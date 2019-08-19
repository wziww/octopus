package opcap

import (
	"bufio"
	"net"
	"octopus/log"
	"strconv"
	"sync"
)

var client map[string]*OConn

// OConn ...
type OConn struct {
	mutex   sync.Mutex
	address string
	fd      net.Conn
	rd      *reader
	wr      *writer
}

type writer struct {
	wr *bufio.Writer
}
type reader struct {
	rd   *bufio.Reader
	_buf []byte
}

var (
	mutex sync.Mutex
	// CRLF ...
	CRLF = []byte("\r\n")
)

func init() {
	client = make(map[string]*OConn)
}

func connGet(conn *OConn) {
	conn.mutex.Lock()
}
func connRelease(conn *OConn) {
	conn.mutex.Unlock()
}

// CreateOrGetClient ...
func CreateOrGetClient(address string) (c *OConn, err error) {
	mutex.Lock()
	defer mutex.Unlock()
	if client[address] == nil {
		fd, e := net.Dial("tcp", address)
		if e == nil {
			client[address] = &OConn{
				address: address,
				fd:      fd,
				wr: &writer{
					wr: bufio.NewWriter(fd),
				},
				rd: &reader{
					rd:   bufio.NewReader(fd),
					_buf: make([]byte, 64),
				},
			}
		}
		return nil, e
	}
	return client[address], nil
}

// RefreshClient ...
func RefreshClient(address string) (c *OConn, err error) {
	mutex.Lock()
	defer mutex.Unlock()
	delete(client, address)
	fd, e := net.Dial("tcp", address)
	if e == nil {
		client[address] = &OConn{
			address: address,
			fd:      fd,
			wr: &writer{
				wr: bufio.NewWriter(fd),
			},
			rd: &reader{
				rd:   bufio.NewReader(fd),
				_buf: make([]byte, 64),
			},
		}
	}
	return client[address], nil
}

// PING ...
func PING(conn *OConn, address string) (str string) {
	if conn == nil {
		return "failed"
	}
	connGet(conn)
	defer connRelease(conn)
	e := conn.write("ping")
	if e != nil {
		log.FMTLog(log.LOGERROR, e.Error())
		return ""
	}
	e = conn.clrf()
	if e != nil {
		log.FMTLog(log.LOGERROR, e.Error())
		return ""
	}
	line := conn.getline()
	if len(line) == 0 {
		return ""
	}
	str = string(line)
	return
}

// Count ...
func Count(conn *OConn, address string) (result []string) {
	if conn == nil {
		return []string{"failed"}
	}
	connGet(conn)
	defer connRelease(conn)
	e := conn.write("get")
	if e != nil {
		log.FMTLog(log.LOGERROR, e.Error())
		return []string{""}
	}
	e = conn.clrf()
	if e != nil {
		log.FMTLog(log.LOGERROR, e.Error())
		return []string{""}
	}
	line := conn.getline()
	if len(line) == 0 {
		return []string{""}
	}
	count, _ := strconv.Atoi(string(line))
	for i := 0; i < count; i++ {
		key := conn.getline()
		value := conn.getline()
		result = append(result, string(key))
		result = append(result, string(value))
	}
	return
}

func (conn *OConn) getline() []byte {
	line, isPrefix, err := conn.rd.rd.ReadLine()
	if err != nil {
		log.FMTLog(log.LOGERROR, err.Error())
		RefreshClient(conn.address)
		return []byte{}
	}
	if isPrefix {
		log.FMTLog(log.LOGERROR, bufio.ErrBufferFull.Error())
		return []byte{}
	}
	if len(line) == 0 {
		return []byte{}
	}
	return line
}

func (conn *OConn) write(arg interface{}) (e error) {
	switch arg.(type) {
	case []byte:
		_, e = conn.wr.wr.Write(arg.([]byte))
	case string:
		_, e = conn.wr.wr.WriteString(arg.(string))
	case byte:
		e = conn.wr.wr.WriteByte(arg.(byte))
	case rune:
		_, e = conn.wr.wr.WriteRune(arg.(rune))
	default:
		break
	}
	if e != nil {
		log.FMTLog(log.LOGERROR, e.Error())
		RefreshClient(conn.address)
	}
	return
}
func (conn *OConn) clrf() (e error) {
	e = conn.write(CRLF[0])
	if e != nil {
		return e
	}
	e = conn.write(CRLF[1])
	if e != nil {
		return e
	}
	return conn.wr.wr.Flush()
}
