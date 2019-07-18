package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"octopus/myredis"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var (
	sendLock sync.Mutex
)

// SafeWrite ws safe write
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
	c, err := upgrader.Upgrade(w, r, nil)
	connID := fmt.Sprintf("%x", md5.Sum([]byte(uuid.New().String())))
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			e := err.(*websocket.CloseError)
			if e.Code == websocket.CloseGoingAway { // conn close
			}
			break
		}
		go func() {
			result := handle(message, connID, c)
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

type socketReturn struct {
	Type string
	Data interface{}
}
type socketRecv struct {
	Func string `json:"Func"`
	Data string `json:"Data"`
}
type newConfig struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Type string `json:"type"`
}

type router func(data string, conns ...*websocket.Conn) []byte

var routerAll map[string]router

func handle(message []byte, connID string, c *websocket.Conn) []byte {
	b := &socketRecv{}
	json.Unmarshal(message, b)
	routerPath := b.Func
	if routerAll[routerPath] != nil {
		return routerAll[routerPath](b.Data, c)
	}
	return []byte("404")
}
func init() {
	routerAll = make(map[string]router)
	Router("/config/redis", func(data string, conns ...*websocket.Conn) []byte {
		bytes, _ := json.Marshal(&socketReturn{
			Type: "/config/redis",
			Data: myredis.GetConfig(),
		})
		return bytes
	})
	Router("/config/redis/slots/migrating", func(data string, conns ...*websocket.Conn) []byte {
		body := &struct {
			ID         string `json:"id"`
			SourceID   string `json:"sourceId"`
			TargetID   string `json:"targetId"`
			SlotsStart int64  `json:"slotsStart"`
			SlotsEnd   int64  `json:"slotsEnd"`
		}{}
		json.Unmarshal([]byte(data), body)
		myredis.ClusterSlotsMigrating(body.ID, body.SourceID, body.TargetID, body.SlotsStart, body.SlotsEnd, func(str string, flag ...int64) {
			if len(conns) > 0 {
				t := "/config/redis/slots/migrating"
				if len(flag) > 0 {
					if flag[0] == 0 {
						t = "/config/redis/slots/migrating/0"
					}
				}
				bts, _ := json.Marshal(&socketReturn{
					Type: t,
					Data: str,
				})
				go SafeWrite(conns[0],
					bts)
			}
		})
		return []byte{}
	})
	Router("/config/redis/clusterSlots", func(data string, conns ...*websocket.Conn) []byte {
		c := &struct {
			ID string `json:"id"`
		}{}
		json.Unmarshal([]byte(data), c)
		bytes, _ := json.Marshal(&socketReturn{
			Type: "/config/redis/clusterSlots",
			Data: myredis.ClusterSlotsStats(c.ID),
		})
		return bytes
	})
	Router("/config/redis/clusterNodes", func(data string, conns ...*websocket.Conn) []byte {
		c := &struct {
			ID string `json:"id"`
		}{}
		json.Unmarshal([]byte(data), c)
		bytes, _ := json.Marshal(&socketReturn{
			Type: "/config/redis/clusterNodes",
			Data: myredis.GetClusterNodes(c.ID),
		})
		return bytes
	})
	Router("/config/redis/setSlots", func(data string, conns ...*websocket.Conn) []byte {
		c := &struct {
			ID    string `json:"id"`
			Host  string `json:"host"`
			Port  string `json:"port"`
			Start int64  `json:"start"`
			End   int64  `json:"end"`
		}{}
		json.Unmarshal([]byte(data), c)
		bytes, _ := json.Marshal(&socketReturn{
			Type: "/config/redis/setSlots",
			Data: myredis.ClusterSlotsSet(c.ID, c.Host, c.Port, c.Start, c.End),
		})
		return bytes
	})
	// Router("/config/redis/delSlots", func(data string) []byte {
	// 	c := &struct {
	// 		ID    string `json:"id"`
	// 		Host  string `json:"host"`
	// 		Port  string `json:"port"`
	// 		Start int64  `json:"start"`
	// 		End   int64  `json:"end"`
	// 	}{}
	// 	json.Unmarshal([]byte(data), c)
	// 	bytes, _ := json.Marshal(&socketReturn{
	// 		Type: "/config/redis/delSlots",
	// 		Data: myredis.ClusterSlotsDel(c.ID, c.Host, c.Port, c.Start, c.End),
	// 	})
	// 	return bytes
	// })
	Router("/config/redis/clusterReplicate", func(data string, conns ...*websocket.Conn) []byte {
		c := &struct {
			ID     string `json:"id"`
			Host   string `json:"host"`
			Port   string `json:"port"`
			NodeID string `json:"nodeid"`
		}{}
		json.Unmarshal([]byte(data), c)
		bytes, _ := json.Marshal(&socketReturn{
			Type: "/config/redis/clusterReplicate",
			Data: myredis.ClusterReplicate(c.ID, c.Host, c.Port, c.NodeID),
		})
		return bytes
	})
	Router("/config/redis/clusterForget", func(data string, conns ...*websocket.Conn) []byte {
		c := &struct {
			ID     string `json:"id"`
			NodeID string `json:"nodeid"`
		}{}
		json.Unmarshal([]byte(data), c)
		bytes, _ := json.Marshal(&socketReturn{
			Type: "/config/redis/clusterForget",
			Data: myredis.ClusterForget(c.ID, c.NodeID),
		})
		return bytes
	})
	Router("/config/redis/clusterMeet", func(data string, conns ...*websocket.Conn) []byte {
		c := &struct {
			ID   string `json:"id"`
			Host string `json:"host"`
			Port string `json:"port"`
		}{}
		json.Unmarshal([]byte(data), c)
		bytes, _ := json.Marshal(&socketReturn{
			Type: "/config/redis/clusterMeet",
			Data: myredis.ClusterMeet(c.ID, c.Host, c.Port),
		})
		return bytes
	})
	Router("/config/redis/detail", func(data string, conns ...*websocket.Conn) []byte {
		c := &struct {
			ID string `json:"id"`
		}{}
		json.Unmarshal([]byte(data), c)
		result := myredis.GetDetail(c.ID)
		bytes, _ := json.Marshal(&socketReturn{
			Type: "/config/redis/detail",
			Data: result,
		})
		return bytes
	})
	Router("/config/redis/del", func(data string, conns ...*websocket.Conn) []byte {
		c := &struct {
			ID string `json:"id"`
		}{}
		json.Unmarshal([]byte(data), c)
		myredis.RemoveSource(c.ID)
		bytes, _ := json.Marshal(&socketReturn{
			Type: "/config/redis/del",
			Data: "success",
		})
		return bytes
	})
	Router("/redis/stats", func(data string, conns ...*websocket.Conn) []byte {
		c := &struct {
			ID string `json:"id"`
		}{}
		json.Unmarshal([]byte(data), c)
		result := myredis.GetStats(c.ID)
		bytes, _ := json.Marshal(&socketReturn{
			Type: "/redis/stats",
			Data: result,
		})
		return bytes
	})
}

// Router ...
func Router(path string, r router) {
	routerAll[path] = r
}
