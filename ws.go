package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"octopus/message"
	"octopus/permission"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// 协议升级
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// websocket 统一数据合法返回结构
type socketReturn struct {
	Type string
	Data string
}

// octopus 内部使用 socket 结构体
type oSocket struct {
	mutex     sync.Mutex
	conn      *websocket.Conn
	user      *permission.User
	namespace string
	connID    string
}

// 当前所有连接
type oSocketAll struct {
	mutex sync.Mutex
	conns []*oSocket
}

func (osall *oSocketAll) remove(id string) {
	osall.mutex.Lock()
	defer osall.mutex.Unlock()
	for i, v := range osall.conns {
		if v.connID == id {
			v.conn.Close()
			osall.conns = append(osall.conns[:i], osall.conns[i+1:]...)
			break
		}
	}
}

var socketAll *oSocketAll

func init() {
	socketAll = &oSocketAll{}
}

// SafeWrite ws 并发控制写入
func SafeWrite(c *oSocket, result []byte, messageType ...int) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	mt := websocket.TextMessage
	if len(messageType) > 0 {
		mt = messageType[0]
	}
	err := c.conn.WriteMessage(mt, result)
	return err
}

func ws(w http.ResponseWriter, r *http.Request) {
	var sc *websocket.Conn
	var err error
	connID := fmt.Sprintf("%x", md5.Sum([]byte(uuid.New().String())))
	sc, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	c := &oSocket{
		conn:   sc,
		connID: connID,
	}
	socketAll.mutex.Lock()
	socketAll.conns = append(socketAll.conns, c)
	socketAll.mutex.Unlock()
	for {
		mt, message, err := sc.ReadMessage()
		if err != nil {
			socketAll.remove(connID)
			break
		}
		go func() {
			result := handle(message, c)
			if result != nil && len(result) > 0 {
				SafeWrite(c, result, mt)
				if err != nil {
					log.Println("write:", err)
				}
			}
		}()
	}
	return
}

func handle(msg []byte, c *oSocket) []byte {
	b := &struct {
		Func string `json:"Func"`
		Data string `json:"Data"`
	}{}
	json.Unmarshal(msg, b)
	routerPath := b.Func
	current := routerAll[routerPath]
	if current != nil {
		var bytes []byte
		if current.permission == 0 {
			goto next
		}
		if c.user == nil || (c.user.Permission&current.permission) != 1 {
			bytes, _ = json.Marshal(&socketReturn{
				Type: routerPath,
				Data: message.Res(403001, ""),
			})
			return bytes
		}
	next:
		bytes, _ = json.Marshal(&socketReturn{
			Type: routerPath,
			Data: current.fn(b.Data, c),
		})
		return bytes
	}
	return []byte("404")
}
