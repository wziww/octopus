package rdb

import (
	"bufio"
	"bytes"
	"container/heap"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"octopus/log"
	"octopus/message"
	"os"
	"sort"
	"strconv"
	"sync/atomic"
	"time"
)

/*
 * //TODO
 * Cause every yupoo redis node's system are little endian,
 * we dont handle with big endian system rdb file.
 */

const (
	RDB_LENERR = math.MaxUint64
	/* The current RDB version. When the format changes in a way that is no longer
	backward compatible this number gets incremented. */
	RDB_VERSION = 9
	/* Module auxiliary data. */
	RDB_OPCODE_MODULE_AUX = 247
	/* LRU idle time. */
	RDB_OPCODE_IDLE = 248
	/* LFU frequency. */
	RDB_OPCODE_FREQ = 249
	/* RDB aux field. */
	RDB_OPCODE_AUX = 250
	/* Hash table resize hint. */
	RDB_OPCODE_RESIZEDB = 251
	/* Expire time in milliseconds. */
	RDB_OPCODE_EXPIRETIME_MS = 252
	/* Old expire time in seconds. */
	RDB_OPCODE_EXPIRETIME = 253
	/* DB number of the following keys. */
	RDB_OPCODE_SELECTDB = 254
	/* End of the RDB file. */
	RDB_OPCODE_EOF = 255
	//-----
	/* Defines related to the dump file format. To store 32 bits lengths for short
	 * keys requires a lot of space, so we check the most significant 2 bits of
	 * the first byte to interpreter the length:
	 *
	 * 00|XXXXXX => if the two MSB are 00 the len is the 6 bits of this byte
	 * 01|XXXXXX XXXXXXXX =>  01, the len is 14 byes, 6 bits + 8 bits of next byte
	 * 10|000000 [32 bit integer] => A full 32 bit len in net byte order will follow
	 * 10|000001 [64 bit integer] => A full 64 bit len in net byte order will follow
	 * 11|OBKIND this means: specially encoded object will follow. The six bits
	 *           number specify the kind of object that follows.
	 *           See the RDB_ENC_* defines.
	 *
	 * Lengths up to 63 are stored using a single byte, most DB keys, and may
	 * values, will fit inside. */
	RDB_6BITLEN  = 0
	RDB_14BITLEN = 1
	RDB_32BITLEN = 0x80
	RDB_64BITLEN = 0x81
	RDB_ENCVAL   = 3
	//---
	// rdbLoad...() functions flags.
	RDB_LOAD_NONE  = 0
	RDB_LOAD_ENC   = (1 << 0)
	RDB_LOAD_PLAIN = (1 << 1)
	RDB_LOAD_SDS   = (1 << 2)
	/* When a length of a string object stored on disk has the first two bits
	 * set, the remaining six bits specify a special encoding for the object
	 * accordingly to the following defines: */
	RDB_ENC_INT8  = 0 /* 8 bit signed integer */
	RDB_ENC_INT16 = 1 /* 16 bit signed integer */
	RDB_ENC_INT32 = 2 /* 32 bit signed integer */
	RDB_ENC_LZF   = 3 /* string compressed with FASTLZ */
	/* Map object types to RDB object types. Macros starting with OBJ_ are for
	 * memory storage and may change. Instead RDB types must be fixed because
	 * we store them on disk. */
	RDB_TYPE_STRING   = 0
	RDB_TYPE_LIST     = 1
	RDB_TYPE_SET      = 2
	RDB_TYPE_ZSET     = 3
	RDB_TYPE_HASH     = 4
	RDB_TYPE_ZSET_2   = 5 /* ZSET version 2 with doubles stored in binary. */
	RDB_TYPE_MODULE   = 6
	RDB_TYPE_MODULE_2 = 7 /* Module value with annotations for parsing without
	   the generating module being loaded. */
	/* NOTE: WHEN ADDING NEW RDB TYPE, UPDATE rdbIsObjectType() BELOW */

	/* Object types for encoded objects. */
	RDB_TYPE_HASH_ZIPMAP      = 9
	RDB_TYPE_LIST_ZIPLIST     = 10
	RDB_TYPE_SET_INTSET       = 11
	RDB_TYPE_ZSET_ZIPLIST     = 12
	RDB_TYPE_HASH_ZIPLIST     = 13
	RDB_TYPE_LIST_QUICKLIST   = 14
	RDB_TYPE_STREAM_LISTPACKS = 15
)

