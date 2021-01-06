package log

import (
	"fmt"
	"octopus/config"
	"os"
	"path"
	"runtime/debug"
	"sync"
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

type _f struct {
	fd     *os.File
	fdlock sync.RWMutex
}

var f *_f

func (file *_f) Set(fd *os.File) {
	file.fdlock.Lock()
	defer file.fdlock.Unlock()
	file.fd = fd
}
func (file *_f) Print(strs ...interface{}) {
	strs = append(make([]interface{}, 1), strs...)
	strs[0] = time.Now().Format("2006-01-02 15:04:05") + ":"
	file.fdlock.Lock()
	defer file.fdlock.Unlock()
	_, e := fmt.Fprintln(file.fd, strs...)
	if e != nil {
		fmt.Println(e)
	}
}

// Init ...
func Init() {
	f = &_f{}
	logLevel = LOGNONE
	dir := config.C.Log.LogPath
	LOGFILETIME := time.Now().Format("2006010215")
	if dir != "" {
		fd, err := os.OpenFile(path.Join(dir, LOGFILETIME), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		f.Set(fd)
		go func() {
			for {
				select {
				case <-time.After(time.Second * 10):
					CURRENTTIME := time.Now().Format("2006010215")
					if CURRENTTIME != LOGFILETIME {
						fd, err := os.OpenFile(path.Join(dir, CURRENTTIME), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
						f.Set(fd)
						if err != nil {
							fmt.Println(err)
							os.Exit(1)
						}
						LOGFILETIME = CURRENTTIME
					}
				}
			}
		}()
	} else {
		f.Set(os.Stdout)
	}
	for _, v := range config.C.Log.LogLevel {
		switch v {
		case "LOGERROR":
			SetLogLevel(LOGERROR)
			FMTLog(LOGWARN, "LOGERROR ENABLE")
		case "LOGWARN":
			SetLogLevel(LOGWARN)
			FMTLog(LOGWARN, "LOGWARN ENABLE")
		case "LOGDEBUG":
			SetLogLevel(LOGDEBUG)
			FMTLog(LOGWARN, "LOGDEBUG ENABLE")
		}
	}
}

// SetLogLevel 设置日志级别，多次调用权限叠加
func SetLogLevel(i int) {
	logLevel |= i
}

// FMTLog ...
func FMTLog(level int, strs ...interface{}) {
	if (logLevel & level) > 0 {
		strs = append(make([]interface{}, 1), strs...)
		switch level {
		case LOGERROR:
			strs[0] = "[LOGERROR]"
			strs = append(strs, string(debug.Stack()))
		case LOGWARN:
			strs[0] = "[LOGWARN]"
		case LOGDEBUG:
			strs[0] = "[LOGDEBUG]"
		}
		f.Print(strs...)
	}
}
