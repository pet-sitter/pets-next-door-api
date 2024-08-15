package media

import (
	"time"

	"github.com/google/uuid"

	databasegen "github.com/pet-sitter/pets-next-door-api/internal/infra/database/gen"
)

type DetailView struct {
	ID        uuid.UUID `json:"id"`
	MediaType Type      `json:"mediaType"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"createdAt"`
}

type ListView []*DetailView

func ToDetailView(media databasegen.FindSingleMediaRow) *DetailView {
	return &DetailView{
		ID:        media.ID,
		MediaType: Type(media.MediaType),
		URL:       media.Url,
		CreatedAt: media.CreatedAt,
	}
}

func ToDetailViewFromCreated(media databasegen.CreateMediaRow) *DetailView {
	return &DetailView{
		ID:        media.ID,
		MediaType: Type(media.MediaType),
		URL:       media.Url,
		CreatedAt: media.CreatedAt,
	}
}

func ToDetailViewFromFindByIDs(media databasegen.FindMediasByIDsRow) *DetailView {
	return &DetailView{
		ID:        int64(media.ID),
		MediaType: Type(media.MediaType),
		URL:       media.Url,
		CreatedAt: media.CreatedAt,
	}
}

func ToDetailViewFromResourceMediaRows(resourceMedia databasegen.FindResourceMediaRow) *DetailView {
	return &DetailView{
		ID:        resourceMedia.MediaID,
		MediaType: Type(resourceMedia.MediaType),
		URL:       resourceMedia.Url,
		CreatedAt: resourceMedia.CreatedAt,
	}
}

func ToDetailViewFromViewForSOSPost(media ViewForSOSPost) *DetailView {
	createdAt, err := time.Parse(time.RFC3339, media.CreatedAt)
	if err != nil {
		createdAt = time.Time{}
	}

	return &DetailView{
		ID:        media.ID,
		MediaType: media.MediaType,
		URL:       media.URL,
		CreatedAt: createdAt,
	}
}

func ToListViewFromResourceMediaRows(resourceMediaList []databasegen.FindResourceMediaRow) ListView {
	mediaViewList := make(ListView, len(resourceMediaList))
	for i, resourceMedia := range resourceMediaList {
		mediaViewList[i] = ToDetailViewFromResourceMediaRows(resourceMedia)
	}
	return mediaViewList
}

func ToListViewFromViewListForSOSPost(mediaList ViewListForSOSPost) ListView {
	mediaViewList := make(ListView, len(mediaList))
	for i, media := range mediaList {
		mediaViewList[i] = ToDetailViewFromViewForSOSPost(
			ViewForSOSPost{
				ID:        media.ID,
				MediaType: media.MediaType,
				URL:       media.URL,
				CreatedAt: media.CreatedAt,
			},
		)
	}
	return mediaViewList
}
