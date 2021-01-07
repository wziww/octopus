package main

import (
	"testing"
)

func TestRouterInit(t *testing.T) {
	if len(routerAll) != 18 {
		t.Fatal("router init error")
	}
}
