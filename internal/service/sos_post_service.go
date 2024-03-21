package service

import (
	"context"
	utils "github.com/pet-sitter/pets-next-door-api/internal/common"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
	"github.com/pet-sitter/pets-next-door-api/internal/postgres"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/media"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/sos_post"
)

type SosPostService struct {
	conn *database.DB
}

func NewSosPostService(conn *database.DB) *SosPostService {
	return &SosPostService{
		conn: conn,
	}
}

func (service *SosPostService) WriteSosPost(ctx context.Context, fbUid string, request *sos_post.WriteSosPostRequest) (*sos_post.WriteSosPostView, *pnd.AppError) {
	tx, err := service.conn.BeginTx(ctx)
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}

	userID, err := postgres.FindUserIDByFbUID(ctx, tx, fbUid)
	if err != nil {
		return nil, err
	}

	sosPost, err := postgres.WriteSosPost(ctx, tx, userID, request)
	if err != nil {
		return nil, err
	}

	mediaData, err := postgres.FindResourceMediaByResourceID(ctx, tx, sosPost.ID, string(media.SosResourceType))
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

	conditions, err := postgres.FindConditionByID(ctx, tx, sosPost.ID)
	if err != nil {
		return nil, err
	}

	var conditionsView []sos_post.ConditionView
	for _, c := range conditions {
		view := sos_post.ConditionView{
			ID:   c.ID,
			Name: c.Name,
		}

		conditionsView = append(conditionsView, view)
	}

	pets, err := postgres.FindPetsByID(ctx, tx, sosPost.ID)
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
			BirthDate:  utils.FormatDate(p.BirthDate),
			WeightInKg: p.WeightInKg,
		}
		petsView = append(petsView, p)
	}

	dates, err := postgres.FindDatesBySosPostID(ctx, tx, sosPost.ID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	var sosDatesView []sos_post.SosDateView
	for _, d := range dates {
		d := sos_post.SosDateView{
			DateStartAt: utils.FormatDate(d.DateStartAt),
			DateEndAt:   utils.FormatDate(d.DateEndAt),
		}
		sosDatesView = append(sosDatesView, d)
	}

	return &sos_post.WriteSosPostView{
		ID:           sosPost.ID,
		AuthorID:     sosPost.AuthorID,
		Title:        sosPost.Title,
		Content:      sosPost.Content,
		Media:        mediaView,
		Conditions:   conditionsView,
		Pets:         petsView,
		Dates:        sosDatesView,
		Reward:       sosPost.Reward,
		CareType:     sosPost.CareType,
		CarerGender:  sosPost.CarerGender,
		RewardAmount: sosPost.RewardAmount,
		ThumbnailID:  sosPost.ThumbnailID,
		CreatedAt:    sosPost.CreatedAt,
		UpdatedAt:    sosPost.UpdatedAt,
	}, nil
}

func (service *SosPostService) FindSosPosts(ctx context.Context, page int, size int, sortBy string) (*sos_post.FindSosPostListView, *pnd.AppError) {
	tx, err := service.conn.BeginTx(ctx)
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}

	sosPosts, err := postgres.FindSosPosts(ctx, tx, page, size, sortBy)
	if err != nil {
		return nil, err
	}

	sosPostViews := sos_post.FromEmptySosPostList(sosPosts)
	for _, sosPost := range sosPosts.Items {
		mediaData, err := postgres.FindResourceMediaByResourceID(ctx, tx, sosPost.ID, string(media.SosResourceType))
		if err != nil {
			return nil, err
		}

		mediaView := make([]media.MediaView, 0)
		for _, m := range mediaData {
			view := media.MediaView{
				ID:        m.ID,
				MediaType: m.MediaType,
				URL:       m.URL,
				CreatedAt: m.CreatedAt,
			}
			mediaView = append(mediaView, view)
		}

		conditions, err := postgres.FindConditionByID(ctx, tx, sosPost.ID)
		if err != nil {
			return nil, err
		}

		pets, err := postgres.FindPetsByID(ctx, tx, sosPost.ID)
		if err != nil {
			return nil, err
		}

		petsView := make([]pet.PetView, 0)
		for _, p := range pets {
			p := pet.PetView{
				ID:         p.ID,
				Name:       p.Name,
				PetType:    p.PetType,
				Sex:        p.Sex,
				Neutered:   p.Neutered,
				Breed:      p.Breed,
				BirthDate:  utils.FormatDate(p.BirthDate),
				WeightInKg: p.WeightInKg,
			}
			petsView = append(petsView, p)
		}

		dates, err := postgres.FindDatesBySosPostID(ctx, tx, sosPost.ID)
		if err != nil {
			return nil, err
		}

		author, err := postgres.FindUserByID(ctx, tx, sosPost.AuthorID, true)
		if err != nil {
			return nil, err
		}

		conditionsView := make([]sos_post.ConditionView, 0)
		for _, c := range conditions {
			view := sos_post.ConditionView{
				ID:   c.ID,
				Name: c.Name,
			}

			conditionsView = append(conditionsView, view)
		}

		sosDatesView := make([]sos_post.SosDateView, 0)
		for _, d := range dates {
			d := sos_post.SosDateView{
				DateStartAt: utils.FormatDate(d.DateStartAt),
				DateEndAt:   utils.FormatDate(d.DateEndAt),
			}
			sosDatesView = append(sosDatesView, d)
		}

		findByAuthorSosPostView := &sos_post.FindSosPostView{
			ID:           sosPost.ID,
			Author:       author.ToUserWithoutPrivateInfo(),
			Title:        sosPost.Title,
			Content:      sosPost.Content,
			Media:        mediaView,
			Conditions:   conditionsView,
			Pets:         petsView,
			Dates:        sosDatesView,
			Reward:       sosPost.Reward,
			CareType:     sosPost.CareType,
			CarerGender:  sosPost.CarerGender,
			RewardAmount: sosPost.RewardAmount,
			ThumbnailID:  sosPost.ThumbnailID,
			CreatedAt:    sosPost.CreatedAt,
			UpdatedAt:    sosPost.UpdatedAt,
		}

		sosPostViews.Items = append(sosPostViews.Items, *findByAuthorSosPostView)
	}

	return sosPostViews, nil
}

