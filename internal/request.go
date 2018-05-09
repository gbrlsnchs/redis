package internal

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
	"time"
)

// Request is a Redis request.
type Request struct {
	net.Conn
}

func NewRequest(conn net.Conn, cmd []byte, args ...interface{}) (*Request, error) {
	var buf bytes.Buffer
	r := &Request{Conn: conn}
	parsed := r.parseArgs(args...)

	buf.WriteString(fmt.Sprintf("*%d\r\n$%d\r\n%s\r\n", len(parsed)+1, len(cmd), cmd))

	for i := range parsed {
		arg := parsed[i]

		buf.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(arg), arg))
	}

	if _, err := r.Write(buf.Bytes()); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Request) parseArgs(args ...interface{}) [][]byte {
	b := make([][]byte, 0)

	for i := range args {
		switch v := args[i].(type) {
		case string:
			b = append(b, []byte(v))

		case []byte:
			b = append(b, v)

		case byte:
			b = append(b, []byte{v})

		case int:
			b = append(b, []byte(strconv.Itoa(v)))

		case time.Time:
			t := int(v.Unix())
			b = append(b, []byte(strconv.Itoa(t)))
		}
	}

	return b
}
