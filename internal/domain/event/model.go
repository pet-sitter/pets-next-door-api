package event

import (
	"github.com/pet-sitter/pets-next-door-api/internal/domain/media"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/user"
	databasegen "github.com/pet-sitter/pets-next-door-api/internal/infra/database/gen"
)

type Event struct {
	Event  databasegen.Event
	Author user.WithoutPrivateInfo
	Media  *media.DetailView
}

func ToEvent(
	eventData databasegen.Event,
	authorData user.WithoutPrivateInfo,
	mediaData *media.DetailView,
) *Event {
	return &Event{
		Event:  eventData,
		Author: authorData,
		Media:  mediaData,
	}
}
