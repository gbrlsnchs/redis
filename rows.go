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
	ErrNoResult  = errors.New("redis: no result")
	ErrNullArray = errors.New("redis: array is null")
)

type Rows struct {
	conn    net.Conn
	buf     *bufio.Reader
	err     error
	value   []byte
	null    bool
	isMulti bool
	closed  bool
	mu      *sync.RWMutex
}

func (rs *Rows) Close() error {
	if rs.isMulti {
		return nil
	}
	if err := rs.conn.Close(); err != nil {
		return err
	}
	rs.mu.Lock()
	rs.closed = true
	rs.mu.Unlock()
	return nil
}

func (rs *Rows) Err() error {
	rs.mu.RLock()
	defer rs.mu.RUnlock()
	return rs.err
}

func (rs *Rows) Next() bool {
	rs.null = false
	// Read response and store its buffer.
	if rs.buf == nil {
		rs.buf = bufio.NewReader(rs.conn)
	}
	var b byte
	b, err := rs.buf.ReadByte()
	// EOF.
	if err != nil {
		return false
	}
	t := Type(b)                   // first byte dictates the response type
	s, _ := rs.buf.ReadBytes('\n') // trust Redis has properly set CRLF
	s = s[:len(s)-2]

	if t == ArrayType {
		rs.null = s[0] == '-'
		if rs.null {
			rs.err = ErrNullArray
			return false
		}
		b, _ = rs.buf.ReadByte()
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
			rs.buf.Read(str)
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
		return errors.New("redis: rows are closed")
	}
	var err error
	switch p := value.(type) {
	case *string:
		if rs.null {
			err = errors.New("RESP response is a null value, not string")
			break
		}
		*p = string(rs.value)
	case *NullString:
		*p = NullString{Valid: !rs.null}
		if p.Valid {
			p.String = string(rs.value)
		}
	case *NullInt64:
		*p = NullInt64{Valid: !rs.null}
		if p.Valid {
			p.Int64, err = strconv.ParseInt(string(rs.value), 10, 64)
		}
	case *int:
		*p, err = strconv.Atoi(string(rs.value))
	case *time.Time:
		var unix int64
		if unix, err = strconv.ParseInt(string(rs.value), 10, 64); err != nil {
			break
		}
		*p = time.Unix(unix, 0)
	case *internal.Dump:
		return nil
	default:
		err = errors.New("redis: type is not scannable")
	}
	if err != nil {
		rs.err = err
	}
	return nil
}

func (rs *Rows) isClosed() bool {
	rs.mu.RLock()
	defer rs.mu.RUnlock()
	return rs.closed
}
