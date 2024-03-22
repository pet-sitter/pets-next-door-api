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

func (u *User) ToRegisterUserView(profileImageURL *string) *RegisterUserView {
	return &RegisterUserView{
		ID:                   u.ID,
		Email:                u.Email,
		Nickname:             u.Nickname,
		Fullname:             u.Fullname,
		ProfileImageURL:      profileImageURL,
		FirebaseProviderType: u.FirebaseProviderType,
		FirebaseUID:          u.FirebaseUID,
	}
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

func (u *UserWithProfileImage) ToFindUserView() *FindUserView {
	return &FindUserView{
		ID:                   u.ID,
		Email:                u.Email,
		Nickname:             u.Nickname,
		Fullname:             u.Fullname,
		ProfileImageURL:      u.ProfileImageURL,
		FirebaseProviderType: u.FirebaseProviderType,
		FirebaseUID:          u.FirebaseUID,
	}
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

func (s *UserStatus) ToUserStatusView() *UserStatusView {
	return &UserStatusView{
		Status:               UserStatusRegistered,
		FirebaseProviderType: s.FirebaseProviderType,
	}
}

type UpdateUserRequest struct {
	Nickname       string `json:"nickname" validate:"required"`
	ProfileImageID *int   `json:"profileImageId" validate:"omitempty"`
}

type UpdateUserView struct {
	ID                   int                  `json:"id"`
	Email                string               `json:"email"`
	Nickname             string               `json:"nickname"`
	Fullname             string               `json:"fullname"`
	ProfileImageURL      *string              `json:"profileImageUrl"`
	FirebaseProviderType FirebaseProviderType `json:"fbProviderType"`
}

func (u *UserWithProfileImage) ToUpdateUserView() *UpdateUserView {
	return &UpdateUserView{
		ID:                   u.ID,
		Email:                u.Email,
		Nickname:             u.Nickname,
		Fullname:             u.Fullname,
		ProfileImageURL:      u.ProfileImageURL,
		FirebaseProviderType: u.FirebaseProviderType,
	}
}
