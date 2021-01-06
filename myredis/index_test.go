package myredis

import "testing"

func TestGetFromRDSStr(t *testing.T) {
	if getFromRDSStr("cluster_state:ok", "cluster_state:") != "ok" {
		t.Fatalf("%s", "cluster_state:ok, cluster_state:, getFromRDSStr error")
	}
}
