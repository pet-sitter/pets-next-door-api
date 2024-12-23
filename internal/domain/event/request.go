package event

import (
	"time"

	"github.com/google/uuid"
)

type CreateRequest struct {
	BaseCreateRequest
	RecurringPeriod *EventRecurringPeriod `json:"recurringPeriod,omitempty"`
}

type BaseCreateRequest struct {
	EventType       EventType    `json:"type"`
	Name            string       `json:"name"`
	Description     string       `json:"description"`
	MediaID         *uuid.UUID   `json:"mediaId,omitempty"`
	Topics          []EventTopic `json:"topics"`
	MaxParticipants *int         `json:"maxParticipants,omitempty"`
	Fee             int          `json:"fee"`
	StartAt         *time.Time   `json:"startAt,omitempty"`
}

type ShortTermCreateRequest struct {
	BaseCreateRequest
}

type RecurringCreateRequest struct {
	BaseCreateRequest
	RecurringPeriod EventRecurringPeriod `json:"recurringPeriod"`
}

type UpdateRequest struct {
	BaseUpdateRequest
	RecurringPeriod *EventRecurringPeriod `json:"recurringPeriod,omitempty"`
}

type BaseUpdateRequest struct {
	Name            string        `json:"name"`
	Description     string        `json:"description"`
	MediaID         uuid.NullUUID `json:"mediaId,omitempty"`
	Topics          []EventTopic  `json:"topics"`
	MaxParticipants *int          `json:"maxParticipants,omitempty"`
	Fee             int           `json:"fee"`
	StartAt         *time.Time    `json:"startAt,omitempty"`
}

type ShortTermUpdateRequest struct {
	BaseUpdateRequest
}

type RecurringUpdateRequest struct {
	BaseUpdateRequest
	ReccuringPeriod EventRecurringPeriod `json:"recurringPeriod"`
}
