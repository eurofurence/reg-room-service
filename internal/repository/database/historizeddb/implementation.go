package historizeddb

import (
	"context"
	"github.com/eurofurence/reg-room-service/internal/entity"
	"github.com/eurofurence/reg-room-service/internal/repository/database"
)

type HistorizingRepository struct {
	wrappedRepository database.Repository
}

func Create(wrappedRepository database.Repository) database.Repository {
	return &HistorizingRepository{wrappedRepository: wrappedRepository}
}

func (r *HistorizingRepository) Open(ctx context.Context) error {
	return r.wrappedRepository.Open(ctx)
}

func (r HistorizingRepository) Close(ctx context.Context) {
	r.wrappedRepository.Close(ctx)
}

func (r HistorizingRepository) Migrate(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (r HistorizingRepository) AddGroup(ctx context.Context, g *entity.Group) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (r HistorizingRepository) UpdateGroup(ctx context.Context, g *entity.Group) error {
	//TODO implement me
	panic("implement me")
}

func (r HistorizingRepository) GetGroupByID(ctx context.Context, id string) (*entity.Group, error) {
	//TODO implement me
	panic("implement me")
}

func (r HistorizingRepository) SoftDeleteGroupByID(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (r HistorizingRepository) UndeleteGroupByID(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (r HistorizingRepository) NewEmptyGroupMembership(ctx context.Context, groupID string, attendeeID uint) *entity.GroupMember {
	//TODO implement me
	panic("implement me")
}

func (r HistorizingRepository) GetGroupMembershipByAttendeeID(ctx context.Context, attendeeID uint) (*entity.GroupMember, error) {
	//TODO implement me
	panic("implement me")
}

func (r HistorizingRepository) GetGroupMembersByGroupID(ctx context.Context, groupID string) ([]entity.GroupMember, error) {
	//TODO implement me
	panic("implement me")
}

func (r HistorizingRepository) AddGroupMembership(ctx context.Context, gm *entity.GroupMember) error {
	//TODO implement me
	panic("implement me")
}

func (r HistorizingRepository) UpdateGroupMembership(ctx context.Context, gm *entity.GroupMember) error {
	//TODO implement me
	panic("implement me")
}

func (r HistorizingRepository) DeleteGroupMembership(ctx context.Context, attendeeID uint) error {
	//TODO implement me
	panic("implement me")
}

func (r HistorizingRepository) AddRoom(ctx context.Context, g *entity.Room) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (r HistorizingRepository) UpdateRoom(ctx context.Context, g *entity.Room) error {
	//TODO implement me
	panic("implement me")
}

func (r HistorizingRepository) GetRoomByID(ctx context.Context, id string) (*entity.Room, error) {
	//TODO implement me
	panic("implement me")
}

func (r HistorizingRepository) SoftDeleteRoomByID(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (r HistorizingRepository) UndeleteRoomByID(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (r HistorizingRepository) NewEmptyRoomMembership(ctx context.Context, roomID string, attendeeID uint) *entity.RoomMember {
	//TODO implement me
	panic("implement me")
}

func (r HistorizingRepository) GetRoomMembershipByAttendeeID(ctx context.Context, attendeeID uint) (*entity.RoomMember, error) {
	//TODO implement me
	panic("implement me")
}

func (r HistorizingRepository) GetRoomMembersByRoomID(ctx context.Context, roomID string) ([]entity.RoomMember, error) {
	//TODO implement me
	panic("implement me")
}

func (r HistorizingRepository) AddRoomMembership(ctx context.Context, gm *entity.RoomMember) error {
	//TODO implement me
	panic("implement me")
}

func (r HistorizingRepository) UpdateRoomMembership(ctx context.Context, gm *entity.RoomMember) error {
	//TODO implement me
	panic("implement me")
}

func (r HistorizingRepository) DeleteRoomMembership(ctx context.Context, attendeeID uint) error {
	//TODO implement me
	panic("implement me")
}

func (r HistorizingRepository) RecordHistory(ctx context.Context, h *entity.History) error {
	//TODO implement me
	panic("implement me")
}
