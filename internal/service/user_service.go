package service

import (
	"context"
	"fmt"
	pnd "github.com/pet-sitter/pets-next-door-api/api"
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
	mediaService := NewMediaService(service.conn, service.s3Client)

	var profileImageURL *string
	if registerUserRequest.ProfileImageID != nil {
		mediaData, err := mediaService.FindMediaByID(ctx, *registerUserRequest.ProfileImageID)
		if err != nil {
			return nil, err
		}

		if mediaData != nil {
			profileImageURL = &mediaData.URL
		}
	}

	tx, err := service.conn.BeginTx(ctx)
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}

	created, err := postgres.CreateUser(ctx, tx, registerUserRequest)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &user.RegisterUserView{
		ID:                   created.ID,
		Email:                created.Email,
		Nickname:             created.Nickname,
		Fullname:             created.Fullname,
		ProfileImageURL:      profileImageURL,
		FirebaseProviderType: created.FirebaseProviderType,
		FirebaseUID:          created.FirebaseUID,
	}, nil
}

func (service *UserService) FindUsers(ctx context.Context, page int, size int, nickname *string) (*user.UserWithoutPrivateInfoList, *pnd.AppError) {
	tx, err := service.conn.BeginTx(ctx)
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}

	userList, err := postgres.FindUsers(ctx, tx, page, size, nickname)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return userList, nil
}

// FindMyProfile은 사용자의 프로필 정보를 조회한다.
// 삭제된 유저의 경우 삭제된 유저 정보를 반환한다.
func (service *UserService) FindPublicUserByID(ctx context.Context, id int) (*user.UserWithoutPrivateInfo, *pnd.AppError) {
	tx, err := service.conn.BeginTx(ctx)
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}

	user, err := postgres.FindUserByID(ctx, tx, id, true)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return user.ToUserWithoutPrivateInfo(), nil
}

func (service *UserService) FindUserByEmail(ctx context.Context, email string) (*user.UserWithProfileImage, *pnd.AppError) {
	tx, err := service.conn.BeginTx(ctx)
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}

	user, err := postgres.FindUserByEmail(ctx, tx, email)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return user, nil
}

func (service *UserService) FindUserByUID(ctx context.Context, uid string) (*user.FindUserView, *pnd.AppError) {
	tx, err := service.conn.BeginTx(ctx)
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}

	foundUser, err := postgres.FindUserByUID(ctx, tx, uid)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &user.FindUserView{
		ID:                   foundUser.ID,
		Email:                foundUser.Email,
		Nickname:             foundUser.Nickname,
		Fullname:             foundUser.Fullname,
		ProfileImageURL:      foundUser.ProfileImageURL,
		FirebaseProviderType: foundUser.FirebaseProviderType,
		FirebaseUID:          foundUser.FirebaseUID,
	}, nil
}

func (service *UserService) ExistsByNickname(ctx context.Context, nickname string) (bool, *pnd.AppError) {
	tx, err := service.conn.BeginTx(ctx)
	defer tx.Rollback()
	if err != nil {
		return false, err
	}

	existsByNickname, err := postgres.ExistsUserByNickname(ctx, tx, nickname)
	if err != nil {
		return existsByNickname, err
	}

	if err := tx.Commit(); err != nil {
		return existsByNickname, err
	}

	return existsByNickname, nil
}

func (service *UserService) FindUserStatusByEmail(ctx context.Context, email string) (*user.UserStatus, *pnd.AppError) {
	tx, err := service.conn.BeginTx(ctx)
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}

	userStatus, err := postgres.FindUserStatusByEmail(ctx, tx, email)
	if err != nil {
		return userStatus, err
	}

	if err := tx.Commit(); err != nil {
		return userStatus, err
	}

	return userStatus, nil
}

func (service *UserService) UpdateUserByUID(ctx context.Context, uid string, nickname string, profileImageID *int) (*user.UserWithProfileImage, *pnd.AppError) {
	mediaService := NewMediaService(service.conn, service.s3Client)

	tx, err := service.conn.BeginTx(ctx)
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}

	updatedUser, err := postgres.UpdateUserByUID(ctx, tx, uid, nickname, profileImageID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	var profileImageURL *string
	if updatedUser.ProfileImageID != nil {
		profileImage, err := mediaService.FindMediaByID(ctx, *updatedUser.ProfileImageID)
		if err != nil {
			return nil, err
		}

		if profileImage != nil {
			profileImageURL = &profileImage.URL
		}
	}

	return &user.UserWithProfileImage{
		ID:                   updatedUser.ID,
		Email:                updatedUser.Email,
		Nickname:             updatedUser.Nickname,
		Fullname:             updatedUser.Fullname,
		ProfileImageURL:      profileImageURL,
		FirebaseProviderType: updatedUser.FirebaseProviderType,
		FirebaseUID:          updatedUser.FirebaseUID,
	}, nil
}

func (service *UserService) DeleteUserByUID(ctx context.Context, uid string) *pnd.AppError {
	tx, err := service.conn.BeginTx(ctx)
	defer tx.Rollback()
	if err != nil {
		return err
	}

	if err := postgres.DeleteUserByUID(ctx, tx, uid); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (service *UserService) AddPetsToOwner(ctx context.Context, uid string, addPetsRequest pet.AddPetsToOwnerRequest) ([]pet.PetView, *pnd.AppError) {
	tx, err := service.conn.BeginTx(ctx)
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}

	user, err := postgres.FindUserByUID(ctx, tx, uid)
	if err != nil {
		return nil, err
	}

	pets := make([]pet.PetWithProfileImage, len(addPetsRequest.Pets))
	for i, item := range addPetsRequest.Pets {
		if item.ProfileImageID != nil {
			if _, err := postgres.FindMediaByID(ctx, tx, *item.ProfileImageID); err != nil {
				return nil, pnd.ErrInvalidBody(fmt.Errorf("존재하지 않는 프로필 이미지 ID입니다. ID: %d", *item.ProfileImageID))
			}
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
			ProfileImageID: item.ProfileImageID,
		}
		createdPet, err := postgres.CreatePet(ctx, tx, &petToCreate)
		if err != nil {
			return nil, err
		}
		pets[i] = *createdPet
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return pet.NewPetViewList(pets), nil
}

func (service *UserService) FindPetsByOwnerUID(ctx context.Context, uid string) (*pet.FindMyPetsView, *pnd.AppError) {
	tx, err := service.conn.BeginTx(ctx)
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}

	user, err := postgres.FindUserByUID(ctx, tx, uid)
	if err != nil {
		return nil, err
	}

	pets, err := postgres.FindPetsByOwnerID(ctx, tx, user.ID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return pet.NewFindMyPetsView(pets), nil
}
