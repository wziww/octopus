package opcap

import (
	"net"
	"sync"
)

var client map[string]net.Conn
var (
	mutex sync.Mutex
	// CRLF ...
	CRLF = []byte("\r\n")
)

func init() {
	client = make(map[string]net.Conn)
}

// CreateOrGetClient ...
func CreateOrGetClient(address string) (conn net.Conn, err error) {
	mutex.Lock()
	defer mutex.Unlock()
	if client[address] == nil {
		conn, err = net.Dial("tcp", address)
		if err == nil {
			client[address] = conn
		}
		return
		// conn.Write([]byte("Hello world!"))
	}
	return client[address], nil
}

// PING ...
func PING(conn net.Conn, address string) string {
	conn.Write([]byte("ping\r\n"))
	var buf = make([]byte, 10)
	_, err := conn.Read(buf)
	if err != nil {
		mutex.Lock()
		defer mutex.Unlock()
		delete(client, address)
		return "failed"
	}
	if len(string(buf)) >= 4 {
		return string(buf)[0:4]
	}
	return string(buf) // return pong
}
