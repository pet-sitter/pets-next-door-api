package user

import (
	"context"
	"database/sql"
	"time"

	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"

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
	CreatedAt            time.Time            `field:"created_at"`
	UpdatedAt            time.Time            `field:"updated_at"`
	DeletedAt            sql.NullTime         `field:"deleted_at"`
}

func (u *User) ToUserWithProfileImage(profileImageURL *string) *UserWithProfileImage {
	return &UserWithProfileImage{
		ID:                   u.ID,
		Email:                u.Email,
		Password:             u.Password,
		Nickname:             u.Nickname,
		Fullname:             u.Fullname,
		ProfileImageURL:      profileImageURL,
		FirebaseProviderType: u.FirebaseProviderType,
		FirebaseUID:          u.FirebaseUID,
		CreatedAt:            u.CreatedAt,
		UpdatedAt:            u.UpdatedAt,
		DeletedAt:            u.DeletedAt,
	}
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
	CreatedAt            time.Time            `field:"created_at"`
	UpdatedAt            time.Time            `field:"updated_at"`
	DeletedAt            sql.NullTime         `field:"deleted_at"`
}

func (u *UserWithProfileImage) ToUserWithoutPrivateInfo() *UserWithoutPrivateInfo {
	if u.DeletedAt.Valid {
		return &UserWithoutPrivateInfo{
			ID:              u.ID,
			Nickname:        "탈퇴한 사용자",
			ProfileImageURL: nil,
		}
	}

	return &UserWithoutPrivateInfo{
		ID:              u.ID,
		Nickname:        u.Nickname,
		ProfileImageURL: u.ProfileImageURL,
	}
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

type UserStore interface {
	CreateUser(ctx context.Context, tx *database.Tx, request *RegisterUserRequest) (*User, *pnd.AppError)
	FindUsers(ctx context.Context, tx *database.Tx, page, size int, nickname *string) (*UserWithoutPrivateInfoList, *pnd.AppError)
	FindUserByID(ctx context.Context, tx *database.Tx, id int, includeDeleted bool) (*UserWithProfileImage, *pnd.AppError)
	FindUserByEmail(ctx context.Context, tx *database.Tx, email string) (*UserWithProfileImage, *pnd.AppError)
	FindUserByUID(ctx context.Context, tx *database.Tx, uid string) (*UserWithProfileImage, *pnd.AppError)
	FindUserIDByFbUID(ctx context.Context, tx *database.Tx, fbUID string) (int, *pnd.AppError)
	ExistsUserByNickname(ctx context.Context, tx *database.Tx, nickname string) (bool, *pnd.AppError)
	FindUserStatusByEmail(ctx context.Context, tx *database.Tx, email string) (*UserStatus, *pnd.AppError)
	UpdateUserByUID(ctx context.Context, tx *database.Tx, uid, nickname string, profileImageID *int) (*User, *pnd.AppError)
	DeleteUserByUID(ctx context.Context, tx *database.Tx, uid string) *pnd.AppError
}
