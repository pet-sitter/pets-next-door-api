package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	utils "github.com/pet-sitter/pets-next-door-api/internal/common"
	"github.com/pet-sitter/pets-next-door-api/internal/datatype"
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
) (*user.InternalView, error) {
	if registerUserRequest.ProfileImageID.Valid {
		_, err := service.mediaService.FindMediaByID(ctx, registerUserRequest.ProfileImageID.UUID)
		if err != nil {
			return nil, err
		}
	}

	tx, err := service.conn.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	_, err = databasegen.New(service.conn).
		WithTx(tx.Tx).
		CreateUser(ctx, registerUserRequest.ToDBParams())
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	row, err := databasegen.New(service.conn).FindUser(ctx, databasegen.FindUserParams{
		FbUid: utils.StrToNullStr(registerUserRequest.FirebaseUID),
	})
	if err != nil {
		return nil, err
	}

	return user.ToWithProfileImage(row).ToInternalView(), nil
}

func (service *UserService) FindUsers(
	ctx context.Context, params user.FindUsersParams,
) (*user.ListWithoutPrivateInfo, error) {
	rows, err := databasegen.New(service.conn).FindUsers(ctx, params.ToDBParams())
	if err != nil {
		return nil, err
	}

	return user.ToListWithoutPrivateInfo(params.Page, params.Size, rows), nil
}

func (service *UserService) FindUser(
	ctx context.Context,
	params user.FindUserParams,
) (*user.WithProfileImage, error) {
	row, err := databasegen.New(service.conn).FindUser(ctx, params.ToDBParams())
	if err != nil {
		return nil, err
	}

	return user.ToWithProfileImage(row), nil
}

func (service *UserService) FindUserProfile(
	ctx context.Context, params user.FindUserParams,
) (*user.ProfileView, error) {
	row, err := databasegen.New(service.conn).FindUser(ctx, params.ToDBParams())
	if err != nil {
		return nil, err
	}

	petParams := pet.FindPetsParams{OwnerID: uuid.NullUUID{UUID: row.ID, Valid: true}}
	pets, err := service.FindPets(ctx, petParams)
	if err != nil {
		return nil, err
	}

	return user.NewProfileView(row, pets), nil
}

func (service *UserService) ExistsByNickname(
	ctx context.Context,
	nickname string,
) (bool, error) {
	existsByNickname, err := databasegen.New(service.conn).ExistsUserByNickname(ctx, nickname)
	if err != nil {
		return existsByNickname, err
	}

	return existsByNickname, nil
}

func (service *UserService) UpdateUserByUID(
	ctx context.Context, uid, nickname string, profileImageID uuid.NullUUID,
) (*user.MyProfileView, error) {
	tx, err := service.conn.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	_, err = databasegen.New(service.conn).
		WithTx(tx.Tx).
		UpdateUserByFbUID(ctx, databasegen.UpdateUserByFbUIDParams{
			Nickname:       nickname,
			ProfileImageID: profileImageID,
			FbUid:          utils.StrToNullStr(uid),
		})
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	refreshedUser, err := databasegen.New(service.conn).FindUser(ctx, databasegen.FindUserParams{
		FbUid: utils.StrToNullStr(uid),
	})
	if err != nil {
		return nil, err
	}

	return user.ToWithProfileImage(refreshedUser).ToMyProfileView(), nil
}

func (service *UserService) DeleteUserByUID(ctx context.Context, uid string) error {
	tx, err := service.conn.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := databasegen.New(service.conn).WithTx(tx.Tx).DeleteUserByFbUID(ctx, utils.StrToNullStr(uid)); err != nil {
		return err
	}

	return tx.Commit()
}

func (service *UserService) FindPet(
	ctx context.Context, params pet.FindPetParams,
) (*pet.WithProfileImage, error) {
	row, err := databasegen.New(service.conn).FindPet(ctx, params.ToDBParams())
	if err != nil {
		return nil, err
	}

	return pet.ToWithProfileImage(row), nil
}

func (service *UserService) FindPets(
	ctx context.Context,
	params pet.FindPetsParams,
) (*pet.ListView, error) {
	rows, err := databasegen.New(service.conn).FindPets(ctx, params.ToDBParams())
	if err != nil {
		return nil, err
	}

	return pet.ToListView(rows), nil
}

