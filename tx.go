package redis

import (
	"net"

	"github.com/gbrlsnchs/redis/internal"
)

type Tx struct {
	conn net.Conn
}

func newTx(row *Row) (*Tx, error) {
	row.rows.isMulti = true
	// Dump "MULTI" command response.
	var dump internal.Dump
	if err := row.Scan(&dump); err != nil {
		return nil, err
	}
	return &Tx{conn: row.rows.conn}, nil
}

func (tx *Tx) Commit() error {
	if err := tx.Exec("EXEC"); err != nil {
		return tx.Rollback()
	}
	return tx.conn.Close()
}

func (tx *Tx) Exec(cmd string, args ...interface{}) error {
	var dump internal.Dump
	return tx.QueryRow(cmd, args...).Scan(&dump)
}

func (tx *Tx) Query(cmd string, args ...interface{}) (*Rows, error) {
	if err := tx.do([]byte(cmd), args...); err != nil {
		return nil, err
	}
	return newRows(tx.conn, true), nil
}

func (tx *Tx) QueryRow(cmd string, args ...interface{}) *Row {
	rows, err := tx.Query(cmd, args...)
	return &Row{
		err:  err,
		rows: rows,
	}
}

func (tx *Tx) Rollback() error {
	if err := tx.Exec("DISCARD"); err != nil {
		return err
	}
	return tx.conn.Close()
}

func (tx *Tx) do(cmd []byte, args ...interface{}) error {
	if _, err := tx.conn.Write(internal.Parse(cmd, args...)); err != nil {
		return err
	}
	return nil
}
