package service

import (
	"context"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/user"
	"time"

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
	var sosPostView *sos_post.WriteSosPostView

	err := database.WithTransaction(ctx, service.conn, func(tx *database.Tx) *pnd.AppError {
		userStore := postgres.NewUserPostgresStore(tx)
		sosPostStore := postgres.NewSosPostPostgresStore(tx)
		resourceMediaStore := postgres.NewResourceMediaPostgresStore(tx)

		userID, err := userStore.FindUserIDByFbUID(ctx, fbUid)
		if err != nil {
			return err
		}

		utcDateStart := request.DateStartAt.UTC().Format(time.RFC3339)
		utcDateEnd := request.DateEndAt.UTC().Format(time.RFC3339)

		sosPost, err := sosPostStore.WriteSosPost(ctx, userID, utcDateStart, utcDateEnd, request)
		if err != nil {
			return err
		}

		mediaData, err := resourceMediaStore.FindResourceMediaByResourceID(ctx, sosPost.ID, string(media.SosResourceType))
		if err != nil {
			return err
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

		conditions, err := sosPostStore.FindConditionByID(ctx, sosPost.ID)
		if err != nil {
			return err
		}

		var conditionsView []sos_post.ConditionView
		for _, c := range conditions {
			view := sos_post.ConditionView{
				ID:   c.ID,
				Name: c.Name,
			}

			conditionsView = append(conditionsView, view)
		}

		pets, err := sosPostStore.FindPetsByID(ctx, sosPost.ID)
		if err != nil {
			return err
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

		sosPostView = &sos_post.WriteSosPostView{
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
			CareType:     sosPost.CareType,
			CarerGender:  sosPost.CarerGender,
			RewardAmount: sosPost.RewardAmount,
			ThumbnailID:  sosPost.ThumbnailID,
			CreatedAt:    sosPost.CreatedAt,
			UpdatedAt:    sosPost.UpdatedAt,
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return sosPostView, nil
}

func (service *SosPostService) FindSosPosts(ctx context.Context, page int, size int, sortBy string) (*sos_post.FindSosPostListView, *pnd.AppError) {
	var sosPostViews *sos_post.FindSosPostListView

	err := database.WithTransaction(ctx, service.conn, func(tx *database.Tx) *pnd.AppError {
		sosPostStore := postgres.NewSosPostPostgresStore(tx)
		resourceMediaStore := postgres.NewResourceMediaPostgresStore(tx)
		userStore := postgres.NewUserPostgresStore(tx)

		sosPosts, err := sosPostStore.FindSosPosts(ctx, page, size, sortBy)
		if err != nil {
			return err
		}

		sosPostViews = sos_post.FromEmptySosPostList(sosPosts)

		for _, sosPost := range sosPosts.Items {
			mediaData, err := resourceMediaStore.FindResourceMediaByResourceID(ctx, sosPost.ID, string(media.SosResourceType))
			if err != nil {
				return err
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

			conditions, err := sosPostStore.FindConditionByID(ctx, sosPost.ID)
			if err != nil {
				return err
			}

			conditionsView := make([]sos_post.ConditionView, 0)
			for _, c := range conditions {
				view := sos_post.ConditionView{
					ID:   c.ID,
					Name: c.Name,
				}

				conditionsView = append(conditionsView, view)
			}

			pets, err := sosPostStore.FindPetsByID(ctx, sosPost.ID)
			if err != nil {
				return err
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

			author, err := userStore.FindUserByID(ctx, sosPost.AuthorID)

			userView := &user.UserWithoutPrivateInfo{
				ID:              author.ID,
				ProfileImageURL: author.ProfileImageURL,
				Nickname:        author.Nickname,
			}

			findByAuthorSosPostView := &sos_post.FindSosPostView{
				ID:           sosPost.ID,
				Author:       userView,
				Title:        sosPost.Title,
				Content:      sosPost.Content,
				Media:        mediaView,
				Conditions:   conditionsView,
				Pets:         petsView,
				Reward:       sosPost.Reward,
				DateStartAt:  sosPost.DateStartAt,
				DateEndAt:    sosPost.DateEndAt,
				CareType:     sosPost.CareType,
				CarerGender:  sosPost.CarerGender,
				RewardAmount: sosPost.RewardAmount,
				ThumbnailID:  sosPost.ThumbnailID,
				CreatedAt:    sosPost.CreatedAt,
				UpdatedAt:    sosPost.UpdatedAt,
			}

			sosPostViews.Items = append(sosPostViews.Items, *findByAuthorSosPostView)

		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return sosPostViews, nil
}

func (service *SosPostService) FindSosPostsByAuthorID(ctx context.Context, authorID int, page int, size int, sortBy string) (*sos_post.FindSosPostListView, *pnd.AppError) {
	var sosPostViews *sos_post.FindSosPostListView

	err := database.WithTransaction(ctx, service.conn, func(tx *database.Tx) *pnd.AppError {
		sosPostStore := postgres.NewSosPostPostgresStore(tx)
		resourceMediaStore := postgres.NewResourceMediaPostgresStore(tx)
		userStore := postgres.NewUserPostgresStore(tx)

		sosPosts, err := sosPostStore.FindSosPostsByAuthorID(ctx, authorID, page, size, sortBy)
		if err != nil {
			return err
		}

		sosPostViews = sos_post.FromEmptySosPostList(sosPosts)
		for _, sosPost := range sosPosts.Items {
			mediaData, err := resourceMediaStore.FindResourceMediaByResourceID(ctx, sosPost.ID, string(media.SosResourceType))
			if err != nil {
				return err
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

			conditions, err := sosPostStore.FindConditionByID(ctx, sosPost.ID)
			if err != nil {
				return err
			}

			var conditionsView []sos_post.ConditionView
			for _, c := range conditions {
				view := sos_post.ConditionView{
					ID:   c.ID,
					Name: c.Name,
				}

				conditionsView = append(conditionsView, view)
			}

			pets, err := sosPostStore.FindPetsByID(ctx, sosPost.ID)
			if err != nil {
				return err
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

			author, err := userStore.FindUserByID(ctx, sosPost.AuthorID)

			userView := &user.UserWithoutPrivateInfo{
				ID:              author.ID,
				ProfileImageURL: author.ProfileImageURL,
				Nickname:        author.Nickname,
			}

			findByAuthorSosPostView := &sos_post.FindSosPostView{
				ID:           sosPost.ID,
				Author:       userView,
				Title:        sosPost.Title,
				Content:      sosPost.Content,
				Media:        mediaView,
				Conditions:   conditionsView,
				Pets:         petsView,
				Reward:       sosPost.Reward,
				DateStartAt:  sosPost.DateStartAt,
				DateEndAt:    sosPost.DateEndAt,
				CareType:     sosPost.CareType,
				CarerGender:  sosPost.CarerGender,
				RewardAmount: sosPost.RewardAmount,
				ThumbnailID:  sosPost.ThumbnailID,
				CreatedAt:    sosPost.CreatedAt,
				UpdatedAt:    sosPost.UpdatedAt,
			}

			sosPostViews.Items = append(sosPostViews.Items, *findByAuthorSosPostView)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return sosPostViews, nil
}

func (service *SosPostService) FindSosPostByID(ctx context.Context, id int) (*sos_post.FindSosPostView, *pnd.AppError) {
	var sosPostView *sos_post.FindSosPostView

	err := database.WithTransaction(ctx, service.conn, func(tx *database.Tx) *pnd.AppError {
		sosPostStore := postgres.NewSosPostPostgresStore(tx)
		resourceMediaStore := postgres.NewResourceMediaPostgresStore(tx)
		userStore := postgres.NewUserPostgresStore(tx)

		sosPost, err := sosPostStore.FindSosPostByID(ctx, id)
		if err != nil {
			return err
		}

		mediaData, err := resourceMediaStore.FindResourceMediaByResourceID(ctx, sosPost.ID, string(media.SosResourceType))
		if err != nil {
			return err
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

		conditions, err := sosPostStore.FindConditionByID(ctx, sosPost.ID)
		if err != nil {
			return err
		}

		var conditionsView []sos_post.ConditionView
		for _, c := range conditions {
			view := sos_post.ConditionView{
				ID:   c.ID,
				Name: c.Name,
			}

			conditionsView = append(conditionsView, view)
		}

		pets, err := sosPostStore.FindPetsByID(ctx, sosPost.ID)
		if err != nil {
			return err
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

		author, err := userStore.FindUserByID(ctx, sosPost.AuthorID)

		userView := &user.UserWithoutPrivateInfo{
			ID:              author.ID,
			ProfileImageURL: author.ProfileImageURL,
			Nickname:        author.Nickname,
		}

		sosPostView = &sos_post.FindSosPostView{
			ID:           sosPost.ID,
			Author:       userView,
			Title:        sosPost.Title,
			Content:      sosPost.Content,
			Media:        mediaView,
			Conditions:   conditionsView,
			Pets:         petsView,
			Reward:       sosPost.Reward,
			DateStartAt:  sosPost.DateStartAt,
			DateEndAt:    sosPost.DateEndAt,
			CareType:     sosPost.CareType,
			CarerGender:  sosPost.CarerGender,
			RewardAmount: sosPost.RewardAmount,
			ThumbnailID:  sosPost.ThumbnailID,
			CreatedAt:    sosPost.CreatedAt,
			UpdatedAt:    sosPost.UpdatedAt,
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return sosPostView, nil
}

func (service *SosPostService) UpdateSosPost(ctx context.Context, request *sos_post.UpdateSosPostRequest) (*sos_post.UpdateSosPostView, *pnd.AppError) {
	var sosPostView *sos_post.UpdateSosPostView

	err := database.WithTransaction(ctx, service.conn, func(tx *database.Tx) *pnd.AppError {
		sosPostStore := postgres.NewSosPostPostgresStore(tx)
		resourceMediaStore := postgres.NewResourceMediaPostgresStore(tx)

		updateSosPost, err := sosPostStore.UpdateSosPost(ctx, request)
		if err != nil {
			return err
		}

		mediaData, err := resourceMediaStore.FindResourceMediaByResourceID(ctx, updateSosPost.ID, string(media.SosResourceType))
		if err != nil {
			return err
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

		conditions, err := sosPostStore.FindConditionByID(ctx, updateSosPost.ID)
		if err != nil {
			return err
		}

		var conditionsView []sos_post.ConditionView
		for _, c := range conditions {
			view := sos_post.ConditionView{
				ID:   c.ID,
				Name: c.Name,
			}

			conditionsView = append(conditionsView, view)
		}

		pets, err := sosPostStore.FindPetsByID(ctx, updateSosPost.ID)
		if err != nil {
			return err
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

		sosPostView = &sos_post.UpdateSosPostView{
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
			CareType:     updateSosPost.CareType,
			CarerGender:  updateSosPost.CarerGender,
			RewardAmount: updateSosPost.RewardAmount,
			ThumbnailID:  updateSosPost.ThumbnailID,
			CreatedAt:    updateSosPost.CreatedAt,
			UpdatedAt:    updateSosPost.UpdatedAt,
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return sosPostView, nil
}

func (service *SosPostService) CheckUpdatePermission(ctx context.Context, fbUid string, sosPostID int) (bool, *pnd.AppError) {
	var userID int
	var sosPost *sos_post.SosPost
	var err *pnd.AppError

	err = database.WithTransaction(ctx, service.conn, func(tx *database.Tx) *pnd.AppError {
		userStore := postgres.NewUserPostgresStore(tx)
		sosPostStore := postgres.NewSosPostPostgresStore(tx)

		userID, err = userStore.FindUserIDByFbUID(ctx, fbUid)
		if err != nil {
			return err
		}

		sosPost, err = sosPostStore.FindSosPostByID(ctx, sosPostID)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return false, err
	}

	return userID == sosPost.AuthorID, nil
}
