package redis

import (
	"context"
	"net"

	"github.com/gbrlsnchs/connpool"
)

type DB struct {
	p *connpool.Pool
}

func Open(address string) (*DB, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, err
	}
	return &DB{p: connpool.New("tcp", tcpAddr.String())}, nil
}

func (db *DB) Multi() (*Tx, error) {
	return db.MultiTx(context.Background())
}

func (db *DB) MultiTx(ctx context.Context) (*Tx, error) {
	conn, err := db.p.GetContext(ctx)
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
	conn, err := db.p.GetContext(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	w := newWriter(conn, 1)
	return w.send(ctx, cmd, args...)
}

func (db *DB) SetMaxIdleConns(maxConns int) {
	db.p.SetMaxIdleConns(maxConns)
}

func (db *DB) SetMaxOpenConns(maxConns int) {
	db.p.SetMaxOpenConns(maxConns)
}
