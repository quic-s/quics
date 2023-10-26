package connection

import (
	"fmt"
	"sync"

	qp "github.com/quic-s/quics-protocol"
)

type Pool struct {
	connsMut sync.RWMutex
	Conns    map[string]*qp.Connection
}

func NewnPool() *Pool {
	return &Pool{
		connsMut: sync.RWMutex{},
		Conns:    map[string]*qp.Connection{},
	}
}

func (cp *Pool) UpdateConnection(uuid string, conn *qp.Connection) error {
	cp.connsMut.Lock()
	defer cp.connsMut.Unlock()
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
	conns := []*qp.Connection{}
	for _, value := range uuid {
		if conn, exists := cp.Conns[value]; exists {
			conns = append(conns, conn)
		}
	}
	return conns, nil
}

func (cp *Pool) DeleteConnection(uuid string) error {
	cp.connsMut.Lock()
	defer cp.connsMut.Unlock()
	delete(cp.Conns, uuid)
	return nil
}
