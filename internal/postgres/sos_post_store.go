package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	utils "github.com/pet-sitter/pets-next-door-api/internal/common"
	"time"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/media"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/sos_post"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
)

func WriteSosPost(ctx context.Context, tx *database.Tx, authorID int, request *sos_post.WriteSosPostRequest) (*sos_post.SosPost, *pnd.AppError) {
	const sql = `
	INSERT INTO
		sos_posts
		(
			author_id,
			title,
			content,
			reward,
			care_type,
		 	carer_gender,
		 	reward_type,
		 	thumbnail_id,
			created_at,
			updated_at
		)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
	RETURNING
		id,
		author_id,
		title,
		content,
		reward,
		care_type,
		carer_gender,
		reward_type,
		thumbnail_id
	`

	sosPost := &sos_post.SosPost{}
	err := tx.QueryRowContext(ctx, sql,
		authorID,
		request.Title,
		request.Content,
		request.Reward,
		request.CareType,
		request.CarerGender,
		request.RewardType,
		request.ImageIDs[0],
	).Scan(&sosPost.ID,
		&sosPost.AuthorID,
		&sosPost.Title,
		&sosPost.Content,
		&sosPost.Reward,
		&sosPost.CareType,
		&sosPost.CarerGender,
		&sosPost.RewardType,
		&sosPost.ThumbnailID)

	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	const sql2 = `
		INSERT INTO
		sos_dates
		(
			date_start_at,
			date_end_at,
			created_at,
			updated_at
		)
		VALUES ($1, $2, NOW(), NOW())
		RETURNING id, date_start_at, date_end_at, created_at, updated_at
		`

	sosDates := []sos_post.SosDates{}

	for _, date := range request.Dates {
		SosDate := sos_post.SosDates{}
		if err := tx.QueryRowContext(ctx, sql2,
			date.DateStartAt,
			date.DateEndAt,
		).Scan(
			&SosDate.ID,
			&SosDate.DateStartAt,
			&SosDate.DateEndAt,
			&SosDate.CreatedAt,
			&SosDate.UpdatedAt,
		); err != nil {
			return nil, pnd.FromPostgresError(err)
		}
		sosDates = append(sosDates, SosDate)
	}

	for _, sosDate := range sosDates {
		if _, err := tx.ExecContext(ctx, `
		INSERT INTO
			sos_posts_dates
			(
				 sos_post_id,
				 sos_dates_id,
				 created_at, 
				 updated_at
			)
		VALUES ($1, $2, NOW(), NOW())
		`,
			sosPost.ID,
			sosDate.ID,
		); err != nil {
			return nil, pnd.FromPostgresError(err)
		}
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

	return sosPost, nil
}

func FindSosPosts(ctx context.Context, tx *database.Tx, page int, size int, sortBy string, filterType string) (*sos_post.SosPostInfoList, *pnd.AppError) {
	var sortString string
	switch sortBy {
	case "newest":
		sortString = "v_sos_posts.created_at DESC"
	case "deadline":
		sortString = "earliest_date_start_at ASC"
	default:
		sortString = "sos_posts.created_at DESC"
	}

	var filterString string
	switch filterType {
	case "dog":
		filterString = "'dog' = ANY(pet_type_list)"
	case "cat":
		filterString = "'cat' = ANY(pet_type_list)"
	case "all":
		filterString = "'dog' = ANY(pet_type_list) " +
			"OR 'cat' = ANY(pet_type_list)"
	}

	query := fmt.Sprintf(`
		SELECT
			v_sos_posts.id,
			v_sos_posts.title,
			v_sos_posts.content,
			v_sos_posts.reward,
			v_sos_posts.reward_type,
			v_sos_posts.care_type,
			v_sos_posts.carer_gender,
			v_sos_posts.thumbnail_id,
			v_sos_posts.author_id,
			v_sos_posts.created_at,
			v_sos_posts.updated_at,
			v_sos_posts.dates,
			v_pets.pets_info,
			v_media.media_info,
			v_conditions.conditions_info
		FROM
			v_sos_posts
				LEFT JOIN v_pets ON v_sos_posts.id = v_pets.sos_post_id
				LEFT JOIN v_media ON v_sos_posts.id = v_media.sos_post_id
				LEFT JOIN v_conditions ON v_sos_posts.id = v_conditions.sos_post_id
		WHERE
		    %s
			AND v_sos_posts.earliest_date_start_at >= '%s'
		ORDER BY
			$1
		LIMIT $2
		OFFSET $3;

	`,
		filterString,
		utils.FormatDate(time.Now().String()),
	)

	rows, err := tx.QueryContext(ctx, query, sortString, size+1, (page-1)*size)
	if err != nil {
		return &sos_post.SosPostInfoList{}, pnd.FromPostgresError(err)
	}
	defer rows.Close()

	sosPostList := sos_post.NewSosPostInfoList(page, size)
	for rows.Next() {
		sosPost := sos_post.SosPostInfo{}
		var datesData, petsData, mediaData, conditionsData string

		if err := rows.Scan(
			&sosPost.ID,
			&sosPost.Title,
			&sosPost.Content,
			&sosPost.Reward,
			&sosPost.RewardType,
			&sosPost.CareType,
			&sosPost.CarerGender,
			&sosPost.ThumbnailID,
			&sosPost.AuthorID,
			&sosPost.CreatedAt,
			&sosPost.UpdatedAt,
			&datesData,
			&petsData,
			&mediaData,
			&conditionsData); err != nil {
			return nil, pnd.FromPostgresError(err)
		}

		if err := json.Unmarshal([]byte(datesData), &sosPost.Dates); err != nil {
			return nil, pnd.FromPostgresError(err)
		}
		if err := json.Unmarshal([]byte(petsData), &sosPost.Pets); err != nil {
			return nil, pnd.FromPostgresError(err)
		}
		if err := json.Unmarshal([]byte(mediaData), &sosPost.Media); err != nil {
			return nil, pnd.FromPostgresError(err)
		}
		if err := json.Unmarshal([]byte(conditionsData), &sosPost.Conditions); err != nil {
			return nil, pnd.FromPostgresError(err)
		}

		sosPostList.Items = append(sosPostList.Items, sosPost)
	}
	if err := rows.Err(); err != nil {
		return nil, pnd.FromPostgresError(err)
	}
	sosPostList.CalcLastPage()

	return sosPostList, nil
}

func FindSosPostsByAuthorID(ctx context.Context, tx *database.Tx, authorID int, page int, size int, sortBy string, filterType string) (*sos_post.SosPostInfoList, *pnd.AppError) {
	var sortString string
	switch sortBy {
	case "newest":
		sortString = "v_sos_posts.created_at DESC"
	case "deadline":
		sortString = "earliest_date_start_at ASC"
	default:
		sortString = "sos_posts.created_at DESC"
	}

	var filterString string
	switch filterType {
	case "dog":
		filterString = "'dog' = ANY(pet_type_list)"
	case "cat":
		filterString = "'cat' = ANY(pet_type_list)"
	case "all":
		filterString = "'dog' = ANY(pet_type_list) " +
			"OR 'cat' = ANY(pet_type_list)"
	}

	query := fmt.Sprintf(`
		SELECT
			v_sos_posts.id,
			v_sos_posts.title,
			v_sos_posts.content,
			v_sos_posts.reward,
			v_sos_posts.reward_type,
			v_sos_posts.care_type,
			v_sos_posts.carer_gender,
			v_sos_posts.thumbnail_id,
			v_sos_posts.author_id,
			v_sos_posts.created_at,
			v_sos_posts.updated_at,
			v_sos_posts.dates,
			v_pets.pets_info,
			v_media.media_info,
			v_conditions.conditions_info
		FROM
			v_sos_posts
				LEFT JOIN v_pets ON v_sos_posts.id = v_pets.sos_post_id
				LEFT JOIN v_media ON v_sos_posts.id = v_media.sos_post_id
				LEFT JOIN v_conditions ON v_sos_posts.id = v_conditions.sos_post_id
		WHERE
		    %s
			AND v_sos_posts.earliest_date_start_at >= '%s'
			AND v_sos_posts.author_id = $1
		ORDER BY
			$2
		LIMIT $3
		OFFSET $4;

	`,
		filterString,
		utils.FormatDate(time.Now().String()),
	)

	rows, err := tx.QueryContext(ctx, query, authorID, sortString, size+1, (page-1)*size)
	if err != nil {
		return &sos_post.SosPostInfoList{}, pnd.FromPostgresError(err)
	}
	defer rows.Close()

	sosPostList := sos_post.NewSosPostInfoList(page, size)
	for rows.Next() {
		sosPost := sos_post.SosPostInfo{}
		var datesData, petsData, mediaData, conditionsData string

		if err := rows.Scan(
			&sosPost.ID,
			&sosPost.Title,
			&sosPost.Content,
			&sosPost.Reward,
			&sosPost.RewardType,
			&sosPost.CareType,
			&sosPost.CarerGender,
			&sosPost.ThumbnailID,
			&sosPost.AuthorID,
			&sosPost.CreatedAt,
			&sosPost.UpdatedAt,
			&datesData,
			&petsData,
			&mediaData,
			&conditionsData); err != nil {
			return nil, pnd.FromPostgresError(err)
		}

		if err := json.Unmarshal([]byte(datesData), &sosPost.Dates); err != nil {
			return nil, pnd.FromPostgresError(err)
		}
		if err := json.Unmarshal([]byte(petsData), &sosPost.Pets); err != nil {
			return nil, pnd.FromPostgresError(err)
		}
		if err := json.Unmarshal([]byte(mediaData), &sosPost.Media); err != nil {
			return nil, pnd.FromPostgresError(err)
		}
		if err := json.Unmarshal([]byte(conditionsData), &sosPost.Conditions); err != nil {
			return nil, pnd.FromPostgresError(err)
		}

		sosPostList.Items = append(sosPostList.Items, sosPost)
	}
	if err := rows.Err(); err != nil {
		return nil, pnd.FromPostgresError(err)
	}
	sosPostList.CalcLastPage()

	return sosPostList, nil
}

func FindSosPostByID(ctx context.Context, tx *database.Tx, id int) (*sos_post.SosPost, *pnd.AppError) {
	const query = `
	SELECT
		id,
		author_id,
		title,
		content,
		reward,
		care_type,
		carer_gender,
		reward_type,
		thumbnail_id,
		created_at,
		updated_at
	FROM
		sos_posts
	WHERE
		id = $1 AND
		deleted_at IS NULL
	`

	sosPost := &sos_post.SosPost{}
	if err := tx.QueryRowContext(ctx, query, id).Scan(
		&sosPost.ID,
		&sosPost.AuthorID,
		&sosPost.Title,
		&sosPost.Content,
		&sosPost.Reward,
		&sosPost.CareType,
		&sosPost.CarerGender,
		&sosPost.RewardType,
		&sosPost.ThumbnailID,
		&sosPost.CreatedAt,
		&sosPost.UpdatedAt,
	); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return sosPost, nil
}

func UpdateSosPost(ctx context.Context, tx *database.Tx, request *sos_post.UpdateSosPostRequest) (*sos_post.SosPost, *pnd.AppError) {
	sosPost := &sos_post.SosPost{}

	if _, err := tx.ExecContext(ctx, `
		UPDATE
			sos_posts_dates
		SET
			deleted_at = NOW()
		WHERE
			sos_post_id = $1
	`, request.ID,
	); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	if _, err := tx.ExecContext(ctx, `
		UPDATE
			resource_media
		SET
			deleted_at = NOW()
		WHERE
			resource_id = $1
	`, request.ID,
	); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	const sql = `
		INSERT INTO
		sos_dates
		(
			date_start_at,
			date_end_at,
			created_at,
			updated_at
		)
		VALUES ($1, $2, NOW(), NOW())
		RETURNING id, date_start_at, date_end_at, created_at, updated_at
		`

	sosDates := []sos_post.SosDates{}

	for _, date := range request.Dates {
		SosDate := sos_post.SosDates{}
		if err := tx.QueryRowContext(ctx, sql,
			date.DateStartAt,
			date.DateEndAt,
		).Scan(
			&SosDate.ID,
			&SosDate.DateStartAt,
			&SosDate.DateEndAt,
			&SosDate.CreatedAt,
			&SosDate.UpdatedAt,
		); err != nil {
			return nil, pnd.FromPostgresError(err)
		}
		sosDates = append(sosDates, SosDate)
	}

	for _, sosDate := range sosDates {
		if _, err := tx.ExecContext(ctx, `
		INSERT INTO
			sos_posts_dates
			(
				 sos_post_id,
				 sos_dates_id,
				 created_at, 
				 updated_at
			)
		VALUES ($1, $2, NOW(), NOW())
		`,
			request.ID,
			sosDate.ID,
		); err != nil {
			return nil, pnd.FromPostgresError(err)
		}
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
		VALUES ($1, $2, $3, NOW(), NOW())
		`,
			imageID,
			request.ID,
			media.SosResourceType,
		); err != nil {
			return nil, pnd.FromPostgresError(err)
		}
	}

	if _, err := tx.ExecContext(ctx, `
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
		VALUES ($1, $2, NOW(), NOW())
		`,
			request.ID,
			conditionID,
		); err != nil {
			return nil, pnd.FromPostgresError(err)
		}
	}

	if _, err := tx.ExecContext(ctx, `
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
		if _, err := tx.ExecContext(ctx, `
		INSERT INTO
			sos_posts_pets
			(
				sos_post_id,
				pet_id,
				created_at,
				updated_at
			)
		VALUES ($1, $2, NOW(), NOW())
		`,
			request.ID,
			petID,
		); err != nil {
			return nil, pnd.FromPostgresError(err)
		}
	}

	const sql2 = `
	UPDATE
		sos_posts
	SET
		title = $1,
		content = $2,
		reward = $3,
		care_type = $4,
		carer_gender = $5,
		reward_type = $6,
		thumbnail_id = $7,
		updated_at = NOW()
	WHERE
		id = $8
	RETURNING
		id,
		author_id,
		title,
		content,
		reward,
		care_type,
		carer_gender,
		reward_type,
		thumbnail_id
	`

	if err := tx.QueryRowContext(ctx, sql2,
		request.Title,
		request.Content,
		request.Reward,
		request.CareType,
		request.CarerGender,
		request.RewardType,
		request.ImageIDs[0],
		request.ID,
	).Scan(
		&sosPost.ID,
		&sosPost.AuthorID,
		&sosPost.Title,
		&sosPost.Content,
		&sosPost.Reward,
		&sosPost.CareType,
		&sosPost.CarerGender,
		&sosPost.RewardType,
		&sosPost.ThumbnailID,
	); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return sosPost, nil
}

func FindConditionByID(ctx context.Context, tx *database.Tx, id int) (*sos_post.ConditionList, *pnd.AppError) {
	const sql = `
	SELECT
		sos_conditions.id,
		sos_conditions.name,
		sos_conditions.created_at,
		sos_conditions.updated_at
	FROM
		sos_conditions
	INNER JOIN
		sos_posts_conditions
	ON
		sos_conditions.id = sos_posts_conditions.sos_condition_id
	WHERE
		sos_posts_conditions.sos_post_id = $1 AND
		sos_posts_conditions.deleted_at IS NULL
	`

	conditions := sos_post.ConditionList{}
	rows, err := tx.QueryContext(ctx, sql, id)
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}
	defer rows.Close()

	for rows.Next() {
		condition := sos_post.Condition{}
		if err := rows.Scan(
			&condition.ID,
			&condition.Name,
			&condition.CreatedAt,
			&condition.UpdatedAt,
		); err != nil {
			return nil, pnd.FromPostgresError(err)
		}
		conditions = append(conditions, &condition)
	}
	if err := rows.Err(); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return &conditions, nil
}

