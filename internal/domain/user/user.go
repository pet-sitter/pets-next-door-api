package user

import pnd "github.com/pet-sitter/pets-next-door-api/api"

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

type UserStatus struct {
	FirebaseProviderType FirebaseProviderType `field:"fb_provider_type"`
}

type UserStore interface {
	CreateUser(request *RegisterUserRequest) (*User, *pnd.AppError)
	FindUsers(page int, size int, nickname *string) ([]*UserWithoutPrivateInfo, *pnd.AppError)
	FindUserByEmail(email string) (*UserWithProfileImage, *pnd.AppError)
	FindUserByUID(uid string) (*UserWithProfileImage, *pnd.AppError)
	FindUserIDByFbUID(fbUid string) (int, *pnd.AppError)
	ExistsByNickname(nickname string) (bool, *pnd.AppError)
	FindUserStatusByEmail(email string) (*UserStatus, *pnd.AppError)
	UpdateUserByUID(uid string, nickname string, profileImageID *int) (*User, *pnd.AppError)
}
