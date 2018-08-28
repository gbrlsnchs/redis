package redis

import (
	"context"
	"net"

	"github.com/gbrlsnchs/redis/internal"
)

func send(ctx context.Context, w *internal.Writer, conn net.Conn, times int, cmd string, args ...interface{}) (*Result, error) {
	var err error
	rc := make(chan *Result, 1)
	ec := make(chan error, 1)
	go func() {
		conn := conn.(*internal.Conn)
		defer conn.Unlock()
		if _, err = w.WriteCmd(cmd, args...); err != nil {
			ec <- err
			return
		}
		var r *Result
		if r, err = read(conn, times); err != nil {
			ec <- err
			return
		}

		rc <- r
	}()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err = <-ec:
		return nil, err
	case r := <-rc:
		return r, nil
	}
}
