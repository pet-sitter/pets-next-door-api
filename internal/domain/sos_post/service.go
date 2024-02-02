package sos_post

import (
	"time"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/media"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/user"
)

type SosPostService struct {
	sosPostStore       SosPostStore
	resourceMediaStore media.ResourceMediaStore
	userStore          user.UserStore
}

func NewSosPostService(sosPostStore SosPostStore, resourceMediaStore media.ResourceMediaStore, userStore user.UserStore) *SosPostService {
	return &SosPostService{
		sosPostStore:       sosPostStore,
		resourceMediaStore: resourceMediaStore,
		userStore:          userStore,
	}
}

func (service *SosPostService) WriteSosPost(fbUid string, request *WriteSosPostRequest) (*WriteSosPostResponse, *pnd.AppError) {
	userID, err := service.userStore.FindUserIDByFbUID(fbUid)
	if err != nil {
		return nil, err
	}

	utcDateStart := request.DateStartAt.UTC().Format(time.RFC3339)
	utcDateEnd := request.DateEndAt.UTC().Format(time.RFC3339)

	sosPost, err := service.sosPostStore.WriteSosPost(userID, utcDateStart, utcDateEnd, request)
	if err != nil {
		return nil, err
	}

	mediaData, err := service.resourceMediaStore.FindResourceMediaByResourceID(sosPost.ID, string(media.SosResourceType))
	if err != nil {
		return nil, err
	}

	var mediaView []media.MediaView

	for _, m := range mediaData {
		view := media.MediaView{
			ID:        m.ID,
			MediaType: m.MediaType,
			URL:       m.URL,
			CreatedAt: m.CreatedAt,
		}
		mediaView = append(mediaView, view)
	}

	conditions, err := service.sosPostStore.FindConditionByID(sosPost.ID)
	if err != nil {
		return nil, err
	}

	var conditionsView []ConditionView
	for _, c := range conditions {
		view := ConditionView{
			ID:   c.ID,
			Name: c.Name,
		}

		conditionsView = append(conditionsView, view)
	}

	pets, err := service.sosPostStore.FindPetsByID(sosPost.ID)
	if err != nil {
		return nil, err
	}

	var petsView []pet.PetView
	for _, p := range pets {
		p := pet.PetView{
			ID:         p.ID,
			Name:       p.Name,
			PetType:    p.PetType,
			Sex:        p.Sex,
			Neutered:   p.Neutered,
			Breed:      p.Breed,
			BirthDate:  datetimeToDay(p.BirthDate),
			WeightInKg: p.WeightInKg,
		}
		petsView = append(petsView, p)
	}

	return &WriteSosPostResponse{
		ID:           sosPost.ID,
		AuthorID:     sosPost.AuthorID,
		Title:        sosPost.Title,
		Content:      sosPost.Content,
		Media:        mediaView,
		Conditions:   conditionsView,
		Pets:         petsView,
		Reward:       sosPost.Reward,
		DateStartAt:  sosPost.DateStartAt,
		DateEndAt:    sosPost.DateEndAt,
		TimeStartAt:  datetimeToTime(sosPost.TimeStartAt),
		TimeEndAt:    datetimeToTime(sosPost.TimeEndAt),
		CareType:     sosPost.CareType,
		CarerGender:  sosPost.CarerGender,
		RewardAmount: sosPost.RewardAmount,
		ThumbnailID:  sosPost.ThumbnailID,
		CreatedAt:    sosPost.CreatedAt,
		UpdatedAt:    sosPost.UpdatedAt,
	}, nil
}

func (service *SosPostService) FindSosPosts(page int, size int, sortBy string) ([]FindSosPostResponse, *pnd.AppError) {
	sosPosts, err := service.sosPostStore.FindSosPosts(page, size, sortBy)
	if err != nil {
		return nil, err
	}

	var FindSosPostResponseList []FindSosPostResponse

	for _, sosPost := range sosPosts {
		mediaData, err := service.resourceMediaStore.FindResourceMediaByResourceID(sosPost.ID, string(media.SosResourceType))
		if err != nil {
			return nil, err
		}

		var mediaView []media.MediaView

		for _, m := range mediaData {
			view := media.MediaView{
				ID:        m.ID,
				MediaType: m.MediaType,
				URL:       m.URL,
				CreatedAt: m.CreatedAt,
			}
			mediaView = append(mediaView, view)
		}

		conditions, err := service.sosPostStore.FindConditionByID(sosPost.ID)
		if err != nil {
			return nil, err
		}

		var conditionsView []ConditionView
		for _, c := range conditions {
			view := ConditionView{
				ID:   c.ID,
				Name: c.Name,
			}

			conditionsView = append(conditionsView, view)
		}

		pets, err := service.sosPostStore.FindPetsByID(sosPost.ID)
		if err != nil {
			return nil, err
		}

		var petsView []pet.PetView
		for _, p := range pets {
			p := pet.PetView{
				ID:         p.ID,
				Name:       p.Name,
				PetType:    p.PetType,
				Sex:        p.Sex,
				Neutered:   p.Neutered,
				Breed:      p.Breed,
				BirthDate:  datetimeToDay(p.BirthDate),
				WeightInKg: p.WeightInKg,
			}
			petsView = append(petsView, p)
		}

		findByAuthorSosPostResponse := &FindSosPostResponse{
			ID:           sosPost.ID,
			AuthorID:     sosPost.AuthorID,
			Title:        sosPost.Title,
			Content:      sosPost.Content,
			Media:        mediaView,
			Conditions:   conditionsView,
			Pets:         petsView,
			Reward:       sosPost.Reward,
			DateStartAt:  sosPost.DateStartAt,
			DateEndAt:    sosPost.DateEndAt,
			TimeStartAt:  datetimeToTime(sosPost.TimeStartAt),
			TimeEndAt:    datetimeToTime(sosPost.TimeEndAt),
			CareType:     sosPost.CareType,
			CarerGender:  sosPost.CarerGender,
			RewardAmount: sosPost.RewardAmount,
			ThumbnailID:  sosPost.ThumbnailID,
			CreatedAt:    sosPost.CreatedAt,
			UpdatedAt:    sosPost.UpdatedAt,
		}

		FindSosPostResponseList = append(FindSosPostResponseList, *findByAuthorSosPostResponse)
	}

	return FindSosPostResponseList, nil
}

