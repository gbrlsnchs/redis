package redis

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"
)

var ErrNullArray = errors.New("array is null")

type Rows struct {
	conn   net.Conn
	buf    *bufio.Reader
	err    error
	value  []byte
	null   bool
	multi  bool
	closed bool
	mu     *sync.RWMutex
}

func (rs *Rows) Close() error {
	if err := rs.conn.Close(); err != nil {
		return err
	}

	rs.mu.Lock()
	defer rs.mu.Unlock()

	rs.closed = true

	return nil
}

func (rs *Rows) Err() error {
	rs.mu.RLock()
	defer rs.mu.RUnlock()

	return rs.err
}

func (rs *Rows) Next() bool {
	rs.null = false

	if rs.buf == nil {
		rs.buf = bufio.NewReader(rs.conn)
	}

	var b byte
	b, err := rs.buf.ReadByte()

	// EOF.
	if err != nil {
		return false
	}

	t := Type(b)
	s, _ := rs.buf.ReadBytes('\n')
	s = s[:len(s)-2]

	if t == ArrayType {
		s, _ = rs.buf.ReadBytes('\n')
		s = s[:len(s)-2]
		rs.null = s[0] == '-'

		if rs.null {
			rs.err = ErrNullArray

			return false
		}

		b, _ = rs.buf.ReadByte()
		t = Type(b)
	}

	switch t {
	case ErrorType:
		rs.err = fmt.Errorf("redis: %s", string(s))

		return false

	case BulkStringType:
		size, _ := strconv.Atoi(string(s))

		if rs.null = size == -1; !rs.null {
			str := make([]byte, size+2)

			rs.buf.Read(str)
			s = str[:len(str)-2]
		}
	}

	rs.value = s

	return true
}

func (rs *Rows) Scan(value interface{}) error {
	rs.mu.RLock()

	if rs.closed {
		return errors.New("redis: rows are closed")
	}

	rs.mu.RUnlock()

	switch p := value.(type) {
	case *string:
		if rs.null {
			return errors.New("RESP response is a null value, not string")
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
			var err error
			p.Int64, err = strconv.ParseInt(string(rs.value), 10, 64)

			return err
		}

	case *int:
		var err error
		*p, err = strconv.Atoi(string(rs.value))

		return err

	case *time.Time:
		unix, err := strconv.ParseInt(string(rs.value), 10, 64)

		if err != nil {
			return err
		}

		*p = time.Unix(unix, 0)

	default:
		return errors.New("redis: type is not scannable")
	}

	return nil
}
