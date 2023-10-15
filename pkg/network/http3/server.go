package http3

import (
	"net/http"

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

func (sh *ServerHandler) SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/server/listen", sh.ListenProtocol)
	mux.HandleFunc("/api/v1/server/stop", sh.StopRestServer)

	return mux
}

func (sh *ServerHandler) ListenProtocol(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		err := sh.ServerService.ListenProtocol()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (sh *ServerHandler) StopRestServer(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		err := sh.ServerService.StopServer()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
