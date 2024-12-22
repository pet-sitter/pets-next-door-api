package service

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	pnd "github.com/pet-sitter/pets-next-door-api/api"
	utils "github.com/pet-sitter/pets-next-door-api/internal/common"
	"github.com/pet-sitter/pets-next-door-api/internal/datatype"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/media"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/resourcemedia"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/soscondition"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/sospost"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/user"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
	databasegen "github.com/pet-sitter/pets-next-door-api/internal/infra/database/gen"
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
) (*sospost.DetailView, error) {
	tx, err := service.conn.BeginTx(ctx)
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}
	q := databasegen.New(tx)

	userData, err := q.FindUser(ctx, databasegen.FindUserParams{
		FbUid: utils.StrToNullStr(fbUID),
	})
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	thumbnailID := setThumbnailID(request.ImageIDs)
	sosPost, err := service.createSOSPost(ctx, q, userData.ID, request, thumbnailID)
	if err != nil {
		return nil, err
	}

	if err := service.saveAllLinks(ctx, q, request, sosPost.ID); err != nil {
		return nil, err
	}

	mediaData, err := q.FindResourceMedia(ctx, databasegen.FindResourceMediaParams{
		ResourceID:   uuid.NullUUID{UUID: sosPost.ID, Valid: true},
		ResourceType: utils.StrToNullStr(resourcemedia.SOSResourceType.String()),
	})
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	conditionList, err := q.FindSOSPostConditions(ctx, databasegen.FindSOSPostConditionsParams{
		SosPostID: sosPost.ID,
	})
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	petRows, err := q.FindPetsBySOSPostID(ctx, sosPost.ID)
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	dates, err := q.FindDatesBySOSPostID(ctx, uuid.NullUUID{UUID: sosPost.ID, Valid: true})
	if err != nil {
		return nil, pnd.FromPostgresError(err)
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
	ctx context.Context, q *databasegen.Queries, authorID uuid.UUID,
	request *sospost.WriteSOSPostRequest, thumbnailID uuid.NullUUID,
) (databasegen.WriteSOSPostRow, error) {
	params := databasegen.WriteSOSPostParams{
		ID:          datatype.NewUUIDV7(),
		AuthorID:    authorID,
		Title:       utils.StrToNullStr(request.Title),
		Content:     utils.StrToNullStr(request.Content),
		Reward:      utils.StrToNullStr(request.Reward),
		CareType:    utils.StrToNullStr(string(request.CareType)),
		CarerGender: utils.StrToNullStr(request.CarerGender.String()),
		RewardType:  utils.StrToNullStr(request.RewardType.String()),
	}

	if thumbnailID.Valid {
		params.ThumbnailID = thumbnailID
	}

	sosPost, err := q.WriteSOSPost(ctx, params)
	if err != nil {
		return databasegen.WriteSOSPostRow{}, pnd.FromPostgresError(err)
	}

	return sosPost, nil
}

