package message

import (
	"testing"
)

func TestRes(t *testing.T) {
	Init()
	if Res(200, "test message") != "{\"code\":200,\"message\":\"test message\"}" {
		t.Fatal("Res test error")
	}
	if Res(200, struct {
		C map[string]string
	}{
		C: map[string]string{
			"a": "b",
		},
	}) != "{\"code\":200,\"message\":{\"C\":{\"a\":\"b\"}}}" {
		t.Fatal("Res test error")
	}
	if Res(404001, "test") != "{\"code\":404001,\"message\":\"非集群模式不可操作\"}" {
		t.Fatal("Res test error")
	}
}
