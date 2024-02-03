package sos_post

import (
	"time"

	"github.com/pet-sitter/pets-next-door-api/internal/domain/media"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"
)

type WriteSosPostRequest struct {
	Title        string       `json:"title" validate:"required"`
	Content      string       `json:"content" validate:"required"`
	ImageIDs     []int        `json:"image_ids" validate:"required"`
	Reward       string       `json:"reward" validate:"required"`
	DateStartAt  time.Time    `json:"date_start_at" validate:"required"`
	DateEndAt    time.Time    `json:"date_end_at" validate:"required"`
	TimeStartAt  string       `json:"time_start_at" validate:"required"`
	TimeEndAt    string       `json:"time_end_at" validate:"required"`
	CareType     CareType     `json:"care_type" validate:"required,oneof= foster visiting"`
	CarerGender  CarerGender  `json:"carer_gender" validate:"required,oneof=male female"`
	RewardAmount RewardAmount `json:"reward_amount" validate:"required,oneof=hour"`
	ConditionIDs []int        `json:"condition_ids"`
	PetIDs       []int        `json:"pet_ids"`
}

type WriteSosPostView struct {
	ID           int               `json:"id"`
	AuthorID     int               `json:"author_id"`
	Title        string            `json:"title"`
	Content      string            `json:"content"`
	Media        []media.MediaView `json:"media"`
	Conditions   []ConditionView   `json:"conditions"`
	Pets         []pet.PetView     `json:"pets"`
	Reward       string            `json:"reward"`
	DateStartAt  string            `json:"date_start_at"`
	DateEndAt    string            `json:"date_end_at"`
	TimeStartAt  string            `json:"time_start_at"`
	TimeEndAt    string            `json:"time_end_at"`
	CareType     CareType          `json:"care_type"`
	CarerGender  CarerGender       `json:"carer_gender"`
	RewardAmount RewardAmount      `json:"reward_amount"`
	ThumbnailID  int               `json:"thumbnail_id"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

type FindSosPostView struct {
	ID           int               `json:"id"`
	AuthorID     int               `json:"author_id"`
	Title        string            `json:"title"`
	Content      string            `json:"content"`
	Media        []media.MediaView `json:"media"`
	Conditions   []ConditionView   `json:"conditions"`
	Pets         []pet.PetView     `json:"pets"`
	Reward       string            `json:"reward"`
	DateStartAt  string            `json:"date_start_at"`
	DateEndAt    string            `json:"date_end_at"`
	TimeStartAt  string            `json:"time_start_at"`
	TimeEndAt    string            `json:"time_end_at"`
	CareType     CareType          `json:"care_type"`
	CarerGender  CarerGender       `json:"carer_gender"`
	RewardAmount RewardAmount      `json:"reward_amount"`
	ThumbnailID  int               `json:"thumbnail_id"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

type UpdateSosPostRequest struct {
	ID           int          `json:"id" validate:"required"`
	Title        string       `json:"title" validate:"required"`
	Content      string       `json:"content" validate:"required"`
	ImageIDs     []int        `json:"image_ids" validate:"required"`
	Reward       string       `json:"reward" validate:"required"`
	DateStartAt  string       `json:"date_start_at" validate:"required"`
	DateEndAt    string       `json:"date_end_at" validate:"required"`
	TimeStartAt  string       `json:"time_start_at" validate:"required"`
	TimeEndAt    string       `json:"time_end_at" validate:"required"`
	CareType     CareType     `json:"care_type" validate:"required,oneof= foster visiting"`
	CarerGender  CarerGender  `json:"carer_gender" validate:"required,oneof=male female"`
	RewardAmount RewardAmount `json:"reward_amount" validate:"required,oneof=hour"`
	ConditionIDs []int        `json:"condition_ids"`
	PetIDs       []int        `json:"pet_ids"`
}

type UpdateSosPostView struct {
	ID           int               `json:"id"`
	AuthorID     int               `json:"author_id"`
	Title        string            `json:"title"`
	Content      string            `json:"content"`
	Media        []media.MediaView `json:"media"`
	Conditions   []ConditionView   `json:"conditions"`
	Pets         []pet.PetView     `json:"pets"`
	Reward       string            `json:"reward"`
	DateStartAt  string            `json:"date_start_at"`
	DateEndAt    string            `json:"date_end_at"`
	TimeStartAt  string            `json:"time_start_at"`
	TimeEndAt    string            `json:"time_end_at"`
	CareType     CareType          `json:"care_type"`
	CarerGender  CarerGender       `json:"carer_gender"`
	RewardAmount RewardAmount      `json:"reward_amount"`
	ThumbnailID  int               `json:"thumbnail_id"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}
