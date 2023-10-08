package http3

import (
	"github.com/gorilla/mux"
	"github.com/quic-s/quics/pkg/core/sync"
)

type SyncHandler struct {
	SyncService sync.Service
}

func NewSyncHandler(syncService sync.Service) *SyncHandler {
	return &SyncHandler{
		SyncService: syncService,
	}
}

func (syncHandler *SyncHandler) SetupRoutes(r *mux.Router) {

}
