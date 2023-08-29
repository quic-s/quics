package metadata

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/gorilla/mux"
)

type Handler struct {
	db              *badger.DB
	metadataService *Service
}

func NewMetadataHandler(db *badger.DB) *Handler {
	metadataRepository := NewMetadataRepository(db)
	metadataService := NewMetadataService(metadataRepository)
	return &Handler{db: db, metadataService: metadataService}
}

func (metadataHandler *Handler) SetupRoutes(r *mux.Router) {

}
