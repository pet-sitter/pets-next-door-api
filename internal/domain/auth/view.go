package auth

import "github.com/pet-sitter/pets-next-door-api/internal/domain/user"

type KakaoCallbackView struct {
	AuthToken            string                    `json:"authToken"`
	FirebaseProviderType user.FirebaseProviderType `json:"fbProviderType"`
	FirebaseUID          string                    `json:"fbUid"`
	Email                string                    `json:"email"`
	PhotoURL             string                    `json:"photoURL"`
}