func (service *SosPostService) FindSosPostsByAuthorID(ctx context.Context, authorID int, page int, size int, sortBy string) (*sos_post.FindSosPostListView, *pnd.AppError) {
	tx, err := service.conn.BeginTx(ctx)
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}

	sosPosts, err := postgres.FindSosPostsByAuthorID(ctx, tx, authorID, page, size, sortBy)
	if err != nil {
		return nil, err
	}

	sosPostViews := sos_post.FromEmptySosPostList(sosPosts)
	for _, sosPost := range sosPosts.Items {
		mediaData, err := postgres.FindResourceMediaByResourceID(ctx, tx, sosPost.ID, string(media.SosResourceType))
		if err != nil {
			return nil, err
		}

		conditions, err := postgres.FindConditionByID(ctx, tx, sosPost.ID)
		if err != nil {
			return nil, err
		}

		pets, err := postgres.FindPetsByID(ctx, tx, sosPost.ID)
		if err != nil {
			return nil, err
		}

		dates, err := postgres.FindDatesBySosPostID(ctx, tx, sosPost.ID)
		if err != nil {
			return nil, err
		}

		author, err := postgres.FindUserByID(ctx, tx, sosPost.AuthorID, true)
		if err != nil {
			return nil, err
		}

		var conditionsView []sos_post.ConditionView
		for _, c := range conditions {
			view := sos_post.ConditionView{
				ID:   c.ID,
				Name: c.Name,
			}
			conditionsView = append(conditionsView, view)
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

		var petsView []pet.PetView
		for _, p := range pets {
			p := pet.PetView{
				ID:         p.ID,
				Name:       p.Name,
				PetType:    p.PetType,
				Sex:        p.Sex,
				Neutered:   p.Neutered,
				Breed:      p.Breed,
				BirthDate:  utils.FormatDate(p.BirthDate),
				WeightInKg: p.WeightInKg,
			}
			petsView = append(petsView, p)
		}

		sosDatesView := make([]sos_post.SosDateView, 0)
		for _, d := range dates {
			d := sos_post.SosDateView{
				DateStartAt: utils.FormatDate(d.DateStartAt),
				DateEndAt:   utils.FormatDate(d.DateEndAt),
			}
			sosDatesView = append(sosDatesView, d)
		}

		findByAuthorSosPostView := &sos_post.FindSosPostView{
			ID:           sosPost.ID,
			Author:       author.ToUserWithoutPrivateInfo(),
			Title:        sosPost.Title,
			Content:      sosPost.Content,
			Media:        mediaView,
			Conditions:   conditionsView,
			Pets:         petsView,
			Dates:        sosDatesView,
			Reward:       sosPost.Reward,
			CareType:     sosPost.CareType,
			CarerGender:  sosPost.CarerGender,
			RewardAmount: sosPost.RewardAmount,
			ThumbnailID:  sosPost.ThumbnailID,
			CreatedAt:    sosPost.CreatedAt,
			UpdatedAt:    sosPost.UpdatedAt,
		}

		sosPostViews.Items = append(sosPostViews.Items, *findByAuthorSosPostView)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return sosPostViews, nil
}

func (service *SosPostService) FindSosPostByID(ctx context.Context, id int) (*sos_post.FindSosPostView, *pnd.AppError) {
	tx, err := service.conn.BeginTx(ctx)
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}

	sosPost, err := postgres.FindSosPostByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	mediaData, err := postgres.FindResourceMediaByResourceID(ctx, tx, sosPost.ID, string(media.SosResourceType))
	if err != nil {
		return nil, err
	}

	conditions, err := postgres.FindConditionByID(ctx, tx, sosPost.ID)
	if err != nil {
		return nil, err
	}

	pets, err := postgres.FindPetsByID(ctx, tx, sosPost.ID)
	if err != nil {
		return nil, err
	}

	dates, err := postgres.FindDatesBySosPostID(ctx, tx, sosPost.ID)
	if err != nil {
		return nil, err
	}

	author, err := postgres.FindUserByID(ctx, tx, sosPost.AuthorID, true)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
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

	var conditionsView []sos_post.ConditionView
	for _, c := range conditions {
		view := sos_post.ConditionView{
			ID:   c.ID,
			Name: c.Name,
		}

		conditionsView = append(conditionsView, view)
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
			BirthDate:  utils.FormatDate(p.BirthDate),
			WeightInKg: p.WeightInKg,
		}
		petsView = append(petsView, p)
	}

	sosDatesView := make([]sos_post.SosDateView, 0)
	for _, d := range dates {
		d := sos_post.SosDateView{
			DateStartAt: utils.FormatDate(d.DateStartAt),
			DateEndAt:   utils.FormatDate(d.DateEndAt),
		}
		sosDatesView = append(sosDatesView, d)
	}

	return &sos_post.FindSosPostView{
		ID:           sosPost.ID,
		Author:       author.ToUserWithoutPrivateInfo(),
		Title:        sosPost.Title,
		Content:      sosPost.Content,
		Media:        mediaView,
		Conditions:   conditionsView,
		Pets:         petsView,
		Dates:        sosDatesView,
		Reward:       sosPost.Reward,
		CareType:     sosPost.CareType,
		CarerGender:  sosPost.CarerGender,
		RewardAmount: sosPost.RewardAmount,
		ThumbnailID:  sosPost.ThumbnailID,
		CreatedAt:    sosPost.CreatedAt,
		UpdatedAt:    sosPost.UpdatedAt,
	}, nil
}

