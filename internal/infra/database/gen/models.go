// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package databasegen

import (
	"database/sql"
	"encoding/json"
	"time"
)

type BasePost struct {
	ID        int32
	Title     sql.NullString
	Content   sql.NullString
	AuthorID  sql.NullInt64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
}

type Breed struct {
	ID        int32
	Name      string
	PetType   string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
}

type Medium struct {
	ID        int32
	MediaType string
	Url       string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
}

type Pet struct {
	ID             int32
	OwnerID        int64
	Name           string
	PetType        string
	Sex            string
	Neutered       bool
	Breed          string
	BirthDate      time.Time
	WeightInKg     string
	AdditionalNote sql.NullString
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      sql.NullTime
	ProfileImageID sql.NullInt64
	Remarks        string
}

type ResourceMedium struct {
	ID           int32
	MediaID      sql.NullInt64
	ResourceID   sql.NullInt64
	ResourceType sql.NullString
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    sql.NullTime
}

type SosCondition struct {
	ID        int32
	Name      sql.NullString
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
}

type SosDate struct {
	ID          int32
	DateStartAt sql.NullTime
	DateEndAt   sql.NullTime
	CreatedAt   sql.NullTime
	UpdatedAt   sql.NullTime
	DeletedAt   sql.NullTime
}

type SosPost struct {
	ID          int32
	Title       sql.NullString
	Content     sql.NullString
	AuthorID    sql.NullInt64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   sql.NullTime
	Reward      sql.NullString
	CareType    sql.NullString
	CarerGender sql.NullString
	RewardType  sql.NullString
	ThumbnailID sql.NullInt64
}

type SosPostsCondition struct {
	ID             int32
	SosPostID      sql.NullInt64
	SosConditionID sql.NullInt64
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      sql.NullTime
}

type SosPostsDate struct {
	ID         int32
	SosPostID  sql.NullInt64
	SosDatesID sql.NullInt64
	CreatedAt  sql.NullTime
	UpdatedAt  sql.NullTime
	DeletedAt  sql.NullTime
}

type SosPostsPet struct {
	ID        int32
	SosPostID sql.NullInt64
	PetID     sql.NullInt64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
}

type User struct {
	ID             int32
	Email          string
	Password       string
	Nickname       string
	Fullname       string
	FbProviderType sql.NullString
	FbUid          sql.NullString
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      sql.NullTime
	ProfileImageID sql.NullInt64
}

type VCondition struct {
	SosPostID      sql.NullInt64
	ConditionsInfo json.RawMessage
}

type VMediaForSosPost struct {
	SosPostID sql.NullInt64
	MediaInfo json.RawMessage
}

type VPetsForSosPost struct {
	SosPostID   sql.NullInt64
	PetTypeList interface{}
	PetsInfo    json.RawMessage
}

type VSosPost struct {
	ID                  int32
	Title               sql.NullString
	Content             sql.NullString
	Reward              sql.NullString
	RewardType          sql.NullString
	CareType            sql.NullString
	CarerGender         sql.NullString
	ThumbnailID         sql.NullInt64
	AuthorID            sql.NullInt64
	CreatedAt           time.Time
	UpdatedAt           time.Time
	EarliestDateStartAt interface{}
	Dates               json.RawMessage
}