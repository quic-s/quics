package server

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/quic-s/quics/pkg/types"
	"log"
)

const (
	PASSWORD string = "PASSWORD"
)

type Repository struct {
	DB *badger.DB
}

type RepositoryInterface interface {
	SetPassword(password string, client types.Client)
	UpdateDefaultPassword(id string) (*types.Client, error)
}

func NewServerRepository(db *badger.DB) *Repository {
	return &Repository{DB: db}
}

func (serverRepository *Repository) SetPassword(password string) {
	err := serverRepository.DB.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(PASSWORD), []byte(password))
		return err
	})
	if err != nil {
		log.Panicln("Error while setting server: ", err)
	}
}

func (serverRepository *Repository) GetPassword() string {
	var password string

	err := serverRepository.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(PASSWORD))
		if err != nil {
			log.Println("quis: Error while getting password: ", err)
		}

		err = item.Value(func(val []byte) error {
			server := &types.Server{}
			if err := server.Decode(val); err != nil {
				return err
			}

			password = server.Password
			return nil
		})

		return nil
	})
	if err != nil {
		log.Fatalln("quics: Error while executing GetPassword() with database.")
	}

	return password
}
