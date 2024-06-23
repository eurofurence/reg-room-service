package inmemorydb

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/eurofurence/reg-room-service/internal/entity"
	"github.com/eurofurence/reg-room-service/internal/repository/database"
)

type IMGroup struct {
	Group   entity.Group         // intentionally not a pointer so assignment makes a copy
	Members []entity.GroupMember // intentionally not pointers so assignment makes a copy
}

type IMRoom struct {
	Room    entity.Room         // intentionally not a pointer so assignment makes a copy
	Members []entity.RoomMember // intentionally not pointers so assignment makes a copy
}

type InMemoryRepository struct {
	groups     map[string]*IMGroup
	rooms      map[string]*IMRoom
	history    map[uint]*entity.History
	idSequence uint32
	Now        func() time.Time
}

func New() database.Repository {
	return &InMemoryRepository{
		Now: time.Now,
	}
}

func (r *InMemoryRepository) Open(_ context.Context) error {
	r.groups = make(map[string]*IMGroup)
	r.rooms = make(map[string]*IMRoom)
	r.history = make(map[uint]*entity.History)
	return nil
}

func (r *InMemoryRepository) Close(_ context.Context) {
	r.groups = nil
	r.rooms = nil
	r.history = nil
}

func (r *InMemoryRepository) Migrate(_ context.Context) error {
	// nothing to do
	return nil
}

// groups

func (r *InMemoryRepository) GetGroups(_ context.Context) ([]*entity.Group, error) {
	result := make([]*entity.Group, 0)
	for _, grp := range r.groups {
		if !grp.Group.DeletedAt.Valid {
			grpCopy := grp.Group
			result = append(result, &grpCopy)
		}
	}
	return result, nil
}

func (r *InMemoryRepository) FindGroups(ctx context.Context, minOccupancy uint, maxOccupancy int, anyOfMemberID []uint) ([]string, error) {
	result := make([]string, 0)
	for _, grp := range r.groups {
		if !grp.Group.DeletedAt.Valid {
			if len(grp.Members) >= int(minOccupancy) &&
				(maxOccupancy == -1 || len(grp.Members) <= maxOccupancy) {
				matches := len(anyOfMemberID) == 0
				for _, wantedID := range anyOfMemberID {
					for _, actualMember := range grp.Members {
						if wantedID == actualMember.ID {
							matches = true
						}
					}
				}
				if matches {
					result = append(result, grp.Group.ID)
				}
			}
		}
	}
	return result, nil
}

func (r *InMemoryRepository) AddGroup(_ context.Context, group *entity.Group) (string, error) {
	group.ID = uuid.NewString()
	r.groups[group.ID] = &IMGroup{Group: *group}
	return group.ID, nil
}

func (r *InMemoryRepository) UpdateGroup(_ context.Context, group *entity.Group) error {
	if _, ok := r.groups[group.ID]; ok {
		r.groups[group.ID] = &IMGroup{Group: *group} // this makes a copy
		return nil
	} else {
		return gorm.ErrRecordNotFound
	}
}

func (r *InMemoryRepository) GetGroupByID(_ context.Context, id string) (*entity.Group, error) {
	// allow deleted so history and undelete work
	if result, ok := r.groups[id]; ok {
		grpCopy := result.Group
		return &grpCopy, nil
	} else {
		return &entity.Group{}, gorm.ErrRecordNotFound
	}
}

func (r *InMemoryRepository) SoftDeleteGroupByID(_ context.Context, id string) error {
	if result, ok := r.groups[id]; ok {
		result.Group.DeletedAt = gorm.DeletedAt{
			Time:  r.Now(),
			Valid: true,
		}
		return nil
	} else {
		return gorm.ErrRecordNotFound
	}
}

func (r *InMemoryRepository) UndeleteGroupByID(_ context.Context, id string) error {
	if result, ok := r.groups[id]; ok {
		result.Group.DeletedAt = gorm.DeletedAt{
			Time:  r.Now(),
			Valid: false,
		}
		return nil
	} else {
		return gorm.ErrRecordNotFound
	}
}

// group members

func (r *InMemoryRepository) NewEmptyGroupMembership(_ context.Context, groupID string, attendeeID uint) *entity.GroupMember {
	var m entity.GroupMember
	m.ID = attendeeID
	m.GroupID = groupID
	m.IsInvite = true // default to invite because that's the usual starting point
	return &m
}

