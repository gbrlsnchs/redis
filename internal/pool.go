package internal

import "net"

// Pool is a connection pool to communicate with a Redis database.
type Pool struct {
	c       chan net.Conn
	tcpAddr *net.TCPAddr
}

// NewPool creates a new connection pool if "maxConns" is greater than 0.
// Otherwise, it will spawn a new connection everytime a client tries to
// retrieve a connection.
func NewPool(address string, maxConns int) (*Pool, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", address)

	if err != nil {
		return nil, err
	}

	p := &Pool{tcpAddr: tcpAddr}

	p.SetMaxIdleConns(maxConns)

	if err = p.Start(); err != nil {
		return nil, err
	}

	return p, nil
}

// Get retrieves a new connection if any is available,
// otherwise it spawns a new connection.
func (p *Pool) Get() (net.Conn, error) {
	select {
	case conn := <-p.c:
		return conn, nil

	default:
		return p.newPoolConn()
	}
}

// SetMaxIdleConns limits the amount of available connections.
func (p *Pool) SetMaxIdleConns(maxConns int) {
	if maxConns > 0 {
		p.c = make(chan net.Conn, maxConns)
	}
}

// Start fills the pool connection.
func (p *Pool) Start() error {
	for i := 0; i < cap(p.c); i++ {
		conn, err := p.newPoolConn()

		if err != nil {
			close(p.c)

			return err
		}

		select {
		case p.c <- conn:

		default:
			return conn.Close()
		}
	}

	return nil
}

func (p *Pool) newPoolConn() (net.Conn, error) {
	conn, err := net.DialTCP("tcp", nil, p.tcpAddr)

	if err != nil {
		return nil, err
	}

	return &poolConn{Conn: conn, poolC: p.c}, nil
}
