package datatype

import (
	"github.com/google/uuid"
	pnd "github.com/pet-sitter/pets-next-door-api/api"
)

func NewV7() (uuid.UUID, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return uuid.Nil, pnd.ErrUnknown(err)
	}

	return id, nil
}

func NewUUIDV7() uuid.UUID {
	id, _ := NewV7()
	return id
}
