package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"octopus/config"
	"sync"
	"time"

	"github.com/google/gopacket/layers"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

var (
	// CRLF Carriage-Return Line-Feed 回车&换行
	CRLF  = []byte("\r\n")
	charx = []byte{42}
	chard = []byte{36}
)

type _commands struct {
	CC    map[string]int64 `json:"cc"`
	mutex sync.RWMutex
}

var commands *_commands

type rdsProtocol struct {
	paramsLen byte // 参数数量
	params    []params
}
type params struct {
	len   byte   // 参数长度
	value []byte // 参数内容
}

func (r *rdsProtocol) parse(bts []byte) bool {
	pms := bytes.Split(bts, CRLF)
	if len(pms) < 3 {
		return false
	}
	for i, v := range pms {
		if i == 0 {
			if v[0] != charx[0] {
				return false
			}
			r.paramsLen = v[1]
		} else if ((i % 2) == 1) && i+1 < len(pms) {
			r.params = append(r.params, params{
				len:   v[1],
				value: pms[i+1],
			})
		}
	}
	return true
}
func init() {
	commands = &_commands{}
	commands.CC = make(map[string]int64, 0)
}
func main() {
	go func() {
		for {
			select {
			case <-time.After(60 * time.Second):
				commands.mutex.Lock()
				commands.CC = make(map[string]int64, 0)
				commands.mutex.Unlock()
			}
		}
	}()
	go func() {
		server := &http.Server{
			Addr:         "0.0.0.0:9712",
			WriteTimeout: 5 * time.Second,
			ReadTimeout:  5 * time.Second,
		}
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			commands.mutex.RLock()
			bts, _ := json.Marshal(commands.CC)
			commands.mutex.RUnlock()
			w.Write(bts)
			return
		})
		err := server.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()
	handle, _ := pcap.OpenLive(
		config.C.Opcap.Device,
		int32(65535),
		false,
		60*time.Second,
	)
	defer handle.Close()
	handle.SetBPFFilter(config.C.Opcap.BPFFilter)
	packetSource := gopacket.NewPacketSource(
		handle,
		handle.LinkType(),
	)
	for packet := range packetSource.Packets() {
		var eth layers.Ethernet
		var ip4 layers.IPv4
		var tcp layers.TCP
		parser := gopacket.NewDecodingLayerParser(
			layers.LayerTypeEthernet, &eth, &ip4, &tcp)
		decodedLayers := []gopacket.LayerType{}
		parser.DecodeLayers(packet.Data(), &decodedLayers)
		if len(tcp.Payload) == 0 {
			continue
		}
		go func(bts []byte) {
			rp := &rdsProtocol{}
			rp.parse(bts)
			if rp.paramsLen > 0 {
				currentCommand := rp.params[0].value
				commands.mutex.Lock()
				commands.CC[string(currentCommand)]++
				commands.mutex.Unlock()
			}
		}(tcp.Payload)
	}
}
