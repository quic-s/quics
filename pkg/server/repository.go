package server

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/quic-s/quics/pkg/register"
)

type Repository struct {
	DB *badger.DB
}

type RepositoryInterface interface {
	SetDefaultPassword(newId []byte, client register.Client)
	UpdateDefaultPassword(id string) (*register.Client, error)
}
