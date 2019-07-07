package main

import (
	"fmt"
	"log"
	"net/http"
	_ "octopus/myredis"
	"strings"
	"syscall"

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
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if r.RequestURI == "/v1/websocket" {
				ws(w, r)
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
		log.Fatal(http.ListenAndServe(C.Server.ListenAddress, nil))
	}()
	return nil
}
func (p *program) Stop() error {
	fmt.Println("exit")
	return nil
}
