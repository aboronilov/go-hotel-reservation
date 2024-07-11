package api

import "github.com/aboronilov/go-hotel-reservation/db"

type RoomHandler struct {
	roomStore db.RoomStore
}

func NewRoomHandler(roomStore db.RoomStore) *RoomHandler {
	return &RoomHandler{
		roomStore: roomStore,
	}
}
