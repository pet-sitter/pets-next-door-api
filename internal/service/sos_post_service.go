package service

import (
	"context"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database/pgx"
	"github.com/pet-sitter/pets-next-door-api/internal/postgres"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/media"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/sos_post"
)

type SosPostService struct {
	conn *pgx.DB
}

func NewSosPostService(conn *pgx.DB) *SosPostService {
	return &SosPostService{
		conn: conn,
	}
}

func (service *SosPostService) WriteSosPost(ctx context.Context, fbUid string, request *sos_post.WriteSosPostRequest) (*sos_post.WriteSosPostView, *pnd.AppError) {
	tx, err := service.conn.BeginPgxTx(ctx)
	defer tx.Rollback(ctx)
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

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return sosPost.ToWriteSosPostView(
		mediaData.ToMediaViewList(),
		conditions.ToConditionViewList(),
		pets.ToPetViewList(),
		dates.ToSosDateViewList(),
	), nil
}

func (service *SosPostService) FindSosPosts(ctx context.Context, page int, size int, sortBy string) (*sos_post.FindSosPostListView, *pnd.AppError) {
	tx, err := service.conn.BeginPgxTx(ctx)
	defer tx.Rollback(ctx)
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

		findByAuthorSosPostView := sosPost.ToFindSosPostView(
			author.ToUserWithoutPrivateInfo(),
			mediaData.ToMediaViewList(),
			conditions.ToConditionViewList(),
			pets.ToPetViewList(),
			dates.ToSosDateViewList(),
		)

		sosPostViews.Items = append(sosPostViews.Items, *findByAuthorSosPostView)
	}

	return sosPostViews, nil
}

func (service *SosPostService) FindSosPostsByAuthorID(ctx context.Context, authorID int, page int, size int, sortBy string) (*sos_post.FindSosPostListView, *pnd.AppError) {
	tx, err := service.conn.BeginPgxTx(ctx)
	defer tx.Rollback(ctx)
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

		findByAuthorSosPostView := sosPost.ToFindSosPostView(
			author.ToUserWithoutPrivateInfo(),
			mediaData.ToMediaViewList(),
			conditions.ToConditionViewList(),
			pets.ToPetViewList(),
			dates.ToSosDateViewList(),
		)

		sosPostViews.Items = append(sosPostViews.Items, *findByAuthorSosPostView)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return sosPostViews, nil
}

func (service *SosPostService) FindSosPostByID(ctx context.Context, id int) (*sos_post.FindSosPostView, *pnd.AppError) {
	tx, err := service.conn.BeginPgxTx(ctx)
	defer tx.Rollback(ctx)
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

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return sosPost.ToFindSosPostView(
		author.ToUserWithoutPrivateInfo(),
		mediaData.ToMediaViewList(),
		conditions.ToConditionViewList(),
		pets.ToPetViewList(),
		dates.ToSosDateViewList(),
	), nil
}

func (service *SosPostService) UpdateSosPost(ctx context.Context, request *sos_post.UpdateSosPostRequest) (*sos_post.UpdateSosPostView, *pnd.AppError) {
	tx, err := service.conn.BeginPgxTx(ctx)
	defer tx.Rollback(ctx)
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

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return updateSosPost.ToUpdateSosPostView(
		mediaData.ToMediaViewList(),
		conditions.ToConditionViewList(),
		pets.ToPetViewList(),
		dates.ToSosDateViewList(),
	), nil
}

func (service *SosPostService) CheckUpdatePermission(ctx context.Context, fbUid string, sosPostID int) (bool, *pnd.AppError) {
	tx, err := service.conn.BeginPgxTx(ctx)
	defer tx.Rollback(ctx)
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

	if err := tx.Commit(ctx); err != nil {
		return false, err
	}

	return userID == sosPost.AuthorID, nil
}
