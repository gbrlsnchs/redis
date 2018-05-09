package redis

import (
	"net"
	"sync"

	"github.com/gbrlsnchs/redis/internal"
)

// Tx is a transaction.
type Tx struct {
	conn net.Conn
}

// Commit atomically executes the transaction.
func (tx *Tx) Commit() error {
	conn, err := tx.do([]byte("EXEC"))

	if err != nil {
		return tx.Rollback()
	}

	return conn.Close()
}

// Exec executes a command without caring about the response content.
func (tx *Tx) Exec(cmd string, args ...interface{}) error {
	if _, err := tx.do([]byte(cmd), args...); err != nil {
		return err
	}

	return nil
}

func (tx *Tx) Query(cmd string, args ...interface{}) (*Rows, error) {
	conn, err := tx.do([]byte(cmd), args...)

	if err != nil {
		return nil, err
	}

	return &Rows{conn: conn, mu: &sync.RWMutex{}, multi: true}, nil
}

func (tx *Tx) QueryRow(cmd string, args ...interface{}) *Row {
	conn, err := tx.do([]byte(cmd), args...)

	return &Row{
		err: err,
		rows: &Rows{
			conn:  conn,
			mu:    &sync.RWMutex{},
			multi: true,
		},
	}
}

// Rollback aborts a transaction.
func (tx *Tx) Rollback() error {
	conn, err := tx.do([]byte("DISCARD"))

	if err != nil {
		return err
	}

	return conn.Close()
}

func (tx *Tx) do(cmd []byte, args ...interface{}) (net.Conn, error) {
	return internal.NewRequest(tx.conn, cmd, args...)
}
