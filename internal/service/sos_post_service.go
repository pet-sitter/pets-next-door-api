package service

import (
	"context"
	"time"

	"github.com/pet-sitter/pets-next-door-api/internal/domain/soscondition"

	"github.com/pet-sitter/pets-next-door-api/internal/domain/resourcemedia"

	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"

	utils "github.com/pet-sitter/pets-next-door-api/internal/common"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/user"
	databasegen "github.com/pet-sitter/pets-next-door-api/internal/infra/database/gen"

	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"

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
) (*sospost.DetailView, *pnd.AppError) {
	tx, err := service.conn.BeginTx(ctx)
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}
	q := databasegen.New(tx)

	userData, err2 := q.FindUser(ctx, databasegen.FindUserParams{
		FbUid: utils.StrToNullStr(fbUID),
	})
	if err2 != nil {
		return nil, pnd.FromPostgresError(err2)
	}

	thumbnailID := setThumbnailID(request.ImageIDs)
	sosPost, err := service.createSOSPost(ctx, q, int64(userData.ID), request, thumbnailID)
	if err != nil {
		return nil, err
	}

	if err := service.saveAllLinks(ctx, q, request, int(sosPost.ID)); err != nil {
		return nil, err
	}

	mediaData, err2 := q.FindResourceMedia(ctx, databasegen.FindResourceMediaParams{
		ResourceID:   utils.IntToNullInt64(int(sosPost.ID)),
		ResourceType: utils.StrToNullStr(resourcemedia.SOSResourceType.String()),
	})
	if err2 != nil {
		return nil, pnd.FromPostgresError(err2)
	}

	conditionList, err2 := q.FindSOSPostConditions(ctx, databasegen.FindSOSPostConditionsParams{
		SosPostID: utils.IntToNullInt64(int(sosPost.ID)),
	})
	if err2 != nil {
		return nil, pnd.FromPostgresError(err2)
	}

	petRows, err2 := q.FindPetsBySOSPostID(ctx, utils.IntToNullInt64(int(sosPost.ID)))
	if err2 != nil {
		return nil, pnd.FromPostgresError(err2)
	}

	dates, err2 := q.FindDatesBySOSPostID(ctx, utils.IntToNullInt64(int(sosPost.ID)))
	if err2 != nil {
		return nil, pnd.FromPostgresError(err2)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return sospost.CreateDetailView(
		sosPost,
		media.ToListViewFromResourceMediaRows(mediaData),
		soscondition.ToListViewFromSOSPostConditions(conditionList),
		pet.ToDetailViewList(petRows),
		sospost.ToListViewFromSOSDateRows(dates),
	), nil
}

func (service *SOSPostService) createSOSPost(
	ctx context.Context, q *databasegen.Queries, authorID int64, request *sospost.WriteSOSPostRequest, thumbnailID *int64,
) (databasegen.WriteSOSPostRow, *pnd.AppError) {
	params := databasegen.WriteSOSPostParams{
		AuthorID:    utils.IntToNullInt64(int(authorID)),
		Title:       utils.StrToNullStr(request.Title),
		Content:     utils.StrToNullStr(request.Content),
		Reward:      utils.StrToNullStr(request.Reward),
		CareType:    utils.StrToNullStr(string(request.CareType)),
		CarerGender: utils.StrToNullStr(request.CarerGender.String()),
		RewardType:  utils.StrToNullStr(request.RewardType.String()),
	}

	if thumbnailID != nil {
		params.ThumbnailID = utils.IntToNullInt64(int(*thumbnailID))
	}

	sosPost, err := q.WriteSOSPost(ctx, params)
	if err != nil {
		return databasegen.WriteSOSPostRow{}, pnd.FromPostgresError(err)
	}

	return sosPost, nil
}

func (service *SOSPostService) saveAllLinks(
	ctx context.Context, q *databasegen.Queries, request *sospost.WriteSOSPostRequest, sosPostID int,
) *pnd.AppError {
	if err := service.SaveSOSDates(ctx, q, request.Dates, sosPostID); err != nil {
		return err
	}
	if err := service.SaveLinkSOSPostImage(ctx, q, request.ImageIDs, sosPostID); err != nil {
		return err
	}
	if err := service.SaveLinkConditions(ctx, q, request.ConditionIDs, sosPostID); err != nil {
		return err
	}
	return service.SaveLinkPets(ctx, q, request.PetIDs, sosPostID)
}

