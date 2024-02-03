package postgres

import (
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

func (s *SosPostPostgresStore) WriteSosPost(authorID int, utcDateStart string, utcDateEnd string, request *sos_post.WriteSosPostRequest) (*sos_post.SosPost, *pnd.AppError) {
	sosPost := &sos_post.SosPost{}

	tx, err := s.db.Begin()

	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	err = tx.QueryRow(`
	INSERT INTO
		sos_posts
		(
			author_id,
			title,
			content,
			reward,
			date_start_at,
			date_end_at,
			time_start_at,
			time_end_at,
			care_type,
		 	carer_gender,
		 	reward_amount,
		 	thumbnail_id,
			created_at,
			updated_at
		)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, NOW(), NOW())
	RETURNING
		id,
		author_id,
		title,
		content,
		reward,
		date_start_at,
		date_end_at,
		time_start_at,
		time_end_at,
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
		request.TimeStartAt,
		request.TimeEndAt,
		request.CareType,
		request.CarerGender,
		request.RewardAmount,
		request.ImageIDs[0],
	).Scan(&sosPost.ID, &sosPost.AuthorID, &sosPost.Title, &sosPost.Content, &sosPost.Reward, &sosPost.DateStartAt, &sosPost.DateEndAt, &sosPost.TimeStartAt, &sosPost.TimeEndAt, &sosPost.CareType, &sosPost.CarerGender, &sosPost.RewardAmount, &sosPost.ThumbnailID)

	if err != nil {
		tx.Rollback()
		return nil, pnd.FromPostgresError(err)
	}

	for _, imageID := range request.ImageIDs {
		_, err = tx.Exec(`
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
		)
		if err != nil {
			tx.Rollback()
			return nil, pnd.FromPostgresError(err)
		}
	}

	for _, conditionID := range request.ConditionIDs {
		_, err = tx.Exec(`
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
		)
		if err != nil {
			tx.Rollback()
			return nil, pnd.FromPostgresError(err)
		}
	}

	for _, petID := range request.PetIDs {
		_, err = tx.Exec(`
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
		)
		if err != nil {
			tx.Rollback()
			return nil, pnd.FromPostgresError(err)
		}
	}

	err = tx.Commit()

	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return sosPost, nil
}

func (s *SosPostPostgresStore) FindSosPosts(page int, size int, sortBy string) (*sos_post.SosPostList, *pnd.AppError) {
	sosPostList := sos_post.NewSosPostList(page, size)

	tx, err := s.db.Begin()
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

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
		time_start_at,
		time_end_at,
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

	for rows.Next() {
		sosPost := sos_post.SosPost{}

		err := rows.Scan(&sosPost.ID, &sosPost.AuthorID, &sosPost.Title, &sosPost.Content, &sosPost.Reward, &sosPost.DateStartAt, &sosPost.DateEndAt, &sosPost.TimeStartAt, &sosPost.TimeEndAt, &sosPost.CareType, &sosPost.CarerGender, &sosPost.RewardAmount, &sosPost.ThumbnailID, &sosPost.CreatedAt, &sosPost.UpdatedAt)
		if err != nil {

			return nil, pnd.FromPostgresError(err)
		}

		sosPostList.Items = append(sosPostList.Items, sosPost)
	}
	sosPostList.CalcLastPage()

	if err := tx.Commit(); err != nil {
		return nil, pnd.FromPostgresError(err)
	}
	if err := rows.Err(); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return sosPostList, nil
}

func (s *SosPostPostgresStore) FindSosPostsByAuthorID(authorID int, page int, size int, sortBy string) (*sos_post.SosPostList, *pnd.AppError) {
	sosPostList := sos_post.NewSosPostList(page, size)

	tx, err := s.db.Begin()
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

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
		time_start_at,
		time_end_at,
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

	for rows.Next() {
		sosPost := sos_post.SosPost{}

		err := rows.Scan(&sosPost.ID, &sosPost.AuthorID, &sosPost.Title, &sosPost.Content, &sosPost.Reward, &sosPost.DateStartAt, &sosPost.DateEndAt, &sosPost.TimeStartAt, &sosPost.TimeEndAt, &sosPost.CareType, &sosPost.CarerGender, &sosPost.RewardAmount, &sosPost.ThumbnailID, &sosPost.CreatedAt, &sosPost.UpdatedAt)
		if err != nil {
			return nil, pnd.FromPostgresError(err)
		}

		sosPostList.Items = append(sosPostList.Items, sosPost)
	}

	err = tx.Commit()
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return sosPostList, nil
}

