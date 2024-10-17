package roomservice

import (
	"context"
	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
)

func (r *roomService) FindRooms(ctx context.Context, params *FindRoomParams) ([]*modelsv1.Room, error) {
	//TODO implement me
	panic("implement me")
}

func (r *roomService) FindMyRoom(ctx context.Context) (*modelsv1.Room, error) {
	//TODO implement me
	panic("implement me")
}

func (r *roomService) GetRoomByID(ctx context.Context, roomID string) (*modelsv1.Room, error) {
	//TODO implement me
	panic("implement me")
}

func (r *roomService) CreateRoom(ctx context.Context, room *modelsv1.RoomCreate) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (r *roomService) UpdateRoom(ctx context.Context, room *modelsv1.Room) error {
	//TODO implement me
	panic("implement me")
}

func (r *roomService) DeleteRoom(ctx context.Context, roomID string) error {
	//TODO implement me
	panic("implement me")
}