func (service *SOSPostService) FindSOSPosts(
	ctx context.Context, page, size int, sortBy, filterType string,
) (*sospost.FindSOSPostListView, *pnd.AppError) {
	tx, err := service.conn.BeginTx(ctx)
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}

	sosPosts, err2 := databasegen.New(tx).FindSOSPosts(ctx, databasegen.FindSOSPostsParams{
		EarliestDateStartAt: utils.FormatDateString(time.Now().String()),
		PetType:             utils.StrToNullStr(filterType),
		SortBy:              utils.StrToNullStr(sortBy),
		Limit:               utils.IntToNullInt32(size + 1),
		Offset:              utils.IntToNullInt32((page - 1) * size),
	})
	if err2 != nil {
		return nil, pnd.FromPostgresError(err2)
	}

	sosPostInfoList := sospost.ToInfoListFromFindRow(sosPosts, page, size)
	sosPostViews := sospost.FromEmptySOSPostInfoList(sosPostInfoList)

	for _, sosPost := range sosPostInfoList.Items {
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
			media.ToListViewFromViewListForSOSPost(sosPost.Media),
			soscondition.ToListViewFromViewForSOSPost(sosPost.Conditions),
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

	sosPosts, err2 := databasegen.New(tx).FindSOSPostsByAuthorID(ctx, databasegen.FindSOSPostsByAuthorIDParams{
		EarliestDateStartAt: utils.FormatDateString(time.Now().String()),
		PetType:             utils.StrToNullStr(filterType),
		AuthorID:            utils.IntToNullInt64(authorID),
		SortBy:              utils.StrToNullStr(sortBy),
		Limit:               utils.IntToNullInt32(size + 1),
		Offset:              utils.IntToNullInt32((page - 1) * size),
	})
	if err2 != nil {
		return nil, pnd.FromPostgresError(err2)
	}

	sosPostInfoList := sospost.ToInfoListFromFindAuthorIDRow(sosPosts, page, size)
	sosPostViews := sospost.FromEmptySOSPostInfoList(sosPostInfoList)

	for _, sosPost := range sosPostInfoList.Items {
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
			media.ToListViewFromViewListForSOSPost(sosPost.Media),
			soscondition.ToListViewFromViewForSOSPost(sosPost.Conditions),
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

	sosPost, err2 := databasegen.New(tx).FindSOSPostByID(ctx, utils.IntToNullInt32(id))
	if err2 != nil {
		return nil, pnd.FromPostgresError(err2)
	}
	sosPostInfo := sospost.ToInfoFromFindByIDRow(sosPost)

	author, err2 := databasegen.New(tx).FindUser(ctx, databasegen.FindUserParams{
		ID:             utils.IntToNullInt32(sosPostInfo.AuthorID),
		IncludeDeleted: true,
	})
	if err2 != nil {
		return nil, pnd.FromPostgresError(err2)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return sosPostInfo.ToFindSOSPostInfoView(
		&user.WithoutPrivateInfo{
			ID:              int64(author.ID),
			Nickname:        author.Nickname,
			ProfileImageURL: utils.NullStrToStrPtr(author.ProfileImageUrl),
		},
		media.ToListViewFromViewListForSOSPost(sosPostInfo.Media),
		soscondition.ToListViewFromViewForSOSPost(sosPostInfo.Conditions),
		sosPostInfo.Pets.ToDetailViewList(),
		sosPostInfo.Dates.ToSOSDateViewList(),
	), nil
}

func (service *SOSPostService) UpdateSOSPost(
	ctx context.Context, request *sospost.UpdateSOSPostRequest,
) (*sospost.DetailView, *pnd.AppError) {
	tx, err := service.conn.BeginTx(ctx)
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}
	q := databasegen.New(tx)

	if err2 := service.updateAllLinks(ctx, q, request); err2 != nil {
		return nil, err2
	}

	thumbnailID := setThumbnailID(request.ImageIDs)
	updateSOSPost, err := service.updateSOSPost(ctx, q, request, thumbnailID)
	if err != nil {
		return nil, err
	}

	mediaData, err2 := databasegen.New(tx).FindResourceMedia(ctx, databasegen.FindResourceMediaParams{
		ResourceID:   utils.IntToNullInt64(int(updateSOSPost.ID)),
		ResourceType: utils.StrToNullStr(resourcemedia.SOSResourceType.String()),
	})
	if err2 != nil {
		return nil, pnd.FromPostgresError(err2)
	}

	conditionList, err2 := databasegen.New(tx).FindSOSPostConditions(ctx, databasegen.FindSOSPostConditionsParams{
		SosPostID: utils.IntToNullInt64(int(updateSOSPost.ID)),
	})
	if err2 != nil {
		return nil, pnd.FromPostgresError(err2)
	}

	petRows, err2 := databasegen.New(tx).FindPetsBySOSPostID(ctx, utils.IntToNullInt64(int(updateSOSPost.ID)))
	if err2 != nil {
		return nil, pnd.FromPostgresError(err2)
	}

	dates, err2 := q.FindDatesBySOSPostID(ctx, utils.IntToNullInt64(request.ID))
	if err2 != nil {
		return nil, pnd.FromPostgresError(err2)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return sospost.UpdateDetailView(
		updateSOSPost,
		media.ToListViewFromResourceMediaRows(mediaData),
		soscondition.ToListViewFromSOSPostConditions(conditionList),
		pet.ToDetailViewList(petRows),
		sospost.ToListViewFromSOSDateRows(dates),
	), nil
}

func (service *SOSPostService) updateSOSPost(
	ctx context.Context, q *databasegen.Queries, request *sospost.UpdateSOSPostRequest, thumbnailID *int64,
) (databasegen.UpdateSOSPostRow, *pnd.AppError) {
	params := databasegen.UpdateSOSPostParams{
		ID:          int32(request.ID),
		Title:       utils.StrToNullStr(request.Title),
		Content:     utils.StrToNullStr(request.Content),
		Reward:      utils.StrToNullStr(request.Reward),
		CareType:    utils.StrToNullStr(string(request.CareType)),
		CarerGender: utils.StrToNullStr(request.CarerGender.String()),
		RewardType:  utils.StrToNullStr(request.RewardType.String()),
	}

	if thumbnailID != nil {
		params.ThumbnailID = utils.IntToNullInt64(int(*thumbnailID))
	}

	updateSOSPost, err := q.UpdateSOSPost(ctx, params)
	if err != nil {
		return databasegen.UpdateSOSPostRow{}, pnd.FromPostgresError(err)
	}

	return updateSOSPost, nil
}

func (service *SOSPostService) updateAllLinks(
	ctx context.Context, q *databasegen.Queries, request *sospost.UpdateSOSPostRequest,
) *pnd.AppError {
	if err := service.DeleteLinkSOSPostDates(ctx, q, request.ID); err != nil {
		return err
	}
	if err := service.SaveSOSDates(ctx, q, request.Dates, request.ID); err != nil {
		return err
	}

	if err := service.DeleteLinkSOSPostImages(ctx, q, request.ID); err != nil {
		return err
	}
	if err := service.SaveLinkSOSPostImage(ctx, q, request.ImageIDs, request.ID); err != nil {
		return err
	}

	if err := service.DeleteLinkSOSPostConditions(ctx, q, request.ID); err != nil {
		return err
	}
	if err := service.SaveLinkConditions(ctx, q, request.ConditionIDs, request.ID); err != nil {
		return err
	}

	if err := service.DeleteLinkSOSPostPets(ctx, q, request.ID); err != nil {
		return err
	}
	return service.SaveLinkPets(ctx, q, request.PetIDs, request.ID)
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

	sosPost, err2 := databasegen.New(tx).FindSOSPostByID(ctx, utils.IntToNullInt32(sosPostID))
	if err2 != nil {
		return false, pnd.FromPostgresError(err2)
	}
	sosPostInfo := sospost.ToInfoFromFindByIDRow(sosPost)

	if err := tx.Commit(); err != nil {
		return false, err
	}

	return int(userData.ID) == sosPostInfo.AuthorID, nil
}

func (service *SOSPostService) SaveSOSDates(
	ctx context.Context, tx *databasegen.Queries, dates []sospost.SOSDateView, sosPostID int,
) *pnd.AppError {
	for _, date := range dates {
		dateStartAt, err := utils.StrToNullTime(date.DateStartAt)
		if err != nil {
			return err
		}
		dateEndAt, err := utils.StrToNullTime(date.DateEndAt)
		if err != nil {
			return err
		}

		d, err2 := databasegen.New(service.conn).InsertSOSDate(ctx, databasegen.InsertSOSDateParams{
			DateStartAt: dateStartAt,
			DateEndAt:   dateEndAt,
		})
		if err2 != nil {
			return pnd.FromPostgresError(err2)
		}

		err3 := tx.LinkSOSPostDate(ctx, databasegen.LinkSOSPostDateParams{
			SosPostID:  utils.IntToNullInt64(sosPostID),
			SosDatesID: utils.IntToNullInt64(int(d.ID)),
		})
		if err3 != nil {
			return pnd.FromPostgresError(err3)
		}
	}
	return nil
}

func (service *SOSPostService) SaveLinkSOSPostImage(
	ctx context.Context, tx *databasegen.Queries, imageIDs []int64, sosPostID int,
) *pnd.AppError {
	for _, mediaID := range imageIDs {
		err := tx.LinkResourceMedia(ctx, databasegen.LinkResourceMediaParams{
			MediaID:      utils.IntToNullInt64(int(mediaID)),
			ResourceID:   utils.IntToNullInt64(sosPostID),
			ResourceType: utils.StrToNullStr(resourcemedia.SOSResourceType.String()),
		})
		if err != nil {
			return pnd.FromPostgresError(err)
		}
	}
	return nil
}

func (service *SOSPostService) SaveLinkConditions(
	ctx context.Context, tx *databasegen.Queries, conditionIDs []int, sosPostID int,
) *pnd.AppError {
	for _, conditionID := range conditionIDs {
		err := tx.LinkSOSPostCondition(ctx, databasegen.LinkSOSPostConditionParams{
			SosPostID:      utils.IntToNullInt64(sosPostID),
			SosConditionID: utils.IntToNullInt64(conditionID),
		})
		if err != nil {
			return pnd.FromPostgresError(err)
		}
	}
	return nil
}

func (service *SOSPostService) SaveLinkPets(
	ctx context.Context, tx *databasegen.Queries, petIDs []int64, sosPostID int,
) *pnd.AppError {
	for _, petID := range petIDs {
		err := tx.LinkSOSPostPet(ctx, databasegen.LinkSOSPostPetParams{
			SosPostID: utils.IntToNullInt64(sosPostID),
			PetID:     utils.IntToNullInt64(int(petID)),
		})
		if err != nil {
			return pnd.FromPostgresError(err)
		}
	}
	return nil
}

func (service *SOSPostService) DeleteLinkSOSPostDates(
	ctx context.Context, tx *databasegen.Queries, sosPostID int,
) *pnd.AppError {
	err := tx.DeleteSOSPostDateBySOSPostID(ctx, utils.IntPtrToNullInt64(&sosPostID))
	if err != nil {
		return pnd.FromPostgresError(err)
	}
	return nil
}

func (service *SOSPostService) DeleteLinkSOSPostImages(
	ctx context.Context, tx *databasegen.Queries, sosPostID int,
) *pnd.AppError {
	err := tx.DeleteResourceMediaByResourceID(ctx, utils.IntToNullInt64(sosPostID))
	if err != nil {
		return pnd.FromPostgresError(err)
	}
	return nil
}

func (service *SOSPostService) DeleteLinkSOSPostConditions(
	ctx context.Context, tx *databasegen.Queries, sosPostID int,
) *pnd.AppError {
	err := tx.DeleteSOSPostConditionBySOSPostID(ctx, utils.IntToNullInt64(sosPostID))
	if err != nil {
		return pnd.FromPostgresError(err)
	}
	return nil
}

func (service *SOSPostService) DeleteLinkSOSPostPets(
	ctx context.Context, tx *databasegen.Queries, sosPostID int,
) *pnd.AppError {
	err := tx.DeleteSOSPostPetBySOSPostID(ctx, utils.IntToNullInt64(sosPostID))
	if err != nil {
		return pnd.FromPostgresError(err)
	}
	return nil
}

func setThumbnailID(imageIDs []int64) *int64 {
	if len(imageIDs) > 0 {
		return &imageIDs[0]
	}
	return nil
}
