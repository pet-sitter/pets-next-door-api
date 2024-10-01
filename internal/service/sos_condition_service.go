package service

import (
	"context"

	"github.com/pet-sitter/pets-next-door-api/internal/datatype"

	utils "github.com/pet-sitter/pets-next-door-api/internal/common"

	"github.com/pet-sitter/pets-next-door-api/internal/domain/soscondition"
	databasegen "github.com/pet-sitter/pets-next-door-api/internal/infra/database/gen"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
)

type SOSConditionService struct {
	conn *database.DB
}

func NewSOSConditionService(conn *database.DB) *SOSConditionService {
	return &SOSConditionService{
		conn: conn,
	}
}

func (service *SOSConditionService) InitConditions(ctx context.Context) (soscondition.ListView, *pnd.AppError) {
	tx, err := service.conn.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	conditionList := make([]databasegen.SosCondition, len(soscondition.AvailableNames))
	for idx, conditionName := range soscondition.AvailableNames {
		created, err := databasegen.New(tx).CreateSOSCondition(ctx, databasegen.CreateSOSConditionParams{
			ID:   datatype.NewUUIDV7(),
			Name: utils.StrToNullStr(conditionName.String()),
		})
		if err != nil {
			return nil, pnd.FromPostgresError(err)
		}

		conditionList[idx] = created
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return soscondition.ToListView(conditionList), nil
}

func (service *SOSConditionService) FindConditions(ctx context.Context) (soscondition.ListView, *pnd.AppError) {
	conditionList, err := databasegen.New(service.conn).FindConditions(ctx, false)
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return soscondition.ToListViewFromRows(conditionList), nil
}
