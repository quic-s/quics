package sync

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/gorilla/mux"
)

type Handler struct {
	db          *badger.DB
	syncService *Service
}

func NewSyncHandler(db *badger.DB) *Handler {
	syncRepository := NewSyncRepository(db)
	syncService := NewSyncService(syncRepository)
	return &Handler{db: db, syncService: syncService}
}

func (syncHandler *Handler) SetupRoutes(r *mux.Router) {

}
