package registration

import (
	"github.com/dgraph-io/badger/v3"
)

type Handler struct {
	DB                  *badger.DB
	RegistrationService *Service
}

func NewRegistrationHandler(db *badger.DB) *Handler {
	registrationRepository := NewRegistrationRepository(db)
	registrationService := NewRegistrationService(registrationRepository)
	return &Handler{DB: db, RegistrationService: registrationService}
}
