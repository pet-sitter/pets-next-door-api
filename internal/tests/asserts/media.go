package asserts

import (
	"testing"

	"github.com/pet-sitter/pets-next-door-api/internal/domain/media"
)

func MediaEquals(t *testing.T, got, want media.ListView) {
	t.Helper()

	for i, mediaData := range want {
		if got[i].ID != mediaData.ID {
			t.Errorf("got %v want %v", got[i].ID, mediaData.ID)
		}
		if got[i].MediaType != mediaData.MediaType {
			t.Errorf("got %v want %v", got[i].MediaType, mediaData.MediaType)
		}
		if got[i].URL != mediaData.URL {
			t.Errorf("got %v want %v", got[i].URL, mediaData.URL)
		}
	}
}
