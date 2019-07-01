package myredis

import (
	"crypto/md5"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/go-redis/redis"
)

type target struct {
	Name        string
	Type        string
	Addrs       []string
	MemoryTotal int64
	Status      string
	self        interface{}
}

var (
	rw sync.RWMutex
)
var redisSources map[string]*target

func init() {
	redisSources = make(map[string]*target)
}

// AddSource 添加监控源
func AddSource(name string, cfg interface{}) error {
	n := fmt.Sprintf("%x", md5.Sum([]byte(name)))
	if redisSources[n] != nil {
		return errors.New(name + " has exists")
	}
	switch cfg.(type) {
	case *redis.Options: // 单机模式
		c := redis.NewClient(cfg.(*redis.Options))
		a, e := c.Ping().Result()
		if e != nil {
			fmt.Println(e)
		}
		if a == "PONG" {
			redisSources[n] = &target{
				Name:  name,
				Type:  "single",
				Addrs: []string{cfg.(*redis.Options).Addr},
				self:  c,
			}
		}
	case *redis.ClusterOptions: // 集群模式
		c := redis.NewClusterClient(cfg.(*redis.ClusterOptions))
		a, e := c.Ping().Result()
		if e != nil {
			fmt.Println(e)
		}
		if a == "PONG" {
			redisSources[n] = &target{
				Name:  name,
				Type:  "cluster",
				Addrs: cfg.(*redis.ClusterOptions).Addrs,
				self:  c,
			}
		}
	default:
		return errors.New("unknow type to init")
	}
	return nil
}

// GetConfig 获取配置好的数据源
func GetConfig() interface{} {
	rw.RLock()
	for _, v := range redisSources {
		switch v.self.(type) {
		case *redis.Client: // 单机模式
			z := v.self.(*redis.Client)
			str, err := z.Ping().Result()
			if err != nil || str != "PONG" {
				v.Status = "failed"
			} else {
				v.Status = "ok"
			}
		case *redis.ClusterClient: // 集群模式
			z := v.self.(*redis.ClusterClient)
			str, e := z.ClusterInfo().Result()
			if e != nil {
				fmt.Println(e)
			} else {
				for _, x := range strings.Split(str, "\n") {
					if len(x) > 14 && x[:14] == "cluster_state:" {
						v.Status = strings.Replace(x[14:], "\r", "", -1)
					}
				}
			}
		}
	}
	rw.RUnlock()
	return redisSources
}

// DetailResult ...
type DetailResult struct {
	ID                string
	ADDR              string
	FOLLOW            string
	ROLE              string
	EPOTH             string
	STATE             string
	SLOT              string
	TotalSystemMemory string
	Maxmemory         string
	UsedMemory        string
}

// GetDetail 获取节点详情
func GetDetail(id string) []*DetailResult {
	if redisSources[id] == nil {
		return nil
	}
	switch redisSources[id].Type {

	case "single":
		z := redisSources[id].self.(*redis.Client)
		str, _ := z.Info("memory").Result()
		strArr := strings.Split(str, "\n")
		v := &DetailResult{
			ADDR: z.Options().Addr,
		}
		for _, z := range strArr {
			if len(z) > len("used_memory:") && z[:len("used_memory:")] == "used_memory:" {
				v.UsedMemory = z[len("used_memory:"):]
				continue
			}
			if len(z) > len("total_system_memory:") && z[:len("total_system_memory:")] == "total_system_memory:" {
				v.TotalSystemMemory = z[len("total_system_memory:"):]
				continue
			}
			if len(z) > len("maxmemory:") && z[:len("maxmemory:")] == "maxmemory:" {
				v.Maxmemory = z[len("maxmemory:"):]
				continue
			}
		}
		return []*DetailResult{v}
	case "cluster":
		z := redisSources[id].self.(*redis.ClusterClient)
		str, e := z.ClusterNodes().Result()
		if e != nil {
			fmt.Println(e)
		} else {
			var result []*DetailResult
			for _, x := range strings.Split(str, "\n") {
				arr := strings.Split(x, " ")
				if len(arr) < 3 {
					continue
				}
				ID := arr[0]
				ADDR := arr[1]
				ROLE := arr[2]
				FOLLOW := arr[3]
				// ping time 4
				// pong time 5
				EPOCH := arr[6]
				STATE := arr[7]
				var slot string
				if len(arr) > 8 {
					index := 8
					for {
						if index > (len(arr) - 1) {
							break
						} else {
							slot += arr[index] + " "
							index++
						}
					}
				}
				result = append(result, &DetailResult{
					ID:     ID,
					ADDR:   ADDR,
					FOLLOW: FOLLOW,
					ROLE:   ROLE,
					EPOTH:  EPOCH,
					STATE:  STATE,
					SLOT:   slot,
				})
			}
			z.ForEachNode(func(c *redis.Client) error {
				str, _ := c.Info("memory").Result()
				for _, v := range result {
					if v.ADDR == c.Options().Addr {
						strArr := strings.Split(str, "\n")
						for _, z := range strArr {
							if len(z) > len("used_memory:") && z[:len("used_memory:")] == "used_memory:" {
								v.UsedMemory = z[len("used_memory:"):]
								continue
							}
							if len(z) > len("total_system_memory:") && z[:len("total_system_memory:")] == "total_system_memory:" {
								v.TotalSystemMemory = z[len("total_system_memory:"):]
								continue
							}
							if len(z) > len("maxmemory:") && z[:len("maxmemory:")] == "maxmemory:" {
								v.Maxmemory = z[len("maxmemory:"):]
								continue
							}
						}
					}
				}
				return nil
			})
			return result
		}
		return nil
	}
	return nil
}

