package user

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

type UserRepo interface {
	CreateUser(user *UserModel) (*UserModel, error)
	FindUserByEmail(email string) (*UserModel, error)
	FindUserByUID(uid string) (*UserModel, error)
	UpdateUserByUID(uid string, nickname string) (*UserModel, error)
}

type UserService struct {
	userRepo UserRepo
}

type UserInMemoryRepo struct {
	Users []UserModel
}

func NewUserInMemoryRepo() *UserInMemoryRepo {
	users := []UserModel{}

	return &UserInMemoryRepo{
		Users: users,
	}
}
