package main

import (
	"encoding/json"
	"octopus/message"
	"octopus/myredis"
	"octopus/permission"

	"github.com/gorilla/websocket"
)

type routerExec func(data string, conns ...*oSocket) string
type router struct {
	fn         routerExec
	permission int
}

var routerAll map[string]*router

// Router ...
func Router(path string, mode int, r routerExec) {
	routerAll[path] = &router{
		fn:         r,
		permission: mode,
	}
}
func init() {
	routerAll = make(map[string]*router)
	Router("token", permission.PERMISSIONNONE, func(data string, conns ...*oSocket) string {
		body := &struct {
			Token string `json:"token"`
		}{}
		json.Unmarshal([]byte(data), body)
		conns[0].user = permission.Get(body.Token)
		return message.Res(200, "success")
	})
	Router("namespace", permission.PERMISSIONNONE, func(data string, conns ...*oSocket) string {
		body := &struct {
			Namespace string `json:"namespace"`
		}{}
		json.Unmarshal([]byte(data), body)
		conns[0].namespace = body.Namespace
		return message.Res(200, "success")
	})
	Router("/login", permission.PERMISSIONNONE, func(data string, conns ...*oSocket) string {
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
	Router("/opcap", permission.PERMISSIONMONIT, func(data string, conns ...*oSocket) string {
		body := &struct {
			Address string `json:"address"`
		}{}
		json.Unmarshal([]byte(data), body)
		return myredis.OpcapCount(body.Address)
	})
	Router("/redis", permission.PERMISSIONMONIT, func(data string, conns ...*oSocket) string {
		return myredis.GetConfig()
	})
	Router("/redis/slots/migrating", permission.PERMISSIONDEV, func(data string, conns ...*oSocket) string {
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
			for _, v := range socketAll.conns {
				if v.namespace == conns[0].namespace {
					SafeWrite(v, bts, websocket.TextMessage)
				}
			}
		})
		return message.Res(200, "success")
	})
	Router("/redis/clusterSlots", permission.PERMISSIONDEV, func(data string, conns ...*oSocket) string {
		c := &struct {
			ID string `json:"id"`
		}{}
		json.Unmarshal([]byte(data), c)
		return myredis.ClusterSlotsStats(c.ID)
	})
	Router("/redis/clusterNodes", permission.PERMISSIONMONIT, func(data string, conns ...*oSocket) string {
		c := &struct {
			ID string `json:"id"`
		}{}
		json.Unmarshal([]byte(data), c)

		return myredis.GetClusterNodes(c.ID)
	})
	Router("/redis/setSlots", permission.PERMISSIONDEV, func(data string, conns ...*oSocket) string {
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
	Router("/redis/clusterReplicate", permission.PERMISSIONDEV, func(data string, conns ...*oSocket) string {
		c := &struct {
			ID     string `json:"id"`
			Host   string `json:"host"`
			Port   string `json:"port"`
			NodeID string `json:"nodeid"`
		}{}
		json.Unmarshal([]byte(data), c)
		return myredis.ClusterReplicate(c.ID, c.Host, c.Port, c.NodeID)
	})
	Router("/redis/clusterForget", permission.PERMISSIONDEV, func(data string, conns ...*oSocket) string {
		c := &struct {
			ID     string `json:"id"`
			NodeID string `json:"nodeid"`
		}{}
		json.Unmarshal([]byte(data), c)
		return myredis.ClusterForget(c.ID, c.NodeID)
	})
	Router("/redis/clusterMeet", permission.PERMISSIONDEV, func(data string, conns ...*oSocket) string {
		c := &struct {
			ID   string `json:"id"`
			Host string `json:"host"`
			Port string `json:"port"`
		}{}
		json.Unmarshal([]byte(data), c)
		return myredis.ClusterMeet(c.ID, c.Host, c.Port)
	})
	Router("/redis/detail", permission.PERMISSIONMONIT, func(data string, conns ...*oSocket) string {
		c := &struct {
			ID string `json:"id"`
		}{}
		json.Unmarshal([]byte(data), c)
		return myredis.GetDetail(c.ID)
	})
	Router("/redis/stats", permission.PERMISSIONMONIT, func(data string, conns ...*oSocket) string {
		c := &struct {
			ID string `json:"id"`
		}{}
		json.Unmarshal([]byte(data), c)
		return myredis.GetStats(c.ID)
	})
}
