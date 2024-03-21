package postgres

import (
	"context"
	pnd "github.com/pet-sitter/pets-next-door-api/api"
	utils "github.com/pet-sitter/pets-next-door-api/internal/common"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
)

func CreatePet(ctx context.Context, tx *database.Tx, pet *pet.Pet) (*pet.PetWithProfileImage, *pnd.AppError) {
	const sql = `
	INSERT INTO
		pets
		(
			owner_id,
			name,
			pet_type,
			sex,
			neutered,
			breed,
			birth_date,
			weight_in_kg,
			profile_image_id,
			created_at,
			updated_at
		)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
	RETURNING id, created_at, updated_at
	`

	if err := tx.QueryRowContext(ctx, sql,
		pet.OwnerID,
		pet.Name,
		pet.PetType,
		pet.Sex,
		pet.Neutered,
		pet.Breed,
		pet.BirthDate,
		pet.WeightInKg,
		pet.ProfileImageID,
	).Scan(&pet.ID, &pet.CreatedAt, &pet.UpdatedAt); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return FindPetByID(ctx, tx, pet.ID)
}

func FindPetByID(ctx context.Context, tx *database.Tx, id int) (*pet.PetWithProfileImage, *pnd.AppError) {
	const sql = `
	SELECT
		pets.id,
		pets.owner_id,
		pets.name,
		pets.pet_type,
		pets.sex,
		pets.neutered,
		pets.breed,
		pets.birth_date,
		pets.weight_in_kg,
		media.url AS profile_image_url,
		pets.created_at,
		pets.updated_at
	FROM
		pets
	LEFT OUTER JOIN
		media
	ON
	    pets.profile_image_id = media.id
	WHERE
		pets.id = $1 AND
		pets.deleted_at IS NULL
	`

	var pet pet.PetWithProfileImage
	if err := tx.QueryRowContext(ctx, sql,
		id,
	).Scan(
		&pet.ID,
		&pet.OwnerID,
		&pet.Name,
		&pet.PetType,
		&pet.Sex,
		&pet.Neutered,
		&pet.Breed,
		&pet.BirthDate,
		&pet.WeightInKg,
		&pet.ProfileImageURL,
		&pet.CreatedAt,
		&pet.UpdatedAt,
	); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	pet.BirthDate = utils.FormatDate(pet.BirthDate)

	return &pet, nil
}

func FindPetsByOwnerID(ctx context.Context, tx *database.Tx, ownerID int) ([]pet.PetWithProfileImage, *pnd.AppError) {
	const sql = `
	SELECT
		pets.id,
		pets.owner_id,
		pets.name,
		pets.pet_type,
		pets.sex,
		pets.neutered,
		pets.breed,
		pets.birth_date,
		pets.weight_in_kg,
		media.url AS profile_image_url,
		pets.created_at,
		pets.updated_at
	FROM
		pets
	LEFT OUTER JOIN
		media
	ON
	    pets.profile_image_id = media.id
	WHERE
		pets.owner_id = $1 AND
		pets.deleted_at IS NULL
	`

	var pets []pet.PetWithProfileImage
	rows, err := tx.QueryContext(ctx, sql,
		ownerID,
	)
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var pet pet.PetWithProfileImage
		if err := rows.Scan(
			&pet.ID,
			&pet.OwnerID,
			&pet.Name,
			&pet.PetType,
			&pet.Sex,
			&pet.Neutered,
			&pet.Breed,
			&pet.BirthDate,
			&pet.WeightInKg,
			&pet.ProfileImageURL,
			&pet.CreatedAt,
			&pet.UpdatedAt,
		); err != nil {
			return nil, pnd.FromPostgresError(err)
		}
		pet.BirthDate = utils.FormatDate(pet.BirthDate)
		pets = append(pets, pet)
	}
	if err := rows.Err(); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return pets, nil
}
