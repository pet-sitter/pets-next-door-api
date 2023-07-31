package models

type User struct {
	ID                   int    `field:"id"`
	Email                string `field:"email"`
	Password             string `field:"password"`
	Nickname             string `field:"nickname"`
	Fullname             string `field:"fullname"`
	FirebaseProviderType string `field:"fb_provider_type"`
	FirebaseUID          string `field:"fb_uid"`
	CreatedAt            string `field:"created_at"`
	UpdatedAt            string `field:"updated_at"`
	DeletedAt            string `field:"deleted_at"`
}
