package internal

import "net"

type poolConn struct {
	net.Conn
	poolC chan<- net.Conn
}

func (pc *poolConn) Close() error {
	select {
	case pc.poolC <- pc:
		return nil

	default:
		return pc.Conn.Close()
	}
}
