package http3

import (
	"github.com/gorilla/mux"
	"github.com/quic-s/quics/pkg/core/sync"
)

type SyncHandler struct {
	syncService sync.Service
}

func NewSyncHandler(syncService sync.Service) *SyncHandler {
	return &SyncHandler{
		syncService: syncService,
	}
}

func (syncHandler *SyncHandler) SetupRoutes(r *mux.Router) {

}
