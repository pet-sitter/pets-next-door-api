package sos_post

import (
	"context"
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

type SosDates struct {
	ID          int       `field:"id"`
	DateStartAt string    `field:"date_start_at"`
	DateEndAt   string    `field:"date_end_at"`
	CreatedAt   time.Time `field:"created_at"`
	UpdatedAt   time.Time `field:"updated_at"`
	DeletedAt   time.Time `field:"deleted_at"`
}

type SosPostStore interface {
	WriteSosPost(ctx context.Context, authorID int, utcDateStart string, utcDateEnd string, request *WriteSosPostRequest) (*SosPost, *pnd.AppError)
	FindSosPosts(ctx context.Context, page int, size int, sortBy string) (*SosPostList, *pnd.AppError)
	FindSosPostsByAuthorID(ctx context.Context, authorID int, page int, size int, sortBy string) (*SosPostList, *pnd.AppError)
	FindSosPostByID(ctx context.Context, id int) (*SosPost, *pnd.AppError)
	UpdateSosPost(ctx context.Context, request *UpdateSosPostRequest) (*SosPost, *pnd.AppError)
	FindConditionByID(ctx context.Context, id int) ([]Condition, *pnd.AppError)
	FindPetsByID(ctx context.Context, id int) ([]pet.Pet, *pnd.AppError)
	WriteDates(ctx context.Context, dates []string, sosPostID int) ([]SosDates, *pnd.AppError)
	FindDatesBySosPostID(ctx context.Context, sosPostID int) (SosDates, *pnd.AppError)
}
