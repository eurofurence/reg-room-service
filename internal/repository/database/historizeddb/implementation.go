package historizeddb

import (
	"context"
	"errors"
	"fmt"

	"github.com/d4l3k/messagediff"

	"github.com/eurofurence/reg-room-service/internal/entity"
	"github.com/eurofurence/reg-room-service/internal/repository/database"
)

// TODO need request id and identity from context for history, see at bottom of file!!!

type HistorizingRepository struct {
	wrappedRepository database.Repository
}

func New(wrappedRepository database.Repository) database.Repository {
	return &HistorizingRepository{wrappedRepository: wrappedRepository}
}

func (r *HistorizingRepository) Open(ctx context.Context) error {
	return r.wrappedRepository.Open(ctx)
}

func (r *HistorizingRepository) Close(ctx context.Context) {
	r.wrappedRepository.Close(ctx)
}

func (r *HistorizingRepository) Migrate(ctx context.Context) error {
	return r.wrappedRepository.Migrate(ctx)
}

type entityType string

const (
	typeGroup       entityType = "Group"
	typeGroupMember entityType = "GroupMember"
	typeRoom        entityType = "Room"
	typeRoomMember  entityType = "RoomMember"
)

type operationType string

const (
	opUpdate   operationType = "update"
	opDelete   operationType = "delete"
	opUndelete operationType = "undelete"
)

// group

func (r *HistorizingRepository) GetGroups(ctx context.Context) ([]*entity.Group, error) {
	return r.wrappedRepository.GetGroups(ctx)
}

func (r *HistorizingRepository) AddGroup(ctx context.Context, group *entity.Group) (string, error) {
	return r.wrappedRepository.AddGroup(ctx, group)
}

func (r *HistorizingRepository) FindGroups(ctx context.Context, minOccupancy uint, maxOccupancy int, anyOfMemberID []uint) ([]string, error) {
	return r.wrappedRepository.FindGroups(ctx, minOccupancy, maxOccupancy, anyOfMemberID)
}

func (r *HistorizingRepository) UpdateGroup(ctx context.Context, group *entity.Group) error {
	oldVersion, err := r.wrappedRepository.GetGroupByID(ctx, group.ID)
	if err != nil {
		return err
	}

	// hide always present diff in times
	oldVersion.CreatedAt = group.CreatedAt
	oldVersion.UpdatedAt = group.UpdatedAt

	histEntry := diffReverse(ctx, oldVersion, group, typeGroup, group.ID, opUpdate)

	err = r.wrappedRepository.RecordHistory(ctx, histEntry)
	if err != nil {
		return err
	}

	return r.wrappedRepository.UpdateGroup(ctx, group)
}

func (r *HistorizingRepository) GetGroupByID(ctx context.Context, id string) (*entity.Group, error) {
	return r.wrappedRepository.GetGroupByID(ctx, id)
}

func (r *HistorizingRepository) DeleteGroupByID(ctx context.Context, id string) error {
	histEntry := noDiffRecord(ctx, typeGroup, id, opDelete)

	if err := r.wrappedRepository.RecordHistory(ctx, histEntry); err != nil {
		return err
	}

	return r.wrappedRepository.DeleteGroupByID(ctx, id)
}

// group members

func (r *HistorizingRepository) NewEmptyGroupMembership(ctx context.Context, groupID string, attendeeID uint, nickname string) *entity.GroupMember {
	return r.wrappedRepository.NewEmptyGroupMembership(ctx, groupID, attendeeID, nickname)
}

func (r *HistorizingRepository) GetGroupMembershipByAttendeeID(ctx context.Context, attendeeID uint) (*entity.GroupMember, error) {
	return r.wrappedRepository.GetGroupMembershipByAttendeeID(ctx, attendeeID)
}

func (r *HistorizingRepository) GetGroupMembersByGroupID(ctx context.Context, groupID string) ([]*entity.GroupMember, error) {
	return r.wrappedRepository.GetGroupMembersByGroupID(ctx, groupID)
}

func (r *HistorizingRepository) AddGroupMembership(ctx context.Context, gm *entity.GroupMember) error {
	return r.wrappedRepository.AddGroupMembership(ctx, gm)
}

func (r *HistorizingRepository) UpdateGroupMembership(ctx context.Context, gm *entity.GroupMember) error {
	oldVersion, err := r.wrappedRepository.GetGroupMembershipByAttendeeID(ctx, gm.ID)
	if err != nil {
		return err
	}

	// hide always present diff in times
	oldVersion.CreatedAt = gm.CreatedAt
	oldVersion.UpdatedAt = gm.UpdatedAt

	histEntry := diffReverse(ctx, oldVersion, gm, typeGroupMember, fmt.Sprintf("%d", gm.ID), opUpdate)

	err = r.wrappedRepository.RecordHistory(ctx, histEntry)
	if err != nil {
		return err
	}

	return r.wrappedRepository.UpdateGroupMembership(ctx, gm)
}

func (r *HistorizingRepository) DeleteGroupMembership(ctx context.Context, attendeeID uint) error {
	histEntry := noDiffRecord(ctx, typeGroupMember, fmt.Sprintf("%d", attendeeID), opDelete)

	if err := r.wrappedRepository.RecordHistory(ctx, histEntry); err != nil {
		return err
	}

	return r.wrappedRepository.DeleteGroupMembership(ctx, attendeeID)
}

