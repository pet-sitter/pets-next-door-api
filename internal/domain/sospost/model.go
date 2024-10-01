package sospost

import (
	"encoding/json"
	"log"
	"time"

	utils "github.com/pet-sitter/pets-next-door-api/internal/common"
	databasegen "github.com/pet-sitter/pets-next-door-api/internal/infra/database/gen"

	"github.com/google/uuid"

	"github.com/pet-sitter/pets-next-door-api/internal/domain/soscondition"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/media"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"
)

type (
	CareType    string
	CarerGender string
	RewardType  string
)

const (
	CareTypeFoster   CareType = "foster"
	CareTypeVisiting CareType = "visiting"
)

const (
	CarerGenderMale   CarerGender = "male"
	CarerGenderFemale CarerGender = "female"
	CarerGenderAll    CarerGender = "all"
)

const (
	RewardTypeFee        RewardType = "fee"
	RewardTypeGifticon   RewardType = "gifticon"
	RewardTypeNegotiable RewardType = "negotiable"
)

const (
	JSONNullString = "null"
	JSONEmptyArray = "[]"
)

func (c *CareType) String() string {
	return string(*c)
}

func (c *CarerGender) String() string {
	return string(*c)
}

func (r *RewardType) String() string {
	return string(*r)
}

type SOSPost struct {
	ID          uuid.UUID     `field:"id"`
	AuthorID    uuid.UUID     `field:"author_id"`
	Title       string        `field:"title"`
	Content     string        `field:"content"`
	Reward      string        `field:"reward"`
	CareType    CareType      `field:"care_type"`
	CarerGender CarerGender   `field:"carer_gender"`
	RewardType  RewardType    `field:"reward_type"`
	ThumbnailID uuid.NullUUID `field:"thumbnail_id"`
	CreatedAt   time.Time     `field:"created_at"`
	UpdatedAt   time.Time     `field:"updated_at"`
	DeletedAt   time.Time     `field:"deleted_at"`
}

type SOSPostList struct {
	*pnd.PaginatedView[SOSPost]
}

