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
	ImageIDs     []int         `json:"image_ids" validate:"required"`
	Reward       string        `json:"reward" validate:"required"`
	Dates        []SosDateView `json:"dates" validate:"required"`
	CareType     CareType      `json:"care_type" validate:"required,oneof= foster visiting"`
	CarerGender  CarerGender   `json:"carer_gender" validate:"required,oneof=male female all"`
	RewardAmount RewardAmount  `json:"reward_amount" validate:"required,oneof=hour"`
	ConditionIDs []int         `json:"condition_ids"`
	PetIDs       []int         `json:"pet_ids"`
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
	CareType     CareType          `json:"care_type"`
	CarerGender  CarerGender       `json:"carer_gender"`
	RewardAmount RewardAmount      `json:"reward_amount"`
	ThumbnailID  int               `json:"thumbnail_id"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
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
	CareType     CareType                     `json:"care_type"`
	CarerGender  CarerGender                  `json:"carer_gender"`
	RewardAmount RewardAmount                 `json:"reward_amount"`
	ThumbnailID  int                          `json:"thumbnail_id"`
	CreatedAt    time.Time                    `json:"created_at"`
	UpdatedAt    time.Time                    `json:"updated_at"`
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
	ImageIDs     []int         `json:"image_ids" validate:"required"`
	Dates        []SosDateView `json:"dates" validate:"required"`
	Reward       string        `json:"reward" validate:"required"`
	CareType     CareType      `json:"care_type" validate:"required,oneof= foster visiting"`
	CarerGender  CarerGender   `json:"carer_gender" validate:"required,oneof=male female all"`
	RewardAmount RewardAmount  `json:"reward_amount" validate:"required,oneof=hour"`
	ConditionIDs []int         `json:"condition_ids"`
	PetIDs       []int         `json:"pet_ids"`
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
	CareType     CareType          `json:"care_type"`
	CarerGender  CarerGender       `json:"carer_gender"`
	RewardAmount RewardAmount      `json:"reward_amount"`
	ThumbnailID  int               `json:"thumbnail_id"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

type SosDateView struct {
	DateStartAt string `field:"dateStartAt"`
	DateEndAt   string `field:"dateEndAt"`
}
