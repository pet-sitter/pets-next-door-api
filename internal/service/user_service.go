package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	utils "github.com/pet-sitter/pets-next-door-api/internal/common"
	databasegen "github.com/pet-sitter/pets-next-door-api/internal/infra/database/gen"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/user"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
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
	ctx context.Context, uid, nickname string, profileImageID *int64,
) (*user.MyProfileView, *pnd.AppError) {
	tx, err := service.conn.BeginTx(ctx)
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}

	_, err2 := databasegen.New(service.conn).WithTx(tx.Tx).UpdateUserByFbUID(ctx, databasegen.UpdateUserByFbUIDParams{
		Nickname:       nickname,
		ProfileImageID: utils.Int64PtrToNullInt64(profileImageID),
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

func (service *UserService) FindPet(
	ctx context.Context, params pet.FindPetParams,
) (*pet.WithProfileImage, *pnd.AppError) {
	row, err := databasegen.New(service.conn).FindPet(ctx, params.ToDBParams())
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return pet.ToWithProfileImage(row), nil
}

func (service *UserService) FindPets(ctx context.Context, params pet.FindPetsParams) (*pet.ListView, *pnd.AppError) {
	rows, err := databasegen.New(service.conn).FindPets(ctx, params.ToDBParams())
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return pet.ToListView(rows), nil
}

func (service *UserService) AddPetsToOwner(
	ctx context.Context, uid string, addPetsRequest pet.AddPetsToOwnerRequest,
) (*pet.ListView, *pnd.AppError) {
	tx, err := service.conn.BeginTx(ctx)
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}

	// 사용자가 존재하는지 확인
	userData, err2 := databasegen.New(service.conn).FindUser(ctx, databasegen.FindUserParams{
		FbUid: utils.StrToNullStr(uid),
	})
	if err2 != nil {
		return nil, pnd.FromPostgresError(err2)
	}

	// 프로필 이미지 ID가 DB에 존재하는지 확인
	for _, item := range addPetsRequest.Pets {
		if item.ProfileImageID != nil {
			if _, err := service.mediaService.FindMediaByID(ctx, *item.ProfileImageID); err != nil {
				return nil, pnd.ErrInvalidBody(fmt.Errorf("존재하지 않는 프로필 이미지 ID입니다. ID: %d", *item.ProfileImageID))
			}
		}
	}

	// 사용자의 반려동물 추가
	petIDs := make([]int32, 0, len(addPetsRequest.Pets))
	for _, item := range addPetsRequest.Pets {
		birthDate, err := time.Parse(time.DateOnly, item.BirthDate.String())
		if err != nil {
			return nil, pnd.ErrInvalidBody(fmt.Errorf("잘못된 생년월일 형식입니다. %s", item.BirthDate.String()))
		}

		petToCreate := databasegen.CreatePetParams{
			OwnerID:        int64(userData.ID),
			Name:           item.Name,
			PetType:        string(item.PetType),
			Sex:            string(item.Sex),
			Neutered:       item.Neutered,
			Breed:          item.Breed,
			BirthDate:      birthDate,
			WeightInKg:     item.WeightInKg.String(),
			Remarks:        item.Remarks,
			ProfileImageID: utils.Int64PtrToNullInt64(item.ProfileImageID),
		}
		row, err := databasegen.New(service.conn).WithTx(tx.Tx).CreatePet(ctx, petToCreate)
		if err != nil {
			return nil, pnd.FromPostgresError(err)
		}
		petIDs = append(petIDs, row.ID)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	rows, err2 := databasegen.New(service.conn).FindPetsByIDs(ctx, databasegen.FindPetsByIDsParams{
		Ids: petIDs,
	})
	if err2 != nil {
		return nil, pnd.FromPostgresError(err2)
	}

	return pet.ToListViewFromIDsRows(rows), nil
}

func (service *UserService) UpdatePet(
	ctx context.Context, uid string, petID int64, updatePetRequest pet.UpdatePetRequest,
) (*pet.DetailView, *pnd.AppError) {
	owner, err := service.FindUser(ctx, user.FindUserParams{FbUID: &uid, IncludeDeleted: false})
	if err != nil {
		return nil, err
	}

	petToUpdate, err := service.FindPet(ctx, pet.FindPetParams{ID: &petID, IncludeDeleted: false})
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

	birthDate, err2 := time.Parse(time.DateOnly, updatePetRequest.BirthDate.String())
	if err2 != nil {
		return nil, pnd.ErrInvalidBody(fmt.Errorf("잘못된 생년월일 형식입니다. %s", updatePetRequest.BirthDate.String()))
	}

	if err := databasegen.New(service.conn).WithTx(tx.Tx).UpdatePet(ctx, databasegen.UpdatePetParams{
		ID:             int32(petID),
		Name:           updatePetRequest.Name,
		Neutered:       updatePetRequest.Neutered,
		Breed:          updatePetRequest.Breed,
		BirthDate:      birthDate,
		WeightInKg:     updatePetRequest.WeightInKg.String(),
		Remarks:        updatePetRequest.Remarks,
		ProfileImageID: utils.Int64PtrToNullInt64(updatePetRequest.ProfileImageID),
	}); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	updatedPet, err := service.FindPet(ctx, pet.FindPetParams{ID: &petID, IncludeDeleted: false})
	if err != nil {
		return nil, err
	}
	return updatedPet.ToDetailView(), nil
}

func (service *UserService) DeletePet(ctx context.Context, uid string, petID int64) *pnd.AppError {
	owner, err := service.FindUser(ctx, user.FindUserParams{FbUID: &uid})
	if err != nil {
		return err
	}

	petToDelete, err := service.FindPet(ctx, pet.FindPetParams{ID: &petID})
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

	if err := databasegen.New(service.conn).WithTx(tx.Tx).DeletePet(ctx, int32(petID)); err != nil {
		return pnd.FromPostgresError(err)
	}

	return tx.Commit()
}
