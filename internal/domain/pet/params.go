package pet

import (
	utils "github.com/pet-sitter/pets-next-door-api/internal/common"
	databasegen "github.com/pet-sitter/pets-next-door-api/internal/infra/database/gen"
)

type FindPetParams struct {
	ID             *int64
	OwnerID        *int64
	IncludeDeleted bool
}

func (p *FindPetParams) ToDBParams() databasegen.FindPetParams {
	return databasegen.FindPetParams{
		ID:             utils.Int64PtrToNullInt32(p.ID),
		OwnerID:        utils.Int64PtrToNullInt64(p.OwnerID),
		IncludeDeleted: p.IncludeDeleted,
	}
}

type FindPetsParams struct {
	Page           int
	Size           int
	ID             *int64
	OwnerID        *int64
	IncludeDeleted bool
}

func (p *FindPetsParams) ToDBParams() databasegen.FindPetsParams {
	pagination := utils.OffsetAndLimit(p.Page, p.Size)
	return databasegen.FindPetsParams{
		Limit:          int32(pagination.Limit),
		Offset:         int32(pagination.Offset),
		ID:             utils.Int64PtrToNullInt32(p.ID),
		OwnerID:        utils.Int64PtrToNullInt64(p.OwnerID),
		IncludeDeleted: p.IncludeDeleted,
	}
}
