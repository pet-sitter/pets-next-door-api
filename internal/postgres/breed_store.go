package postgres

import (
	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
)

type BreedPostgresStore struct {
	db *database.DB
}

func NewBreedPostgresStore(db *database.DB) *BreedPostgresStore {
	return &BreedPostgresStore{db: db}
}

func (s *BreedPostgresStore) FindBreeds(page int, size int, petType *string) ([]*pet.Breed, error) {
	var breeds []*pet.Breed

	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	query := `
	SELECT
		id,
		name,
		pet_type,
		created_at,
		updated_at
	FROM
		breeds
	WHERE
	    (pet_type = $1 OR $1 IS NULL) AND
		deleted_at IS NULL
	ORDER BY id ASC
	LIMIT $2
	OFFSET $3
	`

	rows, err := tx.Query(query, petType, size, (page-1)*size)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		breed := &pet.Breed{}

		err := rows.Scan(&breed.ID, &breed.Name, &breed.PetType, &breed.CreatedAt, &breed.UpdatedAt)
		if err != nil {
			return nil, err
		}

		breeds = append(breeds, breed)
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return breeds, nil
}

func (s *BreedPostgresStore) FindBreedByPetTypeAndName(petType pet.PetType, name string) (*pet.Breed, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	breed := &pet.Breed{}

	err = tx.QueryRow(`
	SELECT
		id,
		name,
		pet_type,
		created_at,
		updated_at
	FROM
		breeds
	WHERE
		pet_type = $1 AND
		name = $2 AND
		deleted_at IS NULL
	`, petType, name).Scan(&breed.ID, &breed.Name, &breed.PetType, &breed.CreatedAt, &breed.UpdatedAt)

	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return nil, err
		}

		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return breed, nil
}

func (s *BreedPostgresStore) CreateBreed(breed *pet.Breed) (*pet.Breed, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	err = tx.QueryRow(`
	INSERT INTO
		breeds
		(id, name, pet_type, created_at, updated_at)
	VALUES
		(DEFAULT, $1, $2, DEFAULT, DEFAULT)
	RETURNING
		id, pet_type, name, created_at, updated_at
	`, breed.Name, breed.PetType).Scan(&breed.ID, &breed.PetType, &breed.Name, &breed.CreatedAt, &breed.UpdatedAt)

	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return nil, err
		}

		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return breed, nil
}
