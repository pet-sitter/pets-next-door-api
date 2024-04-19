package media

import (
	"context"
	"time"

	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
)

type ResourceType string

const (
	SOSResourceType ResourceType = "sos_posts"
)

type ResourceMedia struct {
	ID           int          `field:"id"`
	ResourceType ResourceType `field:"resource_type"`
	ResourceID   int          `field:"resource_id"`
	MediaID      int          `field:"media_id"`
	CreatedAt    time.Time    `field:"created_at"`
	UpdatedAt    time.Time    `field:"updated_at"`
	DeletedAt    time.Time    `field:"deleted_at"`
}

type ResourceMediaStore interface {
	CreateResourceMedia(
		ctx context.Context, tx *database.Tx, resourceID, mediaID int, resourceType string,
	) (*ResourceMedia, *pnd.AppError)
	FindResourceMediaByResourceID(
		ctx context.Context, tx *database.Tx, resourceID int, resourceType string,
	) (*MediaList, *pnd.AppError)
}