var (
	idle    uint32 = 0
	busy    uint32 = 1
	_status uint32 = idle

	_offset    int64
	errRdbLoad error = errors.New("rdb load len error")
)

/*
 * CAS
 */
func lock() bool {
	return atomic.CompareAndSwapUint32(&_status, idle, busy)
}
func release() bool {
	return atomic.CompareAndSwapUint32(&_status, busy, idle)
}
func readFull(r io.Reader, buf []byte) (n int, err error) {
	_offset += int64(len(buf))
	return io.ReadFull(r, buf)
}

/*
 * redis internel check function
 */
func rdbIsObjectType(t byte) bool {
	return (t >= 0 && t <= 7) || (t >= 9 && t <= 15)
}
func checkVer(ver int) bool {
	if ver < 1 || ver > RDB_VERSION {
		return false
	}
	return true
}

/*  ----------------
 * |	TYPE (1 byte) |
 *  ----------------
 */
func loadType(buffer *bufio.Reader) (byte, error) {
	bts := make([]byte, 1)
	_, err := readFull(buffer, bts)
	return bts[0], err
}

/*  ----------------
 * |	TYPE (4 byte) |
 *  ----------------
 * This is only used to load old databases stored with the RDB_OPCODE_EXPIRETIME
 * opcode. New versions of Redis store using the RDB_OPCODE_EXPIRETIME_MS
 * opcode.
 */
func rdbLoadTime(buffer *bufio.Reader) (int64, error) {
	bts := make([]byte, 4)
	_, err := readFull(buffer, bts)
	if err != nil {
		log.FMTLog(log.LOGERROR, err)
		return -1, err
	}
	bytebuf := bytes.NewBuffer(bts)
	var result int64 = -1
	err = binary.Read(bytebuf, binary.LittleEndian, &result)
	return result, err
}

/*  ----------------
 * |	TYPE (8 byte) |
 *  ----------------*/
func rdbLoadMillisecondTime(buffer *bufio.Reader, rdbver int) (int64, error) {
	bts := make([]byte, 8)
	_, err := readFull(buffer, bts)
	if err != nil {
		log.FMTLog(log.LOGERROR, err)
		return -1, err
	}
	bytebuf := bytes.NewBuffer(bts)
	var result int64 = -1
	err = binary.Read(bytebuf, binary.LittleEndian, &result)
	return result, err
}

/* Load an encoded length. If the loaded length is a normal length as stored
 * with rdbSaveLen(), the read length is set to '*lenptr'. If instead the
 * loaded length describes a special encoding that follows, then '*isencoded'
 * is set to 1 and the encoding format is stored at '*lenptr'.
 *
 * See the RDB_ENC_* definitions in rdb.h for more information on special
 * encodings.
 *
 * The function returns -1 on error, 0 on success. */
