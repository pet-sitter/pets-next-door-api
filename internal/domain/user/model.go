package user

import (
	"database/sql"
	"time"

	utils "github.com/pet-sitter/pets-next-door-api/internal/common"
	databasegen "github.com/pet-sitter/pets-next-door-api/internal/infra/database/gen"
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

type WithProfileImage struct {
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

func ToWithProfileImage(row databasegen.FindUserRow) *WithProfileImage {
	return &WithProfileImage{
		ID:                   int(row.ID),
		Email:                row.Email,
		Nickname:             row.Nickname,
		Fullname:             row.Fullname,
		ProfileImageURL:      utils.NullStrToStrPtr(row.ProfileImageUrl),
		FirebaseProviderType: FirebaseProviderType(row.FbProviderType.String),
		FirebaseUID:          row.FbUid.String,
		CreatedAt:            row.CreatedAt,
		UpdatedAt:            row.UpdatedAt,
	}
}

func (u *WithProfileImage) ToInternalView() *InternalView {
	return &InternalView{
		ID:                   u.ID,
		Email:                u.Email,
		Nickname:             u.Nickname,
		Fullname:             u.Fullname,
		ProfileImageURL:      u.ProfileImageURL,
		FirebaseProviderType: u.FirebaseProviderType,
		FirebaseUID:          u.FirebaseUID,
	}
}

func (u *WithProfileImage) ToMyProfileView() *MyProfileView {
	return &MyProfileView{
		ID:                   u.ID,
		Email:                u.Email,
		Nickname:             u.Nickname,
		Fullname:             u.Fullname,
		ProfileImageURL:      u.ProfileImageURL,
		FirebaseProviderType: u.FirebaseProviderType,
	}
}
