package sospost

import (
	"github.com/google/uuid"
	pnd "github.com/pet-sitter/pets-next-door-api/api"
	utils "github.com/pet-sitter/pets-next-door-api/internal/common"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/media"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/soscondition"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/user"
	databasegen "github.com/pet-sitter/pets-next-door-api/internal/infra/database/gen"
)

type ViewParams struct {
	ID          uuid.UUID
	AuthorID    uuid.UUID
	Title       string
	Content     string
	MediaList   media.ListView
	Conditions  soscondition.ListView
	Pets        []pet.DetailView
	Reward      string
	SOSDates    []SOSDateView
	CareType    CareType
	CarerGender CarerGender
	RewardType  RewardType
	ThumbnailID uuid.NullUUID
	CreatedAt   string
	UpdatedAt   string
}

type ViewParamsInput struct {
	ID          uuid.UUID
	AuthorID    uuid.UUID
	Title       string
	Content     string
	MediaList   media.ListView
	Conditions  soscondition.ListView
	Pets        []pet.DetailView
	Reward      string
	SOSDates    []SOSDateView
	CareType    string
	CarerGender string
	RewardType  string
	ThumbnailID uuid.NullUUID
	CreatedAt   string
	UpdatedAt   string
}

type DetailView struct {
	ID          uuid.UUID             `json:"id"`
	AuthorID    uuid.UUID             `json:"authorId"`
	Title       string                `json:"title"`
	Content     string                `json:"content"`
	Media       media.ListView        `json:"media"`
	Conditions  soscondition.ListView `json:"conditions"`
	Pets        []pet.DetailView      `json:"pets"`
	Reward      string                `json:"reward"`
	Dates       []SOSDateView         `json:"dates"`
	CareType    CareType              `json:"careType"`
	CarerGender CarerGender           `json:"carerGender"`
	RewardType  RewardType            `json:"rewardType"`
	ThumbnailID uuid.NullUUID         `json:"thumbnailId"`
	CreatedAt   string                `json:"createdAt"`
	UpdatedAt   string                `json:"updatedAt"`
}

func ToDetailView(params ViewParams) *DetailView {
	return &DetailView{
		ID:          params.ID,
		AuthorID:    params.AuthorID,
		Title:       params.Title,
		Content:     params.Content,
		Media:       params.MediaList,
		Conditions:  params.Conditions,
		Pets:        params.Pets,
		Reward:      params.Reward,
		Dates:       params.SOSDates,
		CareType:    params.CareType,
		CarerGender: params.CarerGender,
		RewardType:  params.RewardType,
		ThumbnailID: params.ThumbnailID,
		CreatedAt:   params.CreatedAt,
		UpdatedAt:   params.UpdatedAt,
	}
}

func CreateViewParams(input ViewParamsInput) ViewParams {
	return ViewParams{
		ID:          input.ID,
		AuthorID:    input.AuthorID,
		Title:       input.Title,
		Content:     input.Content,
		MediaList:   input.MediaList,
		Conditions:  input.Conditions,
		Pets:        input.Pets,
		Reward:      input.Reward,
		SOSDates:    input.SOSDates,
		CareType:    CareType(input.CareType),
		CarerGender: CarerGender(input.CarerGender),
		RewardType:  RewardType(input.RewardType),
		ThumbnailID: input.ThumbnailID,
		CreatedAt:   input.CreatedAt,
		UpdatedAt:   input.UpdatedAt,
	}
}

func CreateDetailView(
	sosPost databasegen.WriteSOSPostRow,
	mediaList media.ListView,
	conditions soscondition.ListView,
	pets []pet.DetailView,
	sosDates []SOSDateView,
) *DetailView {
	input := ViewParamsInput{
		// ID:          sosPost.ID,
		// AuthorID:    sosPost.AuthorID.Int64,
		Title:       utils.NullStrToStr(sosPost.Title),
		Content:     utils.NullStrToStr(sosPost.Content),
		MediaList:   mediaList,
		Conditions:  conditions,
		Pets:        pets,
		Reward:      utils.NullStrToStr(sosPost.Reward),
		SOSDates:    sosDates,
		CareType:    sosPost.CareType.String,
		CarerGender: sosPost.CarerGender.String,
		RewardType:  sosPost.RewardType.String,
		// ThumbnailID: sosPost.ThumbnailID.Int64,
		CreatedAt: utils.FormatTimeFromTime(sosPost.CreatedAt),
		UpdatedAt: utils.FormatTimeFromTime(sosPost.UpdatedAt),
	}
	params := CreateViewParams(input)
	return ToDetailView(params)
}

