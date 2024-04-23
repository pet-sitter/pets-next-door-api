package user

import (
	utils "github.com/pet-sitter/pets-next-door-api/internal/common"
	databasegen "github.com/pet-sitter/pets-next-door-api/internal/infra/database/gen"
)

type FindUserParams struct {
	ID             *int
	Email          *string
	FbUID          *string
	IncludeDeleted bool
}

func (p *FindUserParams) ToDBParams() databasegen.FindUserParams {
	return databasegen.FindUserParams{
		ID:             utils.IntPtrToNullInt32(p.ID),
		Email:          utils.StrPtrToNullStr(p.Email),
		FbUid:          utils.StrPtrToNullStr(p.FbUID),
		IncludeDeleted: p.IncludeDeleted,
	}
}

type FindUsersParams struct {
	Page           int
	Size           int
	Nickname       *string
	IncludeDeleted bool
}

func (p *FindUsersParams) ToDBParams() databasegen.FindUsersParams {
	pagination := utils.OffsetAndLimit(p.Page, p.Size)
	return databasegen.FindUsersParams{
		Limit:          int32(pagination.Limit),
		Offset:         int32(pagination.Offset),
		Nickname:       utils.StrPtrToNullStr(p.Nickname),
		IncludeDeleted: p.IncludeDeleted,
	}
}
