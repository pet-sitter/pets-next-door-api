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
) (*user.RegisterUserView, *pnd.AppError) {
	var profileImageURL *string
	if registerUserRequest.ProfileImageID != nil {
		mediaData, err := service.mediaService.FindMediaByID(ctx, *registerUserRequest.ProfileImageID)
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

	created, err2 := databasegen.New(service.conn).WithTx(tx.Tx).CreateUser(ctx, databasegen.CreateUserParams{
		Email:          registerUserRequest.Email,
		Nickname:       registerUserRequest.Nickname,
		Fullname:       registerUserRequest.Fullname,
		Password:       "",
		ProfileImageID: utils.IntPtrToNullInt64(registerUserRequest.ProfileImageID),
		FbProviderType: registerUserRequest.FirebaseProviderType.NullString(),
		FbUid:          utils.StrToNullStr(registerUserRequest.FirebaseUID),
	})
	if err2 != nil {
		return nil, pnd.FromPostgresError(err2)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return user.ToRegisterUserView(&created, profileImageURL), nil
}

func (service *UserService) FindUsers(
	ctx context.Context, page, size int, nickname *string,
) (*user.UserWithoutPrivateInfoList, *pnd.AppError) {
	pagination := utils.OffsetAndLimit(page, size)
	rows, err := databasegen.New(service.conn).FindUsers(ctx, databasegen.FindUsersParams{
		Limit:          int32(pagination.Limit),
		Offset:         int32(pagination.Offset),
		Nickname:       utils.StrPtrToNullStr(nickname),
		IncludeDeleted: false,
	})
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	userList := user.NewUserWithoutPrivateInfoList(page, size)
	for _, row := range rows {
		userList.Items = append(userList.Items, user.UserWithoutPrivateInfo{
			ID:              int(row.ID),
			Nickname:        row.Nickname,
			ProfileImageURL: utils.NullStrToStrPtr(row.ProfileImageUrl),
		})
	}
	userList.CalcLastPage()

	return userList, nil
}

type FindUserParams struct {
	ID             *int
	Email          *string
	FbUID          *string
	IncludeDeleted bool
}

func (service *UserService) FindUser(
	ctx context.Context,
	params FindUserParams,
) (*user.UserWithProfileImage, *pnd.AppError) {
	row, err := databasegen.New(service.conn).FindUser(ctx, databasegen.FindUserParams{
		ID:             utils.IntPtrToNullInt32(params.ID),
		Email:          utils.StrPtrToNullStr(params.Email),
		FbUid:          utils.StrPtrToNullStr(params.FbUID),
		IncludeDeleted: params.IncludeDeleted,
	})
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	userData := &user.UserWithProfileImage{
		ID:                   int(row.ID),
		Email:                row.Email,
		Nickname:             row.Nickname,
		Fullname:             row.Fullname,
		ProfileImageURL:      utils.NullStrToStrPtr(row.ProfileImageUrl),
		FirebaseProviderType: user.FirebaseProviderType(row.FbProviderType.String),
		FirebaseUID:          row.FbUid.String,
		CreatedAt:            row.CreatedAt,
		UpdatedAt:            row.UpdatedAt,
	}

	return userData, nil
}

func (service *UserService) findUserByUID(ctx context.Context, uid string) (*user.UserWithProfileImage, *pnd.AppError) {
	row, err := databasegen.New(service.conn).FindUser(ctx, databasegen.FindUserParams{
		FbUid:          utils.StrToNullStr(uid),
		IncludeDeleted: false,
	})
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	userData := &user.UserWithProfileImage{
		ID:                   int(row.ID),
		Email:                row.Email,
		Nickname:             row.Nickname,
		Fullname:             row.Fullname,
		ProfileImageURL:      utils.NullStrToStrPtr(row.ProfileImageUrl),
		FirebaseProviderType: user.FirebaseProviderType(row.FbProviderType.String),
		FirebaseUID:          row.FbUid.String,
		CreatedAt:            row.CreatedAt,
		UpdatedAt:            row.UpdatedAt,
	}

	return userData, nil
}

// FindMyProfile은 사용자의 프로필 정보를 조회한다.
// 삭제된 유저의 경우 삭제된 유저 정보를 반환한다.
func (service *UserService) FindPublicUserByID(
	ctx context.Context, id int,
) (*user.UserWithoutPrivateInfo, *pnd.AppError) {
	row, err := databasegen.New(service.conn).FindUser(ctx, databasegen.FindUserParams{
		ID:             utils.IntToNullInt32(id),
		IncludeDeleted: true,
	})
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	userData := &user.UserWithoutPrivateInfo{
		ID:              int(row.ID),
		Nickname:        row.Nickname,
		ProfileImageURL: utils.NullStrToStrPtr(row.ProfileImageUrl),
	}

	return userData, nil
}

func (service *UserService) FindUserByEmail(
	ctx context.Context, email string,
) (*user.UserWithProfileImage, *pnd.AppError) {
	row, err := databasegen.New(service.conn).FindUser(ctx, databasegen.FindUserParams{
		Email:          utils.StrToNullStr(email),
		IncludeDeleted: false,
	})
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	userData := &user.UserWithProfileImage{
		ID:                   int(row.ID),
		Email:                row.Email,
		Nickname:             row.Nickname,
		Fullname:             row.Fullname,
		ProfileImageURL:      utils.NullStrToStrPtr(row.ProfileImageUrl),
		FirebaseProviderType: user.FirebaseProviderType(row.FbProviderType.String),
		FirebaseUID:          row.FbUid.String,
		CreatedAt:            row.CreatedAt,
		UpdatedAt:            row.UpdatedAt,
	}

	return userData, nil
}

func (service *UserService) FindUserByUID(ctx context.Context, uid string) (*user.FindUserView, *pnd.AppError) {
	row, err := databasegen.New(service.conn).FindUser(ctx, databasegen.FindUserParams{
		FbUid:          utils.StrToNullStr(uid),
		IncludeDeleted: false,
	})
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	userData := &user.UserWithProfileImage{
		ID:                   int(row.ID),
		Email:                row.Email,
		Nickname:             row.Nickname,
		Fullname:             row.Fullname,
		ProfileImageURL:      utils.NullStrToStrPtr(row.ProfileImageUrl),
		FirebaseProviderType: user.FirebaseProviderType(row.FbProviderType.String),
		FirebaseUID:          row.FbUid.String,
		CreatedAt:            row.CreatedAt,
		UpdatedAt:            row.UpdatedAt,
	}

	return userData.ToFindUserView(), nil
}

func (service *UserService) ExistsByNickname(ctx context.Context, nickname string) (bool, *pnd.AppError) {
	existsByNickname, err := databasegen.New(service.conn).ExistsUserByNickname(ctx, nickname)
	if err != nil {
		return existsByNickname, pnd.FromPostgresError(err)
	}

	return existsByNickname, nil
}

func (service *UserService) FindUserStatusByEmail(ctx context.Context, email string) (*user.UserStatus, *pnd.AppError) {
	userData, err := databasegen.New(service.conn).FindUser(ctx, databasegen.FindUserParams{
		Email:          utils.StrToNullStr(email),
		IncludeDeleted: false,
	})
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	userStatus := &user.UserStatus{
		FirebaseProviderType: user.FirebaseProviderType(userData.FbProviderType.String),
	}
	return userStatus, nil
}

func (service *UserService) UpdateUserByUID(
	ctx context.Context, uid, nickname string, profileImageID *int,
) (*user.UserWithProfileImage, *pnd.AppError) {
	tx, err := service.conn.BeginTx(ctx)
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}

	row, err2 := databasegen.New(service.conn).WithTx(tx.Tx).UpdateUserByFbUID(ctx, databasegen.UpdateUserByFbUIDParams{
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

	var profileImageURL *string
	if !row.ProfileImageID.Valid {
		profileImage, err := service.mediaService.FindMediaByID(ctx, int(row.ProfileImageID.Int64))
		if err != nil {
			return nil, err
		}

		if profileImage != nil {
			profileImageURL = &profileImage.URL
		}
	}

	updatedUser := user.UserWithProfileImage{
		ID:                   int(row.ID),
		Email:                row.Email,
		Nickname:             row.Nickname,
		Fullname:             row.Fullname,
		ProfileImageURL:      profileImageURL,
		FirebaseProviderType: user.FirebaseProviderType(row.FbProviderType.String),
		FirebaseUID:          row.FbUid.String,
		CreatedAt:            row.CreatedAt,
		UpdatedAt:            row.UpdatedAt,
	}

	return &updatedUser, nil
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
	owner, err := service.findUserByUID(ctx, uid)
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
	owner, err := service.findUserByUID(ctx, uid)
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
