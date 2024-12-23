package event

import (
	"time"

	"github.com/google/uuid"
	utils "github.com/pet-sitter/pets-next-door-api/internal/common"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/media"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/user"
)

type View struct {
	ShortTermView
	RecurringPeriod *EventRecurringPeriod `json:"recurringPeriod,omitempty"`
}

type BaseView struct {
	ID              uuid.UUID               `json:"id"`
	EventType       EventType               `json:"type"`
	Author          user.WithoutPrivateInfo `json:"author"`
	Name            string                  `json:"name"`
	Description     string                  `json:"description"`
	Media           media.DetailView        `json:"media"`
	Topics          []EventTopic            `json:"topics"`
	MaxParticipants *int                    `json:"maxParticipants,omitempty"`
	Fee             int                     `json:"fee"`
	StartAt         *time.Time              `json:"startAt,omitempty"`
	CreatedAt       time.Time               `json:"createdAt"`
	UpdatedAt       time.Time               `json:"updatedAt"`
}

type ShortTermView struct {
	BaseView
}
type RecurringView struct {
	BaseView
	RecurringPeriod EventRecurringPeriod `json:"recurringPeriod"`
}

func ToShortTermView(eventData *Event) ShortTermView {
	topics := make([]EventTopic, len(eventData.Event.Topics))
	for i, topic := range eventData.Event.Topics {
		topics[i] = EventTopic(topic)
	}
	view := ShortTermView{
		BaseView: BaseView{
			ID:              eventData.Event.ID,
			Author:          eventData.Author,
			EventType:       EventType(eventData.Event.EventType),
			Name:            eventData.Event.Name,
			Description:     eventData.Event.Description,
			Media:           *eventData.Media,
			Topics:          topics,
			MaxParticipants: utils.NullInt32ToIntPtr(eventData.Event.MaxParticipants),
			Fee:             int(eventData.Event.Fee),
			StartAt:         utils.NullTimeToTimePtr(eventData.Event.StartAt),
			CreatedAt:       eventData.Event.CreatedAt,
			UpdatedAt:       eventData.Event.UpdatedAt,
		},
	}

	return view
}
