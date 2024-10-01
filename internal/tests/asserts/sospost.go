package asserts

import (
	"reflect"
	"testing"

	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"

	"github.com/pet-sitter/pets-next-door-api/internal/domain/soscondition"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/sospost"
)

func DatesEquals(t *testing.T, want, got []sospost.SOSDateView) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func ConditionIDEquals(t *testing.T, want []uuid.UUID, got soscondition.ListView) {
	t.Helper()

	assert.Equal(t, len(want), len(got))
	for i := range got {
		assert.Equal(t, want[i], got[i].ID)
	}
}