func (r *InMemoryRepository) GetGroupMembershipByAttendeeID(_ context.Context, attendeeID uint) (*entity.GroupMember, error) {
	for _, grp := range r.groups {
		for _, gm := range grp.Members {
			if gm.ID == attendeeID {
				gmCopy := gm
				return &gmCopy, nil
			}
		}
	}
	defaultValue := entity.GroupMember{}
	defaultValue.ID = attendeeID
	return &defaultValue, gorm.ErrRecordNotFound
}

func (r *InMemoryRepository) GetGroupMembersByGroupID(_ context.Context, groupID string) ([]*entity.GroupMember, error) {
	if grp, ok := r.groups[groupID]; ok {
		result := make([]*entity.GroupMember, len(grp.Members))
		for i := range grp.Members {
			cpMem := grp.Members[i]
			result[i] = &cpMem
		}
		return result, nil
	} else {
		return []*entity.GroupMember{}, gorm.ErrRecordNotFound
	}
}

func (r *InMemoryRepository) AddGroupMembership(ctx context.Context, gm *entity.GroupMember) error {
	_, err := r.GetGroupMembershipByAttendeeID(ctx, gm.ID)
	if err != gorm.ErrRecordNotFound {
		return gorm.ErrDuplicatedKey
	}
	if grp, ok := r.groups[gm.GroupID]; ok {
		grp.Members = append(grp.Members, *gm)
		return nil
	} else {
		return gorm.ErrForeignKeyViolated
	}
}

func (r *InMemoryRepository) UpdateGroupMembership(ctx context.Context, gm *entity.GroupMember) error {
	_, err := r.GetGroupMembershipByAttendeeID(ctx, gm.ID)
	if err != nil {
		return err
	}
	if grp, ok := r.groups[gm.GroupID]; ok {
		updatedMembers := make([]entity.GroupMember, len(grp.Members))
		for i, m := range grp.Members {
			if m.ID == gm.ID {
				updatedMembers[i] = *gm
			} else {
				updatedMembers[i] = grp.Members[i]
			}
		}
		grp.Members = updatedMembers
		return nil
	} else {
		return fmt.Errorf("internal error - this should not happen, we just read the group")
	}
}

func (r *InMemoryRepository) DeleteGroupMembership(ctx context.Context, attendeeID uint) error {
	current, err := r.GetGroupMembershipByAttendeeID(ctx, attendeeID)
	if err != nil {
		return err
	}
	if grp, ok := r.groups[current.GroupID]; ok {
		updatedMembers := make([]entity.GroupMember, 0)
		for _, m := range grp.Members {
			if m.ID != attendeeID {
				updatedMembers = append(updatedMembers, m)
			}
		}
		grp.Members = updatedMembers
		return nil
	} else {
		return fmt.Errorf("internal error - this should not happen, we just read the group")
	}
}

// rooms

func (r *InMemoryRepository) GetRooms(ctx context.Context) ([]*entity.Room, error) {
	result := make([]*entity.Room, 0)
	for _, rm := range r.rooms {
		if !rm.Room.DeletedAt.Valid {
			rmCopy := rm.Room
			result = append(result, &rmCopy)
		}
	}
	return result, nil
}

func (r *InMemoryRepository) AddRoom(ctx context.Context, room *entity.Room) (string, error) {
	room.ID = uuid.NewString()
	r.rooms[room.ID] = &IMRoom{Room: *room}
	return room.ID, nil
}

func (r *InMemoryRepository) UpdateRoom(ctx context.Context, room *entity.Room) error {
	if _, ok := r.rooms[room.ID]; ok {
		r.rooms[room.ID] = &IMRoom{Room: *room} // this makes a copy
		return nil
	} else {
		return gorm.ErrRecordNotFound
	}
}

func (r *InMemoryRepository) GetRoomByID(ctx context.Context, id string) (*entity.Room, error) {
	// allow deleted so history and undelete work
	if result, ok := r.rooms[id]; ok {
		grpCopy := result.Room
		return &grpCopy, nil
	} else {
		return &entity.Room{}, gorm.ErrRecordNotFound
	}
}

