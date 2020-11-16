package config

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
)

// C 全局配置
var C *Config

// Config 配置
type Config struct {
	Server     *Server       `toml:"server"`
	Redis      []RedisDetail `toml:"redis"`
	RDB        RDB           `toml:"rdb"`
	DB         DB            `toml:"db"`
	Log        *Log
	Auth       []Auth     `toml:"auth"`
	AuthConfig AuthConfig `toml:"auth-config"`
	Opcap      *Opcap     `toml:"opcap"`
}

// RDB rdb 相关配置
type RDB struct {
	Dir string `toml:"dir"`
}

// DB 落盘配置
type DB struct {
	Address string `toml:"address"`
}

// Log 日志配置
type Log struct {
	LogPath  string   `toml:"log_path"`
	LogLevel []string `toml:"log_level"`
}

// Server 服务端配置
type Server struct {
	ListenAddress string `toml:"listen_address"`
	PidFile       string `toml:"pid_file"`
}

// RedisDetail ...
type RedisDetail struct {
	Name     string   `toml:"name"`
	Address  []string `toml:"address"`
	DB       int      `toml:"db"`
	Password string   `toml:"password"`
}

// Auth ...
type Auth struct {
	User       string   `toml:"user"`
	Password   string   `toml:"password"`
	Permission []string `toml:"permission"`
}

// AuthConfig ...
type AuthConfig struct {
	Key string `toml:"key"`
}

// Opcap ...
type Opcap struct {
	Device    string `toml:"device"`
	BPFFilter string `toml:"BPFFilter"`
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
	_, err := os.Stat(C.RDB.Dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err.Error())
	}
}
