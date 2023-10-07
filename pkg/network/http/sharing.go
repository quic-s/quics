package http

import (
	"github.com/gorilla/mux"
	"github.com/quic-s/quics/pkg/core/sharing"
)

type SharingHandler struct {
	sharingService sharing.Service
}

func NewSharingHandler(sharingService sharing.Service) *SharingHandler {
	return &SharingHandler{
		sharingService: sharingService,
	}
}

func (handler *SharingHandler) SetupRoutes(r *mux.Router) {

}
