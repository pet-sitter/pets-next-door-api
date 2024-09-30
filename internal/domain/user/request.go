package user

import (
	"github.com/google/uuid"
	utils "github.com/pet-sitter/pets-next-door-api/internal/common"
	"github.com/pet-sitter/pets-next-door-api/internal/datatype"
	databasegen "github.com/pet-sitter/pets-next-door-api/internal/infra/database/gen"
)

type RegisterUserRequest struct {
	Email                string               `json:"email" validate:"required,email"`
	Nickname             string               `json:"nickname" validate:"required"`
	Fullname             string               `json:"fullname" validate:"required"`
	ProfileImageID       uuid.NullUUID        `json:"profileImageId"`
	FirebaseProviderType FirebaseProviderType `json:"fbProviderType" validate:"required"`
	FirebaseUID          string               `json:"fbUid" validate:"required"`
}

func (r *RegisterUserRequest) ToDBParams() databasegen.CreateUserParams {
	return databasegen.CreateUserParams{
		ID:             datatype.NewUUIDV7(),
		Email:          r.Email,
		Nickname:       r.Nickname,
		Fullname:       r.Fullname,
		Password:       "",
		ProfileImageID: r.ProfileImageID,
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
	Nickname       string        `json:"nickname" validate:"required"`
	ProfileImageID uuid.NullUUID `json:"profileImageId" validate:"omitempty"`
}