func FindPetsByID(ctx context.Context, tx *database.Tx, id int) (*pet.PetList, *pnd.AppError) {
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
		pets.created_at,
		pets.updated_at
	FROM
		pets
	INNER JOIN
		sos_posts_pets
	ON
		pets.id = sos_posts_pets.pet_id
	WHERE
		sos_posts_pets.sos_post_id = $1 AND
		sos_posts_pets.deleted_at IS NULL
	`

	pets := pet.PetList{}
	rows, err := tx.QueryContext(ctx, sql, id)
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}
	defer rows.Close()

	for rows.Next() {
		pet := pet.Pet{}
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
			&pet.Remarks,
			&pet.CreatedAt,
			&pet.UpdatedAt,
		); err != nil {
			return nil, pnd.FromPostgresError(err)
		}
		pets = append(pets, &pet)
	}
	if err := rows.Err(); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return &pets, nil
}

func FindDatesBySosPostID(ctx context.Context, tx *database.Tx, sosPostID int) (*sos_post.SosDatesList, *pnd.AppError) {
	const sql = `
		SELECT
		    sos_dates.id,
			sos_dates.date_start_at,
			sos_dates.date_end_at,
			sos_dates.created_at,
			sos_dates.updated_at
		FROM
			sos_dates
		INNER JOIN
			sos_posts_dates
		ON sos_dates.id = sos_posts_dates.sos_dates_id
		WHERE 
		    sos_posts_dates.sos_post_id = $1 AND
			sos_posts_dates.deleted_at IS NULL
	`

	var sosDates sos_post.SosDatesList
	rows, err := tx.QueryContext(ctx, sql, sosPostID)
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}
	defer rows.Close()

	for rows.Next() {
		sosDate := sos_post.SosDates{}
		if err := rows.Scan(
			&sosDate.ID,
			&sosDate.DateStartAt,
			&sosDate.DateEndAt,
			&sosDate.CreatedAt,
			&sosDate.DeletedAt,
		); err != nil {
			return nil, pnd.FromPostgresError(err)
		}
		sosDates = append(sosDates, &sosDate)
	}

	return &sosDates, nil
}
