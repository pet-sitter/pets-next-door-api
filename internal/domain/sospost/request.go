package sospost

import "github.com/pet-sitter/pets-next-door-api/internal/domain/commonvo"

type WriteSOSPostRequest struct {
	Title        string               `json:"title" validate:"required"`
	Content      string               `json:"content" validate:"required"`
	ImageIDs     []int64              `json:"imageIds" validate:"required"`
	Reward       string               `json:"reward" validate:"required"`
	Dates        []SOSDateView        `json:"dates" validate:"required,gte=1"`
	CareType     commonvo.CareType    `json:"careType" validate:"required,oneof=foster visiting"`
	CarerGender  commonvo.CarerGender `json:"carerGender" validate:"required,oneof=male female all"`
	RewardType   commonvo.RewardType  `json:"rewardType" validate:"required,oneof=fee gifticon negotiable"`
	ConditionIDs []int                `json:"conditionIds" validate:"required"`
	PetIDs       []int64              `json:"petIds" validate:"required,gte=1"`
}

type UpdateSOSPostRequest struct {
	ID           int                  `json:"id" validate:"required"`
	Title        string               `json:"title" validate:"required"`
	Content      string               `json:"content" validate:"required"`
	ImageIDs     []int64              `json:"imageIds" validate:"required"`
	Dates        []SOSDateView        `json:"dates" validate:"required,gte=1"`
	Reward       string               `json:"reward" validate:"required"`
	CareType     commonvo.CareType    `json:"careType" validate:"required,oneof=foster visiting"`
	CarerGender  commonvo.CarerGender `json:"carerGender" validate:"required,oneof=male female all"`
	RewardType   commonvo.RewardType  `json:"rewardType" validate:"required,oneof=fee gifticon negotiable"`
	ConditionIDs []int                `json:"conditionIds" validate:"required"`
	PetIDs       []int64              `json:"petIds" validate:"required,gte=1"`
}
