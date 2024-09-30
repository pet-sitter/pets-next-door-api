package soscondition

import (
	"github.com/google/uuid"
	utils "github.com/pet-sitter/pets-next-door-api/internal/common"
	databasegen "github.com/pet-sitter/pets-next-door-api/internal/infra/database/gen"
)

type DetailView struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func ToDetailView(row databasegen.SosCondition) *DetailView {
	return &DetailView{
		ID:   row.ID,
		Name: utils.NullStrToStr(row.Name),
	}
}

func ToDetailViewFromRows(row databasegen.FindConditionsRow) *DetailView {
	return &DetailView{
		ID:   row.ID,
		Name: utils.NullStrToStr(row.Name),
	}
}

func ToDetailViewFromSOSPostCondition(row databasegen.FindSOSPostConditionsRow) *DetailView {
	return &DetailView{
		ID:   row.ID,
		Name: utils.NullStrToStr(row.Name),
	}
}

func ToDetailViewFromViewForSOSPost(view ViewForSOSPost) *DetailView {
	return &DetailView{
		ID:   view.ID,
		Name: view.Name,
	}
}

type ListView []*DetailView

func ToListView(row []databasegen.SosCondition) ListView {
	conditionViews := make(ListView, len(row))
	for i, condition := range row {
		conditionViews[i] = ToDetailView(condition)
	}
	return conditionViews
}

func ToListViewFromRows(rows []databasegen.FindConditionsRow) ListView {
	conditionViews := make(ListView, len(rows))
	for i, row := range rows {
		conditionViews[i] = ToDetailViewFromRows(row)
	}
	return conditionViews
}

func ToListViewFromSOSPostConditions(rows []databasegen.FindSOSPostConditionsRow) ListView {
	conditionViews := make(ListView, len(rows))
	for i, row := range rows {
		conditionViews[i] = ToDetailViewFromSOSPostCondition(row)
	}
	return conditionViews
}

func ToListViewFromViewForSOSPost(views ViewListForSOSPost) ListView {
	conditionViews := make(ListView, len(views))
	for i, view := range views {
		conditionViews[i] = ToDetailViewFromViewForSOSPost(*view)
	}
	return conditionViews
}
