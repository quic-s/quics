package http3

import (
	"github.com/gorilla/mux"
	"github.com/quic-s/quics/pkg/core/history"
)

type HistoryHandler struct {
	HistoryService history.Service
}

func NewHistoryHandler(historyService history.Service) *HistoryHandler {
	return &HistoryHandler{
		HistoryService: historyService,
	}
}

func (historyHandler *HistoryHandler) SetupRoutes(r *mux.Router) {

}
