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

var (
	// REDISALREADYEXISTS 该配置已存在
	REDISALREADYEXISTS = -1
	// REDISINITERROR 节点启动失败
	REDISINITERROR = -2
)
var (
	rw sync.RWMutex
)

// 监控的 redis 集合
var redisSources map[string]*target

type target struct {
	Name        string
	Type        string
	Addrs       []string
	MemoryTotal int64
	Status      string
	self        interface{}
}

// Server ...
type Server struct {
	ADDR         string
	RedisVersion string
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

// DetailResult ...
type DetailResult struct {
	ID      string
	ADDR    string
	FOLLOW  string
	ROLE    string
	EPOTH   string
	STATE   string
	SLOT    string
	Type    string
	VERSION string
	Memory
}

func init() {
	redisSources = make(map[string]*target)
}

// trim ...
func trim(str string) string {
	return strings.Replace(
		strings.Replace(str, "\r", "", -1),
		"\n", "", -1)
}

// toLines ...
func toLines(str string) []string {
	return strings.Split(str, "\n")
}

// getFromRDSStr 从 redis string 中解析值
func getFromRDSStr(str1, str2 string) string {
	if len(str1) > len(str2) && str1[:len(str2)] == str2 {
		return str1[len(str2):]
	}
	return ""
}

// ClusterMeet ...
func ClusterMeet(id string, host string, port string) string {
	if redisSources[id] == nil {
		return "error"
	}
	switch redisSources[id].self.(type) {
	case *redis.ClusterClient:
		z := redisSources[id].self.(*redis.ClusterClient)
		str, err := z.ClusterMeet(host, port).Result()
		if err == nil {
			return str
		}
		return err.Error()
	}
	return "error"
}

// ClusterForget ...
func ClusterForget(id string, nodeid string) string {
	if redisSources[id] == nil {
		return "error"
	}
	switch redisSources[id].self.(type) {
	case *redis.ClusterClient:
		z := redisSources[id].self.(*redis.ClusterClient)
		str, err := z.ClusterForget(nodeid).Result()
		if err == nil {
			return str
		}
		return err.Error()
	}
	return "error"
}

// ClusterReplicate ...
func ClusterReplicate(id, host, port, nodeid string) string {
	if redisSources[id] == nil {
		return "error"
	}
	switch redisSources[id].self.(type) {
	case *redis.ClusterClient:
		tmpClient := redis.NewClient(&redis.Options{
			Addr: host + ":" + port,
		})
		defer tmpClient.Close()
		result, err := tmpClient.ClusterReplicate(nodeid).Result()
		if err == nil {
			return result
		}
		return err.Error()
	}
	return "error"
}

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
	for _, v := range toLines(clusterInfoStr) {
		if len(v) > len("cluster_enabled:") && v[:len("cluster_enabled:")] == "cluster_enabled:" &&
			trim(v[len("cluster_enabled:"):]) == "1" {
			REDISTYPE = "cluster"
			c.(*redis.Client).Close()
			c = redis.NewClusterClient(&redis.ClusterOptions{
				Addrs: []string{opt.Addr},
			})
			goto finish
		}
	}
finish:
	if REDISTYPE == "cluster" {
		// zz, ee := c.(*redis.ClusterClient).Get("test").Result()
		// fmt.Println(zz, ee)
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

func getServer(z *redis.Client) *Server {
	str, _ := z.Info("server").Result()
	v := &Server{
		ADDR: z.Options().Addr,
	}
	for _, z := range toLines(str) {
		if value := getFromRDSStr(z, "redis_version:"); value != "" {
			v.RedisVersion = value
			continue
		}
	}
	return v
}

// GetServer 获取服务信息
// redis info server
func GetServer(id string) []*Server {
	if redisSources[id] == nil {
		return nil
	}
	var result []*Server
	var (
		mutex sync.Mutex
	)
	switch redisSources[id].self.(type) {
	case *redis.Client:
		z := redisSources[id].self.(*redis.Client)
		v := getServer(z)
		return []*Server{v}
	case *redis.ClusterClient:
		z := redisSources[id].self.(*redis.ClusterClient)
		z.ForEachNode(func(c *redis.Client) error {
			v := getServer(c)
			mutex.Lock()
			defer mutex.Unlock()
			result = append(result, v)
			return nil
		})
		return result
	}
	return nil
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
				for _, x := range toLines(str) {
					if value := getFromRDSStr(x, "cluster_state:"); value != "" {
						v.Status = trim(value)
					}
				}
			}
		}
	}
	return redisSources
}

