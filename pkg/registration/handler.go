package registration

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dgraph-io/badger/v3"
	"github.com/gorilla/mux"
)

type Handler struct {
	DB                  *badger.DB
	RegistrationService *Service
}

func NewRegistrationHandler(db *badger.DB) *Handler {
	registrationRepository := NewClientRepository(db)
	registrationService := NewRegistrationService(registrationRepository)
	return &Handler{DB: db, RegistrationService: registrationService}
}

func (registrationHandler *Handler) SetupRoutes(r *mux.Router) {
	// 2.1 create (connect) new client
	r.HandleFunc("/", registrationHandler.createClientHandler).Methods(http.MethodPost)

	// 2.2 get the status of a client
	//r.HandleFunc("/{clientId}", clientHandler.getClientHandler).Methods(http.MethodGet)

	// 2.3 get all clients
	//r.HandleFunc("/", clientHandler.getAllClientsHandler).Methods(http.MethodGet)

	// 2.4 disconnect client
	//r.HandleFunc("/{clientId}", clientHandler.disconnectClient).Methods(http.MethodPost)

	// 3.1 registration root directory
	r.HandleFunc("/{clientId}/roots", registrationHandler.registerRootDir).Methods(http.MethodPost)

	// 3.2 get registered root directory
	r.HandleFunc("/{clientId}/roots/{rootId}", registrationHandler.getRegisteredRootDir).Methods(http.MethodGet)

	// 3.3 update reigstered root directory
	r.HandleFunc("/{clientId}/roots/{rootId}", registrationHandler.updateRegisteredRootDir).Methods(http.MethodPatch)

	// 3.4 unregister root directory
	r.HandleFunc("/{clientId}/roots/{rootId}", registrationHandler.deleteRegisteredRootDir).Methods(http.MethodPost)
}

// createClientHandler implements the API of 2.1 create (connect new client)
func (registrationHandler *Handler) createClientHandler(w http.ResponseWriter, r *http.Request) {
	var request RegisterClientRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Faield to create client because of wrong request data", http.StatusBadRequest)
		log.Fatalf("Error while decoding request data: %s", err)
		return
	}

	uuid, err := registrationHandler.RegistrationService.CreateNewClient(&request.Ip)
	if err != nil {
		http.Error(w, "Faield to create client", http.StatusInternalServerError)
		log.Fatalf("Error while returning response data: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(uuid)
}

func (registrationHandler *Handler) registerRootDir(w http.ResponseWriter, r *http.Request) {
	var request RegisterRootDirRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Faield to create client because of wrong request data", http.StatusBadRequest)
		log.Fatalf("Error while decoding request data: %s", err)
		return
	}

	message, err := registrationHandler.RegistrationService.RegisterRootDir(request)
	if err != nil {
		http.Error(w, "Faield to registration root directory", http.StatusInternalServerError)
		log.Fatalf("Error while returning response data: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}

func (registrationHandler *Handler) getRegisteredRootDir(w http.ResponseWriter, r *http.Request) {

}

func (registrationHandler *Handler) updateRegisteredRootDir(w http.ResponseWriter, r *http.Request) {

}

func (registrationHandler *Handler) deleteRegisteredRootDir(w http.ResponseWriter, r *http.Request) {

}
