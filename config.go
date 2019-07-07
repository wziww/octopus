package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"octopus/myredis"
	"os"

	"github.com/go-redis/redis"

	"github.com/BurntSushi/toml"
)

// C 全局配置
var C *Config

// Config 配置
type Config struct {
	Server *ConfigServer `toml:"server"`
	Redis  []RedisDetail `toml:"redis"`
}

// ConfigServer 服务端配置
type ConfigServer struct {
	ListenAddress string `toml:"listen_address"`
	PidFile       string `toml:"pid_file"`
}

// RedisDetail ...
type RedisDetail struct {
	Name    string   `toml:"name"`
	Address []string `toml:"address"`
	DB      int      `toml:"db"`
}

func init() {
	configPath := flag.String("c", "./conf/server.conf.toml", "config path")
	flag.Parse()
	if C == nil {
		fdata, openError := ioutil.ReadFile(*configPath)
		if openError != nil {
			fmt.Fprintf(os.Stderr, "%s", openError.Error())
		}
		C = &Config{}
		toml.Decode(string(fdata), C)
	}
	for _, v := range C.Redis {
		myredis.AddSource(v.Name, &redis.Options{
			Addr: v.Address[0],
		})
	}
}
