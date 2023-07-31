package models

type MediaType string

const (
	IMAGE_MEDIA_TYPE MediaType = "image"
)

type Media struct {
	ID        int       `field:"id"`
	MediaType MediaType `field:"media_type"`
	URL       string    `field:"url"`
	CreatedAt string    `field:"created_at"`
	UpdatedAt string    `field:"updated_at"`
	DeletedAt string    `field:"deleted_at"`
}
