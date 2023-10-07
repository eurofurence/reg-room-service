package database

import (
	"context"
	"github.com/eurofurence/reg-room-service/internal/entity"
)

type Repository interface {
	Open(ctx context.Context) error
	Close(ctx context.Context)
	Migrate(ctx context.Context) error

	AddGroup(ctx context.Context, g *entity.Group) (string, error)
	UpdateGroup(ctx context.Context, g *entity.Group) error
	GetGroupByID(ctx context.Context, id string) (*entity.Group, error)
	SoftDeleteGroupByID(ctx context.Context, id string) error
	UndeleteGroupByID(ctx context.Context, id string) error

	// NewEmptyGroupMembership pre-fills some required and internal fields, including the
	// groupID and attendeeID.
	NewEmptyGroupMembership(ctx context.Context, groupID string, attendeeID uint) *entity.GroupMember
	GetGroupMembershipByAttendeeID(ctx context.Context, attendeeID uint) (*entity.GroupMember, error)
	GetGroupMembersByGroupID(ctx context.Context, groupID string) ([]entity.GroupMember, error)
	AddGroupMembership(ctx context.Context, gm *entity.GroupMember) error
	UpdateGroupMembership(ctx context.Context, gm *entity.GroupMember) error
	DeleteGroupMembership(ctx context.Context, attendeeID uint) error

	AddRoom(ctx context.Context, g *entity.Room) (string, error)
	UpdateRoom(ctx context.Context, g *entity.Room) error
	GetRoomByID(ctx context.Context, id string) (*entity.Room, error)
	SoftDeleteRoomByID(ctx context.Context, id string) error
	UndeleteRoomByID(ctx context.Context, id string) error

	// NewEmptyRoomMembership pre-fills some required and internal fields, including the
	// RoomID and attendeeID.
	NewEmptyRoomMembership(ctx context.Context, roomID string, attendeeID uint) *entity.RoomMember
	GetRoomMembershipByAttendeeID(ctx context.Context, attendeeID uint) (*entity.RoomMember, error)
	GetRoomMembersByRoomID(ctx context.Context, roomID string) ([]entity.RoomMember, error)
	AddRoomMembership(ctx context.Context, gm *entity.RoomMember) error
	UpdateRoomMembership(ctx context.Context, gm *entity.RoomMember) error
	DeleteRoomMembership(ctx context.Context, attendeeID uint) error

	RecordHistory(ctx context.Context, h *entity.History) error
}