type SOSPostInfo struct {
	ID          uuid.UUID                       `field:"id" json:"id"`
	AuthorID    uuid.UUID                       `field:"author" json:"author"`
	Title       string                          `field:"title" json:"title"`
	Content     string                          `field:"content" json:"content"`
	Media       media.ViewListForSOSPost        `field:"media" json:"media"`
	Conditions  soscondition.ViewListForSOSPost `field:"conditions" json:"conditions"`
	Pets        pet.ViewListForSOSPost          `field:"pets" json:"pets"`
	Reward      string                          `field:"reward" json:"reward"`
	Dates       SOSDatesList                    `field:"dates" json:"dates"`
	CareType    CareType                        `field:"careType" json:"careType"`
	CarerGender CarerGender                     `field:"carerGender" json:"carerGender"`
	RewardType  RewardType                      `field:"rewardType" json:"rewardType"`
	ThumbnailID uuid.NullUUID                   `field:"thumbnailId" json:"thumbnailId"`
	CreatedAt   time.Time                       `field:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time                       `field:"updatedAt" json:"updatedAt"`
	DeletedAt   time.Time                       `field:"deletedAt" json:"deletedAt"`
}

type SOSPostInfoList struct {
	*pnd.PaginatedView[SOSPostInfo]
}

func ToInfoFromFindRow(row databasegen.FindSOSPostsRow) *SOSPostInfo {
	return &SOSPostInfo{
		ID:          row.ID,
		AuthorID:    row.AuthorID,
		Title:       utils.NullStrToStr(row.Title),
		Content:     utils.NullStrToStr(row.Content),
		Media:       ParseMediaList(row.MediaInfo.RawMessage),
		Conditions:  ParseConditionsList(row.ConditionsInfo.RawMessage),
		Pets:        ParsePetsList(row.PetsInfo.RawMessage),
		Reward:      utils.NullStrToStr(row.Reward),
		Dates:       ParseSOSDatesList(row.Dates),
		CareType:    CareType(row.CareType.String),
		CarerGender: CarerGender(row.CarerGender.String),
		RewardType:  RewardType(row.RewardType.String),
		ThumbnailID: row.ThumbnailID,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}

func ToInfoListFromFindRow(rows []databasegen.FindSOSPostsRow, page, size int) *SOSPostInfoList {
	sl := NewSOSPostInfoList(page, size)
	for _, row := range rows {
		sl.Items = append(sl.Items, *ToInfoFromFindRow(row))
	}
	sl.CalcLastPage()
	return sl
}

func ToInfoFromFindAuthorIDRow(row databasegen.FindSOSPostsByAuthorIDRow) *SOSPostInfo {
	return &SOSPostInfo{
		ID:          row.ID,
		AuthorID:    row.AuthorID,
		Title:       utils.NullStrToStr(row.Title),
		Content:     utils.NullStrToStr(row.Content),
		Media:       ParseMediaList(row.MediaInfo.RawMessage),
		Conditions:  ParseConditionsList(row.ConditionsInfo.RawMessage),
		Pets:        ParsePetsList(row.PetsInfo.RawMessage),
		Reward:      utils.NullStrToStr(row.Reward),
		Dates:       ParseSOSDatesList(row.Dates),
		CareType:    CareType(row.CareType.String),
		CarerGender: CarerGender(row.CarerGender.String),
		RewardType:  RewardType(row.RewardType.String),
		ThumbnailID: row.ThumbnailID,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}

func ToInfoListFromFindAuthorIDRow(rows []databasegen.FindSOSPostsByAuthorIDRow, page, size int) *SOSPostInfoList {
	sl := NewSOSPostInfoList(page, size)
	for _, row := range rows {
		sl.Items = append(sl.Items, *ToInfoFromFindAuthorIDRow(row))
	}

	sl.CalcLastPage()
	return sl
}

func ToInfoFromFindByIDRow(row databasegen.FindSOSPostByIDRow) *SOSPostInfo {
	return &SOSPostInfo{
		ID:          row.ID,
		AuthorID:    row.AuthorID,
		Title:       utils.NullStrToStr(row.Title),
		Content:     utils.NullStrToStr(row.Content),
		Media:       ParseMediaList(row.MediaInfo.RawMessage),
		Conditions:  ParseConditionsList(row.ConditionsInfo.RawMessage),
		Pets:        ParsePetsList(row.PetsInfo.RawMessage),
		Reward:      utils.NullStrToStr(row.Reward),
		Dates:       ParseSOSDatesList(row.Dates),
		CareType:    CareType(row.CareType.String),
		CarerGender: CarerGender(row.CarerGender.String),
		RewardType:  RewardType(row.RewardType.String),
		ThumbnailID: row.ThumbnailID,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}

func NewSOSPostList(page, size int) *SOSPostList {
	return &SOSPostList{PaginatedView: pnd.NewPaginatedView(
		page, size, false, make([]SOSPost, 0),
	)}
}

func NewSOSPostInfoList(page, size int) *SOSPostInfoList {
	return &SOSPostInfoList{PaginatedView: pnd.NewPaginatedView(
		page, size, false, make([]SOSPostInfo, 0),
	)}
}

func ParseMediaList(rows json.RawMessage) media.ViewListForSOSPost {
	var mediaList media.ViewListForSOSPost
	if len(rows) == 0 || string(rows) == JSONNullString || string(rows) == JSONEmptyArray {
		return mediaList
	}

	if err := json.Unmarshal(rows, &mediaList); err != nil {
		log.Println("Error unmarshalling media:", err)
		return mediaList
	}

	return mediaList
}

func ParseConditionsList(rows json.RawMessage) soscondition.ViewListForSOSPost {
	var conditionsList soscondition.ViewListForSOSPost
	if len(rows) == 0 || string(rows) == JSONNullString || string(rows) == JSONEmptyArray {
		return conditionsList
	}

	if err := json.Unmarshal(rows, &conditionsList); err != nil {
		log.Println("Error unmarshalling conditions:", err)
		return conditionsList
	}

	return conditionsList
}

func ParsePetsList(rows json.RawMessage) pet.ViewListForSOSPost {
	var petList pet.ViewListForSOSPost

	if len(rows) == 0 || string(rows) == JSONNullString || string(rows) == JSONEmptyArray {
		return petList
	}

	if err := json.Unmarshal(rows, &petList); err != nil {
		log.Println("Error unmarshalling pets:", err)
		return petList
	}

	return petList
}

func ParseSOSDatesList(rows json.RawMessage) SOSDatesList {
	var sosDatesList SOSDatesList

	if len(rows) == 0 || string(rows) == JSONNullString || string(rows) == JSONEmptyArray {
		return sosDatesList
	}
	if err := json.Unmarshal(rows, &sosDatesList); err != nil {
		log.Println("Error unmarshalling sosDates:", err)
		return sosDatesList
	}

	return sosDatesList
}

type SOSDates struct {
	ID          uuid.UUID `field:"id" json:"id"`
	DateStartAt string    `field:"date_start_at" json:"date_start_at"`
	DateEndAt   string    `field:"date_end_at" json:"date_end_at"`
	CreatedAt   time.Time `field:"created_at" json:"created_at"`
	UpdatedAt   time.Time `field:"updated_at" json:"updated_at"`
	DeletedAt   time.Time `field:"deleted_at" json:"deleted_at"`
}

type SOSDatesList []*SOSDates
