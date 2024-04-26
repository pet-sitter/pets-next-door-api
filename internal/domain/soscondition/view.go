package soscondition

import (
	utils "github.com/pet-sitter/pets-next-door-api/internal/common"
	databasegen "github.com/pet-sitter/pets-next-door-api/internal/infra/database/gen"
)

type DetailView struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func ToDetailView(row databasegen.SosCondition) *DetailView {
	return &DetailView{
		ID:   int64(row.ID),
		Name: utils.NullStrToStr(row.Name),
	}
}

func ToDetailViewFromRows(row databasegen.SosCondition) *DetailView {
	return &DetailView{
		ID:   int64(row.ID),
		Name: utils.NullStrToStr(row.Name),
	}
}

func ToDetailViewFromSOSPostCondition(row databasegen.FindSOSPostConditionsRow) *DetailView {
	return &DetailView{
		ID:   int64(row.ID),
		Name: utils.NullStrToStr(row.Name),
	}
}

func ToDetailViewFromViewForSOSPost(view ViewForSOSPost) *DetailView {
	return &DetailView{
		ID:   int64(view.ID),
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

func ToListViewFromRows(rows []databasegen.SosCondition) ListView {
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
