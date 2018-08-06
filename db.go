package redis

import (
	"net"

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

func (db *DB) Begin() (*Tx, error) {
	row := db.QueryRow("MULTI")
	return newTx(row)
}

func (db *DB) Exec(cmd string, args ...interface{}) error {
	var dump internal.Dump
	return db.QueryRow(cmd, args...).Scan(&dump)
}

func (db *DB) Ping(args ...interface{}) (string, error) {
	row := db.QueryRow("PING", args...)
	var pong string
	if err := row.Scan(&pong); err != nil {
		return "", err
	}
	return pong, nil
}

func (db *DB) Query(cmd string, args ...interface{}) (*Rows, error) {
	conn, err := db.do([]byte(cmd), args...)
	if err != nil {
		return nil, err
	}
	return newRows(conn, false), nil
}

func (db *DB) QueryRow(cmd string, args ...interface{}) *Row {
	rows, err := db.Query(cmd, args...)
	return &Row{
		err:  err,
		rows: rows,
	}
}

func (db *DB) SetMaxIdleConns(maxConns int) {
	db.pool.SetMaxIdleConns(maxConns)
}

func (db *DB) SetMaxOpenConns(maxConns int) {
	db.pool.SetMaxOpenConns(maxConns)
}

func (db *DB) do(cmd []byte, args ...interface{}) (net.Conn, error) {
	conn, err := db.pool.Get()
	if err != nil {
		return nil, err
	}
	if _, err = conn.Write(internal.Parse(cmd, args...)); err != nil {
		return nil, err
	}
	return conn, nil
}
