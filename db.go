package redis

import (
	"context"
	"errors"
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

func (db *DB) Subscribe(channels ...interface{}) (*Subscription, error) {
	return db.SubscribeContext(context.Background(), channels...)
}

func (db *DB) SubscribeContext(ctx context.Context, channels ...interface{}) (*Subscription, error) {
	if len(channels) == 0 {
		return nil, errors.New("redis: no channels to subscribe")
	}
	// Since subscriptions are meant to be long-running connections,
	// getting a connection from the pool seems unnecessary as it might
	// end up stealing ready-to-go connections from accesses meant to be quick.
	// When a connection spawned by DialContext is closed, the pool tries to store
	// it back in its inner connection pool, so this is not a simple connection at all.
	conn, err := db.p.DialContext(ctx)
	if err != nil {
		return nil, err
	}
	length := len(channels)
	w := newWriter(conn, length)
	r, err := w.send(ctx, "SUBSCRIBE", channels...)
	if err != nil {
		return nil, err
	}

	sub := newSubscription(ctx, length)
	for i := 0; i < len(r.values); i += 3 {
		channel := Value(r.values[i+1]).String()
		sub.channels[channel] = make(chan Value)
	}
	go func() {
		rr := internal.NewReader(conn)
		for {
			values, _ := rr.ReadFull()
			channel := Value(values[1])
			message := Value(values[2])
			sub.channels[channel.String()] <- message
		}
	}()
	return sub, nil
}
