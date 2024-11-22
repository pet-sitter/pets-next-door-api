package event

import (
	"encoding/json"
	"errors"
	"strings"
)

type EventType string

const (
	// 단기
	ShortTerm EventType = "SHORT_TERM"
	// 정기
	Recurring EventType = "RECURRING"
)

func (e EventType) IsValid() bool {
	switch e {
	case ShortTerm, Recurring:
		return true
	}
	return false
}

type EventTopic string

const (
	ETC EventTopic = "ETC"
)

func (e EventTopic) IsValid() bool {
	// TODO: Add more topics
	return e == ETC
}

type EventRecurringPeriod string

const (
	// 매일
	Daily EventRecurringPeriod = "DAILY"
	// 매주
	Weekly EventRecurringPeriod = "WEEKLY"
	// 2주에 한 번
	Biweekly EventRecurringPeriod = "BIWEEKLY"
	// 매달
	Monthly EventRecurringPeriod = "MONTHLY"
)

func (e EventRecurringPeriod) IsValid() bool {
	switch e {
	case Daily, Weekly, Biweekly, Monthly:
		return true
	}
	return false
}

//go:generate stringer -type=GenderCondition
type GenderCondition int

const (
	Male GenderCondition = iota
	Female
	All
)

func (g *GenderCondition) MarshalJSON() ([]byte, error) {
	return json.Marshal(strings.ToLower(g.String()))
}

func (g *GenderCondition) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	switch strings.ToLower(s) {
	case "male":
		*g = Male
	case "female":
		*g = Female
	case "all":
		*g = All
	default:
		return errors.New("invalid GenderCondition")
	}

	return nil
}
