package database

import "github.com/pet-sitter/pets-next-door-api/internal/models"

func (tx *Tx) CreateMedia(media *models.Media) (*models.Media, error) {
	err := tx.QueryRow(`
	INSERT INTO
		media
		(
			media_type,
			url,
			created_at,
			updated_at
		)
	VALUES ($1, $2, NOW(), NOW())
	RETURNING id, created_at, updated_at
	`,
		media.MediaType,
		media.URL,
	).Scan(&media.ID, &media.CreatedAt, &media.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return media, nil
}

func (tx *Tx) FindMediaByID(id int) (*models.Media, error) {
	media := &models.Media{}

	err := tx.QueryRow(`
	SELECT
		id,
		media_type,
		url,
		created_at,
		updated_at
	FROM
		media
	WHERE
		id = $1 AND
		deleted_at IS NULL
	`,
		id,
	).Scan(
		&media.ID,
		&media.MediaType,
		&media.URL,
		&media.CreatedAt,
		&media.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return media, nil
}
