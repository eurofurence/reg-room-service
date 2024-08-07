package database

import (
	"context"

	"github.com/eurofurence/reg-room-service/internal/entity"
)

type Repository interface {
	Open(ctx context.Context) error
	Close(ctx context.Context)
	Migrate(ctx context.Context) error

	// GetGroups returns all non-soft-deleted groups.
	GetGroups(ctx context.Context) ([]*entity.Group, error)
	// FindGroups returns IDs of all non-soft-deleted groups satisfying the criteria.
	//
	// Occupancy is the number of people actually in the group, as opposed to its maximum size.
	// If maxOccupancy is set to -1, it will be ignored as a criterion.
	//
	// A group matches the list of badge numbers in anyOfMemberID if at least one of those badge numbers
	// is a member of the group. An empty list means no condition.
	FindGroups(ctx context.Context, minOccupancy uint, maxOccupancy int, anyOfMemberID []uint) ([]string, error)
	AddGroup(ctx context.Context, group *entity.Group) (string, error)
	UpdateGroup(ctx context.Context, group *entity.Group) error
	GetGroupByID(ctx context.Context, id string) (*entity.Group, error) // may return soft deleted entities!
	SoftDeleteGroupByID(ctx context.Context, id string) error
	UndeleteGroupByID(ctx context.Context, id string) error

	// NewEmptyGroupMembership pre-fills some required and internal fields, including the
	// groupID and attendeeID.
	NewEmptyGroupMembership(ctx context.Context, groupID string, attendeeID uint) *entity.GroupMember
	GetGroupMembershipByAttendeeID(ctx context.Context, attendeeID uint) (*entity.GroupMember, error)
	GetGroupMembersByGroupID(ctx context.Context, groupID string) ([]*entity.GroupMember, error)
	AddGroupMembership(ctx context.Context, gm *entity.GroupMember) error
	UpdateGroupMembership(ctx context.Context, gm *entity.GroupMember) error
	DeleteGroupMembership(ctx context.Context, attendeeID uint) error

	// GetRooms returns all non-soft-deleted rooms.
	GetRooms(ctx context.Context) ([]*entity.Room, error)
	AddRoom(ctx context.Context, room *entity.Room) (string, error)
	UpdateRoom(ctx context.Context, room *entity.Room) error
	GetRoomByID(ctx context.Context, id string) (*entity.Room, error) // may return soft deleted entities!
	SoftDeleteRoomByID(ctx context.Context, id string) error
	UndeleteRoomByID(ctx context.Context, id string) error

	// NewEmptyRoomMembership pre-fills some required and internal fields, including the
	// RoomID and attendeeID.
	NewEmptyRoomMembership(ctx context.Context, roomID string, attendeeID uint) *entity.RoomMember
	GetRoomMembershipByAttendeeID(ctx context.Context, attendeeID uint) (*entity.RoomMember, error)
	GetRoomMembersByRoomID(ctx context.Context, roomID string) ([]*entity.RoomMember, error)
	AddRoomMembership(ctx context.Context, rm *entity.RoomMember) error
	UpdateRoomMembership(ctx context.Context, rm *entity.RoomMember) error
	DeleteRoomMembership(ctx context.Context, attendeeID uint) error

	RecordHistory(ctx context.Context, h *entity.History) error
}
