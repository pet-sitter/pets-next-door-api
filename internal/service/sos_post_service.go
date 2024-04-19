package service

import (
	"context"

	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
	"github.com/pet-sitter/pets-next-door-api/internal/postgres"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/media"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/sospost"
)

type SOSPostService struct {
	conn *database.DB
}

func NewSOSPostService(conn *database.DB) *SOSPostService {
	return &SOSPostService{
		conn: conn,
	}
}

func (service *SOSPostService) WriteSOSPost(
	ctx context.Context, fbUID string, request *sospost.WriteSOSPostRequest,
) (*sospost.WriteSOSPostView, *pnd.AppError) {
	tx, err := service.conn.BeginTx(ctx)
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}

	userID, err := postgres.FindUserIDByFbUID(ctx, tx, fbUID)
	if err != nil {
		return nil, err
	}

	sosPost, err := postgres.WriteSOSPost(ctx, tx, userID, request)
	if err != nil {
		return nil, err
	}

	mediaData, err := postgres.FindResourceMediaByResourceID(ctx, tx, sosPost.ID, string(media.SOSResourceType))
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

	dates, err := postgres.FindDatesBySOSPostID(ctx, tx, sosPost.ID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return sosPost.ToWriteSOSPostView(
		mediaData.ToMediaViewList(),
		conditions.ToConditionViewList(),
		pets.ToPetViewList(),
		dates.ToSOSDateViewList(),
	), nil
}

func (service *SOSPostService) FindSOSPosts(
	ctx context.Context, page, size int, sortBy, filterType string,
) (*sospost.FindSOSPostListView, *pnd.AppError) {
	tx, err := service.conn.BeginTx(ctx)
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}

	sosPosts, err := postgres.FindSOSPosts(ctx, tx, page, size, sortBy, filterType)
	if err != nil {
		return nil, err
	}

	sosPostViews := sospost.FromEmptySOSPostInfoList(sosPosts)

	for _, sosPost := range sosPosts.Items {
		author, err := postgres.FindUserByID(ctx, tx, sosPost.AuthorID, true)
		if err != nil {
			return nil, err
		}
		sosPostView := sosPost.ToFindSOSPostInfoView(
			author.ToUserWithoutPrivateInfo(),
			sosPost.Media.ToMediaViewList(),
			sosPost.Conditions.ToConditionViewList(),
			sosPost.Pets.ToPetViewList(),
			sosPost.Dates.ToSOSDateViewList(),
		)

		sosPostViews.Items = append(sosPostViews.Items, *sosPostView)
	}

	return sosPostViews, nil
}

func (service *SOSPostService) FindSOSPostsByAuthorID(
	ctx context.Context, authorID, page, size int, sortBy, filterType string,
) (*sospost.FindSOSPostListView, *pnd.AppError) {
	tx, err := service.conn.BeginTx(ctx)
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}

	sosPosts, err := postgres.FindSOSPostsByAuthorID(ctx, tx, authorID, page, size, sortBy, filterType)
	if err != nil {
		return nil, err
	}
	sosPostViews := sospost.FromEmptySOSPostInfoList(sosPosts)

	for _, sosPost := range sosPosts.Items {
		author, err := postgres.FindUserByID(ctx, tx, sosPost.AuthorID, true)
		if err != nil {
			return nil, err
		}
		sosPostView := sosPost.ToFindSOSPostInfoView(
			author.ToUserWithoutPrivateInfo(),
			sosPost.Media.ToMediaViewList(),
			sosPost.Conditions.ToConditionViewList(),
			sosPost.Pets.ToPetViewList(),
			sosPost.Dates.ToSOSDateViewList(),
		)

		sosPostViews.Items = append(sosPostViews.Items, *sosPostView)
	}
	return sosPostViews, nil
}

func (service *SOSPostService) FindSOSPostByID(ctx context.Context, id int) (*sospost.FindSOSPostView, *pnd.AppError) {
	tx, err := service.conn.BeginTx(ctx)
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}

	sosPost, err := postgres.FindSOSPostByID(ctx, tx, id)
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

	return sosPost.ToFindSOSPostInfoView(
		author.ToUserWithoutPrivateInfo(),
		sosPost.Media.ToMediaViewList(),
		sosPost.Conditions.ToConditionViewList(),
		sosPost.Pets.ToPetViewList(),
		sosPost.Dates.ToSOSDateViewList(),
	), nil
}

func (service *SOSPostService) UpdateSOSPost(
	ctx context.Context, request *sospost.UpdateSOSPostRequest,
) (*sospost.UpdateSOSPostView, *pnd.AppError) {
	tx, err := service.conn.BeginTx(ctx)
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}

	updateSOSPost, err := postgres.UpdateSOSPost(ctx, tx, request)
	if err != nil {
		return nil, err
	}

	mediaData, err := postgres.FindResourceMediaByResourceID(ctx, tx, updateSOSPost.ID, string(media.SOSResourceType))
	if err != nil {
		return nil, err
	}

	conditions, err := postgres.FindConditionByID(ctx, tx, updateSOSPost.ID)
	if err != nil {
		return nil, err
	}

	pets, err := postgres.FindPetsByID(ctx, tx, updateSOSPost.ID)
	if err != nil {
		return nil, err
	}

	dates, err := postgres.FindDatesBySOSPostID(ctx, tx, request.ID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return updateSOSPost.ToUpdateSOSPostView(
		mediaData.ToMediaViewList(),
		conditions.ToConditionViewList(),
		pets.ToPetViewList(),
		dates.ToSOSDateViewList(),
	), nil
}

func (service *SOSPostService) CheckUpdatePermission(
	ctx context.Context, fbUID string, sosPostID int,
) (bool, *pnd.AppError) {
	tx, err := service.conn.BeginTx(ctx)
	defer tx.Rollback()
	if err != nil {
		return false, err
	}

	userID, err := postgres.FindUserIDByFbUID(ctx, tx, fbUID)
	if err != nil {
		return false, err
	}

	sosPost, err := postgres.FindSOSPostByID(ctx, tx, sosPostID)
	if err != nil {
		return false, err
	}

	if err := tx.Commit(); err != nil {
		return false, err
	}

	return userID == sosPost.AuthorID, nil
}
