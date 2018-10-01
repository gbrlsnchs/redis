package redis

import (
	"context"
	"net"

	"github.com/gbrlsnchs/redis/internal"
)

func send(ctx context.Context, w *internal.Writer, conn net.Conn, times int, cmd string, args ...interface{}) (*Result, error) {
	ch := make(chan interface{})

	// Check if context is done, otherwise wait for the whole connection to be read.
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		go func() {
			if _, err := w.WriteCmd(cmd, args...); err != nil {
				ch <- err
				return
			}
			r, err := read(conn, times)
			if err != nil {
				ch <- err
				return
			}
			ch <- r
		}()
		switch v := (<-ch).(type) {
		case *Result:
			return v, nil
		case error:
			return nil, v
		default:
			return &Result{}, nil // this never happens but compiler is only happy with it included
		}
	}
}
