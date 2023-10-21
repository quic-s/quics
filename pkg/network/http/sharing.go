package http

import (
	"net/http"

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

func (sh *SharingHandler) SetupRoutes(mux *http.ServeMux) http.Handler {
	return mux
}
