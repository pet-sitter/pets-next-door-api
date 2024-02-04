package postgres

import (
	"context"
	"fmt"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/media"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/sos_post"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
)

type SosPostPostgresStore struct {
	db *database.DB
}

func NewSosPostPostgresStore(db *database.DB) *SosPostPostgresStore {
	return &SosPostPostgresStore{
		db: db,
	}
}

func (s *SosPostPostgresStore) WriteSosPost(ctx context.Context, authorID int, utcDateStart string, utcDateEnd string, request *sos_post.WriteSosPostRequest) (*sos_post.SosPost, *pnd.AppError) {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}
	defer tx.Rollback()

	sosPost := &sos_post.SosPost{}
	err = tx.QueryRowContext(ctx, `
	INSERT INTO
		sos_posts
		(
			author_id,
			title,
			content,
			reward,
			date_start_at,
			date_end_at,
			care_type,
		 	carer_gender,
		 	reward_amount,
		 	thumbnail_id,
			created_at,
			updated_at
		)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW())
	RETURNING
		id,
		author_id,
		title,
		content,
		reward,
		date_start_at,
		date_end_at,
		care_type,
		carer_gender,
		reward_amount,
		thumbnail_id`,

		authorID,
		request.Title,
		request.Content,
		request.Reward,
		utcDateStart,
		utcDateEnd,
		request.CareType,
		request.CarerGender,
		request.RewardAmount,
		request.ImageIDs[0],
	).Scan(&sosPost.ID,
		&sosPost.AuthorID,
		&sosPost.Title,
		&sosPost.Content,
		&sosPost.Reward,
		&sosPost.DateStartAt,
		&sosPost.DateEndAt,
		&sosPost.CareType,
		&sosPost.CarerGender,
		&sosPost.RewardAmount,
		&sosPost.ThumbnailID)

	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	for _, imageID := range request.ImageIDs {
		if _, err := tx.ExecContext(ctx, `
		INSERT INTO
			resource_media
			(
				media_id,
				resource_id,
				resource_type,
				created_at,
				updated_at
			)
		VALUES ($1, $2, $3, NOW(), NOW())`,
			imageID,
			sosPost.ID,
			media.SosResourceType,
		); err != nil {
			return nil, pnd.FromPostgresError(err)
		}
	}

	for _, conditionID := range request.ConditionIDs {
		if _, err := tx.ExecContext(ctx, `
		INSERT INTO
			sos_posts_conditions
			(
				sos_post_id,
				sos_condition_id,
				created_at,
				updated_at
			)
		VALUES ($1, $2, NOW(), NOW())`,
			sosPost.ID,
			conditionID,
		); err != nil {
			return nil, pnd.FromPostgresError(err)
		}
	}

	for _, petID := range request.PetIDs {
		if _, err := tx.ExecContext(ctx, `
		INSERT INTO
			sos_posts_pets
			(
				sos_post_id,
				pet_id,
				created_at,
				updated_at
			)
		VALUES ($1, $2, NOW(), NOW())`,
			sosPost.ID,
			petID,
		); err != nil {
			return nil, pnd.FromPostgresError(err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return sosPost, nil
}

func (s *SosPostPostgresStore) FindSosPosts(ctx context.Context, page int, size int, sortBy string) (*sos_post.SosPostList, *pnd.AppError) {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}
	defer tx.Rollback()

	sortColumn := ""
	sortOrder := ""

	switch sortBy {
	case "newest":
		sortColumn = "created_at"
		sortOrder = "DESC"
	case "deadline":
		sortColumn = "date_end_at"
		sortOrder = "ASC"
	default:
		sortColumn = "created_at"
		sortOrder = "DESC"
	}

	query := fmt.Sprintf(`
	SELECT
	    id,
		author_id,
		title,
		content,
		reward,
		date_start_at,
		date_end_at,
		care_type,
		carer_gender,
		reward_amount,
		thumbnail_id,
		created_at,
		updated_at
	FROM
	    sos_posts
	WHERE
	    deleted_at IS NULL
    ORDER BY %s %s
    LIMIT $1 OFFSET $2
    `, sortColumn, sortOrder)

	rows, err := tx.Query(query, size+1, (page-1)*size)

	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}
	defer rows.Close()

	sosPostList := sos_post.NewSosPostList(page, size)
	for rows.Next() {
		sosPost := sos_post.SosPost{}

		err := rows.Scan(
			&sosPost.ID,
			&sosPost.AuthorID,
			&sosPost.Title,
			&sosPost.Content,
			&sosPost.Reward,
			&sosPost.DateStartAt,
			&sosPost.DateEndAt,
			&sosPost.CareType,
			&sosPost.CarerGender,
			&sosPost.RewardAmount,
			&sosPost.ThumbnailID,
			&sosPost.CreatedAt,
			&sosPost.UpdatedAt)

		if err != nil {
			return nil, pnd.FromPostgresError(err)
		}
		sosPostList.Items = append(sosPostList.Items, sosPost)
	}
	if err := rows.Err(); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	if err := tx.Commit(); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	sosPostList.CalcLastPage()
	return sosPostList, nil
}

