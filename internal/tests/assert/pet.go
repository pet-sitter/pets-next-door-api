package assert

import (
	"testing"

	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"
)

func PetRequestAndViewEquals(t *testing.T, expected pet.AddPetRequest, found pet.DetailView) {
	t.Helper()

	if expected.Name != found.Name {
		t.Errorf("got %v want %v", expected.Name, found.Name)
	}

	if expected.PetType != found.PetType {
		t.Errorf("got %v want %v", expected.PetType, found.PetType)
	}

	if expected.Sex != found.Sex {
		t.Errorf("got %v want %v", expected.Sex, found.PetType)
	}

	if expected.Neutered != found.Neutered {
		t.Errorf("got %v want %v", expected.Neutered, found.Neutered)
	}

	if expected.Breed != found.Breed {
		t.Errorf("got %v want %v", expected.Breed, found.Breed)
	}

	if expected.BirthDate != found.BirthDate {
		t.Errorf("got %v want %v", expected.BirthDate, found.BirthDate)
	}

	if expected.WeightInKg.String() != found.WeightInKg.String() {
		t.Errorf("got %v want %v", expected.WeightInKg, found.WeightInKg)
	}
}

func UpdatedPetEquals(t *testing.T, expected pet.UpdatePetRequest, found pet.DetailView) {
	t.Helper()

	if expected.Name != found.Name {
		t.Errorf("got %v want %v", expected.Name, found.Name)
	}

	if expected.Neutered != found.Neutered {
		t.Errorf("got %v want %v", expected.Neutered, found.Neutered)
	}

	if expected.Breed != found.Breed {
		t.Errorf("got %v want %v", expected.Breed, found.Breed)
	}

	if expected.BirthDate != found.BirthDate {
		t.Errorf("got %v want %v", expected.BirthDate, found.BirthDate)
	}

	if expected.WeightInKg.String() != found.WeightInKg.String() {
		t.Errorf("got %v want %v", expected.WeightInKg, found.WeightInKg)
	}
}
