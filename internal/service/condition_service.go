package service

import (
	"context"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/sos_post"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
	"github.com/pet-sitter/pets-next-door-api/internal/postgres"
)

type ConditionService struct {
	conn *database.DB
}

func NewConditionService(conn *database.DB) *ConditionService {
	return &ConditionService{
		conn: conn,
	}
}

func (service *ConditionService) FindConditions(ctx context.Context) ([]sos_post.ConditionView, *pnd.AppError) {
	conditionViews := make([]sos_post.ConditionView, 0)

	err := database.WithTransaction(ctx, service.conn, func(tx *database.Tx) *pnd.AppError {
		conditionStore := postgres.NewConditionPostgresStore(tx)

		conditions, err := conditionStore.FindConditions(ctx)
		if err != nil {
			return err
		}

		for _, v := range conditions {
			conditionView := sos_post.ConditionView{
				ID:   v.ID,
				Name: v.Name,
			}
			conditionViews = append(conditionViews, conditionView)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return conditionViews, nil
}
