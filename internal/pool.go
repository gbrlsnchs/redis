package internal

import (
	"net"
	"sync"
)

const defaultMaxIdleConns = 2

// Pool is a connection pool to communicate with a Redis database.
type Pool struct {
	c            chan net.Conn
	addr         *net.TCPAddr
	maxOpenConns int
	maxIdleConns int
	openConns    int
	mu           *sync.RWMutex
}

// NewPool creates a new connection pool.
func NewPool(address string) (*Pool, error) {
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, err
	}
	p := &Pool{addr: addr, mu: &sync.RWMutex{}}
	p.SetMaxIdleConns(defaultMaxIdleConns)
	return p, nil
}

// Get retrieves a new connection if any is available,
// otherwise it spawns a new connection.
func (p *Pool) Get() (net.Conn, error) {
	select {
	case conn := <-p.c:
		return conn, nil
	default:
		return newConn(p)
	}
}

// SetMaxIdleConns limits the amount of idle connections in the pool.
func (p *Pool) SetMaxIdleConns(maxConns int) {
	p.maxIdleConns = p.resolveConns(maxConns)
	p.resetChan()
}

// SetMaxOpenConns limits the amount of openConns connections.
func (p *Pool) SetMaxOpenConns(maxConns int) {
	p.maxOpenConns = maxConns
	p.maxIdleConns = p.resolveConns(p.maxIdleConns)
	p.resetChan()
}

func (p *Pool) resetChan() {
	if p.maxIdleConns <= 0 {
		p.c = nil
		return
	}
	if p.maxIdleConns != cap(p.c) {
		p.c = make(chan net.Conn, p.maxIdleConns)
	}
}

func (p *Pool) resolveConns(maxConns int) int {
	if p.maxOpenConns > 0 && maxConns > p.maxOpenConns {
		return p.maxOpenConns
	}
	return maxConns
}
