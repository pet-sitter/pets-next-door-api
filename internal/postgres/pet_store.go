package postgres

import (
	pnd "github.com/pet-sitter/pets-next-door-api/api"
	utils "github.com/pet-sitter/pets-next-door-api/internal/common"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
)

type PetPostgresStore struct {
	db *database.DB
}

func NewPetPostgresStore(db *database.DB) *PetPostgresStore {
	return &PetPostgresStore{
		db: db,
	}
}

func (s *PetPostgresStore) CreatePet(pet *pet.Pet) (*pet.Pet, *pnd.AppError) {
	tx, _ := s.db.Begin()
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
	tx.Commit()

	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	pet.BirthDate = utils.FormatDate(pet.BirthDate)
	return pet, nil
}

func (s *PetPostgresStore) FindPetsByOwnerID(ownerID int) ([]pet.Pet, *pnd.AppError) {
	var pets []pet.Pet

	tx, _ := s.db.Begin()
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
		return nil, pnd.FromPostgresError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var pet pet.Pet

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
			return nil, pnd.FromPostgresError(err)
		}

		pet.BirthDate = utils.FormatDate(pet.BirthDate)
		pets = append(pets, pet)
	}
	tx.Commit()

	return pets, nil
}
