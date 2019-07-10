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
			result := handle(message, connID)
			sendLock.Lock()
			defer sendLock.Unlock()
			err = c.WriteMessage(mt, result)
			if err != nil {
				log.Println("write:", err)
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

type router func(data string) []byte

var routerAll map[string]router

func handle(message []byte, connID string) []byte {
	b := &socketRecv{}
	json.Unmarshal(message, b)
	routerPath := b.Func
	if routerAll[routerPath] != nil {
		return routerAll[routerPath](b.Data)
	}
	return []byte("404")
}
func init() {
	routerAll = make(map[string]router)
	Router("/config/redis", func(data string) []byte {
		bytes, _ := json.Marshal(&socketReturn{
			Type: "/config/redis",
			Data: myredis.GetConfig(),
		})
		return bytes
	})
	Router("/config/redis/detail", func(data string) []byte {
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
	Router("/config/redis/del", func(data string) []byte {
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
	Router("/redis/stats", func(data string) []byte {
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
