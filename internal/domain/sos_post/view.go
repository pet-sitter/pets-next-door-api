package sos_post

import (
	pnd "github.com/pet-sitter/pets-next-door-api/api"
	utils "github.com/pet-sitter/pets-next-door-api/internal/common"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/media"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/user"
)

type ConditionView struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ConditionViewList []*ConditionView

func (c *Condition) ToConditionView() *ConditionView {
	return &ConditionView{
		ID:   c.ID,
		Name: c.Name,
	}
}

func (cl *ConditionList) ToConditionViewList() []ConditionView {
	conditionViews := make([]ConditionView, len(*cl))
	for i, c := range *cl {
		conditionViews[i] = *c.ToConditionView()
	}
	return conditionViews
}

type WriteSosPostRequest struct {
	Title        string        `json:"title" validate:"required"`
	Content      string        `json:"content" validate:"required"`
	ImageIDs     []int         `json:"imageIds" validate:"required"`
	Reward       string        `json:"reward"`
	Dates        []SosDateView `json:"dates" validate:"required"`
	CareType     CareType      `json:"careType" validate:"required,oneof=foster visiting"`
	CarerGender  CarerGender   `json:"carerGender" validate:"required,oneof=male female all"`
	RewardType   RewardType    `json:"rewardType" validate:"required,oneof=fee gifticon negotiable"`
	ConditionIDs []int         `json:"conditionIds"`
	PetIDs       []int         `json:"petIds"`
}

type WriteSosPostView struct {
	ID          int                 `json:"id"`
	AuthorID    int                 `json:"authorId"`
	Title       string              `json:"title"`
	Content     string              `json:"content"`
	Media       media.MediaViewList `json:"media"`
	Conditions  []ConditionView     `json:"conditions"`
	Pets        []pet.PetView       `json:"pets"`
	Reward      string              `json:"reward"`
	Dates       []SosDateView       `json:"dates"`
	CareType    CareType            `json:"careType"`
	CarerGender CarerGender         `json:"carerGender"`
	RewardType  RewardType          `json:"rewardType"`
	ThumbnailID int                 `json:"thumbnailId"`
	CreatedAt   string              `json:"createdAt"`
	UpdatedAt   string              `json:"updatedAt"`
}

func (p *SosPost) ToWriteSosPostView(
	media media.MediaViewList,
	conditions []ConditionView,
	pets []pet.PetView,
	sosDates []SosDateView,
) *WriteSosPostView {
	return &WriteSosPostView{
		ID:          p.ID,
		AuthorID:    p.AuthorID,
		Title:       p.Title,
		Content:     p.Content,
		Media:       media,
		Conditions:  conditions,
		Pets:        pets,
		Reward:      p.Reward,
		Dates:       sosDates,
		CareType:    p.CareType,
		CarerGender: p.CarerGender,
		RewardType:  p.RewardType,
		ThumbnailID: p.ThumbnailID,
		CreatedAt:   utils.FormatDateTime(p.CreatedAt),
		UpdatedAt:   utils.FormatDateTime(p.UpdatedAt),
	}
}

type FindSosPostView struct {
	ID          int                          `json:"id"`
	Author      *user.UserWithoutPrivateInfo `json:"author"`
	Title       string                       `json:"title"`
	Content     string                       `json:"content"`
	Media       media.MediaViewList          `json:"media"`
	Conditions  []ConditionView              `json:"conditions"`
	Pets        []pet.PetView                `json:"pets"`
	Reward      string                       `json:"reward"`
	Dates       []SosDateView                `json:"dates"`
	CareType    CareType                     `json:"careType"`
	CarerGender CarerGender                  `json:"carerGender"`
	RewardType  RewardType                   `json:"rewardType"`
	ThumbnailID int                          `json:"thumbnailId"`
	CreatedAt   string                       `json:"createdAt"`
	UpdatedAt   string                       `json:"updatedAt"`
}

func (p *SosPost) ToFindSosPostView(
	author *user.UserWithoutPrivateInfo,
	media media.MediaViewList,
	conditions []ConditionView,
	pets []pet.PetView,
	sosDates []SosDateView,
) *FindSosPostView {
	return &FindSosPostView{
		ID:          p.ID,
		Author:      author,
		Title:       p.Title,
		Content:     p.Content,
		Media:       media,
		Conditions:  conditions,
		Pets:        pets,
		Reward:      p.Reward,
		Dates:       sosDates,
		CareType:    p.CareType,
		CarerGender: p.CarerGender,
		RewardType:  p.RewardType,
		ThumbnailID: p.ThumbnailID,
		CreatedAt:   utils.FormatDateTime(p.CreatedAt),
		UpdatedAt:   utils.FormatDateTime(p.UpdatedAt),
	}
}

