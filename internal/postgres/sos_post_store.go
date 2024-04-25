package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	utils "github.com/pet-sitter/pets-next-door-api/internal/common"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/media"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/sospost"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
)

const writeSOSPostQuery = `
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

func WriteSOSPost(
	ctx context.Context, tx *database.Tx, authorID int, request *sospost.WriteSOSPostRequest,
) (*sospost.SOSPost, *pnd.AppError) {
	sosPost := &sospost.SOSPost{}
	err := tx.QueryRowContext(ctx, writeSOSPostQuery,
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

	sosDates := []sospost.SOSDates{}

	for _, date := range request.Dates {
		sosDate := sospost.SOSDates{}
		if err := tx.QueryRowContext(ctx, sql2,
			date.DateStartAt,
			date.DateEndAt,
		).Scan(
			&sosDate.ID,
			&sosDate.DateStartAt,
			&sosDate.DateEndAt,
			&sosDate.CreatedAt,
			&sosDate.UpdatedAt,
		); err != nil {
			return nil, pnd.FromPostgresError(err)
		}
		sosDates = append(sosDates, sosDate)
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
			media.SOSResourceType,
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

func readSOSPostRows(rows *sql.Rows, page, size int) (*sospost.SOSPostInfoList, *pnd.AppError) {
	sosPostList := sospost.NewSOSPostInfoList(page, size)

	for rows.Next() {
		sosPost := sospost.SOSPostInfo{}
		var datesData, petsData, mediaData, conditionsData []byte
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

		if err := json.Unmarshal(datesData, &sosPost.Dates); err != nil {
			return nil, pnd.FromPostgresError(err)
		}
		if err := json.Unmarshal(petsData, &sosPost.Pets); err != nil {
			return nil, pnd.FromPostgresError(err)
		}
		if err := json.Unmarshal(mediaData, &sosPost.Media); err != nil {
			return nil, pnd.FromPostgresError(err)
		}
		if err := json.Unmarshal(conditionsData, &sosPost.Conditions); err != nil {
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

func FindSOSPosts(
	ctx context.Context, tx *database.Tx, page, size int, sortBy, filterType string,
) (*sospost.SOSPostInfoList, *pnd.AppError) {
	var sortString string
	switch sortBy {
	case "newest":
		sortString = "v_sos_posts.created_at DESC"
	case "deadline":
		sortString = "v_sos_posts.earliest_date_start_at"
	}

	var filterString string
	switch filterType {
	case "dog":
		filterString = "AND " +
			"NOT EXISTS " +
			"(SELECT 1 " +
			"FROM unnest(pet_type_list) AS pet_type " +
			"WHERE pet_type <> 'dog')"
	case "cat":
		filterString = "AND " +
			"NOT EXISTS " +
			"(SELECT 1 " +
			"FROM unnest(pet_type_list) AS pet_type " +
			"WHERE pet_type <> 'cat')"
	case "all":
		filterString = ""
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
			v_pets_for_sos_posts.pets_info,
			v_media_for_sos_posts.media_info,
			v_conditions.conditions_info
		FROM
			v_sos_posts
				LEFT JOIN v_pets_for_sos_posts ON v_sos_posts.id = v_pets_for_sos_posts.sos_post_id
				LEFT JOIN v_media_for_sos_posts ON v_sos_posts.id = v_media_for_sos_posts.sos_post_id
				LEFT JOIN v_conditions ON v_sos_posts.id = v_conditions.sos_post_id
		WHERE
		    v_sos_posts.earliest_date_start_at >= '%s'
			%s
		ORDER BY
			%s
		LIMIT $1
		OFFSET $2;

	`,
		utils.FormatDate(time.Now().String()),
		filterString,
		sortString,
	)

	rows, err := tx.QueryContext(ctx, query, size+1, (page-1)*size)
	if err != nil {
		return &sospost.SOSPostInfoList{}, pnd.FromPostgresError(err)
	}
	defer rows.Close()

	sosPostList, err2 := readSOSPostRows(rows, page, size)
	if err2 != nil {
		return nil, err2
	}

	return sosPostList, nil
}

