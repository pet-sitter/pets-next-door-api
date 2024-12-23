package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"gopkg.in/go-playground/assert.v1"

	"github.com/pet-sitter/pets-next-door-api/internal/domain/event"
	databasegen "github.com/pet-sitter/pets-next-door-api/internal/infra/database/gen"

	"github.com/pet-sitter/pets-next-door-api/internal/domain/media"

	"github.com/pet-sitter/pets-next-door-api/internal/tests"
)

func TestCreateEvent(t *testing.T) {
	t.Run("이벤트를 새로 생성한다", func(t *testing.T) {
		db, tearDown := tests.SetUp(t)
		defer tearDown(t)
		ctx := context.Background()
		mediaService := tests.NewMockMediaService(db)
		userService := tests.NewMockUserService(db)
		eventService := tests.NewMockEventService(db)

		// Given
		profileImage, _ := mediaService.UploadMedia(ctx, nil, media.TypeImage, "profile_image.jpg")
		author, _ := userService.RegisterUser(
			ctx,
			tests.NewDummyRegisterUserRequest(uuid.NullUUID{UUID: profileImage.ID, Valid: true}),
		)
		eventMedia, _ := mediaService.UploadMedia(ctx, nil, media.TypeImage, "event_thumbnail.jpg")

		// When
		now := time.Now().UTC()
		maxParticipants := 10
		created, _ := eventService.CreateEvent(ctx, author.ID, event.CreateRequest{
			BaseCreateRequest: event.BaseCreateRequest{
				EventType:       event.ShortTerm,
				Name:            "테스트 이벤트",
				Description:     "테스트 이벤트 설명",
				Topics:          []event.EventTopic{event.ETC},
				MediaID:         &eventMedia.ID,
				MaxParticipants: &maxParticipants,
				Fee:             3000,
				StartAt:         &now,
			},
		})

		// Then
		found, err := eventService.FindEvent(
			ctx,
			databasegen.FindEventParams{ID: uuid.NullUUID{UUID: created.ID, Valid: true}},
		)
		assert.Equal(t, err, nil)
		assertEventEquals(t, created, found)
	})
}

func TestFindEvents(t *testing.T) {
	t.Run("이벤트 목록을 조회한다", func(t *testing.T) {
		db, tearDown := tests.SetUp(t)
		defer tearDown(t)
		ctx := context.Background()
		mediaService := tests.NewMockMediaService(db)
		userService := tests.NewMockUserService(db)
		eventService := tests.NewMockEventService(db)

		// Given
		profileImage, _ := mediaService.UploadMedia(ctx, nil, media.TypeImage, "profile_image.jpg")
		owner, _ := userService.RegisterUser(
			ctx,
			tests.NewDummyRegisterUserRequest(uuid.NullUUID{UUID: profileImage.ID, Valid: true}),
		)
		eventMedia, _ := mediaService.UploadMedia(ctx, nil, media.TypeImage, "event_thumbnail.jpg")

		// When
		now := time.Now().UTC()
		maxParticipants := 10
		created, _ := eventService.CreateEvent(ctx, owner.ID, event.CreateRequest{
			BaseCreateRequest: event.BaseCreateRequest{
				EventType:       event.ShortTerm,
				Name:            "테스트 이벤트",
				Description:     "테스트 이벤트 설명",
				Topics:          []event.EventTopic{event.ETC},
				MediaID:         &eventMedia.ID,
				MaxParticipants: &maxParticipants,
				Fee:             3000,
				StartAt:         &now,
			},
		})

		// Then
		found, err := eventService.FindEvents(
			ctx,
			databasegen.FindEventsParams{
				Limit:    10,
				AuthorID: uuid.NullUUID{UUID: uuid.Nil, Valid: false},
			},
		)
		assert.Equal(t, err, nil)
		assert.Equal(t, 1, len(found))
		assertEventEquals(t, created, found[0])
	})

	t.Run("작성자 ID로 이벤트 목록을 조회할 수 있다", func(t *testing.T) {
		db, tearDown := tests.SetUp(t)
		defer tearDown(t)
		ctx := context.Background()
		mediaService := tests.NewMockMediaService(db)
		userService := tests.NewMockUserService(db)
		eventService := tests.NewMockEventService(db)

		// Given
		profileImage, _ := mediaService.UploadMedia(ctx, nil, media.TypeImage, "profile_image.jpg")
		author, _ := userService.RegisterUser(
			ctx,
			tests.NewDummyRegisterUserRequest(uuid.NullUUID{UUID: profileImage.ID, Valid: true}),
		)
		eventMedia, _ := mediaService.UploadMedia(ctx, nil, media.TypeImage, "event_thumbnail.jpg")

		// When
		now := time.Now().UTC()
		maxParticipants := 10
		created, _ := eventService.CreateEvent(ctx, author.ID, event.CreateRequest{
			BaseCreateRequest: event.BaseCreateRequest{
				EventType:       event.ShortTerm,
				Name:            "테스트 이벤트",
				Description:     "테스트 이벤트 설명",
				Topics:          []event.EventTopic{event.ETC},
				MediaID:         &eventMedia.ID,
				MaxParticipants: &maxParticipants,
				Fee:             3000,
				StartAt:         &now,
			},
		})

		// Then
		found, err := eventService.FindEvents(
			ctx,
			databasegen.FindEventsParams{
				Limit:    10,
				AuthorID: uuid.NullUUID{UUID: author.ID, Valid: true},
			},
		)
		assert.Equal(t, err, nil)
		assert.Equal(t, 1, len(found))
		assertEventEquals(t, created, found[0])
	})
}

func assertEventEquals(
	t *testing.T,
	expected *databasegen.Event,
	actual *event.Event,
) {
	t.Helper()
	assert.Equal(t, expected.ID, actual.Event.ID)
	assert.Equal(t, expected.AuthorID, actual.Event.AuthorID)
	assert.Equal(t, expected.EventType, actual.Event.EventType)
	assert.Equal(t, expected.Name, actual.Event.Name)
	assert.Equal(t, expected.Description, actual.Event.Description)
	assert.Equal(t, expected.MediaID, actual.Event.MediaID)
	assert.Equal(t, expected.Topics, actual.Event.Topics)
	assert.Equal(t, expected.MaxParticipants, actual.Event.MaxParticipants)
	assert.Equal(t, expected.Fee, actual.Event.Fee)
	assert.Equal(t, expected.StartAt, actual.Event.StartAt)
	assert.Equal(t, expected.CreatedAt, actual.Event.CreatedAt)
	assert.Equal(t, expected.UpdatedAt, actual.Event.UpdatedAt)
	assert.Equal(t, expected.DeletedAt, actual.Event.DeletedAt)
}
