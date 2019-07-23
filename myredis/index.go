package myredis

import (
	"crypto/md5"
	"crypto/tls"
	"fmt"
	"net"
	"octopus/config"
	"octopus/log"
	cluster "octopus/myredis/cluster"
	"strconv"
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

type _redisSources struct {
	RS map[string]*target
	rw sync.RWMutex
}

// 监控的 redis 集合
var redisSources *_redisSources

func (r *_redisSources) Set(k string, t *target) {
	r.rw.Lock()
	defer r.rw.Unlock()
	r.RS[k] = t
}
func (r *_redisSources) Get(k string) *target {
	r.rw.RLock()
	defer r.rw.RUnlock()
	z := r.RS[k]
	return z
}
func (r *_redisSources) Range() map[string]*target {
	r.rw.RLock()
	defer r.rw.RUnlock()
	return r.RS
}
func (r *_redisSources) Delete(k string) {
	r.rw.Lock()
	defer r.rw.Unlock()
	delete(r.RS, k)
}

type target struct {
	Name        string
	Type        string
	Addrs       []string
	MemoryTotal int64
	Status      string
	self        interface{}
}

// Server address => version
type Server struct {
	ADDR         string
	RedisVersion string
}

// Stats redis「info stats」
type Stats struct {
	InstantaneousInputKbps  string
	InstantaneousOutputKbps string
	InstantaneousOpsPerSec  string
	ADDR                    string
}

// Memory redis「info memory」
type Memory struct {
	TotalSystemMemory string
	Maxmemory         string
	UsedMemory        string
}

// DetailResult node detail
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
	redisSources = &_redisSources{
		RS: make(map[string]*target),
	}
	for _, v := range config.C.Redis {
		AddSource(v.Name, &redis.Options{
			Addr: v.Address[0],
		})
	}
}

// strArrToInterface .,..
func strArrToInterface(strArr []string) []interface{} {
	s := make([]interface{}, len(strArr))
	for i, v := range strArr {
		s[i] = v
	}
	return s
}

// Trim ...
func Trim(str string) string {
	return strings.Replace(
		strings.Replace(str, "\r", "", -1),
		"\n", "", -1)
}

// toLines ...
func toLines(str string) []string {
	return strings.Split(str, "\n")
}

// getFromRDSStr get value from redis's return string
func getFromRDSStr(str1, str2 string) string {
	if len(str1) > len(str2) && str1[:len(str2)] == str2 {
		return str1[len(str2):]
	}
	return ""
}

// ClusterMeet ...
func ClusterMeet(id string, host string, port string) string {
	if redisSources.Get(id) == nil || redisSources.Get(id).Type != "cluster" {
		return "error"
	}
	switch redisSources.Get(id).self.(type) {
	case *redis.ClusterClient:
		z := redisSources.Get(id).self.(*redis.ClusterClient)
		str, err := z.ClusterMeet(host, port).Result()
		if err == nil {
			return str
		}
		return err.Error()
	default:
	}
	return "error"
}

// ClusterForget ...
func ClusterForget(id string, nodeid string) string {
	if redisSources.Get(id) == nil || redisSources.Get(id).Type != "cluster" {
		return "error"
	}
	switch redisSources.Get(id).self.(type) {
	case *redis.ClusterClient:
		z := redisSources.Get(id).self.(*redis.ClusterClient)
		str, err := z.ClusterForget(nodeid).Result()
		if err == nil {
			log.FMTLog(log.LOGERROR, err)
			return str
		}
		return err.Error()
	default:
	}
	return "error"
}

// ClusterReplicate ...
func ClusterReplicate(id, host, port, nodeid string) string {
	if redisSources.Get(id) == nil || redisSources.Get(id).Type != "cluster" {
		return "error"
	}
	switch redisSources.Get(id).self.(type) {
	case *redis.ClusterClient:
		tmpClient := redis.NewClient(&redis.Options{
			Addr: host + ":" + port,
		})
		defer tmpClient.Close()
		result, err := tmpClient.ClusterReplicate(nodeid).Result()
		if err == nil {
			log.FMTLog(log.LOGERROR, err)
			return result
		}
		return err.Error()
	default:
	}
	return "error"
}

