package internal

import (
	"net"
)

type conn struct {
	net.Conn
	pool *Pool
}

func newConn(p *Pool) (*conn, error) {
	p.mu.RLock()
	// If there's no room for a new open connection,
	// try to retrieve one until there's room.
	if p.maxOpenConns > 0 && p.openConns >= p.maxOpenConns {
		p.mu.RUnlock()
		return newConn(p)
	}
	p.mu.RUnlock()

	c, err := net.DialTCP("tcp", nil, p.addr)
	if err != nil {
		return nil, err
	}
	p.mu.Lock()
	p.openConns++
	p.mu.Unlock()
	return &conn{Conn: c, pool: p}, nil
}

// Close either returns the connection to the pool or,
// if the pool is full, it simply closes the connection.
func (c *conn) Close() error {
	// Try to send the connection back to the pool,
	// otherwise simply close it.
	select {
	case c.pool.c <- c:
		return nil
	default:
		c.pool.mu.Lock()
		c.pool.openConns--
		c.pool.mu.Unlock()
		return c.Conn.Close()
	}
}
