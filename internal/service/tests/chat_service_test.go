package service

import (
	"context"
	"testing"

	"github.com/pet-sitter/pets-next-door-api/internal/domain/chat"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/media"
	"github.com/pet-sitter/pets-next-door-api/internal/tests"
	"github.com/stretchr/testify/assert"
)

func TestCreateRoom(t *testing.T) {
	t.Run("채팅방을 생성한다", func(t *testing.T) {
		db, tearDown := tests.SetUp(t)
		defer tearDown(t)
		ctx := context.Background()
		chatService := tests.NewMockChatService(db)

		// Given
		roomName := "Test Room"
		roomType := chat.RoomType(chat.RoomTypePersonal)

		// When
		createdRoom, err := chatService.CreateRoom(ctx, roomName, roomType)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, roomName, createdRoom.Name)
		assert.Equal(t, roomType, createdRoom.RoomType)
	})
}

func TestJoinRoom(t *testing.T) {
	t.Run("채팅방에 입장한다", func(t *testing.T) {
		db, tearDown := tests.SetUp(t)
		defer tearDown(t)
		ctx := context.Background()
		chatService := tests.NewMockChatService(db)
		mediaService := tests.NewMockMediaService(db)
		userService := tests.NewMockUserService(db)

		// Given
		profileImage, _ := mediaService.UploadMedia(ctx, nil, media.TypeImage, "profile_image.jpg")

		userRequest := tests.NewDummyRegisterUserRequest(&profileImage.ID)
		createdUser, _ := userService.RegisterUser(ctx, userRequest)

		roomName := "Test Room"
		roomType := chat.RoomType(chat.RoomTypePersonal)
		createdRoom, _ := chatService.CreateRoom(ctx, roomName, roomType)

		// When
		joinRoomView, err := chatService.JoinRoom(ctx, createdRoom.ID, createdUser.FirebaseUID)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, createdRoom.ID, joinRoomView.RoomID)
		assert.Equal(t, createdUser.ID, joinRoomView.UserID)
	})
}

func TestLeaveRoom(t *testing.T) {
	t.Run("채팅방을 떠난다", func(t *testing.T) {
		db, tearDown := tests.SetUp(t)
		defer tearDown(t)
		ctx := context.Background()
		chatService := tests.NewMockChatService(db)
		mediaService := tests.NewMockMediaService(db)
		userService := tests.NewMockUserService(db)

		// Given
		profileImage, _ := mediaService.UploadMedia(ctx, nil, media.TypeImage, "profile_image.jpg")

		userRequest := tests.NewDummyRegisterUserRequest(&profileImage.ID)
		createdUser, _ := userService.RegisterUser(ctx, userRequest)

		roomName := "Test Room"
		roomType := chat.RoomType(chat.RoomTypePersonal)
		createdRoom, _ := chatService.CreateRoom(ctx, roomName, roomType)
		_, _ = chatService.JoinRoom(ctx, createdRoom.ID, createdUser.FirebaseUID)

		// When
		leaveRoomErr := chatService.LeaveRoom(ctx, createdRoom.ID, createdUser.FirebaseUID)

		// Then
		assert.NoError(t, leaveRoomErr)
	})
}

func TestSaveMessage(t *testing.T) {
	t.Run("메시지를 저장한다", func(t *testing.T) {
		db, tearDown := tests.SetUp(t)
		defer tearDown(t)
		ctx := context.Background()
		chatService := tests.NewMockChatService(db)
		mediaService := tests.NewMockMediaService(db)
		userService := tests.NewMockUserService(db)

		// Given
		profileImage, _ := mediaService.UploadMedia(ctx, nil, media.TypeImage, "profile_image.jpg")

		userRequest := tests.NewDummyRegisterUserRequest(&profileImage.ID)
		createdUser, _ := userService.RegisterUser(ctx, userRequest)

		roomName := "Test Room"
		roomType := chat.RoomType(chat.RoomTypePersonal)
		createdRoom, _ := chatService.CreateRoom(ctx, roomName, roomType)
		_, _ = chatService.JoinRoom(ctx, createdRoom.ID, createdUser.FirebaseUID)

		// When
		message := "Hello, World!"
		savedMessage, err := chatService.SaveMessage(
			ctx, createdRoom.ID, createdUser.FirebaseUID, message, chat.MessageTypeNormal,
		)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, createdRoom.ID, savedMessage.RoomID)
		assert.Equal(t, createdUser.ID, savedMessage.UserID)
		assert.Equal(t, message, savedMessage.Content)
	})
}

func TestFindRoomByID(t *testing.T) {
	t.Run("채팅방을 ID로 조회한다", func(t *testing.T) {
		db, tearDown := tests.SetUp(t)
		defer tearDown(t)
		ctx := context.Background()
		chatService := tests.NewMockChatService(db)

		// Given
		roomName := "Test Room"
		roomType := chat.RoomType(chat.RoomTypePersonal)
		createdRoom, _ := chatService.CreateRoom(ctx, roomName, roomType)

		// When
		foundRoom, err := chatService.FindRoomByID(ctx, &createdRoom.ID)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, createdRoom.ID, foundRoom.ID)
		assert.Equal(t, createdRoom.Name, foundRoom.Name)
	})
}

func TestFindUserChatRoom(t *testing.T) {
	t.Run("사용자의 채팅방 참여 목록을 조회한다", func(t *testing.T) {
		db, tearDown := tests.SetUp(t)
		defer tearDown(t)
		ctx := context.Background()
		chatService := tests.NewMockChatService(db)
		mediaService := tests.NewMockMediaService(db)
		userService := tests.NewMockUserService(db)

		// Given
		profileImage, _ := mediaService.UploadMedia(ctx, nil, media.TypeImage, "profile_image.jpg")

		userRequest := tests.NewDummyRegisterUserRequest(&profileImage.ID)
		createdUser, _ := userService.RegisterUser(ctx, userRequest)

		roomName := "Test Room"
		roomType := chat.RoomType(chat.RoomTypePersonal)
		createdRoom, _ := chatService.CreateRoom(ctx, roomName, roomType)

		_, _ = chatService.JoinRoom(ctx, createdRoom.ID, createdUser.FirebaseUID)

		// When
		userChatRooms, err := chatService.FindUserChatRoom(ctx)

		// Then
		assert.NoError(t, err)
		assert.NotEmpty(t, userChatRooms)

		for _, userChatRoom := range userChatRooms {
			if userChatRoom.RoomID == createdRoom.ID && userChatRoom.UserID == createdUser.ID {
				assert.Equal(t, createdRoom.ID, userChatRoom.RoomInfo.ID)
				assert.Equal(t, roomName, userChatRoom.RoomInfo.Name)
				assert.Equal(t, roomType, userChatRoom.RoomInfo.RoomType)
				assert.Equal(t, createdUser.Email, userChatRoom.UserInfo.Email)
				assert.Equal(t, createdUser.Nickname, userChatRoom.UserInfo.Nickname)
				assert.Equal(t, profileImage.URL, *userChatRoom.UserInfo.ProfileImageURL)
			}
		}
	})
}
