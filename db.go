package redis

import (
	"fmt"
	"net"
	"sync"

	"github.com/gbrlsnchs/redis/internal"
)

// DB is a bridge between a Redis database and a connection pool to access it.
type DB struct {
	pool *internal.Pool
}

// Open fills the connection pool.
func Open(host string, port int, maxConns int) (*DB, error) {
	p, err := internal.NewPool(fmt.Sprintf("%s:%d", host, port), maxConns)

	if err != nil {
		return nil, err
	}

	return &DB{pool: p}, nil
}

func (db *DB) Begin() (*Tx, error) {
	conn, err := db.do([]byte("MULTI"))

	if err != nil {
		return nil, err
	}

	return &Tx{conn: conn}, nil
}

// Exec executes a Redis command.
func (db *DB) Exec(cmd string, args ...interface{}) error {
	conn, err := db.do([]byte(cmd), args...)

	if err != nil {
		return err
	}

	return conn.Close()
}

// Ping sends a "PING" signal plus arguments to the Redis database.
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

	return &Rows{conn: conn, mu: &sync.RWMutex{}}, nil
}

func (db *DB) QueryRow(cmd string, args ...interface{}) *Row {
	conn, err := db.do([]byte(cmd), args...)

	return &Row{
		err: err,
		rows: &Rows{
			conn: conn,
			mu:   &sync.RWMutex{},
		},
	}
}

// SetMaxIdleConns limits the amount of available connections.
func (db *DB) SetMaxIdleConns(max int) {
	db.pool.SetMaxIdleConns(max)
}

func (db *DB) do(cmd []byte, args ...interface{}) (net.Conn, error) {
	conn, err := db.pool.Get()

	if err != nil {
		return nil, err
	}

	return internal.NewRequest(conn, []byte(cmd), args...)
}
