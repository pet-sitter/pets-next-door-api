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

func (s *ConditionPostgresStore) InitConditions(condition []sos_post.SosCondition) (string, error) {
	tx, _ := s.db.Begin()
	for _, v := range condition {
		var exists bool

		err := tx.QueryRow(`
		SELECT EXISTS (
			SELECT 1
			FROM sos_conditions
			WHERE name = $1
		)`, v).Scan(&exists)

		if err != nil {
			return "condition init error", err
		}

		if exists {
			continue
		}

		_, err = tx.Exec(`
		INSERT INTO
			sos_conditions
			(
				name,
				created_at,
				updated_at
			)
		VALUES ($1, NOW(), NOW())
		`,
			v,
		)
	}

	tx.Commit()

	return "condition init success", nil
}
