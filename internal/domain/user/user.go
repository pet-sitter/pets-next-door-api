package user

import (
	"context"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
)

type FirebaseProviderType string

const (
	FirebaseProviderTypeEmail  FirebaseProviderType = "email"
	FirebaseProviderTypeGoogle FirebaseProviderType = "google"
	FirebaseProviderTypeApple  FirebaseProviderType = "apple"
	FirebaseProviderTypeKakao  FirebaseProviderType = "kakao"
)

type User struct {
	ID                   int                  `field:"id"`
	Email                string               `field:"email"`
	Password             string               `field:"password"`
	Nickname             string               `field:"nickname"`
	Fullname             string               `field:"fullname"`
	ProfileImageID       *int                 `field:"profile_image_id"`
	FirebaseProviderType FirebaseProviderType `field:"fb_provider_type"`
	FirebaseUID          string               `field:"fb_uid"`
	CreatedAt            string               `field:"created_at"`
	UpdatedAt            string               `field:"updated_at"`
	DeletedAt            string               `field:"deleted_at"`
}

type UserWithProfileImage struct {
	ID                   int                  `field:"id"`
	Email                string               `field:"email"`
	Password             string               `field:"password"`
	Nickname             string               `field:"nickname"`
	Fullname             string               `field:"fullname"`
	ProfileImageURL      *string              `field:"profile_image_url"`
	FirebaseProviderType FirebaseProviderType `field:"fb_provider_type"`
	FirebaseUID          string               `field:"fb_uid"`
	CreatedAt            string               `field:"created_at"`
	UpdatedAt            string               `field:"updated_at"`
	DeletedAt            string               `field:"deleted_at"`
}

type UserWithoutPrivateInfo struct {
	ID              int     `field:"id" json:"id"`
	Nickname        string  `field:"nickname" json:"nickname"`
	ProfileImageURL *string `field:"profile_image_url" json:"profileImageUrl"`
}

type UserWithoutPrivateInfoList struct {
	*pnd.PaginatedView[UserWithoutPrivateInfo]
}

func NewUserWithoutPrivateInfoList(page int, size int) *UserWithoutPrivateInfoList {
	return &UserWithoutPrivateInfoList{PaginatedView: pnd.NewPaginatedView(
		page, size, false, make([]UserWithoutPrivateInfo, 0),
	)}
}

type UserStatus struct {
	FirebaseProviderType FirebaseProviderType `field:"fb_provider_type"`
}

type UserStore interface {
	CreateUser(ctx context.Context, request *RegisterUserRequest) (*User, *pnd.AppError)
	HardDeleteUserByUID(ctx context.Context, uid string) *pnd.AppError
	FindUsers(ctx context.Context, page int, size int, nickname *string) (*UserWithoutPrivateInfoList, *pnd.AppError)
	FindUserByEmail(ctx context.Context, email string) (*UserWithProfileImage, *pnd.AppError)
	FindUserByUID(ctx context.Context, uid string) (*UserWithProfileImage, *pnd.AppError)
	FindUserIDByFbUID(ctx context.Context, fbUid string) (int, *pnd.AppError)
	ExistsByNickname(ctx context.Context, nickname string) (bool, *pnd.AppError)
	FindUserStatusByEmail(ctx context.Context, email string) (*UserStatus, *pnd.AppError)
	UpdateUserByUID(ctx context.Context, uid string, nickname string, profileImageID *int) (*User, *pnd.AppError)
}
