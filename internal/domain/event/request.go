package event

type CreateRequest struct {
	BaseCreateRequest
	RecurringPeriod *EventRecurringPeriod `json:"recurringPeriod,omitempty"`
}

type BaseCreateRequest struct {
	EventType       EventType    `json:"type"`
	AuthorID        int          `json:"authorId"`
	Name            string       `json:"name"`
	Description     string       `json:"description"`
	MediaID         int          `json:"mediaId"`
	Topics          []EventTopic `json:"topics"`
	GenderCondition string       `json:"genderCondition" enums:"male,female,all"`
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
	AuthorID        int          `json:"authorId"`
	Name            string       `json:"name"`
	Description     string       `json:"description"`
	MediaID         int          `json:"mediaId"`
	Topics          []EventTopic `json:"topics"`
	GenderCondition string       `json:"genderCondition" enums:"male,female,all"`
}

type ShortTermUpdateRequest struct {
	BaseUpdateRequest
}

type RecurringUpdateRequest struct {
	BaseUpdateRequest
	ReccuringPeriod EventRecurringPeriod `json:"recurringPeriod"`
}
