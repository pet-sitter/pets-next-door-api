package postgres

import (
	"context"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	utils "github.com/pet-sitter/pets-next-door-api/internal/common"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
)

type PetPostgresStore struct {
	conn *database.Tx
}

func NewPetPostgresStore(conn *database.Tx) *PetPostgresStore {
	return &PetPostgresStore{
		conn: conn,
	}
}

func (s *PetPostgresStore) CreatePet(ctx context.Context, pet *pet.Pet) (*pet.Pet, *pnd.AppError) {
	return (&petQueries{conn: s.conn}).CreatePet(ctx, pet)
}

func (s *PetPostgresStore) FindPetsByOwnerID(ctx context.Context, ownerID int) ([]pet.Pet, *pnd.AppError) {
	return (&petQueries{conn: s.conn}).FindPetsByOwnerID(ctx, ownerID)
}

type petQueries struct {
	conn database.DBTx
}

func (s *petQueries) CreatePet(ctx context.Context, pet *pet.Pet) (*pet.Pet, *pnd.AppError) {
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
			created_at,
			updated_at
		)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
	RETURNING id, created_at, updated_at
	`

	if err := s.conn.QueryRowContext(ctx, sql,
		pet.OwnerID,
		pet.Name,
		pet.PetType,
		pet.Sex,
		pet.Neutered,
		pet.Breed,
		pet.BirthDate,
		pet.WeightInKg,
	).Scan(&pet.ID, &pet.CreatedAt, &pet.UpdatedAt); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	pet.BirthDate = utils.FormatDate(pet.BirthDate)
	return pet, nil
}

func (s *petQueries) FindPetsByOwnerID(ctx context.Context, ownerID int) ([]pet.Pet, *pnd.AppError) {
	const sql = `
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
	`

	var pets []pet.Pet
	rows, err := s.conn.QueryContext(ctx, sql,
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
	if err := rows.Err(); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return pets, nil
}
