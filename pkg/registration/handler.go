package registration

import (
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

}
