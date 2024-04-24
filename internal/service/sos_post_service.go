package service

import (
	"context"

	utils "github.com/pet-sitter/pets-next-door-api/internal/common"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/user"
	databasegen "github.com/pet-sitter/pets-next-door-api/internal/infra/database/gen"

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

	userData, err2 := databasegen.New(tx).FindUser(ctx, databasegen.FindUserParams{
		FbUid: utils.StrToNullStr(fbUID),
	})
	if err2 != nil {
		return nil, pnd.FromPostgresError(err2)
	}

	sosPost, err := postgres.WriteSOSPost(ctx, tx, int(userData.ID), request)
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
		pets.ToDetailViewList(),
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
		author, err := databasegen.New(tx).FindUser(ctx, databasegen.FindUserParams{
			ID:             utils.IntToNullInt32(sosPost.AuthorID),
			IncludeDeleted: true,
		})
		if err != nil {
			return nil, pnd.FromPostgresError(err)
		}

		sosPostView := sosPost.ToFindSOSPostInfoView(
			&user.WithoutPrivateInfo{
				ID:              int64(author.ID),
				Nickname:        author.Nickname,
				ProfileImageURL: utils.NullStrToStrPtr(author.ProfileImageUrl),
			},
			sosPost.Media.ToMediaViewList(),
			sosPost.Conditions.ToConditionViewList(),
			sosPost.Pets.ToDetailViewList(),
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
		author, err := databasegen.New(tx).FindUser(ctx, databasegen.FindUserParams{
			ID:             utils.IntToNullInt32(sosPost.AuthorID),
			IncludeDeleted: true,
		})
		if err != nil {
			return nil, pnd.FromPostgresError(err)
		}

		sosPostView := sosPost.ToFindSOSPostInfoView(
			&user.WithoutPrivateInfo{
				ID:              int64(author.ID),
				Nickname:        author.Nickname,
				ProfileImageURL: utils.NullStrToStrPtr(author.ProfileImageUrl),
			},
			sosPost.Media.ToMediaViewList(),
			sosPost.Conditions.ToConditionViewList(),
			sosPost.Pets.ToDetailViewList(),
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

	author, err2 := databasegen.New(tx).FindUser(ctx, databasegen.FindUserParams{
		ID:             utils.IntToNullInt32(sosPost.AuthorID),
		IncludeDeleted: true,
	})
	if err2 != nil {
		return nil, pnd.FromPostgresError(err2)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return sosPost.ToFindSOSPostInfoView(
		&user.WithoutPrivateInfo{
			ID:              int64(author.ID),
			Nickname:        author.Nickname,
			ProfileImageURL: utils.NullStrToStrPtr(author.ProfileImageUrl),
		},
		sosPost.Media.ToMediaViewList(),
		sosPost.Conditions.ToConditionViewList(),
		sosPost.Pets.ToDetailViewList(),
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
		pets.ToDetailViewList(),
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

	userData, err2 := databasegen.New(tx).FindUser(ctx, databasegen.FindUserParams{
		FbUid: utils.StrToNullStr(fbUID),
	})
	if err2 != nil {
		return false, pnd.FromPostgresError(err2)
	}

	sosPost, err := postgres.FindSOSPostByID(ctx, tx, sosPostID)
	if err != nil {
		return false, err
	}

	if err := tx.Commit(); err != nil {
		return false, err
	}

	return int(userData.ID) == sosPost.AuthorID, nil
}
