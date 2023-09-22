package pet

type Pet struct {
	ID         int     `field:"id"`
	OwnerID    int     `field:"owner_id"`
	Name       string  `field:"name"`
	PetType    PetType `field:"pet_type"`
	Sex        PetSex  `field:"sex"`
	Neutered   bool    `field:"neutered"`
	Breed      string  `field:"breed"`
	BirthDate  string  `field:"birth_date"`
	WeightInKg float64 `field:"weight_in_kg"`
	CreatedAt  string  `field:"created_at"`
	UpdatedAt  string  `field:"updated_at"`
	DeletedAt  string  `field:"deleted_at"`
}

type PetStore interface {
	CreatePet(pet *Pet) (*Pet, error)
	FindPetsByOwnerID(ownerID int) ([]Pet, error)
}

type PetType string

const (
	PetTypeDog PetType = "dog"
	PetTypeCat PetType = "cat"
)

type PetSex string

const (
	PetSexMale   PetSex = "male"
	PetSexFemale PetSex = "female"
)
