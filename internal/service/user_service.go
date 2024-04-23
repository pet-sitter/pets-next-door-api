package service

import (
	"context"
	"errors"
	"fmt"

	utils "github.com/pet-sitter/pets-next-door-api/internal/common"
	databasegen "github.com/pet-sitter/pets-next-door-api/internal/infra/database/gen"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/user"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
	"github.com/pet-sitter/pets-next-door-api/internal/postgres"
)

type UserService struct {
	conn         *database.DB
	mediaService *MediaService
}

func NewUserService(conn *database.DB, mediaService *MediaService) *UserService {
	return &UserService{
		conn:         conn,
		mediaService: mediaService,
	}
}

func (service *UserService) RegisterUser(
	ctx context.Context, registerUserRequest *user.RegisterUserRequest,
) (*user.InternalView, *pnd.AppError) {
	if registerUserRequest.ProfileImageID != nil {
		_, err := service.mediaService.FindMediaByID(ctx, *registerUserRequest.ProfileImageID)
		if err != nil {
			return nil, err
		}
	}

	tx, err := service.conn.BeginTx(ctx)
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}

	_, err2 := databasegen.New(service.conn).WithTx(tx.Tx).CreateUser(ctx, registerUserRequest.ToDBParams())
	if err2 != nil {
		return nil, pnd.FromPostgresError(err2)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	row, err2 := databasegen.New(service.conn).FindUser(ctx, databasegen.FindUserParams{
		FbUid: utils.StrToNullStr(registerUserRequest.FirebaseUID),
	})
	if err2 != nil {
		return nil, pnd.FromPostgresError(err2)
	}

	return user.ToWithProfileImage(row).ToInternalView(), nil
}

func (service *UserService) FindUsers(
	ctx context.Context, params user.FindUsersParams,
) (*user.ListWithoutPrivateInfo, *pnd.AppError) {
	rows, err := databasegen.New(service.conn).FindUsers(ctx, params.ToDBParams())
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return user.ToListWithoutPrivateInfo(params.Page, params.Size, rows), nil
}

func (service *UserService) FindUser(
	ctx context.Context,
	params user.FindUserParams,
) (*user.WithProfileImage, *pnd.AppError) {
	row, err := databasegen.New(service.conn).FindUser(ctx, params.ToDBParams())
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return user.ToWithProfileImage(row), nil
}

func (service *UserService) ExistsByNickname(ctx context.Context, nickname string) (bool, *pnd.AppError) {
	existsByNickname, err := databasegen.New(service.conn).ExistsUserByNickname(ctx, nickname)
	if err != nil {
		return existsByNickname, pnd.FromPostgresError(err)
	}

	return existsByNickname, nil
}

func (service *UserService) UpdateUserByUID(
	ctx context.Context, uid, nickname string, profileImageID *int,
) (*user.MyProfileView, *pnd.AppError) {
	tx, err := service.conn.BeginTx(ctx)
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}

	_, err2 := databasegen.New(service.conn).WithTx(tx.Tx).UpdateUserByFbUID(ctx, databasegen.UpdateUserByFbUIDParams{
		Nickname:       nickname,
		ProfileImageID: utils.IntPtrToNullInt64(profileImageID),
		FbUid:          utils.StrToNullStr(uid),
	})
	if err2 != nil {
		return nil, pnd.FromPostgresError(err2)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	refreshedUser, err2 := databasegen.New(service.conn).FindUser(ctx, databasegen.FindUserParams{
		FbUid: utils.StrToNullStr(uid),
	})
	if err2 != nil {
		return nil, pnd.FromPostgresError(err2)
	}

	return user.ToWithProfileImage(refreshedUser).ToMyProfileView(), nil
}

func (service *UserService) DeleteUserByUID(ctx context.Context, uid string) *pnd.AppError {
	tx, err := service.conn.BeginTx(ctx)
	defer tx.Rollback()
	if err != nil {
		return err
	}

	if err := databasegen.New(service.conn).WithTx(tx.Tx).DeleteUserByFbUID(ctx, utils.StrToNullStr(uid)); err != nil {
		return pnd.FromPostgresError(err)
	}

	return tx.Commit()
}

