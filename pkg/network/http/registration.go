package http

import (
	"github.com/gorilla/mux"
	"github.com/quic-s/quics/pkg/core/registration"
)

type RegistrationHandler struct {
	RegistrationService registration.Service
}

func NewRegistrationHandler(registrationService registration.Service) *RegistrationHandler {
	return &RegistrationHandler{
		RegistrationService: registrationService,
	}
}

func (handler *RegistrationHandler) SetupRoutes(r *mux.Router) {

}
