package internal

import (
	"context"
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
	queue        chan struct{}
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

// DialContext tries to stablish a connection before a context is canceled.
func (p *Pool) DialContext(ctx context.Context) (net.Conn, error) {
	var err error
	if err = p.wait(ctx); err != nil {
		return nil, err
	}
	var d net.Dialer
	c, err := d.DialContext(ctx, "tcp", p.addr.String())
	if err != nil {
		return nil, err
	}
	return &Conn{Conn: c, pool: p.c, queue: p.queue, done: make(chan struct{}, 1)}, nil
}

// Get retrieves a new connection if any is available,
// otherwise it spawns a new connection.
func (p *Pool) Get(ctx context.Context) (net.Conn, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	select {
	case conn := <-p.c:
		return conn, nil
	default:
		return p.DialContext(ctx)
	}
}

// SetMaxIdleConns limits the amount of idle connections in the pool.
func (p *Pool) SetMaxIdleConns(maxConns int) {
	p.maxIdleConns = p.resolveConns(maxConns)
	p.resetChan()
}

// SetMaxOpenConns limits the amount of open connections.
func (p *Pool) SetMaxOpenConns(maxConns int) {
	p.maxOpenConns = maxConns
	p.maxIdleConns = p.resolveConns(p.maxIdleConns)
	p.resetChan()
	if p.maxOpenConns > 0 {
		p.queue = make(chan struct{}, p.maxOpenConns)
		return
	}
	p.queue = nil
}

func (p *Pool) resetChan() {
	// Don't reuse any connections.
	if p.maxIdleConns <= 0 {
		p.c = nil
		return
	}
	// Reset channel only if size has changed.
	if p.maxIdleConns != cap(p.c) {
		if p.c != nil {
			close(p.c)
		}
		p.c = make(chan net.Conn, p.maxIdleConns)
	}
}

func (p *Pool) resolveConns(maxConns int) int {
	if p.maxOpenConns > 0 && maxConns > p.maxOpenConns {
		return p.maxOpenConns
	}
	return maxConns
}

func (p *Pool) wait(ctx context.Context) error {
	if p.maxOpenConns > 0 {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case p.queue <- struct{}{}:
			return nil
		}
	}
	return nil
}