func rdbLoadLenByRef(buffer *bufio.Reader, isencoded *bool, lenptr *uint64) int {
	var _type int
	if isencoded != nil {
		*isencoded = false
	}
	buf := make([]byte, 1)
	_, err := readFull(buffer, buf)
	if err != nil {
		log.FMTLog(log.LOGERROR, err)
		return -1
	}
	_type = (int(buf[0]) & 0xC0) >> 6
	if _type == RDB_ENCVAL {
		/* Read a 6 bit encoding type. */
		if isencoded != nil {
			*isencoded = true
		}
		if lenptr != nil {
			*lenptr = uint64(buf[0]) & 0x3F
		}
	} else if _type == RDB_6BITLEN {
		/* Read a 6 bit len. */
		if lenptr != nil {
			*lenptr = uint64(buf[0]) & 0x3F
		}
	} else if _type == RDB_14BITLEN {
		/* Read a 14 bit len. */
		bts := make([]byte, 1)
		_, err := readFull(buffer, bts)
		if err != nil {
			log.FMTLog(log.LOGERROR, err)
			return -1
		}
		if lenptr != nil {
			*lenptr = ((uint64(buf[0]) & 0x3F) << 8) | uint64(bts[0])
		}
	} else if buf[0] == RDB_32BITLEN {
		/* Read a 32 bit len. */
		// uint32_t len;
		var length uint32
		bts := make([]byte, 4)
		_, err := readFull(buffer, bts)
		if err != nil {
			log.FMTLog(log.LOGERROR, err)
			return -1
		}
		bytebuf := bytes.NewBuffer(bts)
		binary.Read(bytebuf, binary.LittleEndian, &length)
		if lenptr != nil {
			bts2 := make([]byte, 4)
			binary.LittleEndian.PutUint32(bts2, length)
			*lenptr = uint64(binary.BigEndian.Uint32(bts2))
		}

	} else if buf[0] == RDB_64BITLEN {
		/* Read a 64 bit len. */
		var length uint64
		bts := make([]byte, 8)
		_, err := readFull(buffer, bts)
		if err != nil {
			log.FMTLog(log.LOGERROR, err)
			return -1
		}
		bytebuf := bytes.NewBuffer(bts)
		binary.Read(bytebuf, binary.LittleEndian, &length)
		if lenptr != nil {
			*lenptr = length
		}
	} else {
		log.FMTLog(log.LOGERROR, fmt.Sprintf("Unknown length encoding %d in rdbLoadLen()", _type))
		return -1 /* Never reached. */
	}
	return 0
}

/* This is like rdbLoadLenByRef() but directly returns the value read
 * from the RDB stream, signaling an error by returning RDB_LENERR
 * (since it is a too large count to be applicable in any Redis data
 * structure). */
func rdbLoadLen(buffer *bufio.Reader, isencoded *bool) uint64 {
	var length uint64
	if rdbLoadLenByRef(buffer, isencoded, &length) == -1 {
		log.FMTLog(log.LOGERROR, "rdb load len error")
		return RDB_LENERR
	}
	return length
}

/* Load an LZF compressed string in RDB format. The returned value
 * changes according to 'flags'. For more info check the
 * rdbGenericLoadStringObject() function. */
func rdbLoadLzfStringObject(rdb *bufio.Reader, flags int, lenptr *uint64) []byte {
	// plain := flags & RDB_LOAD_PLAIN
	// sds := flags & RDB_LOAD_SDS
	var ulen, clen uint64
	clen = rdbLoadLen(rdb, nil)
	if clen == RDB_LENERR {
		return nil
	}
	ulen = rdbLoadLen(rdb, nil)
	if ulen == RDB_LENERR {
		return nil
	}
	bts := make([]byte, int(clen))
	_, err := readFull(rdb, bts)
	if err != nil {
		log.FMTLog(log.LOGERROR, err)
		return nil
	}
	decompressed := lzfDecompress(bts[:int(clen)], int(ulen))
	if len(decompressed) != int(ulen) {
		log.FMTLog(log.LOGERROR, fmt.Sprintf("decompressed string length %d didn't match expected length %d", len(decompressed), ulen))
		return nil
	}
	return decompressed
}
func lzfDecompress(in []byte, outlen int) []byte {
	out := make([]byte, outlen)
	for i, o := 0, 0; i < len(in); {
		ctrl := int(in[i])
		i++
		if ctrl < 32 {
			for x := 0; x <= ctrl; x++ {
				out[o] = in[i]
				i++
				o++
			}
		} else {
			length := ctrl >> 5
			if length == 7 {
				length = length + int(in[i])
				i++
			}
			ref := o - ((ctrl & 0x1f) << 8) - int(in[i]) - 1
			i++
			for x := 0; x <= length+1; x++ {
				out[o] = out[ref]
				ref++
				o++
			}
		}
	}
	return out
}

