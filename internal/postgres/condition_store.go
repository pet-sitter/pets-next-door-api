package postgres

import (
	"context"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database/pgx"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/sos_post"
)

func InitConditions(ctx context.Context, tx *pgx.PgxTx, conditions []sos_post.SosCondition) (string, *pnd.AppError) {
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
		_, err := tx.Exec(ctx, sql, n+1, string(v))
		if err != nil {
			return "", err
		}
	}

	return "condition init success", nil
}

func FindConditions(ctx context.Context, tx *pgx.PgxTx) (*sos_post.ConditionList, *pnd.AppError) {
	const sql = `
	SELECT
		id,
		name
	FROM
		sos_conditions
	`

	conditions := make(sos_post.ConditionList, 0)
	rows, err := tx.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		condition := sos_post.Condition{}
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
