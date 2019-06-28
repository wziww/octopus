package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "octopus/myredis"
	"syscall"

	"github.com/judwhite/go-svc/svc"
)

var addr = flag.String("addr", "0.0.0.0:8080", "http service address")

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
		http.HandleFunc("/v1/websocket", ws)
		log.Fatal(http.ListenAndServe(*addr, nil))
	}()
	return nil
}
func (p *program) Stop() error {
	fmt.Println("exit")
	return nil
}
