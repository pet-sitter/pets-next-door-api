package commonvo

// Pet
type PetType string

const (
	PetTypeDog PetType = "dog"
	PetTypeCat PetType = "cat"
)

func (p *PetType) String() string {
	return string(*p)
}
