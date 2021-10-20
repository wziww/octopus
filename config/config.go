package config

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
)

const (
	MIN_SLOW_LOG_COUNT = 10
)

// C 全局配置
var C *Config

// Config 配置
type Config struct {
	Server        *Server       `toml:"server"`
	ElasticSearch []ESDetail    `toml:"elasticsearch"`
	Redis         []RedisDetail `toml:"redis"`
	RDB           RDB           `toml:"rdb"`
	DB            DB            `toml:"db"`
	Log           *Log
	Auth          []Auth     `toml:"auth"`
	AuthConfig    AuthConfig `toml:"auth-config"`
	Robin         *Robin     `toml:"robin"`
	Feishu        *Feishu    `toml:"feishu"`
	Luffy         *Luffy     `toml:"luffy"`
}

// Luffy luffy bot
type Luffy struct {
	Username   string `toml:"username"`
	Password   string `toml:"password"`
	APIAddress string `toml:"api_address"`
}

// Feishu 飞书机器人配置
type Feishu struct {
	AppID          string `toml:"app_id"`
	AppSecret      string `toml:"app_secret"`
	ChatID         string `toml:"chat_id"`
	ESURL          string `toml:"es_url"`
	RedisURLMemory string `toml:"redis_url_memory"`
	RedisURLFlow   string `toml:"redis_url_flow"`
}

// RDB rdb 相关配置
type RDB struct {
	Dir string `toml:"dir"`
}

// DB 落盘配置
type DB struct {
	Address string `toml:"address"`
}

// Robin 配置
type Robin struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
	APIKEY   string `toml:"apikey"`
	Host     string `toml:"host"`
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

// ESDetail ...
type ESDetail struct {
	Name              string   `toml:"name"`
	Address           []string `toml:"address"`
	SnapshotsCPULimit int      `toml:"snapshots_cpu_limit"` // 低于多少 cpu 进行快照
	DataDir           string   `toml:"data_dir"`            // 快照保存目录
}

// RedisDetail ...
type RedisDetail struct {
	Name              string   `toml:"name"`
	Address           []string `toml:"address"`
	DB                int      `toml:"db"`
	Password          string   `toml:"password"`
	SlowLogLimitCount int      `toml:"slowlog_limit"` // >= 多少进行快照
	DataDir           string   `toml:"data_dir"`      // 快照保存目录
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

// Init ...
func Init() {
	configPath := flag.String("c", "", "config path")
	flag.Parse()
	if *configPath == "" {
		*configPath = os.Getenv("CONFIG_FILE")
		if *configPath == "" {
			*configPath = "./conf/server.conf.toml"
		}
	}
	if C == nil {
		fdata, openError := ioutil.ReadFile(*configPath)
		if openError != nil {
			fmt.Fprintf(os.Stderr, "%s\n", openError.Error())
			os.Exit(1)
		}
		C = &Config{}
		toml.Decode(string(fdata), C)
	}
	_, err := os.Stat(C.RDB.Dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
	for _, v := range C.Redis {
		if v.SlowLogLimitCount < MIN_SLOW_LOG_COUNT {
			v.SlowLogLimitCount = MIN_SLOW_LOG_COUNT
		}
	}
}