func (service *SosPostService) UpdateSosPost(ctx context.Context, request *sos_post.UpdateSosPostRequest) (*sos_post.UpdateSosPostView, *pnd.AppError) {
	tx, err := service.conn.BeginTx(ctx)
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}

	updateSosPost, err := postgres.UpdateSosPost(ctx, tx, request)
	if err != nil {
		return nil, err
	}

	mediaData, err := postgres.FindResourceMediaByResourceID(ctx, tx, updateSosPost.ID, string(media.SosResourceType))
	if err != nil {
		return nil, err
	}

	conditions, err := postgres.FindConditionByID(ctx, tx, updateSosPost.ID)
	if err != nil {
		return nil, err
	}

	pets, err := postgres.FindPetsByID(ctx, tx, updateSosPost.ID)
	if err != nil {
		return nil, err
	}

	dates, err := postgres.FindDatesBySosPostID(ctx, tx, request.ID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
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

	var conditionsView []sos_post.ConditionView
	for _, c := range conditions {
		view := sos_post.ConditionView{
			ID:   c.ID,
			Name: c.Name,
		}

		conditionsView = append(conditionsView, view)
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
			BirthDate:  utils.FormatDate(p.BirthDate),
			WeightInKg: p.WeightInKg,
		}
		petsView = append(petsView, p)
	}

	sosDatesView := make([]sos_post.SosDateView, 0)
	for _, d := range dates {
		d := sos_post.SosDateView{
			DateStartAt: utils.FormatDate(d.DateStartAt),
			DateEndAt:   utils.FormatDate(d.DateEndAt),
		}
		sosDatesView = append(sosDatesView, d)
	}

	return &sos_post.UpdateSosPostView{
		ID:           updateSosPost.ID,
		AuthorID:     updateSosPost.AuthorID,
		Title:        updateSosPost.Title,
		Content:      updateSosPost.Content,
		Media:        mediaView,
		Conditions:   conditionsView,
		Pets:         petsView,
		Dates:        sosDatesView,
		Reward:       updateSosPost.Reward,
		CareType:     updateSosPost.CareType,
		CarerGender:  updateSosPost.CarerGender,
		RewardAmount: updateSosPost.RewardAmount,
		ThumbnailID:  updateSosPost.ThumbnailID,
		CreatedAt:    updateSosPost.CreatedAt,
		UpdatedAt:    updateSosPost.UpdatedAt,
	}, nil
}

func (service *SosPostService) CheckUpdatePermission(ctx context.Context, fbUid string, sosPostID int) (bool, *pnd.AppError) {
	tx, err := service.conn.BeginTx(ctx)
	defer tx.Rollback()
	if err != nil {
		return false, err
	}

	userID, err := postgres.FindUserIDByFbUID(ctx, tx, fbUid)
	if err != nil {
		return false, err
	}

	sosPost, err := postgres.FindSosPostByID(ctx, tx, sosPostID)
	if err != nil {
		return false, err
	}

	if err := tx.Commit(); err != nil {
		return false, err
	}

	return userID == sosPost.AuthorID, nil
}
