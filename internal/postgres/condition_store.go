package postgres

import (
	"context"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/sospost"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
)

func InitConditions(ctx context.Context, tx *database.Tx, conditions []sospost.SosCondition) (string, *pnd.AppError) {
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
		_, err := tx.ExecContext(ctx, sql, n+1, string(v))
		if err != nil {
			return "", pnd.FromPostgresError(err)
		}
	}

	return "condition init success", nil
}

func FindConditions(ctx context.Context, tx *database.Tx) (*sospost.ConditionList, *pnd.AppError) {
	const sql = `
	SELECT
		id,
		name
	FROM
		sos_conditions
	`

	conditions := make(sospost.ConditionList, 0)
	rows, err := tx.QueryContext(ctx, sql)
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}
	for rows.Next() {
		condition := sospost.Condition{}
		if err := rows.Scan(&condition.ID, &condition.Name); err != nil {
			return nil, pnd.FromPostgresError(err)
		}
		conditions = append(conditions, &condition)
	}
	if err := rows.Err(); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return &conditions, nil
}
