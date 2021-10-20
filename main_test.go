package main

import (
	"octopus/config"
	"octopus/log"
	"octopus/message"
	"octopus/myredis"
	"octopus/permission"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Setenv("CONFIG_FILE", "./conf/test.conf")
	config.Init()
	log.Init()
	message.Init()
	permission.Init()
	RouterInit()
	myredis.Init()
	m.Run()
}
