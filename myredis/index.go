package myredis

import (
	"crypto/md5"
	"crypto/tls"
	"fmt"
	"net"
	"octopus/config"
	"octopus/log"
	"octopus/message"
	cluster "octopus/myredis/cluster"
	"octopus/opcap"
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
	/**/
	typeNone   string = "none"
	typeString string = "string"
	typeList   string = "list"
	typeSet    string = "set"
	typeZset   string = "zset"
	typeHash   string = "hash"
	typeStream string = "stream"
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
	if t != nil {
		log.FMTLog(log.LOGDEBUG, k, t.Addrs, t.Type)
	} else {
		log.FMTLog(log.LOGDEBUG, k, "nil")
	}
	r.RS[k] = t
}
func (r *_redisSources) Get(k string) *target {
	r.rw.RLock()
	defer r.rw.RUnlock()
	z := r.RS[k]
	if z != nil {
		log.FMTLog(log.LOGDEBUG, k, z.Addrs, z.Type)
	} else {
		log.FMTLog(log.LOGDEBUG, k, "nil")
	}
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
	password    string
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
	KeyspaceHits            string
	KeyspaceMisses          string
	//--
	ADDR string
}

// Memory redis「info memory」
type Memory struct {
	TotalSystemMemory string
	Maxmemory         string
	UsedMemory        string
}

// DetailResult node detail
type DetailResult struct {
	ID          string
	ADDR        string
	FOLLOW      string
	ROLE        string
	EPOTH       string
	STATE       string
	SLOT        string
	Type        string
	VERSION     string
	OpcapOnline bool
	Memory
	Stats
}

