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

func (s *PetPostgresStore) CreatePet(ctx context.Context, pet *pet.Pet) (*pet.PetWithProfileImage, *pnd.AppError) {
	return (&petQueries{conn: s.conn}).CreatePet(ctx, pet)
}

func (s *PetPostgresStore) FindPetByID(ctx context.Context, id int) (*pet.PetWithProfileImage, *pnd.AppError) {
	return (&petQueries{conn: s.conn}).FindPetByID(ctx, id)
}

func (s *PetPostgresStore) FindPetsByOwnerID(ctx context.Context, ownerID int) ([]pet.PetWithProfileImage, *pnd.AppError) {
	return (&petQueries{conn: s.conn}).FindPetsByOwnerID(ctx, ownerID)
}

type petQueries struct {
	conn database.DBTx
}

func (s *petQueries) CreatePet(ctx context.Context, pet *pet.Pet) (*pet.PetWithProfileImage, *pnd.AppError) {
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

	if err := s.conn.QueryRowContext(ctx, sql,
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

	return s.FindPetByID(ctx, pet.ID)
}

func (s *petQueries) FindPetByID(ctx context.Context, id int) (*pet.PetWithProfileImage, *pnd.AppError) {
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
	if err := s.conn.QueryRowContext(ctx, sql,
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

func (s *petQueries) FindPetsByOwnerID(ctx context.Context, ownerID int) ([]pet.PetWithProfileImage, *pnd.AppError) {
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
	rows, err := s.conn.QueryContext(ctx, sql,
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
