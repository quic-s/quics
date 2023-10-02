package http3

import (
	"github.com/gorilla/mux"
	"github.com/quic-s/quics/pkg/core/registration"
)

type RegistrationHandler struct {
	registrationService registration.Service
}

func NewRegistrationHandler(registrationService registration.Service) *RegistrationHandler {
	return &RegistrationHandler{
		registrationService: registrationService,
	}
}

func (handler *RegistrationHandler) SetupRoutes(r *mux.Router) {

}
