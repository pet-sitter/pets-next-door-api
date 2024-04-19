package media

import (
	"context"

	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
)

type MediaType string

const (
	MediaTypeImage MediaType = "image"
)

type Media struct {
	ID        int       `field:"id" json:"id"`
	MediaType MediaType `field:"media_type" json:"mediaType"`
	URL       string    `field:"url" json:"url"`
	CreatedAt string    `field:"created_at" json:"createdAt"`
	UpdatedAt string    `field:"updated_at" json:"updatedAt"`
	DeletedAt string    `field:"deleted_at" json:"deletedAt"`
}

type MediaList []*Media

type MediaStore interface {
	CreateMedia(ctx context.Context, tx *database.Tx, media *Media) (*Media, *pnd.AppError)
	FindMediaByID(ctx context.Context, tx *database.Tx, id int) (*Media, *pnd.AppError)
}