func (s *SosPostPostgresStore) FindSosPostsByAuthorID(ctx context.Context, authorID int, page int, size int, sortBy string) (*sos_post.SosPostList, *pnd.AppError) {
	sosPostList := sos_post.NewSosPostList(page, size)

	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}
	defer tx.Rollback()

	sortColumn := ""
	sortOrder := ""

	switch sortBy {
	case "newest":
		sortColumn = "created_at"
		sortOrder = "DESC"
	case "deadline":
		sortColumn = "date_end_at"
		sortOrder = "ASC"
	default:
		sortColumn = "created_at"
		sortOrder = "DESC"
	}

	query := fmt.Sprintf(`
	SELECT
	    id,
		author_id,
		title,
		content,
		reward,
		date_start_at,
		date_end_at,
		care_type,
		carer_gender,
		reward_amount,
		thumbnail_id,
		created_at,
		updated_at
	FROM
	    sos_posts
	WHERE
	    author_id = $1 AND
	    deleted_at IS NULL
	ORDER BY %s %s
    LIMIT $2 OFFSET $3
    `, sortColumn, sortOrder)

	rows, err := tx.Query(query, authorID, size+1, (page-1)*size)

	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}
	defer rows.Close()

	for rows.Next() {
		sosPost := sos_post.SosPost{}

		err := rows.Scan(
			&sosPost.ID,
			&sosPost.AuthorID,
			&sosPost.Title,
			&sosPost.Content,
			&sosPost.Reward,
			&sosPost.DateStartAt,
			&sosPost.DateEndAt,
			&sosPost.CareType,
			&sosPost.CarerGender,
			&sosPost.RewardAmount,
			&sosPost.ThumbnailID,
			&sosPost.CreatedAt,
			&sosPost.UpdatedAt)

		if err != nil {
			return nil, pnd.FromPostgresError(err)
		}
		sosPostList.Items = append(sosPostList.Items, sosPost)
	}

	if err := tx.Commit(); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return sosPostList, nil
}

func (s *SosPostPostgresStore) FindSosPostByID(ctx context.Context, id int) (*sos_post.SosPost, *pnd.AppError) {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}
	defer tx.Rollback()

	sosPost := &sos_post.SosPost{}
	err = tx.QueryRowContext(ctx, `
	SELECT
		id,
		author_id,
		title,
		content,
		reward,
		date_start_at,
		date_end_at,
		care_type,
		carer_gender,
		reward_amount,
		thumbnail_id,
		created_at,
		updated_at
	FROM
	    sos_posts
	WHERE
	    id = $1 AND
	    deleted_at IS NULL
	`, id).Scan(
		&sosPost.ID,
		&sosPost.AuthorID,
		&sosPost.Title,
		&sosPost.Content,
		&sosPost.Reward,
		&sosPost.DateStartAt,
		&sosPost.DateEndAt,
		&sosPost.CareType,
		&sosPost.CarerGender,
		&sosPost.RewardAmount,
		&sosPost.ThumbnailID,
		&sosPost.CreatedAt,
		&sosPost.UpdatedAt)

	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	if err := tx.Commit(); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return sosPost, nil
}