type FindSOSPostView struct {
	ID          uuid.UUID                `json:"id"`
	Author      *user.WithoutPrivateInfo `json:"author"`
	Title       string                   `json:"title"`
	Content     string                   `json:"content"`
	Media       media.ListView           `json:"media"`
	Conditions  soscondition.ListView    `json:"conditions"`
	Pets        []pet.DetailView         `json:"pets"`
	Reward      string                   `json:"reward"`
	Dates       []SOSDateView            `json:"dates"`
	CareType    CareType                 `json:"careType"`
	CarerGender CarerGender              `json:"carerGender"`
	RewardType  RewardType               `json:"rewardType"`
	ThumbnailID uuid.NullUUID            `json:"thumbnailId"`
	CreatedAt   string                   `json:"createdAt"`
	UpdatedAt   string                   `json:"updatedAt"`
}

func (p *SOSPost) ToFindSOSPostView(
	author *user.WithoutPrivateInfo,
	mediaList media.ListView,
	conditions soscondition.ListView,
	pets []pet.DetailView,
	sosDates []SOSDateView,
) *FindSOSPostView {
	return &FindSOSPostView{
		// ID:          p.ID,
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
		// ThumbnailID: p.ThumbnailID,
		CreatedAt: utils.FormatDateTimeFromTime(p.CreatedAt),
		UpdatedAt: utils.FormatDateTimeFromTime(p.UpdatedAt),
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
	conditions soscondition.ListView,
	pets []pet.DetailView,
	sosDates []SOSDateView,
) *FindSOSPostView {
	return &FindSOSPostView{
		// ID:          p.ID,
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
		// ThumbnailID: p.ThumbnailID,
		CreatedAt: utils.FormatDateTimeFromTime(p.CreatedAt),
		UpdatedAt: utils.FormatDateTimeFromTime(p.UpdatedAt),
	}
}

func UpdateDetailView(
	sosPost databasegen.UpdateSOSPostRow,
	mediaList media.ListView,
	conditions soscondition.ListView,
	pets []pet.DetailView,
	sosDates []SOSDateView,
) *DetailView {
	params := ViewParams{
		// ID:          sosPost.ID,
		// AuthorID:    sosPost.AuthorID,
		Title:       utils.NullStrToStr(sosPost.Title),
		Content:     utils.NullStrToStr(sosPost.Content),
		MediaList:   mediaList,
		Conditions:  conditions,
		Pets:        pets,
		Reward:      utils.NullStrToStr(sosPost.Reward),
		SOSDates:    sosDates,
		CareType:    CareType(sosPost.CareType.String),
		CarerGender: CarerGender(sosPost.CarerGender.String),
		RewardType:  RewardType(sosPost.RewardType.String),
		// ThumbnailID: sosPost.ThumbnailID,
		CreatedAt: utils.FormatDateTimeFromTime(sosPost.CreatedAt),
		UpdatedAt: utils.FormatDateTimeFromTime(sosPost.UpdatedAt),
	}
	return ToDetailView(params)
}

type SOSDateView struct {
	DateStartAt string `json:"dateStartAt"`
	DateEndAt   string `json:"dateEndAt"`
}

func (d *SOSDates) ToSOSDateView() SOSDateView {
	return SOSDateView{
		DateStartAt: utils.FormatDateString(d.DateStartAt),
		DateEndAt:   utils.FormatDateString(d.DateEndAt),
	}
}

func ToListViewFromSOSDateRows(rows []databasegen.FindDatesBySOSPostIDRow) []SOSDateView {
	sosDateViews := make([]SOSDateView, len(rows))
	for i, row := range rows {
		date := SOSDates{
			DateStartAt: utils.NullTimeToStr(row.DateStartAt),
			DateEndAt:   utils.NullTimeToStr(row.DateEndAt),
		}
		sosDateViews[i] = date.ToSOSDateView()
	}
	return sosDateViews
}

func (dl *SOSDatesList) ToSOSDateViewList() []SOSDateView {
	sosDateViews := make([]SOSDateView, len(*dl))
	for i, d := range *dl {
		sosDateViews[i] = d.ToSOSDateView()
	}
	return sosDateViews
}
