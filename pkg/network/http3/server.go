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

func (sh *ServerHandler) SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/server/stop", sh.StopRestServer)
	mux.HandleFunc("/api/v1/server/listen", sh.ListenProtocol)
	mux.HandleFunc("/api/v1/server/logs/clients", sh.ShowClientLogs)
	mux.HandleFunc("/api/v1/server/logs/directories", sh.ShowDirLogs)
	mux.HandleFunc("/api/v1/server/logs/files", sh.ShowFileLogs)
	mux.HandleFunc("/api/v1/server/logs/histories", sh.ShowHistoryLogs)
	mux.HandleFunc("/api/v1/server/remove/clients", sh.RemoveClient)
	mux.HandleFunc("/api/v1/server/remove/directories", sh.RemoveDir)
	mux.HandleFunc("/api/v1/server/remove/files", sh.RemoveFile)
	mux.HandleFunc("/api/v1/server/rollback/files", sh.RollbackFile)
	mux.HandleFunc("/api/v1/server/download/files", sh.DownloadFile)
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

func (sh *ServerHandler) ShowClientLogs(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		all := r.URL.Query().Get("all")
		id := r.URL.Query().Get("id")

		err := sh.ServerService.ShowClientLogs(all, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (sh *ServerHandler) ShowDirLogs(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		all := r.URL.Query().Get("all")
		id := r.URL.Query().Get("id")

		err := sh.ServerService.ShowDirLogs(all, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (sh *ServerHandler) ShowFileLogs(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		all := r.URL.Query().Get("all")
		id := r.URL.Query().Get("id")

		err := sh.ServerService.ShowFileLogs(all, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (sh *ServerHandler) ShowHistoryLogs(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		all := r.URL.Query().Get("all")
		id := r.URL.Query().Get("id")

		err := sh.ServerService.ShowHistoryLogs(all, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (sh *ServerHandler) RemoveClient(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		all := r.URL.Query().Get("all")
		id := r.URL.Query().Get("id")

		err := sh.ServerService.RemoveClient(all, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (sh *ServerHandler) RemoveDir(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		all := r.URL.Query().Get("all")
		id := r.URL.Query().Get("id")

		err := sh.ServerService.RemoveDir(all, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (sh *ServerHandler) RemoveFile(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		all := r.URL.Query().Get("all")
		id := r.URL.Query().Get("id")

		err := sh.ServerService.RemoveFile(all, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (sh *ServerHandler) RollbackFile(w http.ResponseWriter, r *http.Request) {
	// switch r.Method {
	// case "POST":
	// 	path := r.URL.Query().Get("path")
	// 	version := r.URL.Query().Get("version")

	// 	err := sh.ServerService.RollbackFile(path, version)
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}
	// }
}

func (sh *ServerHandler) DownloadFile(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		path := r.URL.Query().Get("path")
		version := r.URL.Query().Get("version")
		target := r.URL.Query().Get("target")

		err := sh.ServerService.DownloadFile(path, version, target)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
