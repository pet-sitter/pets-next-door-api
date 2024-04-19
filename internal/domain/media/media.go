package media

import (
	"context"

	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
)

type MediaType string

const (
	IMAGE_MEDIA_TYPE MediaType = "image"
)

type Media struct {
	ID        int       `field:"id" json:"id"`
	MediaType MediaType `field:"media_type" json:"media_type"`
	URL       string    `field:"url" json:"url"`
	CreatedAt string    `field:"created_at" json:"created_at"`
	UpdatedAt string    `field:"updated_at" json:"updated_at"`
	DeletedAt string    `field:"deleted_at" json:"deleted_at"`
}

type MediaList []*Media

type MediaStore interface {
	CreateMedia(ctx context.Context, tx *database.Tx, media *Media) (*Media, *pnd.AppError)
	FindMediaByID(ctx context.Context, tx *database.Tx, id int) (*Media, *pnd.AppError)
}