// ClusterSlotsStats slots 情况
func ClusterSlotsStats(id string) interface{} {
	if redisSources.Get(id) == nil || redisSources.Get(id).Type != "cluster" {
		return []byte("error")
	}
	switch redisSources.Get(id).self.(type) {
	case *redis.ClusterClient:
		z := redisSources.Get(id).self.(*redis.ClusterClient)
		slots, err := z.ClusterSlots().Result()
		if err == nil {
			return slots
		}
		log.FMTLog(log.LOGERROR, err)
		return err.Error()
	default:
	}
	return []byte("error")
}

// ClusterSlotsSet ...
func ClusterSlotsSet(id, host, port string, start, end int64) interface{} {
	if redisSources.Get(id) == nil || redisSources.Get(id).Type != "cluster" {
		return "error"
	}
	switch redisSources.Get(id).self.(type) {
	case *redis.ClusterClient:
		tmpClient := redis.NewClient(&redis.Options{
			Addr: host + ":" + port,
		})
		defer tmpClient.Close()
		result, err := tmpClient.Eval(cluster.AddSlotsLua(start, end), []string{}).Result()
		if err != nil {
			log.FMTLog(log.LOGERROR, err)
			return err.Error()
		}
		return result
	default:
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
	if redisSources.Get(n) != nil {
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
			Trim(v[len("cluster_enabled:"):]) == "1" {
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
		pingStr, pingError = c.(*redis.ClusterClient).Ping().Result()
	} else {
		pingStr, pingError = c.(*redis.Client).Ping().Result()
	}
	if pingError != nil {
		log.FMTLog(log.LOGERROR)
	}
	if pingStr == "PONG" {
		redisSources.Set(n, &target{
			Name:  name,
			Type:  REDISTYPE,
			Addrs: []string{opt.Addr},
			self:  c,
		})
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
	if redisSources.Get(id) == nil {
		return nil
	}
	var result []*Server
	var (
		mutex sync.Mutex
	)
	switch redisSources.Get(id).self.(type) {
	case *redis.Client:
		z := redisSources.Get(id).self.(*redis.Client)
		v := getServer(z)
		return []*Server{v}
	case *redis.ClusterClient:
		z := redisSources.Get(id).self.(*redis.ClusterClient)
		z.ForEachNode(func(c *redis.Client) error {
			v := getServer(c)
			mutex.Lock()
			defer mutex.Unlock()
			result = append(result, v)
			return nil
		})
		return result
	default:
	}
	return nil
}

// GetConfig 获取配置好的数据源
func GetConfig() interface{} {
	for _, v := range redisSources.Range() {
		switch v.self.(type) {
		case *redis.Client: // 单机模式
			z := v.self.(*redis.Client)
			str, err := z.Ping().Result()
			if err != nil || str != "PONG" {
				log.FMTLog(log.LOGERROR, err)
				v.Status = "failed"
			} else {
				v.Status = "ok"
			}
		case *redis.ClusterClient: // 集群模式
			z := v.self.(*redis.ClusterClient)
			str, e := z.ClusterInfo().Result()
			if e != nil {
				log.FMTLog(log.LOGERROR, e)
			} else {
				for _, x := range toLines(str) {
					if value := getFromRDSStr(x, "cluster_state:"); value != "" {
						v.Status = Trim(value)
					}
				}
			}
		}
	}
	return redisSources.Range()
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
	if redisSources.Get(id) == nil {
		return nil
	}
	switch redisSources.Get(id).self.(type) {
	case *redis.Client:
		z := redisSources.Get(id).self.(*redis.Client)
		v := &DetailResult{
			ADDR: z.Options().Addr,
			Type: "single",
		}
		getMemory(z, v)
		return []*DetailResult{v}
	case *redis.ClusterClient:
		z := redisSources.Get(id).self.(*redis.ClusterClient)
		var result []*DetailResult
		var (
			resultAppendLock sync.Mutex
		)
		z.ForEachNode(func(c *redis.Client) error {
			str, nodesError := c.ClusterNodes().Result()
			if nodesError != nil {
				log.FMTLog(log.LOGERROR, nodesError)
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
				break
			}
			for _, v := range result {
				oaddr := c.Options().Addr
				if len(v.ADDR) >= len(oaddr) && len(strings.Split(v.ADDR, oaddr)) > 1 {
					getMemory(c, v)
				}
			}
			resultAppendLock.Unlock()
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
	default:
	}
	return nil
}

// 获取不含 slots 的 master 节点
func getNoSlotsMaster(id string) (rst []*DetailResult) {
	if redisSources.Get(id) == nil || redisSources.Get(id).Type != "cluster" {
		return
	}
	switch redisSources.Get(id).self.(type) {
	case *redis.ClusterClient:
		z := redisSources.Get(id).self.(*redis.ClusterClient)
		cr, _ := z.ClusterNodes().Result()
		crArr := toLines(cr)
		for _, v := range crArr {
			/*
			*continue if v's empty or its not master node or has slots*/
			if len(v) == 0 {
				continue
			}
			vArr := strings.Split(v, " ")
			if len(vArr) != 8 {
				continue
			}
			if strings.IndexAny(vArr[2], "master") == -1 {
				continue
			}
			rst = append(rst, &DetailResult{
				ID:   vArr[0],
				ADDR: vArr[1],
				ROLE: vArr[2],
			})
		}
	default:
	}
	return
}

// ClusterSlotsMigrating slots 迁移
func ClusterSlotsMigrating(id, sourceID, targetID string, slotsStart,
	slotsEnd int64, fn func(string, ...int64)) interface{} {
	if redisSources.Get(id) == nil || redisSources.Get(id).Type != "cluster" {
		return "error"
	}
	switch redisSources.Get(id).self.(type) {
	case *redis.ClusterClient:
		nodesResult := GetDetail(id)
		var (
			sourceNode *DetailResult
			targetNode *DetailResult
		)
		for _, v := range getNoSlotsMaster(id) {
			nodesResult = append(nodesResult, v)
		}
		for _, v := range nodesResult {
			/*
			* must be master node */
			if strings.IndexAny(strings.ToLower(v.ROLE), "master") == -1 {
				continue
			}
			if v.ID == sourceID {
				sourceNode = v
				continue
			}
			if v.ID == targetID {
				targetNode = v
				continue
			}
			if sourceNode != nil && targetNode != nil {
				break
			}
		}
		if sourceNode == nil || targetNode == nil {
			fn("not found source or target")
			return ""
		}
		/*
		* redis v4.x returns the address likes 127.0.0.1:6379@16379
		*       v3.x likes 127.0.0.1:6379 */
		tmpSourceClient := redis.NewClient(&redis.Options{
			Addr: strings.Split(sourceNode.ADDR, "@")[0],
		})
		defer tmpSourceClient.Close()
		tmpTargetClient := redis.NewClient(&redis.Options{
			Addr: strings.Split(targetNode.ADDR, "@")[0], // redis v4.x returns the address likes 127.0.0.1:6379@16379     v3.0 likes 127.0.0.1:6379
		})
		defer tmpTargetClient.Close()
		{
			for i := slotsStart; i <= slotsEnd; i++ {
				step1 := []string{"CLUSTER", "SETSLOT", strconv.FormatInt(i, 10), "MIGRATING", targetID}
				log.FMTLog(log.LOGWARN, strings.Join(step1, " "))
				result, err := tmpSourceClient.Do(strArrToInterface(step1)...).Result()
				if result != nil {
					log.FMTLog(log.LOGWARN, result.(string))
				}
				step2 := []string{"CLUSTER", "SETSLOT", strconv.FormatInt(i, 10), "IMPORTING", sourceID}
				log.FMTLog(log.LOGWARN, strings.Join(step2, " "))
				result2, err2 := tmpTargetClient.Do(strArrToInterface(step2)...).Result()
				if result2 != nil {
					log.FMTLog(log.LOGWARN, result2.(string))
				}
				/*
				* CLUSTER SETSLOT went error  */
				if err != nil || err2 != nil || strings.IndexAny(strings.ToLower(result.(string)), "ok") == -1 || strings.IndexAny(strings.ToLower(result2.(string)), "ok") == -1 {
					if err != nil {
						log.FMTLog(log.LOGERROR, err)
						fn(err.Error())
					}
					if err2 != nil {
						log.FMTLog(log.LOGERROR, err2)
						fn(err2.Error())
					}
					stepErr1 := []string{"CLUSTER", "SETSLOT", strconv.FormatInt(i, 10), "STABLE"}
					log.FMTLog(log.LOGWARN, strings.Join(stepErr1, " "))
					tmpSourceClient.Do(strArrToInterface(stepErr1)...)
					stepErr2 := []string{"CLUSTER", "SETSLOT", strconv.FormatInt(i, 10), "STABLE"}
					log.FMTLog(log.LOGWARN, strings.Join(stepErr2, " "))
					tmpTargetClient.Do(strArrToInterface(stepErr2)...)
					goto fail
				}
				for {
					keys, err := tmpSourceClient.ClusterGetKeysInSlot(int(i), 10).Result()
					if err != nil {
						log.FMTLog(log.LOGERROR, err.Error())
						goto finish
					}
					/*
					* this slot's keys migraton have been finished   */
					if len(keys) == 0 {
						Annouce := []string{"CLUSTER", "SETSLOT", strconv.FormatInt(i, 10), "NODE", targetID}
						log.FMTLog(log.LOGWARN, strings.Join(Annouce, " "))
						tmpSourceClient.Do(strArrToInterface(Annouce)...)
						log.FMTLog(log.LOGWARN, strings.Join(Annouce, " "))
						tmpTargetClient.Do(strArrToInterface(Annouce)...)
						break
					}
					hostNPort := strings.Split(tmpTargetClient.Options().Addr, ":")
					var (
						wg sync.WaitGroup
					)
					for _, v := range keys {
						wg.Add(1)
						go func(v string) {
							log.FMTLog(log.LOGWARN, strings.Join([]string{hostNPort[0], hostNPort[1], v, "0", "10s"}, " "))
							result, err := tmpSourceClient.Migrate(hostNPort[0], hostNPort[1], v, 0, time.Second*10).Result()
							if err != nil {
								log.FMTLog(log.LOGERROR, err)
								fn(err.Error())
							} else {
								log.FMTLog(log.LOGWARN, result)
							}
							wg.Done()
						}(v)
					}
					wg.Wait()
				}
				fn(strconv.FormatInt(slotsEnd-slotsStart, 10)+" "+strconv.FormatInt(i-slotsStart, 10), 0)
			}
		}
	default:
	}
finish:
	fn("slots 迁移完毕")
	return "finish"
fail:
	fn("slots 迁移出错,已终止")
	return "fail"
}

// GetClusterNodes 获取节点详情
func GetClusterNodes(id string) string {
	if redisSources.Get(id) == nil {
		return ""
	}
	switch redisSources.Get(id).Type {
	case "single":
		return ""
	case "cluster":
		z := redisSources.Get(id).self.(*redis.ClusterClient)
		result, _ := z.ClusterNodes().Result()
		return result
	default:
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
		if value := getFromRDSStr(z, "instantaneous_output_kbps:"); value != "" {
			v.InstantaneousOutputKbps = value
			continue
		}
	}
	return v
}

// GetStats 获取节点详情
// redis  info STATS
func GetStats(id string) []*Stats {
	if redisSources.Get(id) == nil {
		return nil
	}
	var result []*Stats
	var (
		mutex sync.Mutex
	)
	switch redisSources.Get(id).Type {
	case "single":
		z := redisSources.Get(id).self.(*redis.Client)
		v := getStats(z)
		return []*Stats{v}
	case "cluster":
		z := redisSources.Get(id).self.(*redis.ClusterClient)
		z.ForEachNode(func(c *redis.Client) error {
			v := getStats(c)
			mutex.Lock()
			result = append(result, v)
			mutex.Unlock()
			return nil
		})
		return result
	default:
	}
	return nil
}

// RemoveSource 节点配置移除
func RemoveSource(id string) error {
	if redisSources.Get(id) == nil {
		return nil
	}
	switch redisSources.Get(id).self.(type) {
	case *redis.Client:
		redisSources.Get(id).self.(*redis.Client).Close()
	case *redis.ClusterClient:
		redisSources.Get(id).self.(*redis.ClusterClient).Close()
	}
	redisSources.Delete(id)
	return nil
}
