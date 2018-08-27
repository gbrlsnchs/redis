package redis

import (
	"context"
	"net"

	"github.com/gbrlsnchs/redis/internal"
)

type Tx struct {
	conn net.Conn
	ctx  context.Context
	w    *internal.Writer
}

func multi(ctx context.Context, conn net.Conn) (*Tx, error) {
	var err error
	w := internal.NewWriter(conn)
	if _, err = w.WriteCmd("MULTI"); err != nil {
		conn.Close()
		return nil, err
	}
	if _, err = read(conn); err != nil {
		return nil, err
	}
	return &Tx{conn, ctx, w}, nil
}

func (tx *Tx) Discard() (*Result, error) {
	defer tx.conn.Close()
	return tx.send(tx.ctx, "DISCARD")
}

func (tx *Tx) Exec() (*Result, error) {
	defer tx.conn.Close()
	return tx.send(tx.ctx, "EXEC")
}

func (tx *Tx) Queue(cmd string, args ...interface{}) (*Result, error) {
	return tx.QueueContext(tx.ctx, cmd, args...)
}

func (tx *Tx) QueueContext(ctx context.Context, cmd string, args ...interface{}) (*Result, error) {
	return tx.send(ctx, cmd, args...)
}

func (tx *Tx) send(ctx context.Context, cmd string, args ...interface{}) (*Result, error) {
	var err error
	rc := make(chan *Result, 1)
	ec := make(chan error, 1)
	go func() {
		if _, err = tx.w.WriteCmd(cmd, args...); err != nil {
			ec <- err
			return
		}
		var r *Result
		if r, err = read(tx.conn); err != nil {
			ec <- err
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
