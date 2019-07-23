package log

import (
	"fmt"
	"octopus/config"
	"os"
	"path"
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
	file.fdlock.Lock()
	defer file.fdlock.Unlock()
	_, e := fmt.Fprintln(file.fd, strs...)
	if e != nil {
		fmt.Println(e)
	}
}

func init() {
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
}

// SetLogLevel 设置日志级别，多次调用权限叠加
func SetLogLevel(i int) {
	logLevel |= i
}

// FMTLog ...
func FMTLog(level int, strs ...interface{}) {
	if (logLevel & level) > 0 {
		f.Print(strs...)
	}
}
