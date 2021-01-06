package myredis

import (
	"crypto/md5"
	"fmt"
	"octopus/config"
	"octopus/log"
	"octopus/message"
	"octopus/permission"
	"os"
	"testing"
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
	if len(_getServer(testID)) != 6 {
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
	if len(GetDetailObj(testID)) != 6 {
		t.Fatalf("%s\n", "GetDetailObj error")
	}
}
