package sos_post

import (
	"github.com/pet-sitter/pets-next-door-api/internal/domain/user"
	"time"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/media"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"
)

type WriteSosPostRequest struct {
	Title        string        `json:"title" validate:"required"`
	Content      string        `json:"content" validate:"required"`
	ImageIDs     []int         `json:"imageIds" validate:"required"`
	Reward       string        `json:"reward" validate:"required"`
	Dates        []SosDateView `json:"dates" validate:"required"`
	CareType     CareType      `json:"careType" validate:"required,oneof= foster visiting"`
	CarerGender  CarerGender   `json:"carerGender" validate:"required,oneof=male female all"`
	RewardAmount RewardAmount  `json:"rewardAmount" validate:"required,oneof=hour"`
	ConditionIDs []int         `json:"conditionIds"`
	PetIDs       []int         `json:"petIds"`
}

type WriteSosPostView struct {
	ID           int               `json:"id"`
	AuthorID     int               `json:"authorId"`
	Title        string            `json:"title"`
	Content      string            `json:"content"`
	Media        []media.MediaView `json:"media"`
	Conditions   []ConditionView   `json:"conditions"`
	Pets         []pet.PetView     `json:"pets"`
	Reward       string            `json:"reward"`
	Dates        []SosDateView     `json:"dates"`
	CareType     CareType          `json:"careType"`
	CarerGender  CarerGender       `json:"carerGender"`
	RewardAmount RewardAmount      `json:"rewardAmount"`
	ThumbnailID  int               `json:"thumbnailId"`
	CreatedAt    time.Time         `json:"createdAt"`
	UpdatedAt    time.Time         `json:"updatedAt"`
}

type FindSosPostView struct {
	ID           int                          `json:"id"`
	Author       *user.UserWithoutPrivateInfo `json:"author"`
	Title        string                       `json:"title"`
	Content      string                       `json:"content"`
	Media        []media.MediaView            `json:"media"`
	Conditions   []ConditionView              `json:"conditions"`
	Pets         []pet.PetView                `json:"pets"`
	Reward       string                       `json:"reward"`
	Dates        []SosDateView                `json:"dates"`
	CareType     CareType                     `json:"careType"`
	CarerGender  CarerGender                  `json:"carerGender"`
	RewardAmount RewardAmount                 `json:"rewardAmount"`
	ThumbnailID  int                          `json:"thumbnailId"`
	CreatedAt    time.Time                    `json:"createdAt"`
	UpdatedAt    time.Time                    `json:"updatedAt"`
}

type FindSosPostListView struct {
	*pnd.PaginatedView[FindSosPostView]
}

func FromEmptySosPostList(sosPosts *SosPostList) *FindSosPostListView {
	return &FindSosPostListView{
		PaginatedView: pnd.NewPaginatedView(
			sosPosts.Page, sosPosts.Size, sosPosts.IsLastPage, make([]FindSosPostView, 0),
		),
	}
}

type UpdateSosPostRequest struct {
	ID           int           `json:"id" validate:"required"`
	Title        string        `json:"title" validate:"required"`
	Content      string        `json:"content" validate:"required"`
	ImageIDs     []int         `json:"imageIds" validate:"required"`
	Dates        []SosDateView `json:"dates" validate:"required"`
	Reward       string        `json:"reward" validate:"required"`
	CareType     CareType      `json:"careType" validate:"required,oneof= foster visiting"`
	CarerGender  CarerGender   `json:"carerGender" validate:"required,oneof=male female all"`
	RewardAmount RewardAmount  `json:"rewardAmount" validate:"required,oneof=hour"`
	ConditionIDs []int         `json:"conditionIds"`
	PetIDs       []int         `json:"petIds"`
}

type UpdateSosPostView struct {
	ID           int               `json:"id"`
	AuthorID     int               `json:"authorId"`
	Title        string            `json:"title"`
	Content      string            `json:"content"`
	Media        []media.MediaView `json:"media"`
	Conditions   []ConditionView   `json:"conditions"`
	Pets         []pet.PetView     `json:"pets"`
	Reward       string            `json:"reward"`
	Dates        []SosDateView     `json:"dates"`
	CareType     CareType          `json:"careType"`
	CarerGender  CarerGender       `json:"carerGender"`
	RewardAmount RewardAmount      `json:"rewardAmount"`
	ThumbnailID  int               `json:"thumbnailId"`
	CreatedAt    time.Time         `json:"createdAt"`
	UpdatedAt    time.Time         `json:"updatedAt"`
}

type SosDateView struct {
	DateStartAt string `field:"dateStartAt"`
	DateEndAt   string `field:"dateEndAt"`
}
