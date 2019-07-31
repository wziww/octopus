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

var (
	sendLock sync.Mutex
)

// SafeWrite ws 并发控制写入
func SafeWrite(c *websocket.Conn, result []byte, messageType ...int) error {
	sendLock.Lock()
	defer sendLock.Unlock()
	mt := websocket.TextMessage
	if len(messageType) > 0 {
		mt = messageType[0]
	}
	err := c.WriteMessage(mt, result)
	return err
}

func ws(w http.ResponseWriter, r *http.Request) {
	var c *websocket.Conn
	var err error
	params := r.URL.Query()
	tokens := params["ot"] // octopus token
	var token string
	if len(tokens) > 0 {
		token = tokens[0]
	}
	if token == "" {
		w.WriteHeader(403)
	}
	var path, clusterID string
	p := params["op"] // octopus path
	if len(p) > 0 {
		path = p[0]
	}
	oc := params["ocid"] // octopus cluster id
	if len(oc) > 0 {
		clusterID = oc[0]
	}
	connID := fmt.Sprintf("%x", md5.Sum([]byte(uuid.New().String())))

	u := permission.Get(token)
	if u != nil {
		c, err = upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		if userConnGroup.checkSet(token, c, path, connID, clusterID) {
		} else {
			bytes, _ := json.Marshal(&socketReturn{
				Type: "",
				Data: message.Res(404002, "error"),
			})
			SafeWrite(c, bytes, websocket.TextMessage)
			return
		}
	} else {
		w.WriteHeader(403)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			userConnGroup.checkDel(token, path, connID)
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

func handle(message []byte, c *websocket.Conn) []byte {
	b := &struct {
		Func string `json:"Func"`
		Data string `json:"Data"`
	}{}
	json.Unmarshal(message, b)
	routerPath := b.Func
	if routerAll[routerPath] != nil {
		bytes, _ := json.Marshal(&socketReturn{
			Type: routerPath,
			Data: routerAll[routerPath](b.Data, c),
		})
		return bytes
	}
	return []byte("404")
}
