package internal

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"unsafe"
)

var (
	// ErrNoResult informs that either a null array or null bulk string has been returned.
	ErrNoResult = errors.New("redis: no result")
	crlf        = []byte("\r\n")
)

// Reader is a RESP protocol reader.
type Reader struct {
	rd     *bufio.Reader
	values [][]byte
	mod    int
	null   bool
}

// NewReader returns a new buffered reader.
func NewReader(rd io.Reader) *Reader {
	return &Reader{rd: bufio.NewReader(rd)}
}

// ReadFull reads the RESP message and, if it's an array,
// it keeps reading until the array is iterated over.
func (r *Reader) ReadFull() ([][]byte, error) {
	size, err := r.ReadRESP()
	if err != nil {
		return nil, err
	}
	if size > 0 && r.values == nil {
		r.values = make([][]byte, 0, size*(r.mod+1))
	}
	for i := 0; i < size; i++ {
		if _, err = r.ReadRESP(); err != nil {
			return nil, err
		}
	}
	return r.values, nil
}

// ReadN runs ReadFull N times in order to process an array of arrays.
func (r *Reader) ReadN(times int) ([][]byte, error) {
	r.mod = times - 1
	var err error
	for i := 0; i < times; i++ {
		if _, err = r.ReadFull(); err != nil {
			return nil, err
		}
	}
	return r.values, nil
}

// ReadRESP reads the response until the whole message is read
// according to the RESP protocol.
func (r *Reader) ReadRESP() (int, error) {
	var (
		b   byte
		err error
	)

	// Try to read RESP type.
	if b, err = r.rd.ReadByte(); err != nil {
		return 0, nil // finish reading
	}
	t := Type(b)

	var bb []byte
	if bb, err = r.rd.ReadBytes('\n'); err != nil {
		return 0, errors.New("redis: incomplete message")
	}
	bb = bytes.TrimSuffix(bb, crlf)

	switch t {
	case ArrayType:
		if bb[0] == '-' {
			r.null = true
			return 0, ErrNoResult
		}
		var size int
		s := *(*string)(unsafe.Pointer(&bb))
		if size, err = strconv.Atoi(s); err != nil {
			return 0, err
		}
		return size, nil
	case ErrorType:
		return 0, fmt.Errorf("redis: %s", bb)
	case BulkStringType:
		var size int
		s := *(*string)(unsafe.Pointer(&bb))
		if size, err = strconv.Atoi(s); err != nil {
			return 0, err
		}
		if size == -1 {
			return 0, ErrNoResult
		}
		bb = make([]byte, size+2) // string length + CRLF
		if _, err = io.ReadFull(r.rd, bb); err != nil {
			return 0, err
		}
		bb = bytes.TrimSuffix(bb, crlf)
	}
	r.values = append(r.values, bb)
	return 0, nil
}
