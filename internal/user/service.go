package user

func NewUserService(userRepo UserRepo) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

type UserServicer interface {
	CreateUser(user *UserModel) (*UserModel, error)
	FindUserByEmail(email string) (*UserModel, error)
	FindUserByUID(uid string) (*UserModel, error)
	UpdateUserByUID(uid string, nickname string) (*UserModel, error)
}

func (service *UserService) CreateUser(user *UserModel) (*UserModel, error) {
	created, err := service.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (service *UserService) FindUserByEmail(email string) (*UserModel, error) {
	user, err := service.userRepo.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (service *UserService) FindUserByUID(uid string) (*UserModel, error) {
	user, err := service.userRepo.FindUserByUID(uid)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (service *UserService) UpdateUserByUID(uid string, nickname string) (*UserModel, error) {
	updated, err := service.userRepo.UpdateUserByUID(uid, nickname)
	if err != nil {
		return nil, err
	}

	return updated, nil
}
