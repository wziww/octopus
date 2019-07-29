package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"octopus/message"
	"octopus/myredis"
	"octopus/permission"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type socketReturn struct {
	Type string
	Data string
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
var userConnGroup *_userConnGroup

type _userConnGroup struct {
	conns []*userConns
	lock  sync.RWMutex
}
type userConns struct {
	user  *permission.User
	conns []*userConnType
}

type userConnType struct {
	conn      *websocket.Conn
	path      string
	clusterID string
	id        string
}

func (c *_userConnGroup) checkSet(token string, conn *websocket.Conn, path, id, clusterID string) bool {
	if path != "dev" && path != "monit" && path != "exec" || clusterID == "" {
		return false
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	for _, v := range c.conns {
		if v.user.Token == token {
			if len(v.conns) > 100 {
				return false
			}
			for _, vv := range v.conns {
				if vv.id == id { // Duplicate id
					return false
				}
			}
			v.conns = append(v.conns, &userConnType{
				conn:      conn,
				path:      path,
				id:        id,
				clusterID: clusterID,
			})
			return true
		}
	}
	c.conns = append(c.conns, &userConns{
		user: permission.Get(token),
		conns: []*userConnType{&userConnType{
			path:      path,
			conn:      conn,
			id:        id,
			clusterID: clusterID,
		}},
	})
	return true
}
func (c *_userConnGroup) checkDel(token string, path string, id string) bool {
	for _, v := range c.conns {
		if v.user.Token == token {
			for i, v2 := range v.conns {
				if v2.path == path && v2.id == id {
					c.lock.Lock()
					v.conns = append(v.conns[:i], v.conns[i+1:]...)
					c.lock.Unlock()
					break
				}
			}
		}
	}
	return false
}
func init() {
	userConnGroup = &_userConnGroup{}
	userConnGroup.conns = make([]*userConns, 0)
	routerAll = make(map[string]router)
	Router("/redis", func(data string, conns ...*websocket.Conn) []byte {
		bytes, _ := json.Marshal(&socketReturn{
			Type: "/redis",
			Data: myredis.GetConfig(),
		})
		return bytes
	})
	Router("/redis/slots/migrating", func(data string, conns ...*websocket.Conn) []byte {
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
				t := "/redis/slots/migrating"
				if len(flag) > 0 {
					if flag[0] == 0 {
						t = "/redis/slots/migrating/0"
					}
				}
				bts, _ := json.Marshal(&socketReturn{
					Type: t,
					Data: message.Res(200, str),
				})
				PATH := "dev"
				go func() {
					userConnGroup.lock.RLock()
					defer userConnGroup.lock.RUnlock()
					for _, v := range userConnGroup.conns {
						if (v.user.Permission & permission.PERMISSIONDEV) > 0 {
							for _, v2 := range v.conns {
								if v2.path == PATH && v2.clusterID == body.ID {
									SafeWrite(v2.conn, bts, websocket.TextMessage)
								}
							}
						}
					}
				}()
			}
		})
		return []byte{}
	})
	Router("/redis/clusterSlots", func(data string, conns ...*websocket.Conn) []byte {
		c := &struct {
			ID string `json:"id"`
		}{}
		json.Unmarshal([]byte(data), c)
		bytes, _ := json.Marshal(&socketReturn{
			Type: "/redis/clusterSlots",
			Data: myredis.ClusterSlotsStats(c.ID),
		})
		return bytes
	})
	Router("/redis/clusterNodes", func(data string, conns ...*websocket.Conn) []byte {
		c := &struct {
			ID string `json:"id"`
		}{}
		json.Unmarshal([]byte(data), c)
		bytes, _ := json.Marshal(&socketReturn{
			Type: "/redis/clusterNodes",
			Data: myredis.GetClusterNodes(c.ID),
		})
		return bytes
	})
	Router("/redis/setSlots", func(data string, conns ...*websocket.Conn) []byte {
		c := &struct {
			ID    string `json:"id"`
			Host  string `json:"host"`
			Port  string `json:"port"`
			Start int64  `json:"start"`
			End   int64  `json:"end"`
		}{}
		json.Unmarshal([]byte(data), c)
		bytes, _ := json.Marshal(&socketReturn{
			Type: "/redis/setSlots",
			Data: myredis.ClusterSlotsSet(c.ID, c.Host, c.Port, c.Start, c.End),
		})
		return bytes
	})
	Router("/redis/clusterReplicate", func(data string, conns ...*websocket.Conn) []byte {
		c := &struct {
			ID     string `json:"id"`
			Host   string `json:"host"`
			Port   string `json:"port"`
			NodeID string `json:"nodeid"`
		}{}
		json.Unmarshal([]byte(data), c)
		bytes, _ := json.Marshal(&socketReturn{
			Type: "/redis/clusterReplicate",
			Data: myredis.ClusterReplicate(c.ID, c.Host, c.Port, c.NodeID),
		})
		return bytes
	})
	Router("/redis/clusterForget", func(data string, conns ...*websocket.Conn) []byte {
		c := &struct {
			ID     string `json:"id"`
			NodeID string `json:"nodeid"`
		}{}
		json.Unmarshal([]byte(data), c)
		bytes, _ := json.Marshal(&socketReturn{
			Type: "/redis/clusterForget",
			Data: myredis.ClusterForget(c.ID, c.NodeID),
		})
		return bytes
	})
	Router("/redis/clusterMeet", func(data string, conns ...*websocket.Conn) []byte {
		c := &struct {
			ID   string `json:"id"`
			Host string `json:"host"`
			Port string `json:"port"`
		}{}
		json.Unmarshal([]byte(data), c)
		bytes, _ := json.Marshal(&socketReturn{
			Type: "/redis/clusterMeet",
			Data: myredis.ClusterMeet(c.ID, c.Host, c.Port),
		})
		return bytes
	})
	Router("/redis/detail", func(data string, conns ...*websocket.Conn) []byte {
		c := &struct {
			ID string `json:"id"`
		}{}
		json.Unmarshal([]byte(data), c)
		bytes, _ := json.Marshal(&socketReturn{
			Type: "/redis/detail",
			Data: myredis.GetDetail(c.ID),
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
			Data: message.Res(200, "success"),
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
	var c *websocket.Conn
	var err error
	params := r.URL.Query()
	tokens := params["octopusToken"]
	var token string
	if len(tokens) > 0 {
		token = tokens[0]
	}
	var path, clusterID string
	p := params["octopusPath"]
	if len(p) > 0 {
		path = p[0]
	}
	oc := params["octopusClusterID"]
	if len(oc) > 0 {
		clusterID = oc[0]
	}
	connID := fmt.Sprintf("%x", md5.Sum([]byte(uuid.New().String())))
	if token != "" {
		u := permission.Get(token)
		if u != nil {
			c, err = upgrader.Upgrade(w, r, nil)
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
			return
		}
	} else {
		return
	}
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			userConnGroup.checkDel(token, path, connID)
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

func handle(message []byte, connID string, c *websocket.Conn) []byte {
	b := &socketRecv{}
	json.Unmarshal(message, b)
	routerPath := b.Func
	if routerAll[routerPath] != nil {
		return routerAll[routerPath](b.Data, c)
	}
	return []byte("404")
}

// Router ...
func Router(path string, r router) {
	routerAll[path] = r
}
