package redis

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/gbrlsnchs/redis/internal"
)

var (
	ErrNoResult   = errors.New("redis: no result")
	ErrNullArray  = errors.New("redis: array is null")
	ErrRowsClosed = errors.New("redis: rows are closed")
)

type Rows struct {
	conn     net.Conn
	rd       *bufio.Reader
	err      error
	value    []byte
	null     bool
	isMulti  bool
	closed   bool
	errMu    *sync.RWMutex
	closedMu *sync.RWMutex
}

func newRows(conn net.Conn, isMulti bool) *Rows {
	return &Rows{
		conn:     conn,
		rd:       bufio.NewReader(conn),
		isMulti:  isMulti,
		errMu:    &sync.RWMutex{},
		closedMu: &sync.RWMutex{},
	}
}

func (rs *Rows) Close() error {
	if rs.isMulti {
		return nil
	}
	if err := rs.conn.Close(); err != nil {
		return err
	}
	rs.closedMu.Lock()
	rs.closed = true
	rs.closedMu.Unlock()
	return nil
}

func (rs *Rows) Err() error {
	rs.errMu.RLock()
	defer rs.errMu.RUnlock()
	return rs.err
}

func (rs *Rows) Next() bool {
	rs.null = false
	b, err := rs.rd.ReadByte()
	// EOF.
	if err != nil {
		return false
	}
	t := Type(b)                  // first byte dictates the response type
	s, _ := rs.rd.ReadBytes('\n') // trust Redis has properly set CRLF
	s = s[:len(s)-2]

	if t == ArrayType {
		rs.null = s[0] == '-'
		if rs.null {
			rs.err = ErrNullArray
			return false
		}
		b, _ = rs.rd.ReadByte()
		t = Type(b)
	}
	if t == ErrorType {
		rs.err = fmt.Errorf("redis: %s", s)
		return false
	}
	if t == BulkStringType {
		size, _ := strconv.Atoi(string(s))
		if rs.null = size == -1; !rs.null {
			str := make([]byte, size+2) // read including CRLF
			rs.rd.Read(str)
			s = str[:len(str)-2] // remove CRLF
		} else {
			rs.err = ErrNoResult
		}
	}
	// In case the value is either a simple string
	// or an integer, it needn't be parsed, simply
	// return it as a byte array.
	rs.value = s
	return true
}

func (rs *Rows) Scan(value interface{}) error {
	if rs.isClosed() {
		return ErrRowsClosed
	}
	var err error
	switch v := value.(type) {
	case *[]byte:
		*v = rs.value
	case *string:
		if rs.null {
			err = errors.New("RESP response is a null value, not string")
			break
		}
		*v = string(rs.value)
	case *NullString:
		*v = NullString{Valid: !rs.null}
		if v.Valid {
			v.String = string(rs.value)
		}
	case *NullInt64:
		*v = NullInt64{Valid: !rs.null}
		if v.Valid {
			v.Int64, err = strconv.ParseInt(string(rs.value), 10, 64)
		}
	case *bool:
		*v, err = strconv.ParseBool(string(rs.value))
	case *int:
		*v, err = strconv.Atoi(string(rs.value))
	case *int8:
		var n int64
		n, err = strconv.ParseInt(string(rs.value), 10, 8)
		*v = int8(n)
	case *int16:
		var n int64
		n, err = strconv.ParseInt(string(rs.value), 10, 16)
		*v = int16(n)
	case *int32:
		var n int64
		n, err = strconv.ParseInt(string(rs.value), 10, 32)
		*v = int32(n)
	case *int64:
		*v, err = strconv.ParseInt(string(rs.value), 10, 64)
	case *uint:
		var n uint64
		n, err = strconv.ParseUint(string(rs.value), 10, 0)
		*v = uint(n)
	case *uint8:
		var n uint64
		n, err = strconv.ParseUint(string(rs.value), 10, 8)
		*v = uint8(n)
	case *uint16:
		var n uint64
		n, err = strconv.ParseUint(string(rs.value), 10, 16)
		*v = uint16(n)
	case *uint32:
		var n uint64
		n, err = strconv.ParseUint(string(rs.value), 10, 32)
		*v = uint32(n)
	case *uint64:
		*v, err = strconv.ParseUint(string(rs.value), 10, 64)
	case *float32:
		var n float64
		n, err = strconv.ParseFloat(string(rs.value), 32)
		*v = float32(n)
	case *float64:
		*v, err = strconv.ParseFloat(string(rs.value), 64)
	case *time.Time:
		var unix int64
		if unix, err = strconv.ParseInt(string(rs.value), 10, 64); err != nil {
			break
		}
		*v = time.Unix(unix, 0)
	case *internal.Dump:
		return nil
	default:
		err = errors.New("redis: type is not scannable")
	}
	if err != nil {
		rs.errMu.Lock()
		rs.err = err
		rs.errMu.Unlock()
	}
	return nil
}

func (rs *Rows) isClosed() bool {
	rs.closedMu.RLock()
	defer rs.closedMu.RUnlock()
	return rs.closed
}
