package postgres

import (
	"context"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
)

type BreedPostgresStore struct {
	db *database.DB
}

func NewBreedPostgresStore(db *database.DB) *BreedPostgresStore {
	return &BreedPostgresStore{db: db}
}

func (s *BreedPostgresStore) FindBreeds(ctx context.Context, page int, size int, petType *string) (*pet.BreedList, *pnd.AppError) {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}
	defer tx.Rollback()

	breedList := pet.NewBreedList(page, size)
	rows, err := tx.QueryContext(ctx, `
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
	`, petType, size+1, (page-1)*size)
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}
	defer rows.Close()

	for rows.Next() {
		breed := &pet.Breed{}
		if err := rows.Scan(&breed.ID, &breed.Name, &breed.PetType, &breed.CreatedAt, &breed.UpdatedAt); err != nil {
			return nil, pnd.FromPostgresError(err)
		}
		breedList.Items = append(breedList.Items, *breed)
	}
	if err := rows.Err(); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	if err := tx.Commit(); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	breedList.CalcLastPage()
	return breedList, nil
}

func (s *BreedPostgresStore) FindBreedByPetTypeAndName(ctx context.Context, petType pet.PetType, name string) (*pet.Breed, *pnd.AppError) {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}
	defer tx.Rollback()

	breed := &pet.Breed{}
	err = tx.QueryRowContext(ctx, `
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
		return nil, pnd.FromPostgresError(err)
	}

	if err := tx.Commit(); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return breed, nil
}

func (s *BreedPostgresStore) CreateBreed(ctx context.Context, breed *pet.Breed) (*pet.Breed, *pnd.AppError) {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}
	defer tx.Rollback()

	err = tx.QueryRowContext(ctx, `
	INSERT INTO
		breeds
		(id, name, pet_type, created_at, updated_at)
	VALUES
		(DEFAULT, $1, $2, DEFAULT, DEFAULT)
	RETURNING
		id, pet_type, name, created_at, updated_at
	`, breed.Name, breed.PetType).Scan(&breed.ID, &breed.PetType, &breed.Name, &breed.CreatedAt, &breed.UpdatedAt)
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	if err := tx.Commit(); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return breed, nil
}
