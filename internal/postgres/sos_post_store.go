package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pet-sitter/pets-next-door-api/internal/domain/media"

	"github.com/pet-sitter/pets-next-door-api/internal/domain/soscondition"

	utils "github.com/pet-sitter/pets-next-door-api/internal/common"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/sospost"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
)

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
		if len(mediaData) > 0 {
			if err := json.Unmarshal(mediaData, &sosPost.Media); err != nil {
				return nil, pnd.FromPostgresError(err)
			}
		} else {
			sosPost.Media = media.ViewListForSOSPost{}
		}
		if len(conditionsData) > 0 {
			if err := json.Unmarshal(conditionsData, &sosPost.Conditions); err != nil {
				return nil, pnd.FromPostgresError(err)
			}
		} else {
			sosPost.Conditions = soscondition.ViewListForSOSPost{}
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

	filterString := buildFilterString(filterType)

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

	filterString := buildFilterString(filterType)

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
	if len(mediaData) > 0 {
		if err := json.Unmarshal(mediaData, &sosPost.Media); err != nil {
			return nil, pnd.FromPostgresError(err)
		}
	} else {
		sosPost.Media = media.ViewListForSOSPost{}
	}
	if len(conditionsData) > 0 {
		if err := json.Unmarshal(conditionsData, &sosPost.Conditions); err != nil {
			return nil, pnd.FromPostgresError(err)
		}
	} else {
		sosPost.Conditions = soscondition.ViewListForSOSPost{}
	}

	return &sosPost, nil
}

func buildFilterString(petType string) string {
	if petType == "all" {
		return ""
	}
	return fmt.Sprintf(`
AND NOT EXISTS 
	(SELECT 1 
	FROM unnest(pet_type_list) AS pet_type 
	WHERE pet_type <> '%s')
`, petType)
}
