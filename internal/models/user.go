package models

type UserModel struct {
	ID        int    `field:"id"`
	Email     string `field:"email"`
	Password  string `field:"password"`
	UID       string `field:"uid"`
	Providers string `field:"providers"`
	Nickname  string `field:"nickname"`
	CreatedAt string `field:"created_at"`
	UpdatedAt string `field:"updated_at"`
	DeletedAt string `field:"deleted_at"`
}