/* Load a string object from an RDB file according to flags:
 *
 * RDB_LOAD_NONE (no flags): load an RDB object, unencoded.
 * RDB_LOAD_ENC: If the returned type is a Redis object, try to
 *               encode it in a special way to be more memory
 *               efficient. When this flag is passed the function
 *               no longer guarantees that obj->ptr is an SDS string.
 * RDB_LOAD_PLAIN: Return a plain string allocated with zmalloc()
 *                 instead of a Redis object with an sds in it.
 * RDB_LOAD_SDS: Return an SDS string instead of a Redis object.
 *
 * On I/O error NULL is returned.
 */
func rdbGenericLoadString(buffer *bufio.Reader, flags int, lenptr *uint64) []byte {
	// var encode int = flags & RDB_LOAD_ENC
	var plain int = flags & RDB_LOAD_PLAIN
	var sds int = flags & RDB_LOAD_SDS
	var isencoded bool
	var length uint64 = rdbLoadLen(buffer, &isencoded)
	if isencoded {
		switch length {
		case RDB_ENC_INT8:
			bts := make([]byte, 1)
			_, err := readFull(buffer, bts)
			if err != nil {
				log.FMTLog(log.LOGERROR, err)
				return nil
			}
			return []byte(strconv.Itoa(int(bts[0])))
		case RDB_ENC_INT16:
			bts := make([]byte, 2)
			_, err := readFull(buffer, bts)
			if err != nil {
				log.FMTLog(log.LOGERROR, err)
				return nil
			}
			return bts
		case RDB_ENC_INT32:
			bts := make([]byte, 4)
			_, err := readFull(buffer, bts)
			if err != nil {
				log.FMTLog(log.LOGERROR, err)
				return nil
			}
			return []byte(
				strconv.Itoa(int(bts[0]) | (int(bts[1]) << 8) | (int(bts[2]) << 16) | (int(bts[3]) << 24)),
			)
		case RDB_ENC_LZF:
			return rdbLoadLzfStringObject(buffer, flags, lenptr)
		default:
			log.FMTLog(log.LOGERROR, fmt.Sprintf("Unknown RDB string encoding type %d", length))
			return nil
		}
	}
	if length == RDB_LENERR {
		return nil
	}
	if plain > 0 || sds > 0 {
		bts := make([]byte, int(length))
		_, err := readFull(buffer, bts)
		if err != nil {
			log.FMTLog(log.LOGERROR, err)
			return nil
		}
		return bts
	}
	bts := make([]byte, int(length))
	_, err := readFull(buffer, bts)
	if err != nil {
		log.FMTLog(log.LOGERROR, err)
		return nil
	}
	return bts
}
func rdbLoadString(buffer *bufio.Reader) []byte {
	return rdbGenericLoadString(buffer, RDB_LOAD_NONE, nil)
}
func rdbLoadBinaryDoubleValue(buffer *bufio.Reader) (result float64, err error) {
	bts := make([]byte, 8)
	_, err = readFull(buffer, bts)
	if err != nil {
		return
	}
	return math.Float64frombits(binary.LittleEndian.Uint64(bts)), nil
}
func rdbLoadDoubleValue(buffer *bufio.Reader) (float64, error) {
	length := rdbLoadLen(buffer, nil)
	switch length {
	case 255:
		// *val = R_NegInf; return 0;
		return 0, nil
	case 254:
		// *val = R_PosInf; return 0;
		return 0, nil
	case 253:
		// *val = R_Nan; return 0;
		return 0, nil
	default:
	}
	bts := make([]byte, length)
	_, err := readFull(buffer, bts)
	if err != nil {
		log.FMTLog(log.LOGERROR, err)
		return -1, err
	}
	return strconv.ParseFloat(string(bts), 64)
}

