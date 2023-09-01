package server

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/quic-s/quics/pkg/client"
)

type Repository struct {
	DB *badger.DB
}

type RepositoryInterface interface {
	SetDefaultPassword(newId []byte, client client.Client)
	UpdateDefaultPassword(id string) (*client.Client, error)
}
