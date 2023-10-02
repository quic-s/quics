package http3

import (
	"github.com/gorilla/mux"
	"github.com/quic-s/quics/pkg/core/server"
)

type ServerHandler struct {
	ServerService server.Service
}

func NewServerHandler(serverService server.Service) *ServerHandler {
	return &ServerHandler{
		ServerService: serverService,
	}
}

func (serverHandler *ServerHandler) SetupRoutes(r *mux.Router) {

}
