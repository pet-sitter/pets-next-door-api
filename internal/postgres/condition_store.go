package postgres

import (
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

func (s *ConditionPostgresStore) InitConditions(conditions []sos_post.SosCondition) (string, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return "", err
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
			return "", err
		}
	}

	err = tx.Commit()
	if err != nil {
		return "", err
	}

	return "condition init success", nil
}

func (s *ConditionPostgresStore) FindConditions() ([]sos_post.Condition, error) {
	conditions := make([]sos_post.Condition, 0)

	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(`
		SELECT 
			id,
			name
		FROM 
			sos_conditions
	`)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		condition := sos_post.Condition{}

		err := rows.Scan(
			&condition.ID,
			&condition.Name,
		)
		if err != nil {
			return nil, err
		}

		conditions = append(conditions, condition)
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return conditions, nil
}
