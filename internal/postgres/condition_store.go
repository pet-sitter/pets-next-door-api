package postgres

import (
	"context"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/sos_post"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
)

type ConditionPostgresStore struct {
	conn *database.Tx
}

func NewConditionPostgresStore(conn *database.Tx) *ConditionPostgresStore {
	return &ConditionPostgresStore{
		conn: conn,
	}
}

func (s *ConditionPostgresStore) InitConditions(ctx context.Context, conditions []sos_post.SosCondition) (string, *pnd.AppError) {
	return (&conditionQueries{conn: s.conn}).InitConditions(ctx, conditions)
}

func (s *ConditionPostgresStore) FindConditions(ctx context.Context) ([]sos_post.Condition, *pnd.AppError) {
	return (&conditionQueries{conn: s.conn}).FindConditions(ctx)
}

type conditionQueries struct {
	conn *database.Tx
}

func (s *conditionQueries) InitConditions(ctx context.Context, conditions []sos_post.SosCondition) (string, *pnd.AppError) {
	const sql = `
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
	`

	for n, v := range conditions {
		_, err := s.conn.ExecContext(ctx, sql, n+1, string(v))
		if err != nil {
			return "", pnd.FromPostgresError(err)
		}
	}

	return "condition init success", nil
}

func (s *conditionQueries) FindConditions(ctx context.Context) ([]sos_post.Condition, *pnd.AppError) {
	const sql = `
	SELECT
		id,
		name
	FROM
		sos_conditions
	`

	conditions := make([]sos_post.Condition, 0)
	rows, err := s.conn.QueryContext(ctx, sql)
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}
	for rows.Next() {
		condition := sos_post.Condition{}
		if err := rows.Scan(&condition.ID, &condition.Name); err != nil {
			return nil, pnd.FromPostgresError(err)
		}
		conditions = append(conditions, condition)
	}
	if err := rows.Err(); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return conditions, nil
}
