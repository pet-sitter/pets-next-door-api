package user

type RegisterUserRequest struct {
	Email                string               `json:"email" validate:"required,email"`
	Nickname             string               `json:"nickname" validate:"required"`
	Fullname             string               `json:"fullname" validate:"required"`
	ProfileImageID       *int                 `json:"profileImageId"`
	FirebaseProviderType FirebaseProviderType `json:"fbProviderType" validate:"required"`
	FirebaseUID          string               `json:"fbUid" validate:"required"`
}

type RegisterUserView struct {
	ID                   int                  `json:"id"`
	Email                string               `json:"email"`
	Nickname             string               `json:"nickname"`
	Fullname             string               `json:"fullname"`
	ProfileImageURL      *string              `json:"profileImageUrl"`
	FirebaseProviderType FirebaseProviderType `json:"fbProviderType"`
	FirebaseUID          string               `json:"fbUid"`
}

type FindUserView struct {
	ID                   int                  `json:"id"`
	Email                string               `json:"email"`
	Nickname             string               `json:"nickname"`
	Fullname             string               `json:"fullname"`
	ProfileImageURL      *string              `json:"profileImageUrl"`
	FirebaseProviderType FirebaseProviderType `json:"fbProviderType"`
	FirebaseUID          string               `json:"fbUid"`
}

func (r *FindUserView) ToMyProfileView() *MyProfileView {
	return &MyProfileView{
		ID:                   r.ID,
		Email:                r.Email,
		Nickname:             r.Nickname,
		Fullname:             r.Fullname,
		ProfileImageURL:      r.ProfileImageURL,
		FirebaseProviderType: r.FirebaseProviderType,
	}
}

type MyProfileView struct {
	ID                   int                  `json:"id"`
	Email                string               `json:"email"`
	Nickname             string               `json:"nickname"`
	Fullname             string               `json:"fullname"`
	ProfileImageURL      *string              `json:"profileImageUrl"`
	FirebaseProviderType FirebaseProviderType `json:"fbProviderType"`
}

type CheckNicknameRequest struct {
	Nickname string `json:"nickname" validate:"required"`
}

type CheckNicknameView struct {
	IsAvailable bool `json:"isAvailable"`
}

type UserStatusRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type UserRegistrationStatus string

const (
	UserStatusNotRegistered UserRegistrationStatus = "NOT_REGISTERED"
	UserStatusRegistered    UserRegistrationStatus = "REGISTERED"
)

type UserStatusView struct {
	Status               UserRegistrationStatus `json:"status"`
	FirebaseProviderType FirebaseProviderType   `json:"fbProviderType,omitempty"`
}

type UpdateUserRequest struct {
	Nickname       string `json:"nickname"`
	ProfileImageID *int   `json:"profileImageId"`
}

type UpdateUserView struct {
	ID                   int                  `json:"id"`
	Email                string               `json:"email"`
	Nickname             string               `json:"nickname"`
	Fullname             string               `json:"fullname"`
	ProfileImageURL      *string              `json:"profileImageUrl"`
	FirebaseProviderType FirebaseProviderType `json:"fbProviderType"`
}
