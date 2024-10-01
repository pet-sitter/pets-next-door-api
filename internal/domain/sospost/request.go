package sospost

import "github.com/google/uuid"

type WriteSOSPostRequest struct {
	Title        string        `json:"title" validate:"required"`
	Content      string        `json:"content" validate:"required"`
	ImageIDs     []uuid.UUID   `json:"imageIds" validate:"required"`
	Reward       string        `json:"reward" validate:"required"`
	Dates        []SOSDateView `json:"dates" validate:"required,gte=1"`
	CareType     CareType      `json:"careType" validate:"required,oneof=foster visiting"`
	CarerGender  CarerGender   `json:"carerGender" validate:"required,oneof=male female all"`
	RewardType   RewardType    `json:"rewardType" validate:"required,oneof=fee gifticon negotiable"`
	ConditionIDs []uuid.UUID   `json:"conditionIds" validate:"required"`
	PetIDs       []uuid.UUID   `json:"petIds" validate:"required,gte=1"`
}

type UpdateSOSPostRequest struct {
	ID           uuid.UUID     `json:"id" validate:"required"`
	Title        string        `json:"title" validate:"required"`
	Content      string        `json:"content" validate:"required"`
	ImageIDs     []uuid.UUID   `json:"imageIds" validate:"required"`
	Dates        []SOSDateView `json:"dates" validate:"required,gte=1"`
	Reward       string        `json:"reward" validate:"required"`
	CareType     CareType      `json:"careType" validate:"required,oneof=foster visiting"`
	CarerGender  CarerGender   `json:"carerGender" validate:"required,oneof=male female all"`
	RewardType   RewardType    `json:"rewardType" validate:"required,oneof=fee gifticon negotiable"`
	ConditionIDs []uuid.UUID   `json:"conditionIds" validate:"required"`
	PetIDs       []uuid.UUID   `json:"petIds" validate:"required,gte=1"`
}
