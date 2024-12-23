package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	pnd "github.com/pet-sitter/pets-next-door-api/api"
	utils "github.com/pet-sitter/pets-next-door-api/internal/common"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/event"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/media"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/user"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
	databasegen "github.com/pet-sitter/pets-next-door-api/internal/infra/database/gen"
)

type EventService struct {
	conn         *database.DB
	mediaService *MediaService
	userService  *UserService
}

func NewEventService(
	conn *database.DB,
	userService *UserService,
	mediaService *MediaService,
) *EventService {
	return &EventService{
		conn:         conn,
		userService:  userService,
		mediaService: mediaService,
	}
}

func (s *EventService) CreateEvent(
	ctx context.Context,
	authorID uuid.UUID,
	req event.CreateRequest,
) (*event.Event, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, pnd.ErrUnknown(errors.New("failed to generate id"))
	}

	q := databasegen.New(s.conn)

	// Check if author exists
	authorData, err := s.userService.FindUserProfile(
		ctx,
		user.FindUserParams{ID: uuid.NullUUID{UUID: authorID, Valid: true}},
	)
	if err != nil {
		return nil, err
	}
	authorView := user.WithoutPrivateInfo{
		ID:              authorData.ID,
		Nickname:        authorData.Nickname,
		ProfileImageURL: authorData.ProfileImageURL,
	}

	// Check if media exists
	if req.MediaID != nil {
		if _, err = s.mediaService.FindMediaByID(ctx, *req.MediaID); err != nil {
			return nil, err
		}
	}
	mediaID := uuid.NullUUID{}
	if req.MediaID != nil {
		mediaID = uuid.NullUUID{UUID: *req.MediaID, Valid: true}
	}
	mediaData, err := s.mediaService.FindMediaByID(ctx, mediaID.UUID)
	if err != nil {
		return nil, err
	}

	// Map topic[] to string[]
	topics := make([]string, len(req.Topics))
	for i, topic := range req.Topics {
		topics[i] = topic.String()
	}

	eventData, err := q.CreateEvent(ctx, databasegen.CreateEventParams{
		ID:              id,
		AuthorID:        authorID,
		EventType:       req.EventType.String(),
		Name:            req.Name,
		Description:     req.Description,
		MediaID:         mediaID,
		Topics:          topics,
		MaxParticipants: utils.IntPtrToNullInt32(req.MaxParticipants),
		Fee:             int32(req.Fee),
		StartAt:         utils.TimePtrToNullTime(req.StartAt),
	})
	if err != nil {
		return nil, err
	}

	return event.ToEvent(eventData, authorView, mediaData), nil
}

func (s *EventService) FindEvent(
	ctx context.Context,
	params databasegen.FindEventParams,
) (*event.Event, error) {
	q := databasegen.New(s.conn)

	eventData, err := q.FindEvent(ctx, params)
	if err != nil {
		return nil, err
	}

	authorData, err := s.userService.FindUserProfile(
		ctx,
		user.FindUserParams{ID: uuid.NullUUID{UUID: eventData.AuthorID, Valid: true}},
	)
	if err != nil {
		return nil, err
	}
	authorView := user.WithoutPrivateInfo{
		ID:              authorData.ID,
		Nickname:        authorData.Nickname,
		ProfileImageURL: authorData.ProfileImageURL,
	}

	var mediaData *media.DetailView
	if eventData.MediaID.Valid {
		mediaData, err = s.mediaService.FindMediaByID(ctx, eventData.MediaID.UUID)
		if err != nil {
			return nil, err
		}
	}

	return event.ToEvent(eventData, authorView, mediaData), nil
}

func (s *EventService) FindEvents(
	ctx context.Context,
	params databasegen.FindEventsParams,
) ([]*event.Event, error) {
	q := databasegen.New(s.conn)

	eventsData, err := q.FindEvents(ctx, params)
	if err != nil {
		return nil, err
	}

	authorIDs := make([]uuid.UUID, len(eventsData))
	for i, eventData := range eventsData {
		authorIDs[i] = eventData.AuthorID
	}
	authorsData, err := s.userService.FindUsersByIDs(ctx, databasegen.FindUsersByIDsParams{
		Ids: authorIDs,
	})
	if err != nil {
		return nil, err
	}
	authorIDsMap := make(map[uuid.UUID]user.WithoutPrivateInfo)
	for _, authorData := range authorsData {
		authorIDsMap[authorData.ID] = authorData
	}

	mediaIDs := make([]uuid.UUID, len(eventsData))
	for i, eventData := range eventsData {
		if eventData.MediaID.Valid {
			mediaIDs[i] = eventData.MediaID.UUID
		}
	}
	mediasData, err := s.mediaService.FindMediasByIDs(ctx, mediaIDs)
	if err != nil {
		return nil, err
	}
	mediasDataMap := make(map[uuid.UUID]media.DetailView)
	for _, mediaData := range mediasData {
		mediasDataMap[mediaData.ID] = mediaData
	}

	events := make([]*event.Event, len(eventsData))
	for i, eventData := range eventsData {
		authorData := authorIDsMap[eventData.AuthorID]
		mediaData := mediasDataMap[eventData.MediaID.UUID]
		events[i] = event.ToEvent(eventData, authorData, &mediaData)
	}
	return events, nil
}