func (service *UserService) AddPetsToOwner(
	ctx context.Context, uid string, addPetsRequest pet.AddPetsToOwnerRequest,
) ([]pet.PetView, *pnd.AppError) {
	tx, err := service.conn.BeginTx(ctx)
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}

	userData, err2 := databasegen.New(service.conn).FindUser(ctx, databasegen.FindUserParams{
		FbUid: utils.StrToNullStr(uid),
	})
	if err2 != nil {
		return nil, pnd.FromPostgresError(err2)
	}

	pets := make(pet.PetWithProfileList, len(addPetsRequest.Pets))
	for i, item := range addPetsRequest.Pets {
		if item.ProfileImageID != nil {
			if _, err := postgres.FindMediaByID(ctx, tx, *item.ProfileImageID); err != nil {
				return nil, pnd.ErrInvalidBody(fmt.Errorf("존재하지 않는 프로필 이미지 ID입니다. ID: %d", *item.ProfileImageID))
			}
		}

		petToCreate := item.ToPet(int(userData.ID))
		createdPet, err := postgres.CreatePet(ctx, tx, petToCreate)
		if err != nil {
			return nil, err
		}
		pets[i] = createdPet
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return pets.ToPetViewList(), nil
}

func (service *UserService) UpdatePet(
	ctx context.Context, uid string, petID int, updatePetRequest pet.UpdatePetRequest,
) (*pet.PetView, *pnd.AppError) {
	owner, err := service.FindUser(ctx, user.FindUserParams{FbUID: &uid, IncludeDeleted: false})
	if err != nil {
		return nil, err
	}

	petToUpdate, err := service.findPetByID(ctx, petID)
	if err != nil {
		return nil, err
	}

	if petToUpdate.OwnerID != owner.ID {
		return nil, pnd.ErrForbidden(errors.New("해당 반려동물을 수정할 권한이 없습니다"))
	}

	if updatePetRequest.ProfileImageID != nil {
		if _, err = service.mediaService.FindMediaByID(ctx, *updatePetRequest.ProfileImageID); err != nil {
			return nil, err
		}
	}

	tx, err := service.conn.BeginTx(ctx)
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}

	err = postgres.UpdatePet(ctx, tx, petID, &updatePetRequest)
	if err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	updatedPet, err := service.findPetByID(ctx, petID)
	if err != nil {
		return nil, err
	}
	return updatedPet.ToPetView(), nil
}

func (service *UserService) DeletePet(ctx context.Context, uid string, petID int) *pnd.AppError {
	owner, err := service.FindUser(ctx, user.FindUserParams{FbUID: &uid, IncludeDeleted: false})
	if err != nil {
		return err
	}

	petToDelete, err := service.findPetByID(ctx, petID)
	if err != nil {
		return err
	}

	if petToDelete.OwnerID != owner.ID {
		return pnd.ErrForbidden(errors.New("해당 반려동물을 삭제할 권한이 없습니다"))
	}

	tx, err := service.conn.BeginTx(ctx)
	defer tx.Rollback()
	if err != nil {
		return err
	}

	if err := postgres.DeletePet(ctx, tx, petID); err != nil {
		return err
	}

	return tx.Commit()
}

func (service *UserService) FindPetsByOwnerUID(ctx context.Context, uid string) (*pet.FindMyPetsView, *pnd.AppError) {
	userData, err2 := databasegen.New(service.conn).FindUser(ctx, databasegen.FindUserParams{
		FbUid:          utils.StrToNullStr(uid),
		IncludeDeleted: false,
	})
	if err2 != nil {
		return nil, pnd.FromPostgresError(err2)
	}

	tx, err := service.conn.BeginTx(ctx)
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}

	pets, err := postgres.FindPetsByOwnerID(ctx, tx, int(userData.ID))
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return pets.ToFindMyPetsView(), nil
}

func (service *UserService) findPetByID(ctx context.Context, petID int) (*pet.PetWithProfileImage, *pnd.AppError) {
	tx, err := service.conn.BeginTx(ctx)
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}

	petData, err := postgres.FindPetByID(ctx, tx, petID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return petData, nil
}
