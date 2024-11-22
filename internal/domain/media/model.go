package media

import "github.com/google/uuid"

type Type string

const (
	TypeImage Type = "image"
)

func (mt Type) String() string {
	return string(mt)
}

type ViewForSOSPost struct {
	ID        uuid.UUID `field:"id"         json:"id"`
	MediaType Type      `field:"media_type" json:"media_type"`
	URL       string    `field:"url"        json:"url"`
	CreatedAt string    `field:"created_at" json:"created_at"`
	UpdatedAt string    `field:"updated_at" json:"updated_at"`
	DeletedAt string    `field:"deleted_at" json:"deleted_at"`
}

type ViewListForSOSPost []*ViewForSOSPost
