package sos_post

import (
	"time"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
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
	CarerGenderAll    CarerGender = "all"
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
	CareType     CareType     `field:"care_type"`
	CarerGender  CarerGender  `field:"carer_gender"`
	RewardAmount RewardAmount `field:"reward_amount"`
	ThumbnailID  int          `field:"thumbnail_id"`
	CreatedAt    time.Time    `field:"created_at"`
	UpdatedAt    time.Time    `field:"updated_at"`
	DeletedAt    time.Time    `field:"deleted_at"`
}

type SosPostList struct {
	*pnd.PaginatedView[SosPost]
}

func NewSosPostList(page int, size int) *SosPostList {
	return &SosPostList{PaginatedView: pnd.NewPaginatedView(
		page, size, false, make([]SosPost, 0),
	)}
}

type SosPostStore interface {
	WriteSosPost(authorID int, utcDateStart string, utcDateEnd string, request *WriteSosPostRequest) (*SosPost, *pnd.AppError)
	FindSosPosts(page int, size int, sortBy string) (*SosPostList, *pnd.AppError)
	FindSosPostsByAuthorID(authorID int, page int, size int, sortBy string) (*SosPostList, *pnd.AppError)
	FindSosPostByID(id int) (*SosPost, *pnd.AppError)
	UpdateSosPost(request *UpdateSosPostRequest) (*SosPost, *pnd.AppError)
	FindConditionByID(id int) ([]Condition, *pnd.AppError)
	FindPetsByID(id int) ([]pet.Pet, *pnd.AppError)
}
