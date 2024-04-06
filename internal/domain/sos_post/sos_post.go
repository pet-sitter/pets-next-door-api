package sos_post

import (
	"context"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
	"time"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"
)

type CareType string
type CarerGender string
type RewardType string

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
	RewardTypeFee        RewardType = "fee"
	RewardTypeGifticon   RewardType = "gifticon"
	RewardTypeNegotiable RewardType = "negotiable"
)

type SosPost struct {
	ID          int         `field:"id"`
	AuthorID    int         `field:"author_id"`
	Title       string      `field:"title"`
	Content     string      `field:"content"`
	Reward      string      `field:"reward"`
	CareType    CareType    `field:"care_type"`
	CarerGender CarerGender `field:"carer_gender"`
	RewardType  RewardType  `field:"reward_type"`
	ThumbnailID int         `field:"thumbnail_id"`
	CreatedAt   time.Time   `field:"created_at"`
	UpdatedAt   time.Time   `field:"updated_at"`
	DeletedAt   time.Time   `field:"deleted_at"`
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

type SosDatesList []*SosDates

type SosPostStore interface {
	WriteSosPost(ctx context.Context, tx database.Tx, authorID int, utcDateStart string, utcDateEnd string, request *WriteSosPostRequest) (*SosPost, *pnd.AppError)
	FindSosPosts(ctx context.Context, tx database.Tx, page int, size int, sortBy string) (*SosPostList, *pnd.AppError)
	FindSosPostsByAuthorID(ctx context.Context, tx database.Tx, authorID int, page int, size int, sortBy string) (*SosPostList, *pnd.AppError)
	FindSosPostByID(ctx context.Context, tx database.Tx, id int) (*SosPost, *pnd.AppError)
	UpdateSosPost(ctx context.Context, tx database.Tx, request *UpdateSosPostRequest) (*SosPost, *pnd.AppError)
	FindConditionByID(ctx context.Context, tx database.Tx, id int) (*ConditionList, *pnd.AppError)
	FindPetsByID(ctx context.Context, tx database.Tx, id int) (*pet.PetList, *pnd.AppError)
	WriteDates(ctx context.Context, tx database.Tx, dates []string, sosPostID int) (*SosDatesList, *pnd.AppError)
	FindDatesBySosPostID(ctx context.Context, tx database.Tx, sosPostID int) (*SosDatesList, *pnd.AppError)
}