// Init ...
func Init() {
	redisSources = &_redisSources{
		RS: make(map[string]*target),
	}
	for _, v := range config.C.Redis {
		AddSource(v.Name, &redis.Options{
			Addr:     v.Address[0],
			Password: v.Password,
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

func checkIsCluster(id string) bool {
	ifcluster := !(redisSources.Get(id) == nil || redisSources.Get(id).Type != "cluster")
	if !ifcluster {
		log.FMTLog(log.LOGERROR, id, "MUST BE CLUSTER MODE")
	}
	return ifcluster
}

// ClusterMeet ...
func ClusterMeet(id string, host string, port string) string {
	if !checkIsCluster(id) {
		return message.Res(404001, "error")
	}
	switch redisSources.Get(id).self.(type) {
	case *redis.ClusterClient:
		z := redisSources.Get(id).self.(*redis.ClusterClient)
		str, err := z.ClusterMeet(host, port).Result()
		if err == nil {
			return message.Res(200, str)
		}
		log.FMTLog(log.LOGERROR, err)
		return message.Res(500, err.Error())
	default:
	}
	return message.Res(500, "error")
}

// ClusterForget ...
func ClusterForget(id string, nodeid string) string {
	if !checkIsCluster(id) {
		return message.Res(404001, "error")
	}
	switch redisSources.Get(id).self.(type) {
	case *redis.ClusterClient:
		z := redisSources.Get(id).self.(*redis.ClusterClient)
		err := z.ForEachNode(func(c *redis.Client) error {
			_, err := c.ClusterForget(nodeid).Result()
			if err != nil {
				log.FMTLog(log.LOGWARN, err)
			}
			return nil
		})
		if err == nil {
			return message.Res(200, "ok")
		}
		log.FMTLog(log.LOGERROR, err)
		return message.Res(500, err.Error())
	default:
	}
	return message.Res(500, "error")
}

// ClusterReplicate ...
func ClusterReplicate(id, host, port, nodeid string) string {
	if !checkIsCluster(id) {
		return message.Res(404001, "error")
	}
	c := redisSources.Get(id)
	switch c.self.(type) {
	case *redis.ClusterClient:
		tmpClient := redis.NewClient(&redis.Options{
			Addr:     host + ":" + port,
			Password: c.password,
		})
		defer tmpClient.Close()
		result, err := tmpClient.ClusterReplicate(nodeid).Result()
		if err == nil {
			log.FMTLog(log.LOGERROR, err)
			return message.Res(200, result)
		}

		return message.Res(500, err.Error())
	default:
	}
	return message.Res(500, "error")
}

// ClusterSlotsStats slots 情况
func ClusterSlotsStats(id string) string {
	if !checkIsCluster(id) {
		return message.Res(404001, "error")
	}
	switch redisSources.Get(id).self.(type) {
	case *redis.ClusterClient:
		z := redisSources.Get(id).self.(*redis.ClusterClient)
		slots, err := z.ClusterSlots().Result()
		if err == nil {
			return message.Res(200, slots)
		}
		log.FMTLog(log.LOGERROR, err)
		return message.Res(500, err.Error())
	default:
	}
	return message.Res(500, "error")
}

// ClusterSlotsSet ...
func ClusterSlotsSet(id, host, port string, start, end int64) string {
	if !checkIsCluster(id) {
		return message.Res(404001, "error")
	}
	c := redisSources.Get(id)
	switch c.self.(type) {
	case *redis.ClusterClient:
		tmpClient := redis.NewClient(&redis.Options{
			Addr:     host + ":" + port,
			Password: c.password,
		})
		defer tmpClient.Close()
		result, err := tmpClient.Eval(cluster.AddSlotsLua(start, end), []string{}).Result()
		if err != nil {
			log.FMTLog(log.LOGERROR, err)
			return message.Res(500, err.Error())
		}
		return message.Res(200, result)
	default:
	}
	return message.Res(500, "error")
}

// AddSource 添加监控源
func AddSource(name string, opt *redis.Options) string {
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
		return message.Res(200, REDISALREADYEXISTS)
	}
	c = redis.NewClient(opt)
	clusterInfoStr, e := c.(*redis.Client).Info("Cluster").Result()
	if e != nil {
		return message.Res(200, REDISINITERROR)
	}
	var pingStr string
	var pingError error
	for _, v := range toLines(clusterInfoStr) {
		if len(v) > len("cluster_enabled:") && v[:len("cluster_enabled:")] == "cluster_enabled:" &&
			Trim(v[len("cluster_enabled:"):]) == "1" {
			REDISTYPE = "cluster"
			c.(*redis.Client).Close()
			c = redis.NewClusterClient(&redis.ClusterOptions{
				Addrs:    []string{opt.Addr},
				Password: opt.Password,
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
		log.FMTLog(log.LOGERROR, pingError)
	}
	if pingStr == "PONG" {
		log.FMTLog(log.LOGWARN, opt.Addr, "REDIS JOINED")
		redisSources.Set(n, &target{
			Name:     name,
			Type:     REDISTYPE,
			Addrs:    []string{opt.Addr},
			self:     c,
			password: opt.Password,
		})
	}
	return message.Res(200, 0)
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

func _getServer(id string) []*Server {
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

// GetServer 获取服务信息
// redis info server
func GetServer(id string) string {
	if redisSources.Get(id) == nil {
		return ""
	}
	return message.Res(200, _getServer(id))
}

// GetConfig 获取配置好的数据源
func GetConfig() string {
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
	return message.Res(200, redisSources.Range())
}

func _getStats(z *redis.Client, v *DetailResult) {
	str, _ := z.Info("stats").Result()
	strArr := strings.Split(str, "\n")
	for _, z := range strArr {
		if value := getFromRDSStr(z, "keyspace_hits:"); value != "" {
			v.KeyspaceHits = value
			continue
		}
		if value := getFromRDSStr(z, "keyspace_misses:"); value != "" {
			v.KeyspaceMisses = value
			continue
		}
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
	return
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

// GetDetailObj ...
func GetDetailObj(id string) []*DetailResult {
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
		address := strings.Split(v.ADDR, ":")[0] + ":9712"
		conn, e := opcap.CreateOrGetClient(address)
		if e != nil {
			v.OpcapOnline = false
			log.FMTLog(log.LOGWARN, "opcap connected error")
			log.FMTLog(log.LOGWARN, e.Error())
		} else {
			str := opcap.PING(conn, address)
			if str == "pong" {
				v.OpcapOnline = true
			}
		}
		getMemory(z, v)
		_getStats(z, v)
		servers := _getServer(id)
		for _, z := range servers {
			if len(strings.Split(v.ADDR, z.ADDR)) > 1 {
				v.VERSION = z.RedisVersion
			}
		}
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
					_getStats(c, v)
				}
			}
			resultAppendLock.Unlock()
			return nil
		})
		servers := _getServer(id)
		for _, v := range result {
			for _, z := range servers {
				if len(strings.Split(v.ADDR, z.ADDR)) > 1 {
					v.VERSION = z.RedisVersion
					v.OpcapOnline = false
					if v.ROLE == "myself,master" {
						address := strings.Split(v.ADDR, ":")[0] + ":9712"
						conn, e := opcap.CreateOrGetClient(address)
						if e != nil {
							log.FMTLog(log.LOGWARN, "opcap connected error")
							log.FMTLog(log.LOGWARN, e.Error())
						} else {
							if opcap.PING(conn, address) == "pong" {
								v.OpcapOnline = true
							}
						}
					}
				}
			}
		}
		return result
	default:
	}
	return nil
}

// SlogLog ...
type SlogLog struct {
	Addr  string
	Count int
}

// GetSlowLogObj ...
func GetSlowLogObj(id string) []*SlogLog {
	if redisSources.Get(id) == nil {
		return nil
	}
	switch redisSources.Get(id).self.(type) {
	case *redis.Client:
		return []*SlogLog{}
	case *redis.ClusterClient:
		z := redisSources.Get(id).self.(*redis.ClusterClient)
		var result []*SlogLog
		var (
			resultAppendLock sync.Mutex
		)
		z.ForEachNode(func(c *redis.Client) error {
			count, nodesError := c.Do("SLOWLOG", "LEN").Int()
			if nodesError != nil {
				log.FMTLog(log.LOGERROR, nodesError)
				return nodesError
			}
			c.Do("SLOWLOG", "RESET")
			resultAppendLock.Lock()
			result = append(result, &SlogLog{
				Addr:  c.Options().Addr,
				Count: count,
			})
			resultAppendLock.Unlock()
			return nil
		})
		return result
	default:
	}
	return nil
}

// GetDetail ...
func GetDetail(id string) string {
	return message.Res(200, GetDetailObj(id))
}

// 获取不含 slots 的 master 节点
func getNoSlotsMaster(id string) (rst []*DetailResult) {
	if !checkIsCluster(id) {
		return nil
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
	slotsEnd int64, fn func(string, ...int64)) string {
	if !checkIsCluster(id) {
		return message.Res(404001, "error")
	}
	c := redisSources.Get(id)
	switch c.self.(type) {
	case *redis.ClusterClient:
		nodesResult := GetDetailObj(id)
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
			Addr:     strings.Split(sourceNode.ADDR, "@")[0],
			Password: c.password,
		})
		defer tmpSourceClient.Close()
		tmpTargetClient := redis.NewClient(&redis.Options{
			Addr:     strings.Split(targetNode.ADDR, "@")[0], // redis v4.x returns the address likes 127.0.0.1:6379@16379     v3.0 likes 127.0.0.1:6379
			Password: c.password,
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
					log.FMTLog(log.LOGDEBUG, result2.(string))
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
								log.FMTLog(log.LOGDEBUG, result)
							}
							wg.Done()
						}(v)
					}
					wg.Wait()
				}
				fn(strconv.FormatInt(slotsEnd-slotsStart, 10)+" "+strconv.FormatInt(i-slotsStart, 10)+" "+sourceID+" "+targetID+" "+
					strconv.FormatInt(slotsStart, 10)+" "+strconv.FormatInt(slotsEnd, 10), 0)
			}
		}
	default:
	}
finish:
	fn("slots 迁移完毕")
	return message.Res(200, "finish")
fail:
	fn("slots 迁移出错,已终止")
	return message.Res(500, "fail")
}

