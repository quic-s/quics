package server

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/gorilla/mux"
)

type Handler struct {
	db            *badger.DB
	ServerService *Service
}

func NewServerHandler(db *badger.DB) *Handler {
	serverRepository := NewServerRepository(db)
	serverService := NewServerService(serverRepository)
	return &Handler{db: db, ServerService: serverService}
}

func (serverHandler *Handler) SetupRoutes(r *mux.Router) {

}
