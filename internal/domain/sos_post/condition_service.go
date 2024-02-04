package sos_post

import (
	"context"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
)

type ConditionService struct {
	conditionStore ConditionStore
}

func NewConditionService(conditionStore ConditionStore) *ConditionService {
	return &ConditionService{
		conditionStore: conditionStore,
	}
}

func (service *ConditionService) FindConditions(ctx context.Context) ([]ConditionView, *pnd.AppError) {
	conditions, err := service.conditionStore.FindConditions(ctx)
	if err != nil {
		return nil, err
	}

	conditionViews := make([]ConditionView, 0)
	for _, v := range conditions {
		conditionView := ConditionView{
			ID:   v.ID,
			Name: v.Name,
		}
		conditionViews = append(conditionViews, conditionView)
	}

	return conditionViews, nil
}
