package sospost

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

type WriteSOSPostRequest struct {
	Title        string        `json:"title" validate:"required"`
	Content      string        `json:"content" validate:"required"`
	ImageIDs     []int64       `json:"imageIds" validate:"required"`
	Reward       string        `json:"reward"`
	Dates        []SOSDateView `json:"dates" validate:"required"`
	CareType     CareType      `json:"careType" validate:"required,oneof=foster visiting"`
	CarerGender  CarerGender   `json:"carerGender" validate:"required,oneof=male female all"`
	RewardType   RewardType    `json:"rewardType" validate:"required,oneof=fee gifticon negotiable"`
	ConditionIDs []int         `json:"conditionIds"`
	PetIDs       []int64       `json:"petIds"`
}

type WriteSOSPostView struct {
	ID          int              `json:"id"`
	AuthorID    int              `json:"authorId"`
	Title       string           `json:"title"`
	Content     string           `json:"content"`
	Media       media.ListView   `json:"media"`
	Conditions  []ConditionView  `json:"conditions"`
	Pets        []pet.DetailView `json:"pets"`
	Reward      string           `json:"reward"`
	Dates       []SOSDateView    `json:"dates"`
	CareType    CareType         `json:"careType"`
	CarerGender CarerGender      `json:"carerGender"`
	RewardType  RewardType       `json:"rewardType"`
	ThumbnailID int64            `json:"thumbnailId"`
	CreatedAt   string           `json:"createdAt"`
	UpdatedAt   string           `json:"updatedAt"`
}

func (p *SOSPost) ToWriteSOSPostView(
	mediaList media.ListView,
	conditions []ConditionView,
	pets []pet.DetailView,
	sosDates []SOSDateView,
) *WriteSOSPostView {
	return &WriteSOSPostView{
		ID:          p.ID,
		AuthorID:    p.AuthorID,
		Title:       p.Title,
		Content:     p.Content,
		Media:       mediaList,
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

type FindSOSPostView struct {
	ID          int                      `json:"id"`
	Author      *user.WithoutPrivateInfo `json:"author"`
	Title       string                   `json:"title"`
	Content     string                   `json:"content"`
	Media       media.ListView           `json:"media"`
	Conditions  []ConditionView          `json:"conditions"`
	Pets        []pet.DetailView         `json:"pets"`
	Reward      string                   `json:"reward"`
	Dates       []SOSDateView            `json:"dates"`
	CareType    CareType                 `json:"careType"`
	CarerGender CarerGender              `json:"carerGender"`
	RewardType  RewardType               `json:"rewardType"`
	ThumbnailID int64                    `json:"thumbnailId"`
	CreatedAt   string                   `json:"createdAt"`
	UpdatedAt   string                   `json:"updatedAt"`
}

func (p *SOSPost) ToFindSOSPostView(
	author *user.WithoutPrivateInfo,
	mediaList media.ListView,
	conditions []ConditionView,
	pets []pet.DetailView,
	sosDates []SOSDateView,
) *FindSOSPostView {
	return &FindSOSPostView{
		ID:          p.ID,
		Author:      author,
		Title:       p.Title,
		Content:     p.Content,
		Media:       mediaList,
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

type FindSOSPostListView struct {
	*pnd.PaginatedView[FindSOSPostView]
}

func FromEmptySOSPostList(sosPosts *SOSPostList) *FindSOSPostListView {
	return &FindSOSPostListView{
		PaginatedView: pnd.NewPaginatedView(
			sosPosts.Page, sosPosts.Size, sosPosts.IsLastPage, make([]FindSOSPostView, 0),
		),
	}
}

func FromEmptySOSPostInfoList(sosPosts *SOSPostInfoList) *FindSOSPostListView {
	return &FindSOSPostListView{
		PaginatedView: pnd.NewPaginatedView(
			sosPosts.Page, sosPosts.Size, sosPosts.IsLastPage, make([]FindSOSPostView, 0),
		),
	}
}

func (p *SOSPostInfo) ToFindSOSPostInfoView(
	author *user.WithoutPrivateInfo,
	mediaList media.ListView,
	conditions []ConditionView,
	pets []pet.DetailView,
	sosDates []SOSDateView,
) *FindSOSPostView {
	return &FindSOSPostView{
		ID:          p.ID,
		Author:      author,
		Title:       p.Title,
		Content:     p.Content,
		Media:       mediaList,
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

type UpdateSOSPostRequest struct {
	ID           int           `json:"id" validate:"required"`
	Title        string        `json:"title" validate:"required"`
	Content      string        `json:"content" validate:"required"`
	ImageIDs     []int64       `json:"imageIds" validate:"required"`
	Dates        []SOSDateView `json:"dates" validate:"required"`
	Reward       string        `json:"reward"`
	CareType     CareType      `json:"careType" validate:"required,oneof=foster visiting"`
	CarerGender  CarerGender   `json:"carerGender" validate:"required,oneof=male female all"`
	RewardType   RewardType    `json:"rewardType" validate:"required,oneof=fee gifticon negotiable"`
	ConditionIDs []int         `json:"conditionIds"`
	PetIDs       []int         `json:"petIds"`
}

type UpdateSOSPostView struct {
	ID          int              `json:"id"`
	AuthorID    int              `json:"authorId"`
	Title       string           `json:"title"`
	Content     string           `json:"content"`
	Media       media.ListView   `json:"media"`
	Conditions  []ConditionView  `json:"conditions"`
	Pets        []pet.DetailView `json:"pets"`
	Reward      string           `json:"reward"`
	Dates       []SOSDateView    `json:"dates"`
	CareType    CareType         `json:"careType"`
	CarerGender CarerGender      `json:"carerGender"`
	RewardType  RewardType       `json:"rewardType"`
	ThumbnailID int64            `json:"thumbnailId"`
	CreatedAt   string           `json:"createdAt"`
	UpdatedAt   string           `json:"updatedAt"`
}

func (p *SOSPost) ToUpdateSOSPostView(
	mediaList media.ListView,
	conditions []ConditionView,
	pets []pet.DetailView,
	sosDates []SOSDateView,
) *UpdateSOSPostView {
	return &UpdateSOSPostView{
		ID:          p.ID,
		AuthorID:    p.AuthorID,
		Title:       p.Title,
		Content:     p.Content,
		Media:       mediaList,
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

type SOSDateView struct {
	DateStartAt string `json:"dateStartAt"`
	DateEndAt   string `json:"dateEndAt"`
}

func (d *SOSDates) ToSOSDateView() SOSDateView {
	return SOSDateView{
		DateStartAt: utils.FormatDate(d.DateStartAt),
		DateEndAt:   utils.FormatDate(d.DateEndAt),
	}
}

func (dl *SOSDatesList) ToSOSDateViewList() []SOSDateView {
	sosDateViews := make([]SOSDateView, len(*dl))
	for i, d := range *dl {
		sosDateViews[i] = d.ToSOSDateView()
	}
	return sosDateViews
}