func (s *SosPostPostgresStore) UpdateSosPost(ctx context.Context, request *sos_post.UpdateSosPostRequest) (*sos_post.SosPost, *pnd.AppError) {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}
	defer tx.Rollback()

	sosPost := &sos_post.SosPost{}
	_, err = tx.ExecContext(ctx, `
        UPDATE
            resource_media
        SET
            deleted_at = NOW()
        WHERE
            resource_id = $1
    `, request.ID)
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	for _, imageID := range request.ImageIDs {
		if _, err := tx.ExecContext(ctx, `
		INSERT INTO
				resource_media
				(
						media_id,
						resource_id,
						resource_type,
						created_at,
						updated_at
				)
		VALUES ($1, $2, $3, NOW(), NOW())`,
			imageID,
			request.ID,
			media.SosResourceType,
		); err != nil {
			return nil, pnd.FromPostgresError(err)
		}
	}

	if _, err = tx.ExecContext(ctx, `
        UPDATE
            sos_posts_conditions
        SET
            deleted_at = NOW()
        WHERE
            sos_post_id = $1
    `, request.ID); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	for _, conditionID := range request.ConditionIDs {
		if _, err := tx.ExecContext(ctx, `
		INSERT INTO
				sos_posts_conditions
				(
						sos_post_id,
						sos_condition_id,
						created_at,
						updated_at
				)
		VALUES ($1, $2, NOW(), NOW())`,
			request.ID,
			conditionID,
		); err != nil {
			return nil, pnd.FromPostgresError(err)
		}
	}

	if _, err = tx.ExecContext(ctx, `
        UPDATE
            sos_posts_pets
        SET
            deleted_at = NOW()
        WHERE
            sos_post_id = $1
    `, request.ID); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	for _, petID := range request.PetIDs {
		_, err := tx.ExecContext(ctx, `
            INSERT INTO
                sos_posts_pets
                (
                    sos_post_id,
                    pet_id,
                    created_at,
                    updated_at
                )
            VALUES ($1, $2, NOW(), NOW())`,
			request.ID,
			petID,
		)
		if err != nil {
			return nil, pnd.FromPostgresError(err)
		}
	}

	err = tx.QueryRowContext(ctx, `
        UPDATE
            sos_posts
        SET
            title = $1,
            content = $2,
            reward = $3,
            date_start_at = $4,
            date_end_at = $5,
            care_type = $6,
            carer_gender = $7,
            reward_amount = $8,
            thumbnail_id = $9,
            updated_at = NOW()
        WHERE
            id = $10
        RETURNING id, author_id, title, content, reward, date_start_at, date_end_at, care_type, carer_gender, reward_amount, thumbnail_id`,
		request.Title,
		request.Content,
		request.Reward,
		request.DateStartAt,
		request.DateEndAt,
		request.CareType,
		request.CarerGender,
		request.RewardAmount,
		request.ImageIDs[0],
		request.ID,
	).Scan(
		&sosPost.ID,
		&sosPost.AuthorID,
		&sosPost.Title,
		&sosPost.Content,
		&sosPost.Reward,
		&sosPost.DateStartAt,
		&sosPost.DateEndAt,
		&sosPost.CareType,
		&sosPost.CarerGender,
		&sosPost.RewardAmount,
		&sosPost.ThumbnailID)

	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	if err = tx.Commit(); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return sosPost, nil
}

func (s *SosPostPostgresStore) FindConditionByID(ctx context.Context, id int) ([]sos_post.Condition, *pnd.AppError) {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}
	defer tx.Rollback()

	conditions := []sos_post.Condition{}
	rows, err := tx.QueryContext(ctx, `
	SELECT
		sos_conditions.id,
		sos_conditions.name,
		sos_conditions.created_at,
		sos_conditions.updated_at
	FROM
		sos_conditions
	INNER JOIN
		sos_posts_conditions
		ON sos_conditions.id = sos_posts_conditions.sos_condition_id
	WHERE
	    sos_posts_conditions.sos_post_id = $1 AND
		sos_posts_conditions.deleted_at IS NULL
	`, id)
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}
	defer rows.Close()

	for rows.Next() {
		condition := sos_post.Condition{}
		err := rows.Scan(
			&condition.ID,
			&condition.Name,
			&condition.CreatedAt,
			&condition.UpdatedAt,
		)
		if err != nil {
			return nil, pnd.FromPostgresError(err)
		}
		conditions = append(conditions, condition)
	}
	if err := rows.Err(); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	if err := tx.Commit(); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return conditions, nil
}

func (s *SosPostPostgresStore) FindPetsByID(ctx context.Context, id int) ([]pet.Pet, *pnd.AppError) {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}
	defer tx.Rollback()

	pets := []pet.Pet{}
	rows, err := tx.QueryContext(ctx, `
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
		pets.created_at,
		pets.updated_at
	FROM
		pets
	INNER JOIN
		sos_posts_pets
		ON pets.id = sos_posts_pets.pet_id
	WHERE
		sos_posts_pets.sos_post_id = $1 AND
		sos_posts_pets.deleted_at IS NULL
	`, id)
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}
	defer rows.Close()

	for rows.Next() {
		pet := pet.Pet{}
		err := rows.Scan(
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
		)
		if err != nil {
			return nil, pnd.FromPostgresError(err)
		}
		pets = append(pets, pet)
	}
	if err := rows.Err(); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	if err := tx.Commit(); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return pets, nil
}
