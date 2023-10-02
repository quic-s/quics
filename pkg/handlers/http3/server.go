package http3

import (
	"github.com/gorilla/mux"
	"github.com/quic-s/quics/pkg/core/server"
)

type ServerHandler struct {
	serverService server.Service
}

func NewServerHandler(serverService server.Service) *ServerHandler {
	return &ServerHandler{
		serverService: serverService,
	}
}

func (serverHandler *ServerHandler) SetupRoutes(r *mux.Router) {

}
