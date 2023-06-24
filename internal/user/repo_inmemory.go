package user

import (
	"fmt"
)

func (repo *UserInMemoryRepo) CreateUser(user *UserModel) (*UserModel, error) {
	nextID := len(repo.Users) + 1
	user.ID = nextID

	repo.Users = append(repo.Users, *user)

	return user, nil
}

func (repo *UserInMemoryRepo) FindUserByEmail(email string) (*UserModel, error) {
	var user *UserModel

	for _, v := range repo.Users {
		if v.Email == email {
			user = &v
		}
	}

	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func (repo *UserInMemoryRepo) FindUserByUID(uid string) (*UserModel, error) {
	var user *UserModel

	for _, v := range repo.Users {
		if v.UID == uid {
			user = &v
		}
	}

	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func (repo *UserInMemoryRepo) UpdateUserByUID(uid string, nickname string) (*UserModel, error) {
	var user *UserModel
	for i, v := range repo.Users {
		if v.UID == uid {
			repo.Users[i].Nickname = nickname
			user = &repo.Users[i]
		}
	}

	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}
