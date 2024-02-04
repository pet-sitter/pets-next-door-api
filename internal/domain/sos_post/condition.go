package sos_post

import (
	"context"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
)

type Condition struct {
	ID        int    `filed:"id"`
	Name      string `filed:"name"`
	CreatedAt string `filed:"created_at"`
	UpdatedAt string `filed:"update_at"`
	DeletedAt string `filed:"deleted_at"`
}

type ConditionView struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type SosCondition string

const (
	CCTVPermission  SosCondition = "CCTV, 펫캠 촬영 동의"
	IDVerification  SosCondition = "신분증 인증"
	PhonePermission SosCondition = "사전 통화 가능 여부"
)

var ConditionName = []SosCondition{CCTVPermission, IDVerification, PhonePermission}

type ConditionStore interface {
	FindConditions(ctx context.Context) ([]Condition, *pnd.AppError)
}