/* Return length of ziplist. */
func ziplistLen(encoded []byte) int64 {
	// TODO max len UINT16_MAX
	// #define ZIPLIST_LENGTH(zl) (*((uint16_t*)((zl)+sizeof(uint32_t)*2)))
	return int64(binary.LittleEndian.Uint16(encoded[8 : 8+2]))
}

/* Load a Redis object of the specified type from the specified file.
 * On success a newly allocated object is returned, otherwise NULL. */
func rdbLoadObject(buffer *bufio.Reader, rdbtype int) []string {
	var result []string
	switch rdbtype {
	case RDB_TYPE_STRING:
		result = append(result, string(rdbLoadString(buffer)))
	case RDB_TYPE_LIST:
		length := rdbLoadLen(buffer, nil)
		if length == RDB_LENERR {
			return nil
		}
		for ; length > 0; length-- {
			result = append(result, string(rdbGenericLoadString(buffer, RDB_LOAD_ENC, nil)))
		}
	case RDB_TYPE_SET:
		var isEncoded bool
		length := rdbLoadLen(buffer, &isEncoded)
		if length == RDB_LENERR {
			return nil
		}
		for i := 0; i < int(length); i++ {
			if isEncoded {
				result = append(result, string(rdbGenericLoadString(buffer, RDB_LOAD_NONE, nil)))
			} else {
				result = append(result, string(rdbGenericLoadString(buffer, RDB_LOAD_ENC, nil)))
			}
		}
	case RDB_TYPE_ZSET_2,
		RDB_TYPE_ZSET:
		var isEncoded bool
		length := rdbLoadLen(buffer, &isEncoded)
		if length == RDB_LENERR {
			return nil
		}
		for i := 0; i < int(length); i++ {
			if isEncoded {
				result = append(result, string(rdbGenericLoadString(buffer, RDB_LOAD_NONE, nil)))
			} else {
				result = append(result, string(rdbGenericLoadString(buffer, RDB_LOAD_ENC, nil)))
			}
			var v float64
			var err error
			if rdbtype == RDB_TYPE_ZSET_2 {
				v, err = rdbLoadBinaryDoubleValue(buffer)
				if err != nil {
					log.FMTLog(log.LOGERROR, err)
					return nil
				}
			} else {
				v, err = rdbLoadDoubleValue(buffer)
				if err != nil {
					log.FMTLog(log.LOGERROR, err)
					return nil
				}
			}
			result = append(result, strconv.FormatFloat(v, 'f', -1, 64))
		}
	case RDB_TYPE_HASH:
		len := rdbLoadLen(buffer, nil)
		for ; len > 0; len-- {
			result = append(result, string(rdbGenericLoadString(buffer, RDB_LOAD_SDS, nil)))
			result = append(result, string(rdbGenericLoadString(buffer, RDB_LOAD_SDS, nil)))
		}
	case RDB_TYPE_LIST_QUICKLIST:
		len := rdbLoadLen(buffer, nil)
		if len == RDB_LENERR {
			return nil
		}
		for i := len; i > 0; i-- {
			_ = rdbGenericLoadString(buffer, RDB_LOAD_PLAIN, nil)
		}
		return make([]string, len) // TODO
	case RDB_TYPE_HASH_ZIPMAP,
		RDB_TYPE_LIST_ZIPLIST,
		RDB_TYPE_SET_INTSET,
		RDB_TYPE_ZSET_ZIPLIST,
		RDB_TYPE_HASH_ZIPLIST:
		encoded := rdbGenericLoadString(buffer, RDB_LOAD_PLAIN, nil)
		if encoded == nil {
			return nil
		}
		switch rdbtype {
		case RDB_TYPE_HASH_ZIPMAP:
		case RDB_TYPE_ZSET_ZIPLIST:
			var kvlen int64
			kvlen = ziplistLen(encoded)
			if kvlen >= math.MaxUint16 { // TODO
				return nil
			}
			return make([]string, kvlen) // TODO
		case RDB_TYPE_HASH_ZIPLIST:
			var kvlen int64
			kvlen = ziplistLen(encoded)
			if kvlen >= math.MaxUint16 { // TODO
				return nil
			}
			return make([]string, kvlen) // TODO
		case RDB_TYPE_SET_INTSET:
			/*
			 * typedef struct intset {
			 *		uint32_t encoding;  // offset 0
			 *		uint32_t length;		// offset 4
			 *		int8_t contents[];  // offset 9
			 * } intset;
			 */
			var len uint32
			len = binary.LittleEndian.Uint32(encoded[4:9])
			_ = len // avoid waring
			return []string{"待实现"}
		default:
			return nil
		}
	case RDB_TYPE_STREAM_LISTPACKS:
	case RDB_TYPE_MODULE,
		RDB_TYPE_MODULE_2:
	default:
		log.FMTLog(log.LOGERROR, fmt.Sprintf("Unknown RDB encoding type %d", rdbtype))
		return nil
	}
	return result
}