func (s *SosPostPostgresStore) FindSosPostByID(id int) (*sos_post.SosPost, *pnd.AppError) {
	sos_post := &sos_post.SosPost{}

	tx, _ := s.db.Begin()

	err := tx.QueryRow(`
	SELECT
		id,
		author_id,
		title,
		content,
		reward,
		date_start_at,
		date_end_at,
		time_start_at,
		time_end_at,
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
		&sos_post.ID,
		&sos_post.AuthorID,
		&sos_post.Title,
		&sos_post.Content,
		&sos_post.Reward,
		&sos_post.DateStartAt,
		&sos_post.DateEndAt,
		&sos_post.TimeStartAt,
		&sos_post.TimeEndAt,
		&sos_post.CareType,
		&sos_post.CarerGender,
		&sos_post.RewardAmount,
		&sos_post.ThumbnailID,
		&sos_post.CreatedAt,
		&sos_post.UpdatedAt)

	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return sos_post, nil
}

func (s *SosPostPostgresStore) UpdateSosPost(request *sos_post.UpdateSosPostRequest) (*sos_post.SosPost, *pnd.AppError) {
	sosPost := &sos_post.SosPost{}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	_, err = tx.Exec(`
        UPDATE
            resource_media
        SET
            deleted_at = NOW()
        WHERE
            resource_id = $1
    `, request.ID)
	if err != nil {
		tx.Rollback()
		return nil, pnd.FromPostgresError(err)
	}

	for _, imageID := range request.ImageIDs {
		_, err := tx.Exec(`
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
		)
		if err != nil {
			tx.Rollback()

			return nil, pnd.FromPostgresError(err)
		}
	}

	_, err = tx.Exec(`
        UPDATE
            sos_posts_conditions
        SET
            deleted_at = NOW()
        WHERE
            sos_post_id = $1
    `, request.ID)
	if err != nil {
		tx.Rollback()
		return nil, pnd.FromPostgresError(err)
	}

	for _, conditionID := range request.ConditionIDs {
		_, err := tx.Exec(`
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
		)
		if err != nil {
			tx.Rollback()

			return nil, pnd.FromPostgresError(err)
		}
	}

	_, err = tx.Exec(`
        UPDATE
            sos_posts_pets
        SET
            deleted_at = NOW()
        WHERE
            sos_post_id = $1
    `, request.ID)
	if err != nil {
		tx.Rollback()
		return nil, pnd.FromPostgresError(err)
	}

	for _, petID := range request.PetIDs {
		_, err := tx.Exec(`
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
			tx.Rollback()

			return nil, pnd.FromPostgresError(err)
		}
	}

	err = tx.QueryRow(`
        UPDATE
            sos_posts
        SET
            title = $1,
            content = $2,
            reward = $3,
            date_start_at = $4,
            date_end_at = $5,
            time_start_at = $6,
            time_end_at = $7,
            care_type = $8,
            carer_gender = $9,
            reward_amount = $10,
            thumbnail_id = $11,
            updated_at = NOW()
        WHERE
            id = $12
        RETURNING id, author_id, title, content, reward, date_start_at, date_end_at, time_start_at, time_end_at, care_type, carer_gender, reward_amount, thumbnail_id`,
		request.Title,
		request.Content,
		request.Reward,
		request.DateStartAt,
		request.DateEndAt,
		request.TimeStartAt,
		request.TimeEndAt,
		request.CareType,
		request.CarerGender,
		request.RewardAmount,
		request.ImageIDs[0],
		request.ID,
	).Scan(&sosPost.ID, &sosPost.AuthorID, &sosPost.Title, &sosPost.Content, &sosPost.Reward, &sosPost.DateStartAt, &sosPost.DateEndAt, &sosPost.TimeStartAt, &sosPost.TimeEndAt, &sosPost.CareType, &sosPost.CarerGender, &sosPost.RewardAmount, &sosPost.ThumbnailID)

	if err != nil {
		tx.Rollback()
		return nil, pnd.FromPostgresError(err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return sosPost, nil
}

func (s *SosPostPostgresStore) FindConditionByID(id int) ([]sos_post.Condition, *pnd.AppError) {
	conditions := []sos_post.Condition{}

	tx, _ := s.db.Begin()

	rows, err := tx.Query(`
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

	err = tx.Commit()

	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return conditions, nil
}

func (s *SosPostPostgresStore) FindPetsByID(id int) ([]pet.Pet, *pnd.AppError) {
	pets := []pet.Pet{}

	tx, _ := s.db.Begin()

	rows, err := tx.Query(`
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

	err = tx.Commit()

	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return pets, nil
}
