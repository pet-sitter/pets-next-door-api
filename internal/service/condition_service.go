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
	tx, err := service.conn.BeginSqlTx(ctx)
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

	return conditions.ToConditionViewList(), nil
}
