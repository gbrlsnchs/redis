package redis

import (
	"context"

	"github.com/gbrlsnchs/redis/internal"
)

type DB struct {
	pool *internal.Pool
}

func Open(addr string) (*DB, error) {
	p, err := internal.NewPool(addr)
	if err != nil {
		return nil, err
	}
	return &DB{pool: p}, nil
}

func (db *DB) Multi() (*Tx, error) {
	return db.MultiTx(context.Background())
}

func (db *DB) MultiTx(ctx context.Context) (*Tx, error) {
	conn, err := db.pool.Get(ctx)
	if err != nil {
		return nil, err
	}
	tx, err := multi(ctx, conn)
	if err != nil {
		conn.Close()
		return nil, err
	}
	return tx, nil
}

func (db *DB) Ping(args ...interface{}) (string, error) {
	r, err := db.Send("PING", args...)
	if err != nil {
		return "", err
	}
	return r.String(), nil
}

func (db *DB) Send(cmd string, args ...interface{}) (*Result, error) {
	return db.SendContext(context.Background(), cmd, args...)
}

func (db *DB) SendContext(ctx context.Context, cmd string, args ...interface{}) (*Result, error) {
	conn, err := db.pool.Get(ctx)
	defer conn.Close()
	if err != nil {
		return nil, err
	}
	w := internal.NewWriter(conn)
	return send(ctx, w, conn, 1, cmd, args...)
}

func (db *DB) SetMaxIdleConns(maxConns int) {
	db.pool.SetMaxIdleConns(maxConns)
}

func (db *DB) SetMaxOpenConns(maxConns int) {
	db.pool.SetMaxOpenConns(maxConns)
}