func (service *SOSPostService) saveAllLinks(
	ctx context.Context,
	q *databasegen.Queries,
	request *sospost.WriteSOSPostRequest,
	sosPostID uuid.UUID,
) error {
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
) (*sospost.FindSOSPostListView, error) {
	tx, err := service.conn.BeginTx(ctx)
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}

	sosPosts, err := databasegen.New(tx).FindSOSPosts(ctx, databasegen.FindSOSPostsParams{
		EarliestDateStartAt: utils.FormatDateString(time.Now().String()),
		PetType:             utils.StrToNullStr(filterType),
		SortBy:              utils.StrToNullStr(sortBy),
		Limit:               utils.IntToNullInt32(size + 1),
		Offset:              utils.IntToNullInt32((page - 1) * size),
	})
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	sosPostInfoList := sospost.ToInfoListFromFindRow(sosPosts, page, size)
	sosPostViews := sospost.FromEmptySOSPostInfoList(sosPostInfoList)

	for _, sosPost := range sosPostInfoList.Items {
		author, err := databasegen.New(tx).FindUser(ctx, databasegen.FindUserParams{
			ID:             uuid.NullUUID{UUID: sosPost.AuthorID, Valid: true},
			IncludeDeleted: true,
		})
		if err != nil {
			return nil, pnd.FromPostgresError(err)
		}
		sosPostView := sosPost.ToFindSOSPostInfoView(
			&user.WithoutPrivateInfo{
				ID:              author.ID,
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
	ctx context.Context, authorID uuid.UUID, page, size int, sortBy, filterType string,
) (*sospost.FindSOSPostListView, error) {
	tx, err := service.conn.BeginTx(ctx)
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}

	sosPosts, err := databasegen.New(tx).
		FindSOSPostsByAuthorID(ctx, databasegen.FindSOSPostsByAuthorIDParams{
			EarliestDateStartAt: utils.FormatDateString(time.Now().String()),
			PetType:             utils.StrToNullStr(filterType),
			AuthorID:            uuid.NullUUID{UUID: authorID, Valid: true},
			SortBy:              utils.StrToNullStr(sortBy),
			Limit:               utils.IntToNullInt32(size + 1),
			Offset:              utils.IntToNullInt32((page - 1) * size),
		})
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	sosPostInfoList := sospost.ToInfoListFromFindAuthorIDRow(sosPosts, page, size)
	sosPostViews := sospost.FromEmptySOSPostInfoList(sosPostInfoList)

	for _, sosPost := range sosPostInfoList.Items {
		author, err := databasegen.New(tx).FindUser(ctx, databasegen.FindUserParams{
			ID:             uuid.NullUUID{UUID: sosPost.AuthorID, Valid: true},
			IncludeDeleted: true,
		})
		if err != nil {
			return nil, pnd.FromPostgresError(err)
		}

		sosPostView := sosPost.ToFindSOSPostInfoView(
			&user.WithoutPrivateInfo{
				ID:              author.ID,
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

func (service *SOSPostService) FindSOSPostByID(
	ctx context.Context, id uuid.UUID,
) (*sospost.FindSOSPostView, error) {
	tx, err := service.conn.BeginTx(ctx)
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}

	sosPost, err := databasegen.New(tx).FindSOSPostByID(ctx, uuid.NullUUID{UUID: id, Valid: true})
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}
	sosPostInfo := sospost.ToInfoFromFindByIDRow(sosPost)

	author, err := databasegen.New(tx).FindUser(ctx, databasegen.FindUserParams{
		ID:             uuid.NullUUID{UUID: sosPostInfo.AuthorID, Valid: true},
		IncludeDeleted: true,
	})
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return sosPostInfo.ToFindSOSPostInfoView(
		&user.WithoutPrivateInfo{
			ID:              author.ID,
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
) (*sospost.DetailView, error) {
	tx, err := service.conn.BeginTx(ctx)
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}
	q := databasegen.New(tx)

	if err = service.updateAllLinks(ctx, q, request); err != nil {
		return nil, err
	}

	thumbnailID := setThumbnailID(request.ImageIDs)
	updateSOSPost, err := service.updateSOSPost(ctx, q, request, thumbnailID)
	if err != nil {
		return nil, err
	}

	mediaData, err := databasegen.New(tx).
		FindResourceMedia(ctx, databasegen.FindResourceMediaParams{
			ResourceID:   uuid.NullUUID{UUID: updateSOSPost.ID, Valid: true},
			ResourceType: utils.StrToNullStr(resourcemedia.SOSResourceType.String()),
		})
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	conditionList, err := databasegen.New(tx).
		FindSOSPostConditions(ctx, databasegen.FindSOSPostConditionsParams{
			SosPostID: updateSOSPost.ID,
		})
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	petRows, err := databasegen.New(tx).FindPetsBySOSPostID(ctx, updateSOSPost.ID)
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	dates, err := q.FindDatesBySOSPostID(ctx, uuid.NullUUID{UUID: request.ID, Valid: true})
	if err != nil {
		return nil, pnd.FromPostgresError(err)
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
	ctx context.Context,
	q *databasegen.Queries,
	request *sospost.UpdateSOSPostRequest,
	thumbnailID uuid.NullUUID,
) (databasegen.UpdateSOSPostRow, error) {
	params := databasegen.UpdateSOSPostParams{
		ID:          request.ID,
		Title:       utils.StrToNullStr(request.Title),
		Content:     utils.StrToNullStr(request.Content),
		Reward:      utils.StrToNullStr(request.Reward),
		CareType:    utils.StrToNullStr(string(request.CareType)),
		CarerGender: utils.StrToNullStr(request.CarerGender.String()),
		RewardType:  utils.StrToNullStr(request.RewardType.String()),
	}

	if thumbnailID.Valid {
		params.ThumbnailID = thumbnailID
	}

	updateSOSPost, err := q.UpdateSOSPost(ctx, params)
	if err != nil {
		return databasegen.UpdateSOSPostRow{}, pnd.FromPostgresError(err)
	}

	return updateSOSPost, nil
}

func (service *SOSPostService) updateAllLinks(
	ctx context.Context, q *databasegen.Queries, request *sospost.UpdateSOSPostRequest,
) error {
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
	ctx context.Context, fbUID string, sosPostID uuid.UUID,
) (bool, error) {
	tx, err := service.conn.BeginTx(ctx)
	defer tx.Rollback()
	if err != nil {
		return false, err
	}

	userData, err := databasegen.New(tx).FindUser(ctx, databasegen.FindUserParams{
		FbUid: utils.StrToNullStr(fbUID),
	})
	if err != nil {
		return false, pnd.FromPostgresError(err)
	}

	sosPost, err := databasegen.New(tx).
		FindSOSPostByID(ctx, uuid.NullUUID{UUID: sosPostID, Valid: true})
	if err != nil {
		return false, pnd.FromPostgresError(err)
	}
	sosPostInfo := sospost.ToInfoFromFindByIDRow(sosPost)

	if err := tx.Commit(); err != nil {
		return false, err
	}

	return userData.ID == sosPostInfo.AuthorID, nil
}

func (service *SOSPostService) SaveSOSDates(
	ctx context.Context, tx *databasegen.Queries, dates []sospost.SOSDateView, sosPostID uuid.UUID,
) error {
	for _, date := range dates {
		dateStartAt, err := utils.StrToNullTime(date.DateStartAt)
		if err != nil {
			return err
		}
		dateEndAt, err := utils.StrToNullTime(date.DateEndAt)
		if err != nil {
			return err
		}

		d, err := databasegen.New(service.conn).InsertSOSDate(ctx, databasegen.InsertSOSDateParams{
			ID:          datatype.NewUUIDV7(),
			DateStartAt: dateStartAt,
			DateEndAt:   dateEndAt,
		})
		if err != nil {
			return pnd.FromPostgresError(err)
		}

		if err := tx.LinkSOSPostDate(ctx, databasegen.LinkSOSPostDateParams{
			ID:         datatype.NewUUIDV7(),
			SosPostID:  sosPostID,
			SosDatesID: d.ID,
		}); err != nil {
			return pnd.FromPostgresError(err)
		}
	}
	return nil
}

func (service *SOSPostService) SaveLinkSOSPostImage(
	ctx context.Context, tx *databasegen.Queries, imageIDs []uuid.UUID, sosPostID uuid.UUID,
) error {
	for _, mediaID := range imageIDs {
		log.Default().Println("mediaID", mediaID)

		if err := tx.LinkResourceMedia(ctx, databasegen.LinkResourceMediaParams{
			ID:           datatype.NewUUIDV7(),
			MediaID:      mediaID,
			ResourceID:   sosPostID,
			ResourceType: utils.StrToNullStr(resourcemedia.SOSResourceType.String()),
		}); err != nil {
			return pnd.FromPostgresError(err)
		}
	}
	return nil
}

func (service *SOSPostService) SaveLinkConditions(
	ctx context.Context, tx *databasegen.Queries, conditionIDs []uuid.UUID, sosPostID uuid.UUID,
) error {
	for _, conditionID := range conditionIDs {
		log.Default().Println("conditionID", conditionID)

		if err := tx.LinkSOSPostCondition(ctx, databasegen.LinkSOSPostConditionParams{
			ID:             datatype.NewUUIDV7(),
			SosPostID:      sosPostID,
			SosConditionID: uuid.NullUUID{UUID: conditionID, Valid: true},
		}); err != nil {
			return pnd.FromPostgresError(err)
		}
	}
	return nil
}

func (service *SOSPostService) SaveLinkPets(
	ctx context.Context, tx *databasegen.Queries, petIDs []uuid.UUID, sosPostID uuid.UUID,
) error {
	for _, petID := range petIDs {
		log.Default().Println("petID", petID)

		if err := tx.LinkSOSPostPet(ctx, databasegen.LinkSOSPostPetParams{
			ID:        datatype.NewUUIDV7(),
			SosPostID: sosPostID,
			PetID:     petID,
		}); err != nil {
			return pnd.FromPostgresError(err)
		}
	}
	return nil
}

func (service *SOSPostService) DeleteLinkSOSPostDates(
	ctx context.Context, tx *databasegen.Queries, sosPostID uuid.UUID,
) error {
	if err := tx.DeleteSOSPostDateBySOSPostID(ctx, sosPostID); err != nil {
		return pnd.FromPostgresError(err)
	}
	return nil
}

func (service *SOSPostService) DeleteLinkSOSPostImages(
	ctx context.Context, tx *databasegen.Queries, sosPostID uuid.UUID,
) error {
	if err := tx.DeleteResourceMediaByResourceID(ctx, sosPostID); err != nil {
		return pnd.FromPostgresError(err)
	}
	return nil
}

func (service *SOSPostService) DeleteLinkSOSPostConditions(
	ctx context.Context, tx *databasegen.Queries, sosPostID uuid.UUID,
) error {
	if err := tx.DeleteSOSPostConditionBySOSPostID(ctx, sosPostID); err != nil {
		return pnd.FromPostgresError(err)
	}
	return nil
}

func (service *SOSPostService) DeleteLinkSOSPostPets(
	ctx context.Context, tx *databasegen.Queries, sosPostID uuid.UUID,
) error {
	if err := tx.DeleteSOSPostPetBySOSPostID(ctx, sosPostID); err != nil {
		return pnd.FromPostgresError(err)
	}
	return nil
}

func setThumbnailID(imageIDs []uuid.UUID) uuid.NullUUID {
	if len(imageIDs) > 0 {
		return uuid.NullUUID{UUID: imageIDs[0], Valid: true}
	}
	return uuid.NullUUID{UUID: uuid.Nil, Valid: false}
}
