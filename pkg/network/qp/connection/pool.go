package connection

import (
	"fmt"

	qp "github.com/quic-s/quics-protocol"
)

type Pool struct {
	Conns map[string]*qp.Connection
}

func NewnPool() *Pool {
	return &Pool{}
}

func (cp *Pool) UpdateConnection(uuid string, conn *qp.Connection) error {
	cp.Conns[uuid] = conn
	return nil
}

func (cp *Pool) GetConnection(uuid string) (*qp.Connection, error) {
	if conn, exists := cp.Conns[uuid]; exists {
		return conn, nil
	}
	return nil, fmt.Errorf("connection does not exist")
}

func (cp *Pool) GetConnections(uuid []string) ([]*qp.Connection, error) {
	conns := make([]*qp.Connection, 0)
	for _, value := range uuid {
		if conn, exists := cp.Conns[value]; exists {
			conns = append(conns, conn)
		}
	}
	return conns, nil
}