type FindSosPostListView struct {
	*pnd.PaginatedView[FindSosPostView]
}

func FromEmptySosPostList(sosPosts *SosPostList) *FindSosPostListView {
	return &FindSosPostListView{
		PaginatedView: pnd.NewPaginatedView(
			sosPosts.Page, sosPosts.Size, sosPosts.IsLastPage, make([]FindSosPostView, 0),
		),
	}
}

func FromEmptySosPostInfoList(sosPosts *SosPostInfoList) *FindSosPostListView {
	return &FindSosPostListView{
		PaginatedView: pnd.NewPaginatedView(
			sosPosts.Page, sosPosts.Size, sosPosts.IsLastPage, make([]FindSosPostView, 0),
		),
	}
}

func (p *SosPostInfo) ToFindSosPostInfoView(
	author *user.UserWithoutPrivateInfo,
	media media.MediaViewList,
	conditions []ConditionView,
	pets []pet.PetView,
	sosDates []SosDateView,
) *FindSosPostView {
	return &FindSosPostView{
		ID:          p.ID,
		Author:      author,
		Title:       p.Title,
		Content:     p.Content,
		Media:       media,
		Conditions:  conditions,
		Pets:        pets,
		Reward:      p.Reward,
		Dates:       sosDates,
		CareType:    p.CareType,
		CarerGender: p.CarerGender,
		RewardType:  p.RewardType,
		ThumbnailID: p.ThumbnailID,
		CreatedAt:   utils.FormatDateTime(p.CreatedAt),
		UpdatedAt:   utils.FormatDateTime(p.UpdatedAt),
	}
}

type UpdateSosPostRequest struct {
	ID           int           `json:"id" validate:"required"`
	Title        string        `json:"title" validate:"required"`
	Content      string        `json:"content" validate:"required"`
	ImageIDs     []int         `json:"imageIds" validate:"required"`
	Dates        []SosDateView `json:"dates" validate:"required"`
	Reward       string        `json:"reward"`
	CareType     CareType      `json:"careType" validate:"required,oneof=foster visiting"`
	CarerGender  CarerGender   `json:"carerGender" validate:"required,oneof=male female all"`
	RewardType   RewardType    `json:"rewardType" validate:"required,oneof=fee gifticon negotiable"`
	ConditionIDs []int         `json:"conditionIds"`
	PetIDs       []int         `json:"petIds"`
}

type UpdateSosPostView struct {
	ID          int                 `json:"id"`
	AuthorID    int                 `json:"authorId"`
	Title       string              `json:"title"`
	Content     string              `json:"content"`
	Media       media.MediaViewList `json:"media"`
	Conditions  []ConditionView     `json:"conditions"`
	Pets        []pet.PetView       `json:"pets"`
	Reward      string              `json:"reward"`
	Dates       []SosDateView       `json:"dates"`
	CareType    CareType            `json:"careType"`
	CarerGender CarerGender         `json:"carerGender"`
	RewardType  RewardType          `json:"rewardType"`
	ThumbnailID int                 `json:"thumbnailId"`
	CreatedAt   string              `json:"createdAt"`
	UpdatedAt   string              `json:"updatedAt"`
}

func (p *SosPost) ToUpdateSosPostView(
	media media.MediaViewList,
	conditions []ConditionView,
	pets []pet.PetView,
	sosDates []SosDateView,
) *UpdateSosPostView {
	return &UpdateSosPostView{
		ID:          p.ID,
		AuthorID:    p.AuthorID,
		Title:       p.Title,
		Content:     p.Content,
		Media:       media,
		Conditions:  conditions,
		Pets:        pets,
		Reward:      p.Reward,
		Dates:       sosDates,
		CareType:    p.CareType,
		CarerGender: p.CarerGender,
		RewardType:  p.RewardType,
		ThumbnailID: p.ThumbnailID,
		CreatedAt:   utils.FormatDateTime(p.CreatedAt),
		UpdatedAt:   utils.FormatDateTime(p.UpdatedAt),
	}
}

type SosDateView struct {
	DateStartAt string `json:"dateStartAt"`
	DateEndAt   string `json:"dateEndAt"`
}

func (d *SosDates) ToSosDateView() SosDateView {
	return SosDateView{
		DateStartAt: utils.FormatDate(d.DateStartAt),
		DateEndAt:   utils.FormatDate(d.DateEndAt),
	}
}

func (dl *SosDatesList) ToSosDateViewList() []SosDateView {
	sosDateViews := make([]SosDateView, len(*dl))
	for i, d := range *dl {
		sosDateViews[i] = d.ToSosDateView()
	}
	return sosDateViews
}
