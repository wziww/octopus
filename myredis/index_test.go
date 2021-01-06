package myredis

import (
	"crypto/md5"
	"fmt"
	"octopus/config"
	"octopus/log"
	"octopus/message"
	"octopus/permission"
	"os"
	"strings"
	"testing"
	"time"
)

func TestGetFromRDSStr(t *testing.T) {
	if getFromRDSStr("cluster_state:ok", "cluster_state:") != "ok" {
		t.Fatalf("%s", "cluster_state:ok, cluster_state:, getFromRDSStr error")
	}
	if getFromRDSStr("cluster_state:ok", "cluster_state:okk") != "" {
		t.Fatalf("%s", "cluster_state:ok, cluster_state:, getFromRDSStr error")
	}
}

func TestTolines(t *testing.T) {
	if len(toLines(`1
	2
	3
	4
	5`)) != 5 {
		t.Fatal("toLines error")
	}
}
func TestTrim(t *testing.T) {
	if Trim("s\rsdfsdf\r\n\asd\n") != "ssdfsdfsd" {
		t.Fatal("Trim error")
	}
}

func TestStrArrToInterface(t *testing.T) {
	if len(strArrToInterface([]string{"a", "b", "c"})) != 3 {
		t.Fatal("strArrToInterface error")
	}
}

var (
	nodeNum int = 6
)

func TestCluster(t *testing.T) {
	os.Setenv("CONFIG_FILE", "../conf/test.conf")
	config.Init()
	log.Init()
	message.Init()
	permission.Init()
	Init()
	testID := fmt.Sprintf("%x", md5.Sum([]byte(config.C.Redis[0].Name)))
	if !checkIsCluster(testID) {
		t.Fatalf("%s\n", "redis-server not running in cluster mod")
	}
	if len(_getServer(testID)) != nodeNum {
		t.Fatalf("%s\n", "_getServer error")
	}
	if len(_getServer("testID")) != 0 {
		t.Fatalf("%s\n", "_getServer error")
	}
	if GetServer(testID) == "" {
		t.Fatalf("%s\n", "GetServer error")
	}
	if GetServer("testID") != "" {
		t.Fatalf("%s\n", "GetServer error")
	}
	if ClusterSlotsStats(testID) == "" {
		t.Fatalf("%s\n", "ClusterSlotsStats error")
	}
	if len(GetDetailObj(testID)) != nodeNum {
		t.Fatalf("%s\n", "GetDetailObj error")
	}
	z := GetDetailObj(testID)
	var tmp *DetailResult
	for _, v := range z {
		if strings.Index(v.ROLE, "master") == -1 { // 使用 slave 节点进行单元测试
			tmp = v
			goto next
		}
	}
	t.Fatal("cant find slave node to forget")
next:
	t.Log(ClusterForget(testID, tmp.ID), "\n")
	if len(GetDetailObj(testID)) != nodeNum-1 {
		t.Fatalf("%s\n", "ClusterForget error")
	}
	addr := strings.Split(tmp.ADDR, ":")
	t.Log(ClusterMeet(testID, addr[0], strings.Split(addr[1], "@")[0]), "\n") // 高版本兼容
	select {
	case <-time.After(time.Second * 5): // wait cluster info sync
	}
	if len(GetDetailObj(testID)) != nodeNum {
		t.Fatalf("%s\n", "ClusterMeet error")
	}
	// 集群加入非法节点测试
	if ClusterMeet("t", "", "") != message.Res(404001, "error") {
		t.Fatal("ClusterMeet 非法集群操作异常")
	}
	ClusterMeet(testID, "0.0.0.0", "9999")
	if len(GetDetailObj(testID)) != nodeNum {
		t.Fatal("ClusterMeet 非法节点操作异常")
	}
}
