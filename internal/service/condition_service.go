package service

import (
	"context"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database/pgx"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/sos_post"
	"github.com/pet-sitter/pets-next-door-api/internal/postgres"
)

type ConditionService struct {
	conn *pgx.DB
}

func NewConditionService(conn *pgx.DB) *ConditionService {
	return &ConditionService{
		conn: conn,
	}
}

func (service *ConditionService) FindConditions(ctx context.Context) ([]sos_post.ConditionView, *pnd.AppError) {
	tx, err := service.conn.BeginPgxTx(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		return nil, err
	}

	conditions, err := postgres.FindConditions(ctx, tx)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return conditions.ToConditionViewList(), nil
}
