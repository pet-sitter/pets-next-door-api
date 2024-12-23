package user

import (
	"github.com/google/uuid"
	pnd "github.com/pet-sitter/pets-next-door-api/api"
	utils "github.com/pet-sitter/pets-next-door-api/internal/common"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"
	databasegen "github.com/pet-sitter/pets-next-door-api/internal/infra/database/gen"
)

type InternalView struct {
	ID                   uuid.UUID            `json:"id"`
	Email                string               `json:"email"`
	Nickname             string               `json:"nickname"`
	Fullname             string               `json:"fullname"`
	ProfileImageURL      *string              `json:"profileImageUrl"`
	FirebaseProviderType FirebaseProviderType `json:"fbProviderType"`
	FirebaseUID          string               `json:"fbUid"`
}

func (r *InternalView) ToMyProfileView() *MyProfileView {
	return &MyProfileView{
		ID:                   r.ID,
		Email:                r.Email,
		Nickname:             r.Nickname,
		Fullname:             r.Fullname,
		ProfileImageURL:      r.ProfileImageURL,
		FirebaseProviderType: r.FirebaseProviderType,
	}
}

type MyProfileView struct {
	ID                   uuid.UUID            `json:"id"`
	Email                string               `json:"email"`
	Nickname             string               `json:"nickname"`
	Fullname             string               `json:"fullname"`
	ProfileImageURL      *string              `json:"profileImageUrl"`
	FirebaseProviderType FirebaseProviderType `json:"fbProviderType"`
}

type ProfileView struct {
	ID              uuid.UUID        `json:"id"`
	Nickname        string           `json:"nickname"`
	ProfileImageURL *string          `json:"profileImageUrl"`
	Pets            []pet.DetailView `json:"pets"`
}

func NewProfileView(
	user databasegen.FindUserRow,
	pets *pet.ListView,
) *ProfileView {
	return &ProfileView{
		ID:              user.ID,
		Nickname:        user.Nickname,
		ProfileImageURL: utils.NullStrToStrPtr(user.ProfileImageUrl),
		Pets:            pets.Pets,
	}
}

type CheckNicknameView struct {
	IsAvailable bool `json:"isAvailable"`
}

type RegistrationStatus string

const (
	StatusNotRegistered RegistrationStatus = "NOT_REGISTERED"
	StatusRegistered    RegistrationStatus = "REGISTERED"
)

type StatusView struct {
	Status               RegistrationStatus   `json:"status"`
	FirebaseProviderType FirebaseProviderType `json:"fbProviderType,omitempty"`
}

func NewStatusView(providerType FirebaseProviderType) *StatusView {
	return &StatusView{
		Status:               StatusRegistered,
		FirebaseProviderType: providerType,
	}
}

type WithoutPrivateInfo struct {
	ID              uuid.UUID `json:"id"`
	Nickname        string    `json:"nickname"`
	ProfileImageURL *string   `json:"profileImageUrl"`
}

func ToWithoutPrivateInfo(row databasegen.FindUserRow) *WithoutPrivateInfo {
	return &WithoutPrivateInfo{
		ID:              row.ID,
		Nickname:        row.Nickname,
		ProfileImageURL: utils.NullStrToStrPtr(row.ProfileImageUrl),
	}
}

type ListWithoutPrivateInfo struct {
	*pnd.PaginatedView[WithoutPrivateInfo]
}

func ToListWithoutPrivateInfo(
	page, size int,
	rows []databasegen.FindUsersRow,
) *ListWithoutPrivateInfo {
	ul := &ListWithoutPrivateInfo{PaginatedView: pnd.NewPaginatedView(
		page, size, false, make([]WithoutPrivateInfo, 0),
	)}

	for _, row := range rows {
		ul.Items = append(ul.Items, WithoutPrivateInfo{
			ID:              row.ID,
			Nickname:        row.Nickname,
			ProfileImageURL: utils.NullStrToStrPtr(row.ProfileImageUrl),
		})
	}

	ul.CalcLastPage()
	return ul
}

func ToListWithoutPrivateInfoFromFindByIDs(
	rows []databasegen.FindUsersByIDsRow,
) []WithoutPrivateInfo {
	var items []WithoutPrivateInfo
	for _, row := range rows {
		items = append(items, WithoutPrivateInfo{
			ID:              row.ID,
			Nickname:        row.Nickname,
			ProfileImageURL: utils.NullStrToStrPtr(row.ProfileImageUrl),
		})
	}

	return items
}
