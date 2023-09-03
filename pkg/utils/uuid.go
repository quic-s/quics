package utils

import (
	"github.com/google/uuid"
)

// CreateUuid returns uuid after creating new uuid
// TODO: If possible, add the validation of uuid for not duplicating
func CreateUuid() string {
	u := uuid.New()
	return u.String()
}
