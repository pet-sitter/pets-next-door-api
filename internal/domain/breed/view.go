package breed

import (
	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/commonvo"
	databasegen "github.com/pet-sitter/pets-next-door-api/internal/infra/database/gen"
)

type DetailView struct {
	ID      int64            `json:"id"`
	PetType commonvo.PetType `json:"petType"`
	Name    string           `json:"name"`
}

func ToDetailViewFromRows(row databasegen.FindBreedsRow) *DetailView {
	return &DetailView{
		ID:      int64(row.ID),
		PetType: commonvo.PetType(row.PetType),
		Name:    row.Name,
	}
}

type ListView struct {
	*pnd.PaginatedView[*DetailView]
}

func ToListViewFromRows(page, size int, rows []databasegen.FindBreedsRow) *ListView {
	bl := &ListView{PaginatedView: pnd.NewPaginatedView(page, size, false, make([]*DetailView, len(rows)))}
	for i, row := range rows {
		bl.Items[i] = ToDetailViewFromRows(row)
	}

	bl.CalcLastPage()
	return bl
}
