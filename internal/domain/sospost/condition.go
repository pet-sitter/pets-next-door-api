package sospost

import (
	"context"

	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
)

type Condition struct {
	ID        int    `field:"id"`
	Name      string `field:"name"`
	CreatedAt string `field:"created_at"`
	UpdatedAt string `field:"update_at"`
	DeletedAt string `field:"deleted_at"`
}

type ConditionList []*Condition

type SOSCondition string

const (
	CCTVPermission  SOSCondition = "CCTV, 펫캠 촬영 동의"
	IDVerification  SOSCondition = "신분증 인증"
	PhonePermission SOSCondition = "사전 통화 가능 여부"
)

var ConditionName = []SOSCondition{CCTVPermission, IDVerification, PhonePermission}

type ConditionStore interface {
	InitConditions(ctx context.Context, tx *database.Tx, conditions []SOSCondition) (string, *pnd.AppError)
	FindConditions(ctx context.Context, tx *database.Tx) (*ConditionList, *pnd.AppError)
}
