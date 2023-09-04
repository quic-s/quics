package server

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/quic-s/quics/pkg/registration"
)

type Repository struct {
	DB *badger.DB
}

type RepositoryInterface interface {
	SetDefaultPassword(newId []byte, client registration.Client)
	UpdateDefaultPassword(id string) (*registration.Client, error)
}
