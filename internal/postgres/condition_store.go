package postgres

import (
	"context"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/sos_post"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
)

type ConditionPostgresStore struct {
	db *database.DB
}

func NewConditionPostgresStore(db *database.DB) *ConditionPostgresStore {
	return &ConditionPostgresStore{
		db: db,
	}
}

func (s *ConditionPostgresStore) InitConditions(ctx context.Context, conditions []sos_post.SosCondition) (string, *pnd.AppError) {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return "", pnd.FromPostgresError(err)
	}

	for n, v := range conditions {
		_, err := tx.Exec(`
			INSERT INTO sos_conditions
				(
				 	id,
					name,
					created_at,
					updated_at
				)
				SELECT $1, $2, now(), now()
				WHERE NOT EXISTS (
					SELECT
						1
					FROM
						sos_conditions
					WHERE
						name = $2::VARCHAR(50)
				);
		`, n+1, string(v))
		if err != nil {
			tx.Rollback()
			return "", pnd.FromPostgresError(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return "", pnd.FromPostgresError(err)
	}

	return "condition init success", nil
}

func (s *ConditionPostgresStore) FindConditions(ctx context.Context) ([]sos_post.Condition, *pnd.AppError) {
	conditions := make([]sos_post.Condition, 0)

	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	rows, err := tx.Query(`
		SELECT
			id,
			name
		FROM
			sos_conditions
	`)
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	for rows.Next() {
		condition := sos_post.Condition{}

		err := rows.Scan(
			&condition.ID,
			&condition.Name,
		)
		if err != nil {
			return nil, pnd.FromPostgresError(err)
		}

		conditions = append(conditions, condition)
	}

	err = tx.Commit()
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return conditions, nil
}
