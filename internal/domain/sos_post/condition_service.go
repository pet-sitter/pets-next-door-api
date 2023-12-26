package sos_post

type ConditionService struct {
	conditionStore ConditionStore
}

func NewConditionService(conditionStore ConditionStore) *ConditionService {
	return &ConditionService{
		conditionStore: conditionStore,
	}
}

func (service *ConditionService) FindConditions() ([]ConditionView, error) {
	conditions, err := service.conditionStore.FindConditions()

	if err != nil {
		return nil, err
	}

	var conditionViews []ConditionView

	for _, v := range conditions {
		conditionView := ConditionView{
			ID:   v.ID,
			Name: v.Name,
		}
		conditionViews = append(conditionViews, conditionView)
	}

	return conditionViews, nil
}
