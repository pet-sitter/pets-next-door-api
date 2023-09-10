package views

import "github.com/pet-sitter/pets-next-door-api/internal/models"

type RegisterUserRequest struct {
	Email                string                      `json:"email" validate:"required,email"`
	Nickname             string                      `json:"nickname" validate:"required"`
	Fullname             string                      `json:"fullname" validate:"required"`
	ProfileImageID       int                         `json:"profileImageId" validate:"required"`
	FirebaseProviderType models.FirebaseProviderType `json:"fbProviderType" validate:"required"`
	FirebaseUID          string                      `json:"fbUid" validate:"required"`
}

type RegisterUserResponse struct {
	ID                   int                         `json:"id"`
	Email                string                      `json:"email"`
	Nickname             string                      `json:"nickname"`
	Fullname             string                      `json:"fullname"`
	ProfileImageURL      string                      `json:"profileImageUrl"`
	FirebaseProviderType models.FirebaseProviderType `json:"fbProviderType"`
	FirebaseUID          string                      `json:"fbUid"`
}

type FindUserResponse struct {
	ID                   int                         `json:"id"`
	Email                string                      `json:"email"`
	Nickname             string                      `json:"nickname"`
	Fullname             string                      `json:"fullname"`
	ProfileImageURL      string                      `json:"profileImageUrl"`
	FirebaseProviderType models.FirebaseProviderType `json:"fbProviderType"`
	FirebaseUID          string                      `json:"fbUid"`
}

type UserStatusRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type UserRegistrationStatus string

const (
	UserStatusNotRegistered UserRegistrationStatus = "NOT_REGISTERED"
	UserStatusRegistered    UserRegistrationStatus = "REGISTERED"
)

type UserStatusView struct {
	Status               UserRegistrationStatus      `json:"status"`
	FirebaseProviderType models.FirebaseProviderType `json:"fbProviderType,omitempty"`
}

type UpdateUserRequest struct {
	Nickname       string `json:"nickname"`
	ProfileImageID int    `json:"profileImageId"`
}

type UpdateUserResponse struct {
	ID                   int                         `json:"id"`
	Email                string                      `json:"email"`
	Nickname             string                      `json:"nickname"`
	Fullname             string                      `json:"fullname"`
	ProfileImageURL      string                      `json:"profileImageUrl"`
	FirebaseProviderType models.FirebaseProviderType `json:"fbProviderType"`
	FirebaseUID          string                      `json:"fbUid"`
}
