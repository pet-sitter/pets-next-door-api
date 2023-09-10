package media

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

type MediaStore interface {
	CreateMedia(media *Media) (*Media, error)
	FindMediaByID(id int) (*Media, error)
}

type MediaView struct {
	ID        int       `json:"id"`
	MediaType MediaType `json:"mediaType"`
	URL       string    `json:"url"`
	CreatedAt string    `json:"createdAt"`
}
