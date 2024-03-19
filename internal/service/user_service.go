package service

import (
	"context"
	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/media"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/user"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
	s3infra "github.com/pet-sitter/pets-next-door-api/internal/infra/s3"
	"github.com/pet-sitter/pets-next-door-api/internal/postgres"
)

type UserService struct {
	conn     *database.DB
	s3Client *s3infra.S3Client
}

func NewUserService(conn *database.DB, s3Client *s3infra.S3Client) *UserService {
	return &UserService{
		conn:     conn,
		s3Client: s3Client,
	}
}

func (service *UserService) RegisterUser(ctx context.Context, registerUserRequest *user.RegisterUserRequest) (*user.RegisterUserView, *pnd.AppError) {
	var userView *user.RegisterUserView
	var err *pnd.AppError

	err = database.WithTransaction(ctx, service.conn, func(tx *database.Tx) *pnd.AppError {
		userStore := postgres.NewUserPostgresStore(tx)
		mediaService := NewMediaService(service.conn, service.s3Client)

		var mediaData *media.Media
		if registerUserRequest.ProfileImageID != nil {
			mediaData, err = mediaService.FindMediaByID(ctx, *registerUserRequest.ProfileImageID)
			if err != nil {
				return err
			}
		}

		var profileImageURL *string
		if mediaData != nil {
			profileImageURL = &mediaData.URL
		}

		created, err := userStore.CreateUser(ctx, registerUserRequest)
		if err != nil {
			return err
		}

		userView = &user.RegisterUserView{
			ID:                   created.ID,
			Email:                created.Email,
			Nickname:             created.Nickname,
			Fullname:             created.Fullname,
			ProfileImageURL:      profileImageURL,
			FirebaseProviderType: created.FirebaseProviderType,
			FirebaseUID:          created.FirebaseUID,
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return userView, nil
}

func (service *UserService) FindUsers(ctx context.Context, page int, size int, nickname *string) (*user.UserWithoutPrivateInfoList, *pnd.AppError) {
	var userList *user.UserWithoutPrivateInfoList
	var err *pnd.AppError

	err = database.WithTransaction(ctx, service.conn, func(tx *database.Tx) *pnd.AppError {
		userStore := postgres.NewUserPostgresStore(tx)

		userList, err = userStore.FindUsers(ctx, page, size, nickname)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return userList, nil
}

// FindMyProfile은 사용자의 프로필 정보를 조회한다.
// 삭제된 유저의 경우 삭제된 유저 정보를 반환한다.
func (service *UserService) FindPublicUserByID(ctx context.Context, id int) (*user.UserWithoutPrivateInfo, *pnd.AppError) {
	var err *pnd.AppError

	var user *user.UserWithProfileImage
	err = database.WithTransaction(ctx, service.conn, func(tx *database.Tx) *pnd.AppError {
		userStore := postgres.NewUserPostgresStore(tx)

		user, err = userStore.FindUserByID(ctx, id, true)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return user.ToUserWithoutPrivateInfo(), nil
}

func (service *UserService) FindUserByEmail(ctx context.Context, email string) (*user.UserWithProfileImage, *pnd.AppError) {
	var user *user.UserWithProfileImage
	var err *pnd.AppError

	err = database.WithTransaction(ctx, service.conn, func(tx *database.Tx) *pnd.AppError {
		userStore := postgres.NewUserPostgresStore(tx)

		user, err = userStore.FindUserByEmail(ctx, email)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (service *UserService) FindUserByUID(ctx context.Context, uid string) (*user.FindUserView, *pnd.AppError) {
	var userView *user.FindUserView

	err := database.WithTransaction(ctx, service.conn, func(tx *database.Tx) *pnd.AppError {
		userStore := postgres.NewUserPostgresStore(tx)

		foundUser, err := userStore.FindUserByUID(ctx, uid)
		if err != nil {
			return err
		}

		userView = &user.FindUserView{
			ID:                   foundUser.ID,
			Email:                foundUser.Email,
			Nickname:             foundUser.Nickname,
			Fullname:             foundUser.Fullname,
			ProfileImageURL:      foundUser.ProfileImageURL,
			FirebaseProviderType: foundUser.FirebaseProviderType,
			FirebaseUID:          foundUser.FirebaseUID,
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return userView, nil
}

func (service *UserService) ExistsByNickname(ctx context.Context, nickname string) (bool, *pnd.AppError) {
	var existsByNickname bool
	var err *pnd.AppError

	err = database.WithTransaction(ctx, service.conn, func(tx *database.Tx) *pnd.AppError {
		userStore := postgres.NewUserPostgresStore(tx)
		existsByNickname, err = userStore.ExistsByNickname(ctx, nickname)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return existsByNickname, err
	}

	return existsByNickname, nil
}

func (service *UserService) FindUserStatusByEmail(ctx context.Context, email string) (*user.UserStatus, *pnd.AppError) {
	var userStatus *user.UserStatus
	var err *pnd.AppError

	err = database.WithTransaction(ctx, service.conn, func(tx *database.Tx) *pnd.AppError {
		userStore := postgres.NewUserPostgresStore(tx)

		userStatus, err = userStore.FindUserStatusByEmail(ctx, email)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return userStatus, err
	}

	return userStatus, nil
}

func (service *UserService) UpdateUserByUID(ctx context.Context, uid string, nickname string, profileImageID *int) (*user.UserWithProfileImage, *pnd.AppError) {
	var userView *user.UserWithProfileImage
	var profileImage *media.Media

	err := database.WithTransaction(ctx, service.conn, func(tx *database.Tx) *pnd.AppError {
		userStore := postgres.NewUserPostgresStore(tx)
		mediaService := NewMediaService(service.conn, service.s3Client)

		updatedUser, err := userStore.UpdateUserByUID(ctx, uid, nickname, profileImageID)
		if err != nil {
			return err
		}

		if updatedUser.ProfileImageID != nil {
			profileImage, err = mediaService.FindMediaByID(ctx, *updatedUser.ProfileImageID)
			if err != nil {
				return err
			}
		}

		userView = &user.UserWithProfileImage{
			ID:                   updatedUser.ID,
			Email:                updatedUser.Email,
			Nickname:             updatedUser.Nickname,
			Fullname:             updatedUser.Fullname,
			ProfileImageURL:      &profileImage.URL,
			FirebaseProviderType: updatedUser.FirebaseProviderType,
			FirebaseUID:          updatedUser.FirebaseUID,
		}

		return nil
	})
	if err != nil {
		return userView, err
	}

	return userView, nil
}

func (service *UserService) DeleteUserByUID(ctx context.Context, uid string) *pnd.AppError {
	err := database.WithTransaction(ctx, service.conn, func(tx *database.Tx) *pnd.AppError {
		userStore := postgres.NewUserPostgresStore(tx)

		if err := userStore.DeleteUserByUID(ctx, uid); err != nil {
			return err
		}

		return nil
	})

	return err
}

func (service *UserService) AddPetsToOwner(ctx context.Context, uid string, addPetsRequest pet.AddPetsToOwnerRequest) ([]pet.PetView, *pnd.AppError) {
	var petViews []pet.PetView

	err := database.WithTransaction(ctx, service.conn, func(tx *database.Tx) *pnd.AppError {
		userStore := postgres.NewUserPostgresStore(tx)
		petStore := postgres.NewPetPostgresStore(tx)
		mediaStore := postgres.NewMediaPostgresStore(tx)

		user, err := userStore.FindUserByUID(ctx, uid)
		if err != nil {
			return err
		}

		pets := make([]pet.PetWithProfileImage, len(addPetsRequest.Pets))
		for i, item := range addPetsRequest.Pets {
			media, err := mediaStore.FindMediaByID(ctx, *item.ProfileImageID)
			if err != nil {
				return err
			}

			petToCreate := pet.Pet{
				BasePet: pet.BasePet{
					OwnerID:    user.ID,
					Name:       item.Name,
					PetType:    item.PetType,
					Sex:        item.Sex,
					Neutered:   item.Neutered,
					Breed:      item.Breed,
					BirthDate:  item.BirthDate,
					WeightInKg: item.WeightInKg,
				},
				ProfileImageID: &media.ID,
			}
			createdPet, err := petStore.CreatePet(ctx, &petToCreate)
			pets[i] = *createdPet

			if err != nil {
				return err
			}
		}

		petViews = pet.NewPetViewList(pets)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return petViews, nil
}

func (service *UserService) FindPetsByOwnerUID(ctx context.Context, uid string) (*pet.FindMyPetsView, *pnd.AppError) {
	var petListView *pet.FindMyPetsView

	err := database.WithTransaction(ctx, service.conn, func(tx *database.Tx) *pnd.AppError {
		userStore := postgres.NewUserPostgresStore(tx)
		petStore := postgres.NewPetPostgresStore(tx)

		user, err := userStore.FindUserByUID(ctx, uid)
		if err != nil {
			return err
		}

		pets, err := petStore.FindPetsByOwnerID(ctx, user.ID)
		if err != nil {
			return err
		}

		petListView = pet.NewFindMyPetsView(pets)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return petListView, nil
}