type Stats struct {
	InstantaneousInputKbps  string
	InstantaneousOutputKbps string
	InstantaneousOpsPerSec  string
	ADDR                    string
}

// GetSTATS 获取节点详情
// redis  info STATS
func GetSTATS(id string) []*Stats {
	if redisSources[id] == nil {
		return nil
	}
	var result []*Stats
	var (
		mutex sync.Mutex
	)
	switch redisSources[id].Type {

	case "single":
		z := redisSources[id].self.(*redis.Client)
		str, _ := z.Info("stats").Result()
		strArr := strings.Split(str, "\n")
		v := &Stats{
			ADDR: z.Options().Addr,
		}
		for _, z := range strArr {
			if len(z) > len("instantaneous_input_kbps:") && z[:len("instantaneous_input_kbps:")] == "instantaneous_input_kbps:" {
				v.InstantaneousInputKbps = z[len("instantaneous_input_kbps:"):]
				continue
			}
			if len(z) > len("instantaneous_ops_per_sec:") && z[:len("instantaneous_ops_per_sec:")] == "instantaneous_ops_per_sec:" {
				v.InstantaneousOpsPerSec = z[len("instantaneous_ops_per_sec:"):]
				continue
			}
			if len(z) > len("InstantaneousOutputKbps:") && z[:len("InstantaneousOutputKbps:")] == "InstantaneousOutputKbps:" {
				v.InstantaneousOutputKbps = z[len("InstantaneousOutputKbps:"):]
				continue
			}
		}
		return []*Stats{v}
	case "cluster":
		z := redisSources[id].self.(*redis.ClusterClient)
		z.ForEachNode(func(c *redis.Client) error {
			str, e := c.Info("stats").Result()
			if e != nil {
				fmt.Println(e)
				return nil
			}
			v := &Stats{}
			for _, z := range strings.Split(str, "\n") {
				v.ADDR = c.Options().Addr
				if len(z) > len("instantaneous_input_kbps:") && z[:len("instantaneous_input_kbps:")] == "instantaneous_input_kbps:" {
					v.InstantaneousInputKbps = z[len("instantaneous_input_kbps:"):]
					continue
				}
				if len(z) > len("instantaneous_ops_per_sec:") && z[:len("instantaneous_ops_per_sec:")] == "instantaneous_ops_per_sec:" {
					v.InstantaneousOpsPerSec = z[len("instantaneous_ops_per_sec:"):]
					continue
				}
				if len(z) > len("instantaneous_output_kbps:") && z[:len("instantaneous_output_kbps:")] == "instantaneous_output_kbps:" {
					v.InstantaneousOutputKbps = z[len("instantaneous_output_kbps:"):]
					continue
				}
			}
			mutex.Lock()
			result = append(result, v)
			mutex.Unlock()
			return nil
		})
		return result
	}
	return nil
}

// RemoveSource 节点配置移除
func RemoveSource(id string) error {
	if redisSources[id] == nil {
		return nil
	}
	switch redisSources[id].self.(type) {
	case *redis.ClusterClient:
		redisSources[id].self.(*redis.ClusterClient).Close()
		delete(redisSources, id)
	}
	return nil
}
