package database

import "github.com/pet-sitter/pets-next-door-api/internal/models"

func (tx *Tx) CreatePet(pet *models.Pet) (*models.Pet, error) {
	err := tx.QueryRow(`
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
			created_at,
			updated_at
		)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
	RETURNING id, created_at, updated_at
	`,
		pet.OwnerID,
		pet.Name,
		pet.PetType,
		pet.Sex,
		pet.Neutered,
		pet.Breed,
		pet.BirthDate,
		pet.WeightInKg,
	).Scan(&pet.ID, &pet.CreatedAt, &pet.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return pet, nil
}

func (tx *Tx) FindPetsByOwnerID(ownerID int) ([]models.Pet, error) {
	var pets []models.Pet

	rows, err := tx.Query(`
	SELECT
		id,
		owner_id,
		name,
		pet_type,
		sex,
		neutered,
		breed,
		birth_date,
		weight_in_kg,
		created_at,
		updated_at
	FROM
		pets
	WHERE
		owner_id = $1 AND
		deleted_at IS NULL
	`,
		ownerID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var pet models.Pet

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
			&pet.CreatedAt,
			&pet.UpdatedAt,
		); err != nil {
			return nil, err
		}

		pets = append(pets, pet)
	}

	return pets, nil
}
