package main

import (
	"encoding/json"
	"octopus/message"
	"octopus/myredis"
	"octopus/permission"
	"sync"

	"github.com/gorilla/websocket"
)

type _router struct {
	Path          string
	Permision     int
	HasPermission int
}

func (r *_router) check() int {
	return r.Permision & r.HasPermission
}

type router func(data string, conns ...*websocket.Conn) string

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
	if path == "login" {
		goto next
	}
	if path != "dev" && path != "monit" && path != "exec" || clusterID == "" {
		return false
	}
next:
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

// Router ...
func Router(path string, r router) {
	routerAll[path] = r
}
func init() {
	userConnGroup = &_userConnGroup{}
	userConnGroup.conns = make([]*userConns, 0)
	routerAll = make(map[string]router)
	Router("/login", func(data string, conns ...*websocket.Conn) string {
		body := &struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}{}
		json.Unmarshal([]byte(data), body)
		tok, per := permission.Login(body.Username, body.Password)
		d := &struct {
			Token      string `json:"token"`
			Permission int    `json:"permission"`
		}{
			Token:      tok,
			Permission: per,
		}
		bts, _ := json.Marshal(d)
		return string(bts)
	})
	Router("/redis", func(data string, conns ...*websocket.Conn) string {
		return myredis.GetConfig()
	})
	Router("/redis/slots/migrating", func(data string, conns ...*websocket.Conn) string {
		body := &struct {
			ID         string `json:"id"`
			SourceID   string `json:"sourceId"`
			TargetID   string `json:"targetId"`
			SlotsStart int64  `json:"slotsStart"`
			SlotsEnd   int64  `json:"slotsEnd"`
		}{}
		json.Unmarshal([]byte(data), body)
		myredis.ClusterSlotsMigrating(body.ID, body.SourceID, body.TargetID, body.SlotsStart, body.SlotsEnd, func(str string, flag ...int64) {
			if len(conns) == 0 {
				return
			}
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
			go func() {
				userConnGroup.lock.RLock()
				defer userConnGroup.lock.RUnlock()
				for _, v := range userConnGroup.conns {
					if (v.user.Permission & permission.PERMISSIONDEV) == 0 {
						return
					}
					for _, v2 := range v.conns {
						if v2.path != "dev" || v2.clusterID == body.ID {
							continue
						}
						SafeWrite(v2.conn, bts, websocket.TextMessage)
					}
				}
			}()
		})
		return ""
	})
	Router("/redis/clusterSlots", func(data string, conns ...*websocket.Conn) string {
		c := &struct {
			ID string `json:"id"`
		}{}
		json.Unmarshal([]byte(data), c)
		return myredis.ClusterSlotsStats(c.ID)
	})
	Router("/redis/clusterNodes", func(data string, conns ...*websocket.Conn) string {
		c := &struct {
			ID string `json:"id"`
		}{}
		json.Unmarshal([]byte(data), c)

		return myredis.GetClusterNodes(c.ID)
	})
	Router("/redis/setSlots", func(data string, conns ...*websocket.Conn) string {
		c := &struct {
			ID    string `json:"id"`
			Host  string `json:"host"`
			Port  string `json:"port"`
			Start int64  `json:"start"`
			End   int64  `json:"end"`
		}{}
		json.Unmarshal([]byte(data), c)
		return myredis.ClusterSlotsSet(c.ID, c.Host, c.Port, c.Start, c.End)
	})
	Router("/redis/clusterReplicate", func(data string, conns ...*websocket.Conn) string {
		c := &struct {
			ID     string `json:"id"`
			Host   string `json:"host"`
			Port   string `json:"port"`
			NodeID string `json:"nodeid"`
		}{}
		json.Unmarshal([]byte(data), c)
		return myredis.ClusterReplicate(c.ID, c.Host, c.Port, c.NodeID)
	})
	Router("/redis/clusterForget", func(data string, conns ...*websocket.Conn) string {
		c := &struct {
			ID     string `json:"id"`
			NodeID string `json:"nodeid"`
		}{}
		json.Unmarshal([]byte(data), c)
		return myredis.ClusterForget(c.ID, c.NodeID)
	})
	Router("/redis/clusterMeet", func(data string, conns ...*websocket.Conn) string {
		c := &struct {
			ID   string `json:"id"`
			Host string `json:"host"`
			Port string `json:"port"`
		}{}
		json.Unmarshal([]byte(data), c)
		return myredis.ClusterMeet(c.ID, c.Host, c.Port)
	})
	Router("/redis/detail", func(data string, conns ...*websocket.Conn) string {
		c := &struct {
			ID string `json:"id"`
		}{}
		json.Unmarshal([]byte(data), c)
		return myredis.GetDetail(c.ID)
	})
	Router("/redis/stats", func(data string, conns ...*websocket.Conn) string {
		c := &struct {
			ID string `json:"id"`
		}{}
		json.Unmarshal([]byte(data), c)
		return myredis.GetStats(c.ID)
	})
}
