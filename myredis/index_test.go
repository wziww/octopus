package myredis

import (
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
