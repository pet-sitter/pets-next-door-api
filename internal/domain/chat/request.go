package chat

// FIXME : 정책상 채팅방 생성시 필요 정보가 있을 경우 추가 해야함
type CreateRoomRequest struct {
	RoomName string `json:"roomName" validate:"required"`
	// FIXME : Model 항목을 상속받고 있는데 위계질서에 어긋남 수정 필요
	RoomType RoomType `json:"roomType" validate:"required"`
}
