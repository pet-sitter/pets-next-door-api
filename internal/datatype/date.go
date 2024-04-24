package datatype

import (
	"database/sql"
	"database/sql/driver"
	"time"

	"cloud.google.com/go/civil"
)

type Date civil.Date

func (date *Date) Scan(value interface{}) (err error) {
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(value)
	*date = Date(civil.DateOf(nullTime.Time))
	return err
}

func (date Date) Value() (driver.Value, error) {
	return time.Date(date.Year, date.Month, date.Day, 0, 0, 0, 0, time.UTC), nil
}

func (date Date) MarshalJSON() ([]byte, error) {
	marshalled := make([]byte, 0)
	text, err := civil.Date(date).MarshalText()
	marshalled = append(marshalled, byte('"'))
	marshalled = append(marshalled, text...)
	marshalled = append(marshalled, byte('"'))
	return marshalled, err
}

func (date *Date) UnmarshalJSON(b []byte) error {
	c := civil.Date{}
	err := c.UnmarshalText(b[1 : len(b)-1])
	*date = Date(c)
	return err
}

func DateOf(t time.Time) Date {
	return Date(civil.DateOf(t))
}

func ParseDate(s string) (Date, error) {
	c, err := civil.ParseDate(s)
	return Date(c), err
}

func (date Date) String() string {
	return civil.Date(date).String()
}
