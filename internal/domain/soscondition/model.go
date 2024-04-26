package soscondition

type ViewForSOSPost struct {
	ID        int    `field:"id"`
	Name      string `field:"name"`
	CreatedAt string `field:"created_at"`
	UpdatedAt string `field:"update_at"`
	DeletedAt string `field:"deleted_at"`
}

type ViewListForSOSPost []*ViewForSOSPost

type AvailableName string

const (
	CCTVPermission  AvailableName = "CCTV, 펫캠 촬영 동의"
	IDVerification  AvailableName = "신분증 인증"
	PhonePermission AvailableName = "사전 통화 가능 여부"
)

func (c AvailableName) String() string {
	return string(c)
}

var AvailableNames = []AvailableName{CCTVPermission, IDVerification, PhonePermission}