// Analyze ...
func Analyze(filePath string, offsetSize, childSize, count int64, fn func(string, ...int64)) string {
	if !lock() {
		return message.Res(500, "busy now")
	}
	_offset = 0
	defer func() {
		if !release() {
			panic("CAS ERROR")
		}
	}()
	result := CreateResult(offsetSize, childSize, count)
	fd, err := os.OpenFile(filePath, os.O_RDONLY, 0444)
	if err != nil {
		log.FMTLog(log.LOGERROR, err.Error())
		return message.Res(500, err.Error())
	}
	fileInfo, err := fd.Stat()
	if err != nil {
		log.FMTLog(log.LOGERROR, err.Error())
		return message.Res(500, err.Error())
	}
	totalSize := fileInfo.Size()
	defer fd.Close()
	// fd.Stat()
	buffer := bufio.NewReader(fd)
	/*  ---------------------
	 * |	REDIS+VER (9 byte) |
	 *  ---------------------
	 */
	bts := make([]byte, 9)
	_, err = readFull(buffer, bts)
	if err != nil {
		log.FMTLog(log.LOGERROR, err)
		return message.Res(500, err.Error())
	}
	if string(bts[:5]) != "REDIS" {
		log.FMTLog(log.LOGERROR, "Wrong signature trying to load DB from file")
		return message.Res(500, "Wrong signature trying to load DB from file")
	}
	var rdbver int
	rdbver, err = strconv.Atoi(string(bts[5:]))
	if err != nil {
		log.FMTLog(log.LOGERROR, err)
		return message.Res(500, err.Error())
	}
	if !checkVer(rdbver) {
		errStr := fmt.Sprintf("Can't handle RDB format version %d", rdbver)
		log.FMTLog(log.LOGERROR, errStr)
		return message.Res(500, errStr)
	}
	now := time.Now().UnixNano() / 1e6
	var expiretime int64 = -1
	var percent float32 = 0.00
	for {
		_type, err := loadType(buffer)
		if err != nil {
			log.FMTLog(log.LOGERROR, err)
			return message.Res(500, err.Error())
		}
		switch _type {
		case RDB_OPCODE_EXPIRETIME:
			expiretime, err = rdbLoadTime(buffer)
			if err != nil {
				log.FMTLog(log.LOGERROR, err)
				return message.Res(500, err.Error())
			}
			expiretime *= 1000
			continue /* read next opcode */
		case RDB_OPCODE_EXPIRETIME_MS:
			expiretime, err = rdbLoadMillisecondTime(buffer, rdbver)
			if err != nil {
				log.FMTLog(log.LOGERROR, err)
				return message.Res(500, err.Error())
			}
			continue /* read next opcode */
		case RDB_OPCODE_FREQ:
			var _byte uint8
			bts := make([]byte, 1)
			_, err := readFull(buffer, bts)
			if err != nil {
				log.FMTLog(log.LOGERROR, err)
				return message.Res(500, err.Error())
			}
			_byte = bts[0]
			_ = _byte // avoid waring
			continue  /* Read next opcode. */
		case RDB_OPCODE_IDLE:
			if rdbLoadLen(buffer, nil) == RDB_LENERR {
				return message.Res(500, "rdb load len error")
			}
			continue /* Read next opcode. */
		case RDB_OPCODE_EOF:
			/* EOF: End of file, exit the main loop. */
			goto end
		case RDB_OPCODE_SELECTDB:
			var dbid uint64
			if dbid = rdbLoadLen(buffer, nil); dbid == RDB_LENERR {
				return message.Res(500, "rdb load len error")
			}
			// fmt.Printf("Selecting DB ID %d\n", dbid)
			continue
		case RDB_OPCODE_RESIZEDB:
			var dbSize, expiresSize uint64
			dbSize = rdbLoadLen(buffer, nil)
			if dbSize == RDB_LENERR {
				return message.Res(500, "rdb load len error")
			}
			expiresSize = rdbLoadLen(buffer, nil)
			if expiresSize == RDB_LENERR {
				return message.Res(500, "rdb load len error")
			}
			// fmt.Println(fmt.Sprintf("resizedb from %d to %d", dbSize, expiresSize))
			continue
		case RDB_OPCODE_AUX:
			/* AUX: generic string-string fields. Use to add state to RDB
			 * which is backward compatible. Implementations of RDB loading
			 * are requierd to skip AUX fields they don't understand.
			 *
			 * An AUX field is composed of two strings: key and value. */
			var auxkey, auxval []byte
			auxkey = rdbLoadString(buffer)
			if auxkey == nil {
				return message.Res(500, "rdb load string object error")
			}
			auxval = rdbLoadString(buffer)
			if auxval == nil {
				return message.Res(500, "rdb load string object error")
			}
			if string(auxkey) == "lua" {
				result.LuaNums++
			}
			continue
		default:
			if !rdbIsObjectType(_type) {
				errStr := fmt.Sprintf("Invalid object type: %d", _type)
				log.FMTLog(log.LOGERROR, errStr)
				return message.Res(500, errStr)
			}
		}
		var auxkey string
		var auxval []string
		before := _offset
		if auxkeybts := rdbLoadString(buffer); auxkeybts == nil {
			return message.Res(500, "rdb load string object error")
		} else {
			auxkey = string(auxkeybts)
		}
		auxval = rdbLoadObject(buffer, int(_type))
		if auxval == nil {
			return message.Res(500, auxkey+"   -----  rdb load object error")
		}
		if expiretime != -1 {
			result.Expires++
			if expiretime < now {
				result.AlreadyExpired++
			}
		}
		expiretime = -1
		if diff := _offset - before; diff >= result.OffSetSize {
			heap.Push(result.OffSetLog, Node{
				Key: auxkey,
				Val: diff,
			})
			result.OffSetCount++
			if result.OffSetCount > result.Count {
				result.OffSetCount--
				heap.Pop(result.OffSetLog)
			}
		}
		if len(auxval) > int(result.ChildSize) {
			heap.Push(result.ChildLog, Node{
				Key: auxkey,
				Val: int64(len(auxval)),
			})
			result.ChildCount++
			if result.ChildCount > result.Count {
				result.ChildCount--
				heap.Pop(result.ChildLog)
			}
		}
		result.TotalNums++
		if float32(_offset) >= (float32(totalSize) * percent) {
			percent += 0.01
			fn(strconv.FormatInt(_offset*100/totalSize, 10), 0)
		}
	}
end:
	sort.Sort(result.ChildLog)
	var i, j int = 0, result.ChildLog.Len() - 1
	for i < j {
		result.ChildLog.Swap(i, j)
		i++
		j--
	}
	sort.Sort(result.OffSetLog)
	i, j = 0, result.OffSetLog.Len()-1
	for i < j {
		result.OffSetLog.Swap(i, j)
		i++
		j--
	}
	fn("100", 0)
	return message.Res(200, result)
}
func init() {
}
