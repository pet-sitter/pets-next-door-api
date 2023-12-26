package sos_post

type Condition struct {
	ID        int    `filed:"id"`
	Name      string `filed:"name"`
	CreatedAt string `filed:"created_at"`
	UpdatedAt string `filed:"update_at"`
	DeletedAt string `filed:"deleted_at"`
}

type ConditionView struct {
	ID   int    `filed:"id"`
	Name string `filed:"name"`
}

type SosCondition string

const (
	CCTVPermission  SosCondition = "CCTV, 펫캠 촬영 동의"
	PhonePermission SosCondition = "사전 통화 가능 여부"
	PetRegistration SosCondition = "반려 동물 등록 여부"
)

var ConditionName = []SosCondition{CCTVPermission, PhonePermission, PetRegistration}

type ConditionStore interface {
	FindConditions() ([]Condition, error)
}