func (service *SosPostService) FindSosPostsByAuthorID(authorID int, page int, size int) ([]FindSosPostResponse, *pnd.AppError) {
	sosPosts, err := service.sosPostStore.FindSosPostsByAuthorID(authorID, page, size)
	if err != nil {
		return nil, err
	}

	var FindSosPostResponseList []FindSosPostResponse

	for _, sosPost := range sosPosts {
		mediaData, err := service.resourceMediaStore.FindResourceMediaByResourceID(sosPost.ID, string(media.SosResourceType))
		if err != nil {
			return nil, err
		}

		var mediaView []media.MediaView

		for _, m := range mediaData {
			view := media.MediaView{
				ID:        m.ID,
				MediaType: m.MediaType,
				URL:       m.URL,
				CreatedAt: m.CreatedAt,
			}
			mediaView = append(mediaView, view)
		}

		conditions, err := service.sosPostStore.FindConditionByID(sosPost.ID)
		if err != nil {
			return nil, err
		}

		var conditionsView []ConditionView
		for _, c := range conditions {
			view := ConditionView{
				ID:   c.ID,
				Name: c.Name,
			}

			conditionsView = append(conditionsView, view)
		}

		pets, err := service.sosPostStore.FindPetsByID(sosPost.ID)
		if err != nil {
			return nil, err
		}

		var petsView []pet.PetView
		for _, p := range pets {
			p := pet.PetView{
				ID:         p.ID,
				Name:       p.Name,
				PetType:    p.PetType,
				Sex:        p.Sex,
				Neutered:   p.Neutered,
				Breed:      p.Breed,
				BirthDate:  datetimeToDay(p.BirthDate),
				WeightInKg: p.WeightInKg,
			}
			petsView = append(petsView, p)
		}

		findByAuthorSosPostResponse := &FindSosPostResponse{
			ID:           sosPost.ID,
			AuthorID:     sosPost.AuthorID,
			Title:        sosPost.Title,
			Content:      sosPost.Content,
			Media:        mediaView,
			Conditions:   conditionsView,
			Pets:         petsView,
			Reward:       sosPost.Reward,
			DateStartAt:  sosPost.DateStartAt,
			DateEndAt:    sosPost.DateEndAt,
			TimeStartAt:  datetimeToTime(sosPost.TimeStartAt),
			TimeEndAt:    datetimeToTime(sosPost.TimeEndAt),
			CareType:     sosPost.CareType,
			CarerGender:  sosPost.CarerGender,
			RewardAmount: sosPost.RewardAmount,
			ThumbnailID:  sosPost.ThumbnailID,
			CreatedAt:    sosPost.CreatedAt,
			UpdatedAt:    sosPost.UpdatedAt,
		}

		FindSosPostResponseList = append(FindSosPostResponseList, *findByAuthorSosPostResponse)
	}

	return FindSosPostResponseList, nil
}

