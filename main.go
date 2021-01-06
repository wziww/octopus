package main

import (
	"fmt"
	"io"
	"net/http"
	_ "net/http/pprof"
	"octopus/config"
	"octopus/log"
	_ "octopus/myredis"
	_ "octopus/myredis/rdb"
	"os"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/judwhite/go-svc/svc"
)

type program struct {
}

func main() {
	prg := &program{}
	if err := svc.Run(prg, syscall.SIGINT, syscall.SIGTERM); err != nil {
		log.FMTLog(log.LOGERROR, err)
	}
}
func (p *program) Init(e svc.Environment) error {
	return nil
}
func (p *program) Start() error {
	log.FMTLog(log.LOGWARN, "octopus start")
	go func() {
		log.FMTLog(log.LOGWARN, "HTTP start at "+config.C.Server.ListenAddress)
		server := &http.Server{
			Addr:         config.C.Server.ListenAddress,
			WriteTimeout: 60 * time.Second,
			ReadTimeout:  60 * time.Second,
		}
		http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Access-Control-Allow-Origin", "*")
			w.Header().Add("Access-Control-Allow-Headers", "*")
			w.Header().Add("Access-Control-Allow-Methods", "HEAD,PUT,POST,GET,DELETE,OPTIONS")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			r.ParseMultipartForm(30 << 10)
			file, header, err := r.FormFile("file")
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte(err.Error()))
				return
			}
			defer file.Close()
			filePath := path.Join(config.C.RDB.Dir, path.Base(header.Filename))
			if path.Dir(filePath) != path.Clean(config.C.RDB.Dir) {
				w.WriteHeader(500)
				w.Write([]byte("file path not allowed"))
				return
			}
			f, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0666)
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte(err.Error()))
				return
			}
			defer f.Close()
			_, err = io.Copy(f, file)
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte(err.Error()))
				return
			}
			w.WriteHeader(200)
			w.Write([]byte("success"))
			return
		})
		http.HandleFunc("/v1/websocket", func(w http.ResponseWriter, r *http.Request) {
			ws(w, r)
			return
		})
		http.HandleFunc("/prometheus/", func(w http.ResponseWriter, r *http.Request) {
			httprouter(w, r)
			return
		})
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			params := strings.Split(r.URL.Path, "/")
			for _, v := range params[len(params)-1:] {
				for _, z := range v {
					if z == []rune(".")[0] {
						http.FileServer(http.Dir("./src/dist")).ServeHTTP(w, r)
						return
					}
				}
			}
			http.ServeFile(w, r, "./src/dist/index.html")
			return
		})
		err := server.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()
	return nil
}
func (p *program) Stop() error {
	fmt.Println("exit")
	return nil
}