func FindSOSPostsByAuthorID(
	ctx context.Context, tx *database.Tx, authorID, page, size int, sortBy, filterType string,
) (*sospost.SOSPostInfoList, *pnd.AppError) {
	var sortString string
	switch sortBy {
	case "newest":
		sortString = "v_sos_posts.created_at DESC"
	case "deadline":
		sortString = "v_sos_posts.earliest_date_start_at"
	}

	var filterString string
	switch filterType {
	case "dog":
		filterString = "AND " +
			"NOT EXISTS " +
			"(SELECT 1 " +
			"FROM unnest(pet_type_list) AS pet_type " +
			"WHERE pet_type <> 'dog')"
	case "cat":
		filterString = "AND " +
			"NOT EXISTS " +
			"(SELECT 1 " +
			"FROM unnest(pet_type_list) AS pet_type " +
			"WHERE pet_type <> 'cat')"
	case "all":
		filterString = ""
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
			v_pets_for_sos_posts.pets_info,
			v_media_for_sos_posts.media_info,
			v_conditions.conditions_info
		FROM
			v_sos_posts
				LEFT JOIN v_pets_for_sos_posts ON v_sos_posts.id = v_pets_for_sos_posts.sos_post_id
				LEFT JOIN v_media_for_sos_posts ON v_sos_posts.id = v_media_for_sos_posts.sos_post_id
				LEFT JOIN v_conditions ON v_sos_posts.id = v_conditions.sos_post_id
		WHERE
			v_sos_posts.earliest_date_start_at >= '%s'
			AND v_sos_posts.author_id = $1
			%s
		ORDER BY
			%s
		LIMIT $2
		OFFSET $3;

	`,
		utils.FormatDate(time.Now().String()),
		filterString,
		sortString,
	)

	rows, err := tx.QueryContext(ctx, query, authorID, size+1, (page-1)*size)
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}
	defer rows.Close()

	sosPostList, err2 := readSOSPostRows(rows, page, size)
	if err2 != nil {
		return nil, err2
	}

	return sosPostList, nil
}

func FindSOSPostByID(ctx context.Context, tx *database.Tx, id int) (*sospost.SOSPostInfo, *pnd.AppError) {
	query := `
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
				v_pets_for_sos_posts.pets_info,
				v_media_for_sos_posts.media_info,
				v_conditions.conditions_info
			FROM
				v_sos_posts
					LEFT JOIN v_pets_for_sos_posts ON v_sos_posts.id = v_pets_for_sos_posts.sos_post_id
					LEFT JOIN v_media_for_sos_posts ON v_sos_posts.id = v_media_for_sos_posts.sos_post_id
					LEFT JOIN v_conditions ON v_sos_posts.id = v_conditions.sos_post_id
			WHERE
				v_sos_posts.id = $1;
		`

	row := tx.QueryRowContext(ctx, query, id)

	sosPost := sospost.SOSPostInfo{}

	var datesData, petsData, mediaData, conditionsData []byte
	if err := row.Scan(
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

	if err := json.Unmarshal(datesData, &sosPost.Dates); err != nil {
		return nil, pnd.FromPostgresError(err)
	}
	if err := json.Unmarshal(petsData, &sosPost.Pets); err != nil {
		return nil, pnd.FromPostgresError(err)
	}
	if err := json.Unmarshal(mediaData, &sosPost.Media); err != nil {
		return nil, pnd.FromPostgresError(err)
	}
	if err := json.Unmarshal(conditionsData, &sosPost.Conditions); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return &sosPost, nil
}

func UpdateSOSPost(
	ctx context.Context, tx *database.Tx, request *sospost.UpdateSOSPostRequest,
) (*sospost.SOSPost, *pnd.AppError) {
	sosPost := &sospost.SOSPost{}

	if err := updateSOSPostsDates(ctx, tx, request.ID, request.Dates); err != nil {
		return nil, err
	}

	if err := updateSOSPostsMedia(ctx, tx, request.ID, request.ImageIDs); err != nil {
		return nil, err
	}

	if err := updateSOSPostsConditions(ctx, tx, request.ID, request.ConditionIDs); err != nil {
		return nil, err
	}

	if err := updateSOSPostsPets(ctx, tx, request.ID, request.PetIDs); err != nil {
		return nil, err
	}

	const query = `
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

	if err := tx.QueryRowContext(ctx, query,
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

func updateSOSPostsDates(ctx context.Context, tx *database.Tx, postID int, dates []sospost.SOSDateView) *pnd.AppError {
	if _, err := tx.ExecContext(ctx, `
		UPDATE
			sos_posts_dates
		SET
			deleted_at = NOW()
		WHERE
			sos_post_id = $1
	`, postID); err != nil {
		return pnd.FromPostgresError(err)
	}
	const query = `
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

	sosDates := []sospost.SOSDates{}
	for _, date := range dates {
		sosDate := sospost.SOSDates{}
		if err := tx.QueryRowContext(ctx, query,
			date.DateStartAt,
			date.DateEndAt,
		).Scan(
			&sosDate.ID,
			&sosDate.DateStartAt,
			&sosDate.DateEndAt,
			&sosDate.CreatedAt,
			&sosDate.UpdatedAt,
		); err != nil {
			return pnd.FromPostgresError(err)
		}
		sosDates = append(sosDates, sosDate)
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
			postID,
			sosDate.ID,
		); err != nil {
			return pnd.FromPostgresError(err)
		}
	}

	return nil
}

func updateSOSPostsMedia(ctx context.Context, tx *database.Tx, postID int, mediaIDs []int) *pnd.AppError {
	if _, err := tx.ExecContext(ctx, `
		UPDATE
			resource_media
		SET
			deleted_at = NOW()
		WHERE
			resource_id = $1
	`, postID,
	); err != nil {
		return pnd.FromPostgresError(err)
	}

	for _, mediaID := range mediaIDs {
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
			mediaID,
			postID,
			media.SOSResourceType,
		); err != nil {
			return pnd.FromPostgresError(err)
		}
	}

	return nil
}

func updateSOSPostsConditions(ctx context.Context, tx *database.Tx, postID int, conditionIDs []int) *pnd.AppError {
	if _, err := tx.ExecContext(ctx, `
		UPDATE
			sos_posts_conditions
		SET
			deleted_at = NOW()
		WHERE
			sos_post_id = $1
	`, postID,
	); err != nil {
		return pnd.FromPostgresError(err)
	}

	for _, conditionID := range conditionIDs {
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
			postID,
			conditionID,
		); err != nil {
			return pnd.FromPostgresError(err)
		}
	}

	return nil
}

func updateSOSPostsPets(ctx context.Context, tx *database.Tx, postID int, petIDs []int) *pnd.AppError {
	if _, err := tx.ExecContext(ctx, `
		UPDATE
			sos_posts_pets
		SET
			deleted_at = NOW()
		WHERE
			sos_post_id = $1
	`, postID,
	); err != nil {
		return pnd.FromPostgresError(err)
	}

	for _, petID := range petIDs {
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
			postID,
			petID,
		); err != nil {
			return pnd.FromPostgresError(err)
		}
	}

	return nil
}

func FindConditionByID(ctx context.Context, tx *database.Tx, id int) (*sospost.ConditionList, *pnd.AppError) {
	const query = `
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

	conditions := sospost.ConditionList{}
	rows, err := tx.QueryContext(ctx, query, id)
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}
	defer rows.Close()

	for rows.Next() {
		condition := sospost.Condition{}
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

func FindDatesBySOSPostID(ctx context.Context, tx *database.Tx, sosPostID int) (*sospost.SOSDatesList, *pnd.AppError) {
	const query = `
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

	var sosDates sospost.SOSDatesList
	rows, err := tx.QueryContext(ctx, query, sosPostID)
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}
	defer rows.Close()

	for rows.Next() {
		sosDate := sospost.SOSDates{}
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
