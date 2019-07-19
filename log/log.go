package log

import (
	"fmt"
	"octopus/config"
	"os"
	"path"
	"time"
)

var (
	// LOGNONE 禁用日志
	LOGNONE = 1 << 0
	// LOGERROR 错误级别
	LOGERROR = 1 << 1
	// LOGWARN 警告级别
	LOGWARN = 1 << 2
	// LOGDEBUG 调试级别
	LOGDEBUG = 1 << 3
)
var logLevel int
var fd *os.File

func init() {
	logLevel = LOGNONE
	dir := config.C.Log.LogPath
	var err error
	if dir != "" {
		fd, err = os.OpenFile(path.Join(dir, time.Now().Format("20060102")), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0777)
	} else {
		fd = os.Stdout
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// SetLogLevel 设置日志级别，多次调用权限叠加
func SetLogLevel(i int) {
	logLevel |= i
}

// FMTLog ...
func FMTLog(level int, strs ...interface{}) {
	if (logLevel & level) > 0 {
		_, e := fmt.Fprintln(fd, strs...)
		if e != nil {
			fmt.Println(e)
		}
	}
}