func (r *InMemoryRepository) SoftDeleteRoomByID(ctx context.Context, id string) error {
	if result, ok := r.rooms[id]; ok {
		result.Room.DeletedAt = gorm.DeletedAt{
			Time:  r.Now(),
			Valid: true,
		}
		return nil
	} else {
		return gorm.ErrRecordNotFound
	}
}

func (r *InMemoryRepository) UndeleteRoomByID(ctx context.Context, id string) error {
	if result, ok := r.rooms[id]; ok {
		result.Room.DeletedAt = gorm.DeletedAt{
			Time:  r.Now(),
			Valid: false,
		}
		return nil
	} else {
		return gorm.ErrRecordNotFound
	}
}

// room members

func (r *InMemoryRepository) NewEmptyRoomMembership(_ context.Context, roomID string, attendeeID uint) *entity.RoomMember {
	var m entity.RoomMember
	m.ID = attendeeID
	m.RoomID = roomID
	return &m
}

func (r *InMemoryRepository) GetRoomMembershipByAttendeeID(_ context.Context, attendeeID uint) (*entity.RoomMember, error) {
	for _, room := range r.rooms {
		for _, mem := range room.Members {
			if mem.ID == attendeeID {
				rmCopy := mem
				return &rmCopy, nil
			}
		}
	}
	defaultValue := entity.RoomMember{}
	defaultValue.ID = attendeeID
	return &defaultValue, gorm.ErrRecordNotFound
}

func (r *InMemoryRepository) GetRoomMembersByRoomID(ctx context.Context, roomID string) ([]*entity.RoomMember, error) {
	if rm, ok := r.rooms[roomID]; ok {
		result := make([]*entity.RoomMember, len(rm.Members))
		for i := range rm.Members {
			cpRoom := rm.Members[i]
			result[i] = &cpRoom
		}
		return result, nil
	} else {
		return []*entity.RoomMember{}, gorm.ErrRecordNotFound
	}
}

func (r *InMemoryRepository) AddRoomMembership(ctx context.Context, rm *entity.RoomMember) error {
	_, err := r.GetRoomMembershipByAttendeeID(ctx, rm.ID)
	if err != gorm.ErrRecordNotFound {
		return gorm.ErrDuplicatedKey
	}
	if room, ok := r.rooms[rm.RoomID]; ok {
		room.Members = append(room.Members, *rm)
		return nil
	} else {
		return gorm.ErrForeignKeyViolated
	}
}

func (r *InMemoryRepository) UpdateRoomMembership(ctx context.Context, rm *entity.RoomMember) error {
	_, err := r.GetRoomMembershipByAttendeeID(ctx, rm.ID)
	if err != nil {
		return err
	}
	if room, ok := r.rooms[rm.RoomID]; ok {
		updatedMembers := make([]entity.RoomMember, len(room.Members))
		for i, m := range room.Members {
			if m.ID == rm.ID {
				updatedMembers[i] = *rm
			} else {
				updatedMembers[i] = room.Members[i]
			}
		}
		room.Members = updatedMembers
		return nil
	} else {
		return fmt.Errorf("internal error - this should not happen, we just read the room")
	}
}

func (r *InMemoryRepository) DeleteRoomMembership(ctx context.Context, attendeeID uint) error {
	current, err := r.GetRoomMembershipByAttendeeID(ctx, attendeeID)
	if err != nil {
		return err
	}
	if grp, ok := r.rooms[current.RoomID]; ok {
		updatedMembers := make([]entity.RoomMember, 0)
		for _, m := range grp.Members {
			if m.ID != attendeeID {
				updatedMembers = append(updatedMembers, m)
			}
		}
		grp.Members = updatedMembers
		return nil
	} else {
		return fmt.Errorf("internal error - this should not happen, we just read the room")
	}
}

// history

func (r *InMemoryRepository) RecordHistory(_ context.Context, h *entity.History) error {
	newID := uint(atomic.AddUint32(&r.idSequence, 1))
	h.ID = newID
	r.history[newID] = h
	return nil
}

// GetHistoryByID is only offered for testing, and only on the in memory db.
func (r *InMemoryRepository) GetHistoryByID(_ context.Context, id uint) (*entity.History, error) {
	if h, ok := r.history[id]; ok {
		return h, nil
	} else {
		return &entity.History{}, fmt.Errorf("cannot get history entry %d - not present", id)
	}
}
