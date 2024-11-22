package event

import (
	"time"

	"github.com/google/uuid"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/media"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/user"
)

type View struct {
	ShortTermView
	RecurringPeriod *EventRecurringPeriod `json:"recurringPeriod,omitempty"`
}

type BaseView struct {
	ID          uuid.UUID               `json:"id"`
	EventType   EventType               `json:"type"`
	Author      user.WithoutPrivateInfo `json:"author"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Media       media.DetailView        `json:"media"`
	Topics      []EventTopic            `json:"topics"`
	CreatedAt   time.Time               `json:"createdAt"`
	UpdatedAt   time.Time               `json:"updatedAt"`
}

type ShortTermView struct {
	BaseView
}
type RecurringView struct {
	BaseView
	RecurringPeriod EventRecurringPeriod `json:"recurringPeriod"`
}
