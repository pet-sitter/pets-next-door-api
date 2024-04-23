package user

import (
	"database/sql"
	"time"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
)

type FirebaseProviderType string

const (
	FirebaseProviderTypeEmail  FirebaseProviderType = "email"
	FirebaseProviderTypeGoogle FirebaseProviderType = "google"
	FirebaseProviderTypeApple  FirebaseProviderType = "apple"
	FirebaseProviderTypeKakao  FirebaseProviderType = "kakao"
)

func (f FirebaseProviderType) String() string {
	return string(f)
}

func (f FirebaseProviderType) NullString() sql.NullString {
	return sql.NullString{String: string(f), Valid: true}
}

type UserWithProfileImage struct {
	ID                   int
	Email                string
	Password             string
	Nickname             string
	Fullname             string
	ProfileImageURL      *string
	FirebaseProviderType FirebaseProviderType
	FirebaseUID          string
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            sql.NullTime
}

type UserWithoutPrivateInfo struct {
	ID              int     `field:"id" json:"id"`
	Nickname        string  `field:"nickname" json:"nickname"`
	ProfileImageURL *string `field:"profile_image_url" json:"profileImageUrl"`
}

type UserWithoutPrivateInfoList struct {
	*pnd.PaginatedView[UserWithoutPrivateInfo]
}

func NewUserWithoutPrivateInfoList(page, size int) *UserWithoutPrivateInfoList {
	return &UserWithoutPrivateInfoList{PaginatedView: pnd.NewPaginatedView(
		page, size, false, make([]UserWithoutPrivateInfo, 0),
	)}
}

type UserStatus struct {
	FirebaseProviderType FirebaseProviderType `field:"fb_provider_type"`
}
