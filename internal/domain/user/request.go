package user

import (
	utils "github.com/pet-sitter/pets-next-door-api/internal/common"
	databasegen "github.com/pet-sitter/pets-next-door-api/internal/infra/database/gen"
)

type RegisterUserRequest struct {
	Email                string               `json:"email" validate:"required,email"`
	Nickname             string               `json:"nickname" validate:"required"`
	Fullname             string               `json:"fullname" validate:"required"`
	ProfileImageID       *int                 `json:"profileImageId"`
	FirebaseProviderType FirebaseProviderType `json:"fbProviderType" validate:"required"`
	FirebaseUID          string               `json:"fbUid" validate:"required"`
}

func (r *RegisterUserRequest) ToDBParams() databasegen.CreateUserParams {
	return databasegen.CreateUserParams{
		Email:          r.Email,
		Nickname:       r.Nickname,
		Fullname:       r.Fullname,
		Password:       "",
		ProfileImageID: utils.IntPtrToNullInt64(r.ProfileImageID),
		FbProviderType: r.FirebaseProviderType.NullString(),
		FbUid:          utils.StrToNullStr(r.FirebaseUID),
	}
}

type CheckNicknameRequest struct {
	Nickname string `json:"nickname" validate:"required"`
}

type UserStatusRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type UpdateUserRequest struct {
	Nickname       string `json:"nickname" validate:"required"`
	ProfileImageID *int   `json:"profileImageId" validate:"omitempty"`
}
