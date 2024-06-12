package service

import (
	"context"

	"github.com/pet-sitter/pets-next-door-api/internal/domain/soscondition"

	"github.com/pet-sitter/pets-next-door-api/internal/domain/resourcemedia"

	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"

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
	q := databasegen.New(tx)

	userData, err2 := q.FindUser(ctx, databasegen.FindUserParams{
		FbUid: utils.StrToNullStr(fbUID),
	})
	if err2 != nil {
		return nil, pnd.FromPostgresError(err2)
	}

	thumbnailID := setThumbnailID(request.ImageIDs)

	// SOSPost 저장
	sosPost := databasegen.WriteSOSPostRow{}
	var err3 error
	if thumbnailID == nil {
		sosPost, err3 = q.WriteSOSPost(ctx, databasegen.WriteSOSPostParams{
			AuthorID:    utils.IntToNullInt64(int(userData.ID)),
			Title:       utils.StrToNullStr(request.Title),
			Content:     utils.StrToNullStr(request.Content),
			Reward:      utils.StrToNullStr(request.Reward),
			CareType:    utils.StrToNullStr(string(request.CareType)),
			CarerGender: utils.StrToNullStr(request.CarerGender.String()),
			RewardType:  utils.StrToNullStr(request.RewardType.String()),
		})
		if err3 != nil {
			return nil, pnd.FromPostgresError(err3)
		}
	}
	if thumbnailID != nil {
		sosPost, err3 = q.WriteSOSPost(ctx, databasegen.WriteSOSPostParams{
			AuthorID:    utils.IntToNullInt64(int(userData.ID)),
			Title:       utils.StrToNullStr(request.Title),
			Content:     utils.StrToNullStr(request.Content),
			Reward:      utils.StrToNullStr(request.Reward),
			CareType:    utils.StrToNullStr(string(request.CareType)),
			CarerGender: utils.StrToNullStr(request.CarerGender.String()),
			RewardType:  utils.StrToNullStr(request.RewardType.String()),
			ThumbnailID: utils.IntToNullInt64(int(*thumbnailID)),
		})
		if err3 != nil {
			return nil, pnd.FromPostgresError(err3)
		}
	}

	// 날짜 리스트 저장
	err4 := service.SaveSOSDates(ctx, q, request.Dates, int(sosPost.ID))
	if err4 != nil {
		return nil, err4
	}

	// 이미지와 SOSPost 다대다 저장
	err5 := service.SaveLinkSOSPostImage(ctx, q, request.ImageIDs, int(sosPost.ID))
	if err5 != nil {
		return nil, err5
	}

	// 조건과 SOSPost 다대다 저장
	err6 := service.SaveLinkConditions(ctx, q, request.ConditionIDs, int(sosPost.ID))
	if err6 != nil {
		return nil, err6
	}

	// 펫과 SOSPost 다대다 저장
	err7 := service.SaveLinkPets(ctx, q, request.PetIDs, int(sosPost.ID))
	if err7 != nil {
		return nil, err7
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

	dates, err3 := q.FindDatesBySOSPostID(ctx, utils.IntToNullInt64(int(sosPost.ID)))
	if err3 != nil {
		return nil, pnd.FromPostgresError(err3)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return sospost.ToWriteSOSPostView(
		sosPost,
		media.ToListViewFromResourceMediaRows(mediaData),
		soscondition.ToListViewFromSOSPostConditions(conditionList),
		pet.ToDetailViewList(petRows),
		sospost.ToListViewFromSOSDateRows(dates),
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
		media.ToListViewFromViewListForSOSPost(sosPost.Media),
		soscondition.ToListViewFromViewForSOSPost(sosPost.Conditions),
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
	q := databasegen.New(tx)
	// 날짜 업데이트
	// 날짜 삭제
	err2 := service.DeleteLinkSOSPostDates(ctx, q, request.ID)
	if err2 != nil {
		return nil, err2
	}
	// 날짜 저장
	err3 := service.SaveSOSDates(ctx, q, request.Dates, request.ID)
	if err3 != nil {
		return nil, err3
	}
	// 이미지 업데이트
	// 이미지 삭제
	err4 := service.DeleteLinkSOSPostImages(ctx, q, request.ID)
	if err4 != nil {
		return nil, err4
	}
	// 이미지 저장
	err5 := service.SaveLinkSOSPostImage(ctx, q, request.ImageIDs, request.ID)
	if err5 != nil {
		return nil, err5
	}
	// 조건 업데이트
	// 조건 삭제
	err6 := service.DeleteLinkSOSPostConditions(ctx, q, request.ID)
	if err6 != nil {
		return nil, err6
	}
	// 조건 저장
	err7 := service.SaveLinkConditions(ctx, q, request.ConditionIDs, request.ID)
	if err7 != nil {
		return nil, err7
	}
	// 펫 업데이트
	// 펫 삭제
	err8 := service.DeleteLinkSOSPostPets(ctx, q, request.ID)
	if err8 != nil {
		return nil, err8
	}
	// 펫 저장
	err9 := service.SaveLinkPets(ctx, q, request.PetIDs, request.ID)
	if err9 != nil {
		return nil, err9
	}

	// SOSPost 업데이트
	thumbnailID := setThumbnailID(request.ImageIDs)
	updateSOSPost := databasegen.UpdateSOSPostRow{}
	var err10 error
	if thumbnailID == nil {
		updateSOSPost, err10 = q.UpdateSOSPost(ctx, databasegen.UpdateSOSPostParams{
			ID:          int32(request.ID),
			Title:       utils.StrToNullStr(request.Title),
			Content:     utils.StrToNullStr(request.Content),
			Reward:      utils.StrToNullStr(request.Reward),
			CareType:    utils.StrToNullStr(string(request.CareType)),
			CarerGender: utils.StrToNullStr(request.CarerGender.String()),
			RewardType:  utils.StrToNullStr(request.RewardType.String()),
		})
		if err10 != nil {
			return nil, pnd.FromPostgresError(err10)
		}
	}
	if thumbnailID != nil {
		updateSOSPost, err10 = q.UpdateSOSPost(ctx, databasegen.UpdateSOSPostParams{
			ID:          int32(request.ID),
			Title:       utils.StrToNullStr(request.Title),
			Content:     utils.StrToNullStr(request.Content),
			Reward:      utils.StrToNullStr(request.Reward),
			CareType:    utils.StrToNullStr(string(request.CareType)),
			CarerGender: utils.StrToNullStr(request.CarerGender.String()),
			RewardType:  utils.StrToNullStr(request.RewardType.String()),
			ThumbnailID: utils.IntToNullInt64(int(*thumbnailID)),
		})
		if err10 != nil {
			return nil, pnd.FromPostgresError(err10)
		}
	}

	mediaData, err11 := databasegen.New(tx).FindResourceMedia(ctx, databasegen.FindResourceMediaParams{
		ResourceID:   utils.IntToNullInt64(int(updateSOSPost.ID)),
		ResourceType: utils.StrToNullStr(resourcemedia.SOSResourceType.String()),
	})
	if err2 != nil {
		return nil, pnd.FromPostgresError(err11)
	}

	conditionList, err11 := databasegen.New(tx).FindSOSPostConditions(ctx, databasegen.FindSOSPostConditionsParams{
		SosPostID: utils.IntToNullInt64(int(updateSOSPost.ID)),
	})
	if err11 != nil {
		return nil, pnd.FromPostgresError(err11)
	}

	petRows, err12 := databasegen.New(tx).FindPetsBySOSPostID(ctx, utils.IntToNullInt64(int(updateSOSPost.ID)))
	if err12 != nil {
		return nil, pnd.FromPostgresError(err12)
	}

	dates, err13 := q.FindDatesBySOSPostID(ctx, utils.IntToNullInt64(request.ID))
	if err3 != nil {
		return nil, pnd.FromPostgresError(err13)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return sospost.ToUpdateSOSPostView(
		updateSOSPost,
		media.ToListViewFromResourceMediaRows(mediaData),
		soscondition.ToListViewFromSOSPostConditions(conditionList),
		pet.ToDetailViewList(petRows),
		sospost.ToListViewFromSOSDateRows(dates),
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

func (service *SOSPostService) SaveSOSDates(
	ctx context.Context, tx *databasegen.Queries, Dates []sospost.SOSDateView, sosPostID int,
) *pnd.AppError {
	for _, date := range Dates {
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

		// 날짜와 SOSPost 다대다 저장
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
	ctx context.Context, tx *databasegen.Queries, ImageIDs []int64, sosPostID int,
) *pnd.AppError {
	for _, mediaID := range ImageIDs {
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
	ctx context.Context, tx *databasegen.Queries, ConditionIDs []int, sosPostID int,
) *pnd.AppError {
	for _, conditionID := range ConditionIDs {
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
	ctx context.Context, tx *databasegen.Queries, PetIDs []int64, sosPostID int,
) *pnd.AppError {
	for _, petID := range PetIDs {
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
