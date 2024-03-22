package media

type MediaView struct {
	ID        int       `json:"id"`
	MediaType MediaType `json:"mediaType"`
	URL       string    `json:"url"`
	CreatedAt string    `json:"createdAt"`
}

type MediaViewList []*MediaView

func (media *Media) ToMediaView() *MediaView {
	return &MediaView{
		ID:        media.ID,
		MediaType: media.MediaType,
		URL:       media.URL,
		CreatedAt: media.CreatedAt,
	}
}

func (mediaList *MediaList) ToMediaViewList() MediaViewList {
	mediaViewList := make(MediaViewList, len(*mediaList))
	for i, media := range *mediaList {
		mediaViewList[i] = media.ToMediaView()
	}
	return mediaViewList
}

type ResourceMediaView struct {
	ID           int          `field:"id"`
	ResourceType ResourceType `field:"resource_type"`
	ResourceID   int          `field:"resource_id"`
	MediaID      int          `field:"media_id"`
}
