package postgres

import (
	"github.com/pet-sitter/pets-next-door-api/internal/domain/media"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
)

type ResourceMediaPostgresStore struct {
	db *database.DB
}

func NewResourceMediaPostgresStore(db *database.DB) *ResourceMediaPostgresStore {
	return &ResourceMediaPostgresStore{
		db: db,
	}
}

func (s *ResourceMediaPostgresStore) CreateResourceMedia(resourceID int, mediaID int, resourceType string) (*media.ResourceMedia, error) {
	resourceMedia := &media.ResourceMedia{}

	tx, _ := s.db.Begin()
	err := tx.QueryRow(`
	INSERT INTO
		resource_media
		(
			resource_id,
			media_id,
			resource_type,
			created_at,
			updated_at
		)
	VALUES ($1, $2, $3, NOW(), NOW())
	RETURNING id, resource_id, media_id, created_at, updated_at
	`,
		resourceID,
		mediaID,
		resourceType,
	).Scan(&resourceMedia.ID, &resourceMedia.ResourceID, &resourceMedia.MediaID, &resourceMedia.CreatedAt, &resourceMedia.UpdatedAt)
	tx.Commit()

	if err != nil {
		return nil, err
	}

	return resourceMedia, nil
}

func (s *ResourceMediaPostgresStore) FindResourceMediaByResourceID(resourceID int, resourceType string) ([]media.Media, error) {
	var mediaList []media.Media

	tx, _ := s.db.Begin()
	rows, err := tx.Query(`
	SELECT
		m.id,
		m.media_type,
		m.url,
		m.created_at,
		m.updated_at
	FROM
		resource_media rm
	INNER JOIN
		media m
	ON
		rm.media_id = m.id
	WHERE
		rm.resource_id = $1 AND
		rm.resource_type = $2 AND
		rm.deleted_at IS NULL
	`,
		resourceID,
		resourceType,
	)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		mediaItem := media.Media{}

		err := rows.Scan(&mediaItem.ID, &mediaItem.MediaType, &mediaItem.URL, &mediaItem.CreatedAt, &mediaItem.UpdatedAt)
		if err != nil {
			return nil, err
		}

		mediaList = append(mediaList, mediaItem)
	}

	return mediaList, nil
}
