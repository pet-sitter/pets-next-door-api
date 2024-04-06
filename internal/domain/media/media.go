package media

import (
	"context"
	"database/sql"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
	"time"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
)

type MediaType string

const (
	IMAGE_MEDIA_TYPE MediaType = "image"
)

type Media struct {
	ID        int          `field:"id"`
	MediaType MediaType    `field:"media_type"`
	URL       string       `field:"url"`
	CreatedAt time.Time    `field:"created_at"`
	UpdatedAt time.Time    `field:"updated_at"`
	DeletedAt sql.NullTime `field:"deleted_at"`
}

type MediaList []*Media

type MediaStore interface {
	CreateMedia(ctx context.Context, tx database.Tx, media *Media) (*Media, *pnd.AppError)
	FindMediaByID(ctx context.Context, tx database.Tx, id int) (*Media, *pnd.AppError)
}