func (service *SosPostService) FindSosPostByID(id int) (*FindSosPostResponse, *pnd.AppError) {
	sosPost, err := service.sosPostStore.FindSosPostByID(id)
	if err != nil {
		return nil, err
	}

	mediaData, err := service.resourceMediaStore.FindResourceMediaByResourceID(sosPost.ID, string(media.SosResourceType))
	if err != nil {
		return nil, err
	}

	var mediaView []media.MediaView

	for _, m := range mediaData {
		view := media.MediaView{
			ID:        m.ID,
			MediaType: m.MediaType,
			URL:       m.URL,
			CreatedAt: m.CreatedAt,
		}
		mediaView = append(mediaView, view)
	}

	conditions, err := service.sosPostStore.FindConditionByID(sosPost.ID)
	if err != nil {
		return nil, err
	}

	var conditionsView []ConditionView
	for _, c := range conditions {
		view := ConditionView{
			ID:   c.ID,
			Name: c.Name,
		}

		conditionsView = append(conditionsView, view)
	}

	pets, err := service.sosPostStore.FindPetsByID(sosPost.ID)
	if err != nil {
		return nil, err
	}

	var petsView []pet.PetView
	for _, p := range pets {
		p := pet.PetView{
			ID:         p.ID,
			Name:       p.Name,
			PetType:    p.PetType,
			Sex:        p.Sex,
			Neutered:   p.Neutered,
			Breed:      p.Breed,
			BirthDate:  datetimeToDay(p.BirthDate),
			WeightInKg: p.WeightInKg,
		}
		petsView = append(petsView, p)
	}

	return &FindSosPostResponse{
		ID:           sosPost.ID,
		AuthorID:     sosPost.AuthorID,
		Title:        sosPost.Title,
		Content:      sosPost.Content,
		Media:        mediaView,
		Conditions:   conditionsView,
		Pets:         petsView,
		Reward:       sosPost.Reward,
		DateStartAt:  sosPost.DateStartAt,
		DateEndAt:    sosPost.DateEndAt,
		TimeStartAt:  datetimeToTime(sosPost.TimeStartAt),
		TimeEndAt:    datetimeToTime(sosPost.TimeEndAt),
		CareType:     sosPost.CareType,
		CarerGender:  sosPost.CarerGender,
		RewardAmount: sosPost.RewardAmount,
		ThumbnailID:  sosPost.ThumbnailID,
		CreatedAt:    sosPost.CreatedAt,
		UpdatedAt:    sosPost.UpdatedAt,
	}, nil
}

func (service *SosPostService) UpdateSosPost(request *UpdateSosPostRequest) (*UpdateSosPostResponse, *pnd.AppError) {
	updateSosPost, err := service.sosPostStore.UpdateSosPost(request)
	if err != nil {
		return nil, err
	}

	mediaData, err := service.resourceMediaStore.FindResourceMediaByResourceID(updateSosPost.ID, string(media.SosResourceType))
	if err != nil {
		return nil, err
	}

	var mediaView []media.MediaView

	for _, m := range mediaData {
		view := media.MediaView{
			ID:        m.ID,
			MediaType: m.MediaType,
			URL:       m.URL,
			CreatedAt: m.CreatedAt,
		}
		mediaView = append(mediaView, view)
	}

	conditions, err := service.sosPostStore.FindConditionByID(updateSosPost.ID)
	if err != nil {
		return nil, err
	}

	var conditionsView []ConditionView
	for _, c := range conditions {
		view := ConditionView{
			ID:   c.ID,
			Name: c.Name,
		}

		conditionsView = append(conditionsView, view)
	}

	pets, err := service.sosPostStore.FindPetsByID(updateSosPost.ID)
	if err != nil {
		return nil, err
	}

	var petsView []pet.PetView
	for _, p := range pets {
		p := pet.PetView{
			ID:         p.ID,
			Name:       p.Name,
			PetType:    p.PetType,
			Sex:        p.Sex,
			Neutered:   p.Neutered,
			Breed:      p.Breed,
			BirthDate:  datetimeToDay(p.BirthDate),
			WeightInKg: p.WeightInKg,
		}
		petsView = append(petsView, p)
	}

	return &UpdateSosPostResponse{
		ID:           updateSosPost.ID,
		AuthorID:     updateSosPost.AuthorID,
		Title:        updateSosPost.Title,
		Content:      updateSosPost.Content,
		Media:        mediaView,
		Conditions:   conditionsView,
		Pets:         petsView,
		Reward:       updateSosPost.Reward,
		DateStartAt:  updateSosPost.DateStartAt,
		DateEndAt:    updateSosPost.DateEndAt,
		TimeStartAt:  datetimeToTime(updateSosPost.TimeStartAt),
		TimeEndAt:    datetimeToTime(updateSosPost.TimeEndAt),
		CareType:     updateSosPost.CareType,
		CarerGender:  updateSosPost.CarerGender,
		RewardAmount: updateSosPost.RewardAmount,
		ThumbnailID:  updateSosPost.ThumbnailID,
		CreatedAt:    updateSosPost.CreatedAt,
		UpdatedAt:    updateSosPost.UpdatedAt,
	}, nil
}

func (service *SosPostService) CheckUpdatePermission(fbUid string, sosPostID int) (bool, *pnd.AppError) {
	userID, err := service.userStore.FindUserIDByFbUID(fbUid)
	if err != nil {
		return false, pnd.ErrUnknown(err)
	}

	sosPost, err := service.sosPostStore.FindSosPostByID(sosPostID)
	if err != nil {
		return false, pnd.ErrUnknown(err)
	}

	return userID == sosPost.AuthorID, nil
}

func datetimeToTime(datetime time.Time) string {
	return datetime.Format("15:04")
}

func datetimeToDay(datetime string) string {
	return datetime[:10]
}
