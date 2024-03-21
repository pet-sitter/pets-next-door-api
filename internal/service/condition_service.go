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
	tx, err := service.conn.BeginTx(ctx)
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}

	conditions, err := postgres.FindConditions(ctx, tx)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	conditionViews := make([]sos_post.ConditionView, 0)
	for _, v := range conditions {
		conditionView := sos_post.ConditionView{
			ID:   v.ID,
			Name: v.Name,
		}
		conditionViews = append(conditionViews, conditionView)
	}

	return conditionViews, nil
}
