package redis

import (
	"context"
	"net"
)

type Tx struct {
	conn net.Conn
	ctx  context.Context
	w    *writer
}

func multi(ctx context.Context, conn net.Conn) (*Tx, error) {
	w := newWriter(conn, 1)
	if _, err := w.send(ctx, "MULTI"); err != nil {
		return nil, err
	}
	return &Tx{conn, ctx, w}, nil
}

func (tx *Tx) Discard() (*Result, error) {
	defer tx.conn.Close()
	return tx.w.send(tx.ctx, "DISCARD")
}

func (tx *Tx) Exec() (*Result, error) {
	defer tx.conn.Close()
	return tx.w.send(tx.ctx, "EXEC")
}

func (tx *Tx) Queue(cmd string, args ...interface{}) (*Result, error) {
	return tx.QueueContext(tx.ctx, cmd, args...)
}

func (tx *Tx) QueueContext(ctx context.Context, cmd string, args ...interface{}) (*Result, error) {
	return tx.w.send(ctx, cmd, args...)
}
