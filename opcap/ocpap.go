package opcap

import (
	"net"
	"octopus/log"
	"strconv"
	"sync"
)

var client map[string]*OConn

// OConn ...
type OConn struct {
	mutex sync.Mutex
	conn  net.Conn
}

var (
	mutex sync.Mutex
	// CRLF ...
	CRLF = []byte("\r\n")
	// MAXSIZE tcp max size
	MAXSIZE = 65535
)

func init() {
	client = make(map[string]*OConn)
}

// CreateOrGetClient ...
func CreateOrGetClient(address string) (c *OConn, err error) {
	mutex.Lock()
	defer mutex.Unlock()
	if client[address] == nil {
		conn, e := net.Dial("tcp", address)
		if e == nil {
			client[address] = &OConn{
				conn: conn,
			}
		}
		return nil, e
	}
	return client[address], nil
}

// PING ...
func PING(conn *OConn, address string) string {
	if conn == nil {
		return "failed"
	}
	e := Write(conn, address, []byte("ping\r\n"))
	if e != nil {
		return "failed"
	}
	head := readCRLF(conn)
	if len(string(head)) >= 4 {
		return string(head)[0:4]
	}
	return string(head) // return pong
}

func readCRLF(conn *OConn) []byte {
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	var index int
	var all []byte
	for i := 0; i < MAXSIZE; i++ {
		var headByte = make([]byte, 1)
		conn.conn.Read(headByte)
		if headByte[0] == CRLF[0] { // CRLF
			index++
		} else if headByte[0] == CRLF[1] && index == 1 {
			break
		} else {
			index = 0
			if headByte[0] != '\x00' { // empty ascii
				all = append(all, headByte[0])
			}
		}
	}
	return all
}

// Write ...
func Write(conn *OConn, address string, bts []byte) (e error) {
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	_, e = conn.conn.Write(bts)
	if e != nil {
		delete(client, address)
	}
	return
}

// Count ...
func Count(conn *OConn, address string) (result []string) {
	if conn == nil {
		return []string{"failed"}
	}
	e := Write(conn, address, []byte("get\r\n"))
	if e != nil {
		return
	}
	head := readCRLF(conn)
	if len(head) == 0 {
		return
	}
	len, e := strconv.Atoi(string(head))
	if e != nil {
		log.FMTLog(log.LOGERROR, e.Error())
	}
	for i := 0; i < len*2; i += 2 {
		result = append(result, string(readCRLF(conn)))
		result = append(result, string(readCRLF(conn)))
	}
	return
}