func getMemory(z *redis.Client, v *DetailResult) {
	str, _ := z.Info("memory").Result()
	strArr := toLines(str)
	if len(strArr) > 5 {
		v.STATE = "connected"
	}
	for _, z := range strArr {
		if value := getFromRDSStr(z, "used_memory:"); value != "" {
			v.UsedMemory = value
			continue
		}
		if value := getFromRDSStr(z, "total_system_memory:"); value != "" {
			v.TotalSystemMemory = value
			continue
		}
		if value := getFromRDSStr(z, "maxmemory:"); value != "" {
			v.Maxmemory = value
			continue
		}
	}
}

// GetDetail 获取节点详情
func GetDetail(id string) []*DetailResult {
	if redisSources[id] == nil {
		return nil
	}
	switch redisSources[id].self.(type) {
	case *redis.Client:
		z := redisSources[id].self.(*redis.Client)
		v := &DetailResult{
			ADDR: z.Options().Addr,
			Type: "single",
		}
		getMemory(z, v)
		return []*DetailResult{v}
	case *redis.ClusterClient:
		z := redisSources[id].self.(*redis.ClusterClient)
		var result []*DetailResult
		var (
			resultAppendLock sync.Mutex
		)
		z.ForEachNode(func(c *redis.Client) error {
			str, nodesError := c.ClusterNodes().Result()
			if nodesError != nil {
				return nodesError
			}
			for _, x := range toLines(str) {
				arr := strings.Split(x, " ")
				if len(arr) < 3 || strings.Index(x, "myself") == -1 {
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
				resultAppendLock.Lock()
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
				resultAppendLock.Unlock()
			}
			for _, v := range result {
				oaddr := c.Options().Addr
				if len(v.ADDR) > len(oaddr) && len(strings.Split(v.ADDR, oaddr)) > 1 {
					getMemory(c, v)
				}
			}
			return nil
		})
		servers := GetServer(id)
		for _, v := range result {
			for _, z := range servers {
				if len(strings.Split(v.ADDR, z.ADDR)) > 1 {
					v.VERSION = z.RedisVersion
				}
			}
		}
		return result
	}
	return nil
}

// GetClusterNodes 获取节点详情
func GetClusterNodes(id string) string {
	if redisSources[id] == nil {
		return ""
	}
	switch redisSources[id].Type {
	case "single":
		return ""
	case "cluster":
		z := redisSources[id].self.(*redis.ClusterClient)
		result, _ := z.ClusterNodes().Result()
		return result
	}
	return ""
}
func getStats(z *redis.Client) *Stats {
	str, _ := z.Info("stats").Result()
	strArr := strings.Split(str, "\n")
	v := &Stats{
		ADDR: z.Options().Addr,
	}
	for _, z := range strArr {
		if value := getFromRDSStr(z, "instantaneous_input_kbps:"); value != "" {
			v.InstantaneousInputKbps = value
			continue
		}
		if value := getFromRDSStr(z, "instantaneous_ops_per_sec:"); value != "" {
			v.InstantaneousOpsPerSec = value
			continue
		}
		if value := getFromRDSStr(z, "InstantaneousOutputKbps:"); value != "" {
			v.InstantaneousOutputKbps = value
			continue
		}
	}
	return v
}

// GetStats 获取节点详情
// redis  info STATS
func GetStats(id string) []*Stats {
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
		v := getStats(z)
		return []*Stats{v}
	case "cluster":
		z := redisSources[id].self.(*redis.ClusterClient)
		z.ForEachNode(func(c *redis.Client) error {
			v := getStats(c)
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
