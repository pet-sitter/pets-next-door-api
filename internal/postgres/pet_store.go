package postgres

import (
	"context"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	utils "github.com/pet-sitter/pets-next-door-api/internal/common"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
)

func CreatePet(ctx context.Context, tx *database.Tx, petData *pet.Pet) (*pet.PetWithProfileImage, *pnd.AppError) {
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
		 	remarks,
			profile_image_id,
			created_at,
			updated_at
		)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW())
	RETURNING id, created_at, updated_at
	`

	if err := tx.QueryRowContext(ctx, sql,
		petData.OwnerID,
		petData.Name,
		petData.PetType,
		petData.Sex,
		petData.Neutered,
		petData.Breed,
		petData.BirthDate,
		petData.WeightInKg,
		petData.Remarks,
		petData.ProfileImageID,
	).Scan(&petData.ID, &petData.CreatedAt, &petData.UpdatedAt); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return FindPetByID(ctx, tx, petData.ID)
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
		pets.remarks,
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

	var petData pet.PetWithProfileImage
	if err := tx.QueryRowContext(ctx, sql,
		id,
	).Scan(
		&petData.ID,
		&petData.OwnerID,
		&petData.Name,
		&petData.PetType,
		&petData.Sex,
		&petData.Neutered,
		&petData.Breed,
		&petData.BirthDate,
		&petData.WeightInKg,
		&petData.Remarks,
		&petData.ProfileImageURL,
		&petData.CreatedAt,
		&petData.UpdatedAt,
	); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	petData.BirthDate = utils.FormatDate(petData.BirthDate)

	return &petData, nil
}

func FindPetsByOwnerID(ctx context.Context, tx *database.Tx, ownerID int) (*pet.PetWithProfileList, *pnd.AppError) {
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
		pets.remarks,
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

	var pets pet.PetWithProfileList
	rows, err := tx.QueryContext(ctx, sql,
		ownerID,
	)
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var petData pet.PetWithProfileImage
		if err := rows.Scan(
			&petData.ID,
			&petData.OwnerID,
			&petData.Name,
			&petData.PetType,
			&petData.Sex,
			&petData.Neutered,
			&petData.Breed,
			&petData.BirthDate,
			&petData.WeightInKg,
			&petData.Remarks,
			&petData.ProfileImageURL,
			&petData.CreatedAt,
			&petData.UpdatedAt,
		); err != nil {
			return nil, pnd.FromPostgresError(err)
		}
		petData.BirthDate = utils.FormatDate(petData.BirthDate)
		pets = append(pets, &petData)
	}
	if err := rows.Err(); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return &pets, nil
}

func UpdatePet(ctx context.Context, tx *database.Tx, petID int, updatePetRequest *pet.UpdatePetRequest) *pnd.AppError {
	const sql = `
	UPDATE
		pets
	SET
		name = $1,
		neutered = $2,
		breed = $3,
		birth_date = $4,
		weight_in_kg = $5,
		remarks = $6,
		profile_image_id = $7,
		updated_at = NOW()
	WHERE
		id = $8
	`

	if _, err := tx.ExecContext(ctx, sql,
		updatePetRequest.Name,
		updatePetRequest.Neutered,
		updatePetRequest.Breed,
		updatePetRequest.BirthDate,
		updatePetRequest.WeightInKg,
		updatePetRequest.Remarks,
		updatePetRequest.ProfileImageID,
		petID,
	); err != nil {
		return pnd.FromPostgresError(err)
	}

	return nil
}

func DeletePet(ctx context.Context, tx *database.Tx, petID int) *pnd.AppError {
	const sql = `
	UPDATE
		pets
	SET
		deleted_at = NOW()
	WHERE
		id = $1
	`

	if _, err := tx.ExecContext(ctx, sql,
		petID,
	); err != nil {
		return pnd.FromPostgresError(err)
	}

	return nil
}
