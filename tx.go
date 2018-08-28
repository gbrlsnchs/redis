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
	w := internal.NewWriter(conn)
	if _, err := send(ctx, w, conn, 1, "MULTI"); err != nil {
		return nil, err
	}
	return &Tx{conn, ctx, w}, nil
}

func (tx *Tx) Discard() (*Result, error) {
	defer tx.conn.Close()
	return send(tx.ctx, tx.w, tx.conn, 1, "DISCARD")
}

func (tx *Tx) Exec() (*Result, error) {
	defer tx.conn.Close()
	return send(tx.ctx, tx.w, tx.conn, 1, "EXEC")
}

func (tx *Tx) Queue(cmd string, args ...interface{}) (*Result, error) {
	return tx.QueueContext(tx.ctx, cmd, args...)
}

func (tx *Tx) QueueContext(ctx context.Context, cmd string, args ...interface{}) (*Result, error) {
	return send(ctx, tx.w, tx.conn, 1, cmd, args...)
}
