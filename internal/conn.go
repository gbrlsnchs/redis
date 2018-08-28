package internal

import (
	"net"
)

type Conn struct {
	net.Conn
	queue <-chan struct{}
	pool  chan<- net.Conn
	done  chan struct{}
}

// Close either returns the connection to the pool or,
// if the pool is full, it simply closes the connection.
//
// If it can't return the connection to the pool,
// it tries to dequeue the connection in order to respect
// the max open connections limit.
func (c *Conn) Close() error {
	<-c.done
	// Try to send the connection back to the pool,
	// otherwise simply close it.
	select {
	case c.pool <- c:
		return nil
	default:
	}
	select {
	case <-c.queue:
	default:
	}
	return c.Conn.Close()
}

func (c *Conn) Unlock() {
	c.done <- struct{}{}
}
