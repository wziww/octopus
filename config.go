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
	Server *ConfigServer            `toml:"server"`
	Redis  map[string][]RedisDetail `toml:"redis"`
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
}

func init() {
	configPath := flag.String("c", "./conf/server.conf.toml", "config path")
	if C == nil {
		fdata, openError := ioutil.ReadFile(*configPath)
		if openError != nil {
			fmt.Fprintf(os.Stderr, "%s", openError.Error())
		}
		C = &Config{}
		toml.Decode(string(fdata), C)
	}
	if C.Redis["cluster"] != nil {
		for _, v := range C.Redis["cluster"] {
			myredis.AddSource(v.Name, &redis.ClusterOptions{
				Addrs: v.Address,
			})
		}
	}
}