// GetClusterNodes 获取节点详情
func GetClusterNodes(id string) string {
	if redisSources.Get(id) == nil {
		return message.Res(404001, "error")
	}
	switch redisSources.Get(id).Type {
	case "cluster":
		z := redisSources.Get(id).self.(*redis.ClusterClient)
		result, _ := z.ClusterNodes().Result()
		return message.Res(200, result)
	default:
	}
	return message.Res(404001, "error")
}
func getStats(z *redis.Client) *Stats {
	str, _ := z.Info("stats").Result()
	strArr := strings.Split(str, "\n")
	v := &Stats{
		ADDR: z.Options().Addr,
	}
	for _, z := range strArr {
		if value := getFromRDSStr(z, "keyspace_hits:"); value != "" {
			v.KeyspaceHits = value
			continue
		}
		if value := getFromRDSStr(z, "keyspace_misses:"); value != "" {
			v.KeyspaceMisses = value
			continue
		}
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

// GetStatsObj 获取节点详情
// redis  info STATS
func GetStatsObj(id string) []*Stats {
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

// GetStats 获取节点详情
// redis  info STATS
func GetStats(id string) string {
	if redisSources.Get(id) == nil {
		return message.Res(404001, "error")
	}
	var result []*Stats
	var (
		mutex sync.Mutex
	)
	switch redisSources.Get(id).Type {
	case "single":
		z := redisSources.Get(id).self.(*redis.Client)
		v := getStats(z)
		return message.Res(200, []*Stats{v})
	case "cluster":
		z := redisSources.Get(id).self.(*redis.ClusterClient)
		z.ForEachNode(func(c *redis.Client) error {
			v := getStats(c)
			mutex.Lock()
			result = append(result, v)
			mutex.Unlock()
			return nil
		})
		return message.Res(200, result)
	default:
	}
	return message.Res(200, nil)
}

// RemoveSource 节点配置移除
func RemoveSource(id string) string {
	if redisSources.Get(id) == nil {
		return message.Res(200, nil)
	}
	switch redisSources.Get(id).self.(type) {
	case *redis.Client:
		redisSources.Get(id).self.(*redis.Client).Close()
	case *redis.ClusterClient:
		redisSources.Get(id).self.(*redis.ClusterClient).Close()
	}
	redisSources.Delete(id)
	return message.Res(200, nil)
}

// CheckConnect 查看嗅探插件是否在线
func CheckConnect(address string) string {
	opcap.CreateOrGetClient(address)
	return message.Res(200, nil)
}

// OpcapCount 嗅探插件统计结果
func OpcapCount(address string) string {
	address += ":9712"
	conn, e := opcap.CreateOrGetClient(address)
	if e != nil {
		log.FMTLog(log.LOGWARN, "opcap connected error")
		log.FMTLog(log.LOGWARN, e.Error())
	} else {
		return message.Res(200, strings.Join(opcap.Count(conn, address), "_"))
	}
	return ""
}

// redis debug module

// DebugHtstats 获取指定 db 详情
// redis  info STATS
func DebugHtstats(address string, db int) string {
	// return message.Res(200, result)
	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return message.Res(400, err.Error())
	}
	tmpClient := redis.NewClient(&redis.Options{
		Network: tcpAddr.Network(),
		Addr:    address,
		DB:      db,
	})
	defer tmpClient.Close()
	result, err := tmpClient.Do("DEBUG", "HTSTATS", strconv.Itoa(db)).String()
	if err != nil {
		return message.Res(500, err.Error())
	}
	return message.Res(200, result)
}

// SafeDel ...
func SafeDel(address, key string, db int, fn func(string, ...int64)) string {
	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		fn(err.Error())
		return message.Res(500, err.Error())
	}
	tmpClient := redis.NewClient(&redis.Options{
		Network: tcpAddr.Network(),
		Addr:    address,
		DB:      db,
	})
	defer tmpClient.Close()
	// check type
	t, err := tmpClient.Type(key).Result()
	if err != nil {
		fn(err.Error())
		return message.Res(500, err.Error())
	}
	log.FMTLog(log.LOGWARN, fmt.Sprintf("try to del key 「%s」 with type of 「%s」", key, t))
	switch t {
	case typeNone:
		return message.Res(200, "empty key, dont need to del")
	case typeString:
		// string 类型暴力删除
		_, err := tmpClient.Del(key).Result()
		if err != nil {
			return message.Res(500, err.Error())
		}
		return message.Res(200, "succes")
	case typeHash:
		n, err := tmpClient.HLen(key).Result()
		if err != nil {
			fn(err.Error())
			return message.Res(500, err.Error())
		}
		fn(fmt.Sprintf("about %d keys in hash key「%s」to del, job starting...", n, key))
		var keys []string
		fn(fmt.Sprintf("%d / %d keys to delete", n, n))
		var start uint64
		for {
			keys, start, err = tmpClient.HScan(key, start, "*", 10).Result()
			if err != nil {
				return message.Res(500, err.Error())
			}
			for _, v := range keys {
				_, err = tmpClient.HDel(key, v).Result()
				if err != nil {
					fn(err.Error())
					return message.Res(500, err.Error())
				}
				log.FMTLog(log.LOGDEBUG, fmt.Sprintf("del 「%s」「%s」", key, v))
			}
			n2, err := tmpClient.HLen(key).Result()
			if err != nil {
				fn(err.Error())
				return message.Res(500, err.Error())
			}
			fn(fmt.Sprintf("%d / %d keys to delete", n2, n))
			select {
			case <-time.After(time.Microsecond * 100):
			}
			if n2 == 0 {
				break
			}
		}
		// 当为 hash type 的时候， 如果内置键删除完毕了，redis 会自动删除 key，无需额外再次操作
		return message.Res(200, "success")
	case typeZset:
		n, err := tmpClient.ZCard(key).Result()
		if err != nil {
			fn(err.Error())
			return message.Res(500, err.Error())
		}
		fn(fmt.Sprintf("about %d keys in hash key「%s」to del, job starting...", n, key))
		var keys []string
		fn(fmt.Sprintf("%d / %d keys to delete", n, n))
		var start uint64
		for {
			keys, start, err = tmpClient.ZScan(key, start, "*", 10).Result()
			if err != nil {
				fn(err.Error())
				return message.Res(500, err.Error())
			}
			for _, v := range keys {
				_, err = tmpClient.ZRem(key, v).Result()
				if err != nil {
					fn(err.Error())
					return message.Res(500, err.Error())
				}
				log.FMTLog(log.LOGDEBUG, fmt.Sprintf("del 「%s」「%s」", key, v))
			}
			n2, err := tmpClient.ZCard(key).Result()
			if err != nil {
				fn(err.Error())
				return message.Res(500, err.Error())
			}
			fn(fmt.Sprintf("%d / %d keys to delete", n2, n))
			select {
			case <-time.After(time.Microsecond * 100):
			}
			if n2 == 0 {
				break
			}
		}
		// 当为 zset 的时候， 如果内置键删除完毕了，redis 会自动删除 key，无需额外再次操作
		return message.Res(200, "success")
	case typeSet:
		n, err := tmpClient.SCard(key).Result()
		if err != nil {
			fn(err.Error())
			return message.Res(500, err.Error())
		}
		fn(fmt.Sprintf("about %d keys in hash key「%s」to del, job starting...", n, key))
		var keys []string
		fn(fmt.Sprintf("%d / %d keys to delete", n, n))
		var start uint64
		for {
			keys, start, err = tmpClient.SScan(key, start, "*", 10).Result()
			if err != nil {
				fn(err.Error())
				return message.Res(500, err.Error())
			}
			for _, v := range keys {
				_, err = tmpClient.SRem(key, v).Result()
				if err != nil {
					fn(err.Error())
					return message.Res(500, err.Error())
				}
				log.FMTLog(log.LOGDEBUG, fmt.Sprintf("del 「%s」「%s」", key, v))
			}
			n2, err := tmpClient.SCard(key).Result()
			if err != nil {
				return message.Res(500, err.Error())
			}
			fn(fmt.Sprintf("%d / %d keys to delete", n2, n))
			select {
			case <-time.After(time.Microsecond * 100):
			}
			if n2 == 0 {
				break
			}
		}
		// 当为 set 的时候， 如果内置键删除完毕了，redis 会自动删除 key，无需额外再次操作
		return message.Res(200, "success")
	case typeList:
		n, err := tmpClient.LLen(key).Result()
		if err != nil {
			fn(err.Error())
			return message.Res(500, err.Error())
		}
		fn(fmt.Sprintf("about %d keys in hash key「%s」to del, job starting...", n, key))
		fn(fmt.Sprintf("%d / %d keys to delete", n, n))
		for {
			for i := 0; i < 10; i++ {
				_, err := tmpClient.LPop(key).Result()
				if err != nil {
					if err == redis.Nil {
						break
					} else {
						fn(err.Error())
						return message.Res(500, err.Error())
					}
				}
			}
			select {
			case <-time.After(time.Microsecond * 100):
			}
			n2, err := tmpClient.LLen(key).Result()
			if err != nil {
				fn(err.Error())
				return message.Res(500, err.Error())
			}
			fn(fmt.Sprintf("%d / %d keys to delete", n2, n))
			if n2 == 0 {
				break
			}
		}
		// 当为 list 的时候， 如果内置键删除完毕了，redis 会自动删除 key，无需额外再次操作
		return message.Res(200, "success")

	default:
		fn("failed unknow type: " + t)
		return message.Res(500, t)
	}
}
