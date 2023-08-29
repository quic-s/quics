package history

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/gorilla/mux"
)

type Handler struct {
	db             *badger.DB
	historyService *Service
}

func NewHistoryHandler(db *badger.DB) *Handler {
	historyRepository := NewHistoryRepository(db)
	historyService := NewHistoryService(historyRepository)
	return &Handler{db: db, historyService: historyService}
}

func (historyHandler *Handler) SetupRoutes(r *mux.Router) {

}
