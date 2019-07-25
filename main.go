package main

import (
	"fmt"
	"net/http"
	"octopus/config"
	"octopus/log"
	_ "octopus/myredis"
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
		fmt.Println(err)
	}
}
func (p *program) Init(e svc.Environment) error {
	return nil
}
func (p *program) Start() error {
	log.FMTLog(log.LOGWARN, "octopus start")
	go func() {
		server := &http.Server{
			Addr:         config.C.Server.ListenAddress,
			WriteTimeout: 5 * time.Second,
			ReadTimeout:  5 * time.Second,
		}
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if r.RequestURI == "/v1/websocket" {
				ws(w, r)
				return
			}
			if len(r.RequestURI) >= len("/prometheus") && r.RequestURI[:len("/prometheus")] == "/prometheus" {
				httprouter(w, r)
				return
			}
			params := strings.Split(r.RequestURI, "/")
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