func (service *UserService) AddPetsToOwner(
	ctx context.Context, uid string, addPetsRequest pet.AddPetsToOwnerRequest,
) (*pet.ListView, error) {
	tx, err := service.conn.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// 사용자가 존재하는지 확인
	userData, err := databasegen.New(service.conn).FindUser(ctx, databasegen.FindUserParams{
		FbUid: utils.StrToNullStr(uid),
	})
	if err != nil {
		return nil, err
	}

	// 프로필 이미지 ID가 DB에 존재하는지 확인
	for _, item := range addPetsRequest.Pets {
		if item.ProfileImageID.Valid {
			if _, err := service.mediaService.FindMediaByID(ctx, item.ProfileImageID.UUID); err != nil {
				return nil, pnd.ErrInvalidBody(
					fmt.Errorf("존재하지 않는 프로필 이미지 ID입니다. ID: %s", item.ProfileImageID.UUID),
				)
			}
		}
	}

	// 사용자의 반려동물 추가
	petIDs := make([]uuid.UUID, 0, len(addPetsRequest.Pets))
	for _, item := range addPetsRequest.Pets {
		birthDate, err := datatype.ParseDateToTime(item.BirthDate)
		if err != nil {
			return nil, pnd.ErrInvalidBody(fmt.Errorf("잘못된 생년월일 형식입니다. %s", item.BirthDate))
		}

		petToCreate := databasegen.CreatePetParams{
			ID:             datatype.NewUUIDV7(),
			OwnerID:        userData.ID,
			Name:           item.Name,
			PetType:        string(item.PetType),
			Sex:            string(item.Sex),
			Neutered:       item.Neutered,
			Breed:          item.Breed,
			BirthDate:      birthDate,
			WeightInKg:     item.WeightInKg.String(),
			Remarks:        item.Remarks,
			ProfileImageID: item.ProfileImageID,
		}
		row, err := databasegen.New(service.conn).WithTx(tx.Tx).CreatePet(ctx, petToCreate)
		if err != nil {
			return nil, err
		}
		petIDs = append(petIDs, row.ID)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	rows, err := databasegen.New(service.conn).FindPetsByIDs(ctx, databasegen.FindPetsByIDsParams{
		Ids: petIDs,
	})
	if err != nil {
		return nil, err
	}

	return pet.ToListViewFromIDsRows(rows), nil
}

func (service *UserService) UpdatePet(
	ctx context.Context, uid string, petID uuid.UUID, updatePetRequest pet.UpdatePetRequest,
) (*pet.DetailView, error) {
	owner, err := service.FindUser(ctx, user.FindUserParams{FbUID: &uid, IncludeDeleted: false})
	if err != nil {
		return nil, err
	}

	petToUpdate, err := service.FindPet(
		ctx,
		pet.FindPetParams{ID: uuid.NullUUID{UUID: petID, Valid: true}, IncludeDeleted: false},
	)
	if err != nil {
		return nil, err
	}

	if petToUpdate.OwnerID != owner.ID {
		return nil, pnd.ErrForbidden(errors.New("해당 반려동물을 수정할 권한이 없습니다"))
	}

	if updatePetRequest.ProfileImageID.Valid {
		if _, err = service.mediaService.FindMediaByID(ctx, updatePetRequest.ProfileImageID.UUID); err != nil {
			return nil, err
		}
	}

	tx, err := service.conn.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	birthDate, err := datatype.ParseDateToTime(updatePetRequest.BirthDate)
	if err != nil {
		return nil, pnd.ErrInvalidBody(fmt.Errorf("잘못된 생년월일 형식입니다. %s", updatePetRequest.BirthDate))
	}

	if err := databasegen.New(service.conn).WithTx(tx.Tx).UpdatePet(ctx, databasegen.UpdatePetParams{
		ID:             petID,
		Name:           updatePetRequest.Name,
		Neutered:       updatePetRequest.Neutered,
		Breed:          updatePetRequest.Breed,
		BirthDate:      birthDate,
		WeightInKg:     updatePetRequest.WeightInKg.String(),
		Remarks:        updatePetRequest.Remarks,
		ProfileImageID: updatePetRequest.ProfileImageID,
	}); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	updatedPet, err := service.FindPet(
		ctx,
		pet.FindPetParams{ID: uuid.NullUUID{UUID: petID, Valid: true}, IncludeDeleted: false},
	)
	if err != nil {
		return nil, err
	}
	return updatedPet.ToDetailView(), nil
}

func (service *UserService) DeletePet(
	ctx context.Context,
	uid string,
	petID uuid.UUID,
) error {
	owner, err := service.FindUser(ctx, user.FindUserParams{FbUID: &uid})
	if err != nil {
		return err
	}

	petToDelete, err := service.FindPet(
		ctx,
		pet.FindPetParams{ID: uuid.NullUUID{UUID: petID, Valid: true}},
	)
	if err != nil {
		return err
	}

	if petToDelete.OwnerID != owner.ID {
		return pnd.ErrForbidden(errors.New("해당 반려동물을 삭제할 권한이 없습니다"))
	}

	tx, err := service.conn.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := databasegen.New(service.conn).WithTx(tx.Tx).DeletePet(ctx, petID); err != nil {
		return err
	}

	return tx.Commit()
}
