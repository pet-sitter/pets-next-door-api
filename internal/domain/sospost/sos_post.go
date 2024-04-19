package sospost

import (
	"context"
	"time"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/media"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
)

type (
	CareType    string
	CarerGender string
	RewardType  string
)

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

type SOSPost struct {
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

type SOSPostList struct {
	*pnd.PaginatedView[SOSPost]
}

type SOSPostInfo struct {
	ID          int                    `field:"id" json:"id"`
	AuthorID    int                    `field:"author" json:"author"`
	Title       string                 `field:"title" json:"title"`
	Content     string                 `field:"content" json:"content"`
	Media       media.MediaList        `field:"media" json:"media"`
	Conditions  ConditionList          `field:"conditions" json:"conditions"`
	Pets        pet.PetWithProfileList `field:"pets" json:"pets"`
	Reward      string                 `field:"reward" json:"reward"`
	Dates       SOSDatesList           `field:"dates" json:"dates"`
	CareType    CareType               `field:"careType" json:"careType"`
	CarerGender CarerGender            `field:"carerGender" json:"carerGender"`
	RewardType  RewardType             `field:"rewardType" json:"rewardType"`
	ThumbnailID int                    `field:"thumbnailId" json:"thumbnailId"`
	CreatedAt   time.Time              `field:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time              `field:"updatedAt" json:"updatedAt"`
	DeletedAt   time.Time              `field:"deletedAt" json:"deletedAt"`
}

type SOSPostInfoList struct {
	*pnd.PaginatedView[SOSPostInfo]
}

func NewSOSPostList(page, size int) *SOSPostList {
	return &SOSPostList{PaginatedView: pnd.NewPaginatedView(
		page, size, false, make([]SOSPost, 0),
	)}
}

func NewSOSPostInfoList(page, size int) *SOSPostInfoList {
	return &SOSPostInfoList{PaginatedView: pnd.NewPaginatedView(
		page, size, false, make([]SOSPostInfo, 0),
	)}
}

type SOSDates struct {
	ID          int       `field:"id" json:"id"`
	DateStartAt string    `field:"date_start_at" json:"date_start_at"`
	DateEndAt   string    `field:"date_end_at" json:"date_end_at"`
	CreatedAt   time.Time `field:"created_at" json:"created_at"`
	UpdatedAt   time.Time `field:"updated_at" json:"updated_at"`
	DeletedAt   time.Time `field:"deleted_at" json:"deleted_at"`
}

type SOSDatesList []*SOSDates

type SOSPostStore interface {
	WriteSOSPost(
		ctx context.Context,
		tx *database.Tx,
		authorID int,
		utcDateStart string,
		utcDateEnd string,
		request *WriteSOSPostRequest,
	) (*SOSPost, *pnd.AppError)
	FindSOSPosts(
		ctx context.Context,
		tx *database.Tx,
		page int,
		size int,
		sortBy string,
	) (*SOSPostInfoList, *pnd.AppError)
	FindSOSPostsByAuthorID(
		ctx context.Context,
		tx *database.Tx,
		authorID int,
		page int,
		size int,
		sortBy string,
	) (*SOSPostInfoList, *pnd.AppError)
	FindSOSPostByID(ctx context.Context, tx *database.Tx, id int) (*SOSPost, *pnd.AppError)
	UpdateSOSPost(ctx context.Context, tx *database.Tx, request *UpdateSOSPostRequest) (*SOSPost, *pnd.AppError)
	FindConditionByID(ctx context.Context, tx *database.Tx, id int) (*ConditionList, *pnd.AppError)
	FindPetsByID(ctx context.Context, tx *database.Tx, id int) (*pet.PetList, *pnd.AppError)
	WriteDates(ctx context.Context, tx *database.Tx, dates []string, sosPostID int) (*SOSDatesList, *pnd.AppError)
	FindDatesBySOSPostID(ctx context.Context, tx *database.Tx, sosPostID int) (*SOSDatesList, *pnd.AppError)
}
