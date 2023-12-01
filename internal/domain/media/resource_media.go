package media

import "time"

type ResourceType string

const (
	SosResourceType ResourceType = "sos_posts"
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

type ResourceMediaView struct {
	ID           int          `field:"id"`
	ResourceType ResourceType `field:"resource_type"`
	ResourceID   int          `field:"resource_id"`
	MediaID      int          `field:"media_id"`
}

type ResourceMediaStore interface {
	CreateResourceMedia(resourceID int, mediaID int, resourceType string) (*ResourceMedia, error)
	FindResourceMediaByResourceID(resourceID int, resourceType string) ([]Media, error)
}
