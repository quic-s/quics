package http3

import (
	"github.com/gorilla/mux"
	"github.com/quic-s/quics/pkg/core/history"
)

type HistoryHandler struct {
	historyService history.Service
}

func NewHistoryHandler(historyService history.Service) *HistoryHandler {
	return &HistoryHandler{
		historyService: historyService,
	}
}

func (historyHandler *HistoryHandler) SetupRoutes(r *mux.Router) {

}
