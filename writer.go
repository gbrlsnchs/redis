package redis

import (
	"context"
	"net"

	"github.com/gbrlsnchs/redis/internal"
)

type writer struct {
	*internal.Writer
	conn  net.Conn
	times int
}

func newWriter(conn net.Conn, times int) *writer {
	return &writer{
		Writer: internal.NewWriter(conn),
		conn:   conn,
		times:  times,
	}
}

func (w *writer) read() (*Result, error) {
	rr := internal.NewReader(w.conn)
	values, err := rr.ReadN(w.times)
	if err != nil {
		return nil, err
	}
	return &Result{values: values, length: len(values)}, nil
}

func (w *writer) send(ctx context.Context, cmd string, args ...interface{}) (*Result, error) {
	ch := make(chan interface{})

	// Check if context is done, otherwise wait for the whole connection to be read.
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		if _, err := w.WriteCmd(cmd, args...); err != nil {
			return nil, err
		}
		r, err := w.read()
		if err != nil {
			return nil, err
		}
		return r, nil
	}
}
