package chat

type CreateRoomRequest struct {
	RoomName string `json:"roomName" validate:"required"`
	RoomType string `json:"roomType" validate:"required"`
}
