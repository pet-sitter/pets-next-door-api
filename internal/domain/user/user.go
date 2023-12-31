package user

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
	ProfileImageID       int                  `field:"profile_image_id"`
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
	ProfileImageURL      string               `field:"profile_image_url"`
	FirebaseProviderType FirebaseProviderType `field:"fb_provider_type"`
	FirebaseUID          string               `field:"fb_uid"`
	CreatedAt            string               `field:"created_at"`
	UpdatedAt            string               `field:"updated_at"`
	DeletedAt            string               `field:"deleted_at"`
}

type UserWithoutPrivateInfo struct {
	ID              int    `field:"id" json:"id"`
	Nickname        string `field:"nickname" json:"nickname"`
	ProfileImageURL string `field:"profile_image_url" json:"profileImageUrl"`
}

type UserStatus struct {
	FirebaseProviderType FirebaseProviderType `field:"fb_provider_type"`
}

type UserStore interface {
	CreateUser(request *RegisterUserRequest) (*User, error)
	FindUsers(page int, size int, nickname *string) ([]*UserWithoutPrivateInfo, error)
	FindUserByEmail(email string) (*UserWithProfileImage, error)
	FindUserByUID(uid string) (*UserWithProfileImage, error)
	FindUserIDByFbUID(fbUid string) (int, error)
	ExistsByNickname(nickname string) (bool, error)
	FindUserStatusByEmail(email string) (*UserStatus, error)
	UpdateUserByUID(uid string, nickname string, profileImageID int) (*User, error)
}
