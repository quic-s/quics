package sharing

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/gorilla/mux"
)

type Handler struct {
	db             *badger.DB
	sharingService *Service
}

func NewSharingHandler(db *badger.DB) *Handler {
	sharingRepository := NewSharingRepository(db)
	sharingService := NewSharingService(sharingRepository)
	return &Handler{db: db, sharingService: sharingService}
}

func (sharingHandler *Handler) SetupRoutes(r *mux.Router) {

}
