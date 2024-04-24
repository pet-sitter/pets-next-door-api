package postgres

import (
	"context"

	"github.com/pet-sitter/pets-next-door-api/internal/domain/breed"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/commonvo"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
)

func FindBreeds(
	ctx context.Context, tx *database.Tx, page, size int, petType *string) (*breed.BreedList, *pnd.AppError,
) {
	const sql = `
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

	breedList := breed.NewBreedList(page, size)
	rows, err := tx.QueryContext(ctx, sql, petType, size+1, (page-1)*size)
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}
	defer rows.Close()

	for rows.Next() {
		breedData := &breed.Breed{}
		if err := rows.Scan(
			&breedData.ID, &breedData.Name, &breedData.PetType, &breedData.CreatedAt, &breedData.UpdatedAt,
		); err != nil {
			return nil, pnd.FromPostgresError(err)
		}
		breedList.Items = append(breedList.Items, *breedData)
	}
	if err := rows.Err(); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	breedList.CalcLastPage()
	return breedList, nil
}

func FindBreedByPetTypeAndName(
	ctx context.Context, tx *database.Tx, petType commonvo.PetType, name string,
) (*breed.Breed, *pnd.AppError) {
	const sql = `
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
	`

	breedData := &breed.Breed{}
	if err := tx.QueryRowContext(ctx, sql,
		petType,
		name,
	).Scan(
		&breedData.ID,
		&breedData.Name,
		&breedData.PetType,
		&breedData.CreatedAt,
		&breedData.UpdatedAt,
	); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return breedData, nil
}

func CreateBreed(ctx context.Context, tx *database.Tx, breedData *breed.Breed) (*breed.Breed, *pnd.AppError) {
	const sql = `
	INSERT INTO
		breeds
		(
			id,
			name,
			pet_type,
			created_at,
			updated_at
		)
	VALUES
		(DEFAULT, $1, $2, DEFAULT, DEFAULT)
	RETURNING
		id, pet_type, name, created_at, updated_at
	`

	if err := tx.QueryRowContext(ctx, sql,
		breedData.Name,
		breedData.PetType,
	).Scan(
		&breedData.ID,
		&breedData.PetType,
		&breedData.Name,
		&breedData.CreatedAt,
		&breedData.UpdatedAt,
	); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return breedData, nil
}
