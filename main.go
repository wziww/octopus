package main

import (
	"fmt"
	"log"
	"net/http"
	_ "octopus/myredis"
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
		http.HandleFunc("/v1/websocket", ws)
		log.Fatal(http.ListenAndServe(C.Server.ListenAddress, nil))
	}()
	return nil
}
func (p *program) Stop() error {
	fmt.Println("exit")
	return nil
}
