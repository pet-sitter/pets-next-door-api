package event

import "github.com/google/uuid"

type DetailView struct {
	ID uuid.UUID `json:"id"`
}

type ListView struct {
	ID uuid.UUID `json:"id"`
}
