package myredis

import (
	"crypto/md5"
	"crypto/tls"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

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

// Stats redis info stats
type Stats struct {
	InstantaneousInputKbps  string
	InstantaneousOutputKbps string
	InstantaneousOpsPerSec  string
	ADDR                    string
}

// Memory redis info memory
type Memory struct {
	TotalSystemMemory string
	Maxmemory         string
	UsedMemory        string
}

var (
	rw sync.RWMutex
)
var redisSources map[string]*target

func init() {
	redisSources = make(map[string]*target)
}

var (
	// REDISALREADYEXISTS 该配置已存在
	REDISALREADYEXISTS = -1
	// REDISINITERROR 节点启动失败
	REDISINITERROR = -2
)

// AddSource 添加监控源
func AddSource(name string, opt *redis.Options) int {
	opt.Dialer = func() (net.Conn, error) {
		netDialer := &net.Dialer{
			Timeout:   opt.DialTimeout,
			KeepAlive: 5 * time.Minute,
		}
		if opt.TLSConfig == nil {
			return netDialer.Dial(opt.Network, opt.Addr)
		}
		return tls.DialWithDialer(netDialer, opt.Network, opt.Addr, opt.TLSConfig)
	}
	n := fmt.Sprintf("%x", md5.Sum([]byte(name)))
	REDISTYPE := "single"
	var c interface{}
	if redisSources[n] != nil {
		return REDISALREADYEXISTS
	}
	c = redis.NewClient(opt)
	clusterInfoStr, e := c.(*redis.Client).Info("Cluster").Result()
	if e != nil {
		return REDISINITERROR
	}
	var pingStr string
	var pingError error
	for _, v := range strings.Split(clusterInfoStr, "\n") {
		if len(v) > len("cluster_enabled:") && v[:len("cluster_enabled:")] == "cluster_enabled:" &&
			strings.Replace(v[len("cluster_enabled:"):], "\r", "", -1) == "1" {
			REDISTYPE = "cluster"
			c.(*redis.Client).Close()
			c = redis.NewClusterClient(&redis.ClusterOptions{
				Addrs: []string{opt.Addr},
				// Dialer: opt.Dialer,
			})
			goto finish
		}
	}
finish:
	if REDISTYPE == "cluster" {
		pingStr, pingError = c.(*redis.ClusterClient).Ping().Result()
	} else {
		pingStr, pingError = c.(*redis.Client).Ping().Result()
	}
	if pingError != nil {
		fmt.Println(pingError)
	}
	if pingStr == "PONG" {
		rw.Lock()
		redisSources[n] = &target{
			Name:  name,
			Type:  REDISTYPE,
			Addrs: []string{opt.Addr},
			self:  c,
		}
		rw.Unlock()
	}
	return 0
}

// GetConfig 获取配置好的数据源
func GetConfig() interface{} {
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
	return redisSources
}

// DetailResult ...
type DetailResult struct {
	ID     string
	ADDR   string
	FOLLOW string
	ROLE   string
	EPOTH  string
	STATE  string
	SLOT   string
	Type   string
	Memory
}

// GetDetail 获取节点详情
func GetDetail(id string) []*DetailResult {
	if redisSources[id] == nil {
		return nil
	}
	switch redisSources[id].self.(type) {
	case *redis.Client:
		z := redisSources[id].self.(*redis.Client)
		str, _ := z.Info("memory").Result()
		strArr := strings.Split(str, "\n")
		v := &DetailResult{
			ADDR: z.Options().Addr,
			Type: "single",
		}
		if len(strArr) > 5 {
			v.STATE = "connected"
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
	case *redis.ClusterClient:
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
					Type:   "cluster",
				})
			}
			z.ForEachNode(func(c *redis.Client) error {
				str, _ := c.Info("memory").Result()
				for _, v := range result {
					oaddr := c.Options().Addr
					if (len(v.ADDR) >= len(oaddr)) && v.ADDR[:len(oaddr)] == oaddr {
						strArr := strings.Split(str, "\r")
						for _, z := range strArr {
							z = strings.Replace(z, "\n", "", -1)
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
	case *redis.Client:
		redisSources[id].self.(*redis.Client).Close()
	case *redis.ClusterClient:
		redisSources[id].self.(*redis.ClusterClient).Close()
	}
	rw.Lock()
	delete(redisSources, id)
	rw.Unlock()
	return nil
}
