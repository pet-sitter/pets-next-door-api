package breed

import (
	utils "github.com/pet-sitter/pets-next-door-api/internal/common"
	databasegen "github.com/pet-sitter/pets-next-door-api/internal/infra/database/gen"
)

type FindBreedsParams struct {
	Page           int
	Size           int
	PetType        *string
	IncludeDeleted bool
}

func (p *FindBreedsParams) ToDBParams() databasegen.FindBreedsParams {
	pagination := utils.OffsetAndLimit(p.Page, p.Size)
	return databasegen.FindBreedsParams{
		Limit:          int32(pagination.Limit),
		Offset:         int32(pagination.Offset),
		PetType:        utils.StrPtrToNullStr(p.PetType),
		IncludeDeleted: p.IncludeDeleted,
	}
}