// room

func (r *HistorizingRepository) GetRooms(ctx context.Context) ([]*entity.Room, error) {
	return r.wrappedRepository.GetRooms(ctx)
}

func (r *HistorizingRepository) AddRoom(ctx context.Context, room *entity.Room) (string, error) {
	return r.wrappedRepository.AddRoom(ctx, room)
}

func (r *HistorizingRepository) UpdateRoom(ctx context.Context, room *entity.Room) error {
	oldVersion, err := r.wrappedRepository.GetRoomByID(ctx, room.ID)
	if err != nil {
		return err
	}

	// hide always present diff in times
	oldVersion.CreatedAt = room.CreatedAt
	oldVersion.UpdatedAt = room.UpdatedAt

	histEntry := diffReverse(ctx, oldVersion, room, typeRoom, room.ID, opUpdate)

	err = r.wrappedRepository.RecordHistory(ctx, histEntry)
	if err != nil {
		return err
	}

	return r.wrappedRepository.UpdateRoom(ctx, room)
}

func (r *HistorizingRepository) GetRoomByID(ctx context.Context, id string) (*entity.Room, error) {
	return r.wrappedRepository.GetRoomByID(ctx, id)
}

func (r *HistorizingRepository) DeleteRoomByID(ctx context.Context, id string) error {
	histEntry := noDiffRecord(ctx, typeRoom, id, opDelete)

	if err := r.wrappedRepository.RecordHistory(ctx, histEntry); err != nil {
		return err
	}

	return r.wrappedRepository.DeleteRoomByID(ctx, id)
}

// room members

func (r *HistorizingRepository) NewEmptyRoomMembership(ctx context.Context, roomID string, attendeeID uint) *entity.RoomMember {
	return r.wrappedRepository.NewEmptyRoomMembership(ctx, roomID, attendeeID)
}

func (r *HistorizingRepository) GetRoomMembershipByAttendeeID(ctx context.Context, attendeeID uint) (*entity.RoomMember, error) {
	return r.wrappedRepository.GetRoomMembershipByAttendeeID(ctx, attendeeID)
}

func (r *HistorizingRepository) GetRoomMembersByRoomID(ctx context.Context, roomID string) ([]*entity.RoomMember, error) {
	return r.wrappedRepository.GetRoomMembersByRoomID(ctx, roomID)
}

func (r *HistorizingRepository) AddRoomMembership(ctx context.Context, rm *entity.RoomMember) error {
	return r.wrappedRepository.AddRoomMembership(ctx, rm)
}

func (r *HistorizingRepository) UpdateRoomMembership(ctx context.Context, rm *entity.RoomMember) error {
	oldVersion, err := r.wrappedRepository.GetRoomMembershipByAttendeeID(ctx, rm.ID)
	if err != nil {
		return err
	}

	// hide always present diff in times
	oldVersion.CreatedAt = rm.CreatedAt
	oldVersion.UpdatedAt = rm.UpdatedAt

	histEntry := diffReverse(ctx, oldVersion, rm, typeRoomMember, fmt.Sprintf("%d", rm.ID), opUpdate)

	err = r.wrappedRepository.RecordHistory(ctx, histEntry)
	if err != nil {
		return err
	}

	return r.wrappedRepository.UpdateRoomMembership(ctx, rm)
}

func (r *HistorizingRepository) DeleteRoomMembership(ctx context.Context, attendeeID uint) error {
	histEntry := noDiffRecord(ctx, typeRoomMember, fmt.Sprintf("%d", attendeeID), opDelete)

	if err := r.wrappedRepository.RecordHistory(ctx, histEntry); err != nil {
		return err
	}

	return r.wrappedRepository.DeleteRoomMembership(ctx, attendeeID)
}

// --- history ---

func (r *HistorizingRepository) RecordHistory(ctx context.Context, h *entity.History) error {
	// it is an error to call this from the outside. From the inside use wrappedRepository.RecordHistory to bypass the error
	return errors.New("not allowed to directly manipulate history")
}

func diffReverse[T any](_ context.Context, oldVersion *T, newVersion *T, entityName entityType, entityID string, operation operationType) *entity.History {
	// we diff reverse so the OLD value is printed in the diffs. The new value is in the database now.
	histEntry := &entity.History{
		Entity:    string(entityName),
		EntityId:  entityID,
		Operation: string(operation),
		RequestId: "", // TODO ctxvalues.RequestId(ctx),
		Identity:  "", // TODO ctxvalues.Subject(ctx),
	}
	diff, _ := messagediff.PrettyDiff(*newVersion, *oldVersion)
	histEntry.Diff = diff
	return histEntry
}

func noDiffRecord(_ context.Context, entityName entityType, entityID string, operation operationType) *entity.History {
	return &entity.History{
		Entity:    string(entityName),
		EntityId:  entityID,
		Operation: string(operation),
		RequestId: "", // TODO ctxvalues.RequestId(ctx),
		Identity:  "", // TODO ctxvalues.Subject(ctx),
	}
}
