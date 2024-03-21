package media

type MediaView struct {
	ID        int       `json:"id"`
	MediaType MediaType `json:"mediaType"`
	URL       string    `json:"url"`
	CreatedAt string    `json:"createdAt"`
}

func NewMediaView(media *Media) MediaView {
	return MediaView{
		ID:        media.ID,
		MediaType: media.MediaType,
		URL:       media.URL,
		CreatedAt: media.CreatedAt,
	}
}

func NewMediaListView(mediaList []Media) []MediaView {
	mediaListView := make([]MediaView, 0)
	for _, media := range mediaList {
		mediaListView = append(mediaListView, NewMediaView(&media))
	}
	return mediaListView
}

type ResourceMediaView struct {
	ID           int          `field:"id"`
	ResourceType ResourceType `field:"resource_type"`
	ResourceID   int          `field:"resource_id"`
	MediaID      int          `field:"media_id"`
}
