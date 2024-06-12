package commonvo

// Pet
type PetType string

const (
	PetTypeDog PetType = "dog"
	PetTypeCat PetType = "cat"
)

// SOSPost
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

func (p *PetType) String() string {
	return string(*p)
}

func (c *CareType) String() string {
	return string(*c)
}

func (c *CarerGender) String() string {
	return string(*c)
}

func (r *RewardType) String() string {
	return string(*r)
}
