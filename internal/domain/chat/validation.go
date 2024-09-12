package chat

import "errors"

// Validate to validate CreateRoomRequest
func (r CreateRoomRequest) RoomTypeValidate() error {
	// RoomType이 Model에 정의된 값인지 확인
	switch r.RoomType {
	case EVENT_ROOM_TYPE:
		return nil
	default:
		return errors.New("invalid room type. please check room type")
	}
}
