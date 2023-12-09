package sos_post

import (
	"time"

	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"
)

type CareType string
type CarerGender string
type RewardAmount string

const (
	CareTypeFoster   CareType = "foster"
	CareTypeVisiting CareType = "visiting"
)

const (
	CarerGenderMale   CarerGender = "male"
	CarerGenderFemale CarerGender = "female"
)

const (
	RewardAmountHour RewardAmount = "hour"
)

type SosPost struct {
	ID           int          `field:"id"`
	AuthorID     int          `field:"author_id"`
	Title        string       `field:"title"`
	Content      string       `field:"content"`
	Reward       string       `field:"reward"`
	DateStartAt  string       `field:"date_start_at"`
	DateEndAt    string       `field:"date_end_at"`
	TimeStartAt  string       `field:"time_start_at"`
	TimeEndAt    string       `field:"time_end_at"`
	CareType     CareType     `field:"care_type"`
	CarerGender  CarerGender  `field:"carer_gender"`
	RewardAmount RewardAmount `field:"reward_amount"`
	ThumbnailID  int          `field:"thumbnail_id"`
	CreatedAt    time.Time    `field:"created_at"`
	UpdatedAt    time.Time    `field:"updated_at"`
	DeletedAt    time.Time    `field:"deleted_at"`
}

type SosPostStore interface {
	WriteSosPost(authorID int, request *WriteSosPostRequest) (*SosPost, error)
	FindSosPosts(page int, size int, sortBy string) ([]SosPost, error)
	FindSosPostsByAuthorID(authorID int, page int, size int) ([]SosPost, error)
	FindSosPostByID(id int) (*SosPost, error)
	UpdateSosPost(request *UpdateSosPostRequest) (*SosPost, error)
	FindConditionByID(id int) ([]Condition, error)
	FindPetsByID(id int) ([]pet.Pet, error)
}
