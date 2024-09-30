package pet

import (
	"github.com/google/uuid"
	utils "github.com/pet-sitter/pets-next-door-api/internal/common"
	databasegen "github.com/pet-sitter/pets-next-door-api/internal/infra/database/gen"
)

type FindPetParams struct {
	ID             uuid.NullUUID
	OwnerID        uuid.NullUUID
	IncludeDeleted bool
}

func (p *FindPetParams) ToDBParams() databasegen.FindPetParams {
	return databasegen.FindPetParams{
		ID:             p.ID,
		OwnerID:        p.OwnerID,
		IncludeDeleted: p.IncludeDeleted,
	}
}

type FindPetsParams struct {
	Page           int
	Size           int
	ID             uuid.NullUUID
	OwnerID        uuid.NullUUID
	IncludeDeleted bool
}

func (p *FindPetsParams) ToDBParams() databasegen.FindPetsParams {
	pagination := utils.OffsetAndLimit(p.Page, p.Size)
	return databasegen.FindPetsParams{
		Limit:          int32(pagination.Limit),
		Offset:         int32(pagination.Offset),
		ID:             p.ID,
		OwnerID:        p.OwnerID,
		IncludeDeleted: p.IncludeDeleted,
	}
}
