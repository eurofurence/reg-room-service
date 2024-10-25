package roomservice

import (
	"context"
	"errors"
	"fmt"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
	"github.com/eurofurence/reg-room-service/internal/application/common"
	"github.com/eurofurence/reg-room-service/internal/controller/v1/util"
	"github.com/eurofurence/reg-room-service/internal/entity"
	"github.com/eurofurence/reg-room-service/internal/repository/config"
	"github.com/eurofurence/reg-room-service/internal/service/rbac"
	"gorm.io/gorm"
	"net/url"
	"slices"
	"sort"
	"strings"
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
	validator, err := rbac.NewValidator(ctx)
	if err != nil {
		aulogging.ErrorErrf(ctx, err, "Could not retrieve RBAC validator from context. [error]: %v", err)
		return nil, errCouldNotGetValidator(ctx)
	}

	if validator.IsAdmin() || validator.IsAPITokenCall() {
		room, err := r.DB.GetRoomByID(ctx, roomID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errRoomNotFound(ctx)
			}

			return nil, errRoomRead(ctx, err.Error())
		}

		roomMembers, err := r.DB.GetRoomMembersByRoomID(ctx, roomID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// acceptable, empty room
			} else {
				return nil, errRoomRead(ctx, err.Error())
			}
		}

		return &modelsv1.Room{
			ID:        room.ID,
			Name:      room.Name,
			Flags:     aggregateFlags(room.Flags),
			Comments:  common.ToOmitEmpty(room.Comments),
			Size:      room.Size,
			Occupants: toOccupants(roomMembers),
		}, nil
	} else {
		return nil, errNotAdminOrApiToken(ctx, roomID, "(not loaded)")
	}
}

func (r *roomService) CreateRoom(ctx context.Context, room *modelsv1.RoomCreate) (string, error) {
	validator, err := rbac.NewValidator(ctx)
	if err != nil {
		aulogging.ErrorErrf(ctx, err, "Could not retrieve RBAC validator from context. [error]: %v", err)
		return "", errCouldNotGetValidator(ctx)
	}

	if validator.IsAdmin() || validator.IsAPITokenCall() {
		validation := validateRoomCreate(room)
		if len(validation) > 0 {
			return "", common.NewBadRequest(ctx, common.RoomDataInvalid, validation)
		}

		roomID, err := r.DB.AddRoom(ctx, &entity.Room{
			Name:     room.Name,
			Flags:    fmt.Sprintf(",%s,", strings.Join(room.Flags, ",")),
			Comments: common.Deref(room.Comments),
			Size:     room.Size,
		})

		if err != nil {
			return "", err
		}

		return roomID, nil
	} else {
		return "", errNotAdminOrApiToken(ctx, "(new)", room.Name)
	}
}

func (r *roomService) UpdateRoom(ctx context.Context, room *modelsv1.Room) error {
	validator, err := rbac.NewValidator(ctx)
	if err != nil {
		aulogging.ErrorErrf(ctx, err, "Could not retrieve RBAC validator from context. [error]: %v", err)
		return errCouldNotGetValidator(ctx)
	}

	if validator.IsAdmin() || validator.IsAPITokenCall() {
		dbRoom, err := r.DB.GetRoomByID(ctx, room.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errRoomNotFound(ctx)
			} else {
				return errRoomRead(ctx, err.Error())
			}
		}

		validation := validateRoom(room)
		if len(validation) > 0 {
			return common.NewBadRequest(ctx, common.RoomDataInvalid, validation)
		}

		occupants, err := r.DB.GetRoomMembersByRoomID(ctx, room.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				occupants = make([]*entity.RoomMember, 0)
			} else {
				return errRoomRead(ctx, err.Error())
			}
		}
		if int(room.Size) < len(occupants) {
			return common.NewConflict(ctx, common.RoomSizeTooSmall, common.Details("the room cannot be resized, too many occupants for new size"))
		}

		// check for name conflicts
		if dbRoom.Name != room.Name {
			matchingIDs, err := r.DB.FindRooms(ctx, room.Name, 0, -1, 0, 0, nil)
			if err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					return errRoomRead(ctx, err.Error())
				}
			}

			if len(matchingIDs) > 0 {
				return common.NewConflict(ctx, common.RoomDataDuplicate, common.Details("another room with this name already exists"))
			}
		}

		// do not touch fields that we do not wish to change, like createdAt or referenced occupants
		dbRoom.Name = room.Name
		dbRoom.Flags = fmt.Sprintf(",%s,", strings.Join(room.Flags, ","))
		dbRoom.Comments = common.Deref(room.Comments)
		dbRoom.Size = room.Size

		return r.DB.UpdateRoom(ctx, dbRoom)
	} else {
		return errNotAdminOrApiToken(ctx, room.ID, "(not loaded)")
	}
}

func (r *roomService) DeleteRoom(ctx context.Context, roomID string) error {
	validator, err := rbac.NewValidator(ctx)
	if err != nil {
		aulogging.ErrorErrf(ctx, err, "Could not retrieve RBAC validator from context. [error]: %v", err)
		return errCouldNotGetValidator(ctx)
	}

	if validator.IsAdmin() || validator.IsAPITokenCall() {
		_, err := r.DB.GetRoomByID(ctx, roomID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errRoomNotFound(ctx)
			}

			aulogging.Warnf(ctx, "failed to read room %s from db: %s", url.PathEscape(roomID), err.Error())
			return errRoomRead(ctx, "error retrieving room - see logs for details")
		}

		members, err := r.DB.GetRoomMembersByRoomID(ctx, roomID)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return errInternal(ctx, "failed to read room members during delete")
			}
			// empty room is expected
		}
		if len(members) > 0 {
			aulogging.Infof(ctx, "attempt to delete non-empty room %s - rejected", url.PathEscape(roomID))
			return errRoomNotEmpty(ctx)
		}

		if err := r.DB.DeleteRoomByID(ctx, roomID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// may have been deleted concurrently - give correct error
				return errRoomNotFound(ctx)
			}

			aulogging.ErrorErrf(ctx, err, "unexpected error while deleting room %s: %s", url.PathEscape(roomID), err.Error())
			return errInternal(ctx, "unexpected error occurred during deletion of group")
		}

		return nil
	} else {
		return errNotAdminOrApiToken(ctx, roomID, "(not loaded)")
	}
}

// --- helpers ---

func validateRoomCreate(room *modelsv1.RoomCreate) url.Values {
	return validate(room.Name, room.Flags)
}

func validateRoom(room *modelsv1.Room) url.Values {
	return validate(room.Name, room.Flags)
}

func validate(name string, flags []string) url.Values {
	result := url.Values{}
	if len(name) == 0 {
		result.Set("name", "room name cannot be empty")
	}
	if len(name) > 50 {
		result.Set("name", "room name too long, max 50 characters")
	}
	allowed := allowedFlags()
	for _, flag := range flags {
		if !util.SliceContains(flag, allowed) {
			result.Set("flags", fmt.Sprintf("no such flag '%s'", url.PathEscape(flag)))
		}
	}
	return result
}

func allowedFlags() []string {
	conf, err := config.GetApplicationConfig()
	if err != nil {
		panic("configuration not loaded before call to allowedFlags() - this is a bug")
	}
	return conf.Service.RoomFlags
}

func aggregateFlags(input string) []string {
	tags := strings.Split(input, ",")
	tags = slices.DeleteFunc(tags, func(s string) bool {
		return s == ""
	})

	if len(tags) == 0 {
		return make([]string, 0)
	}

	slices.Sort(tags)
	return tags
}

func toOccupants(roomMembers []*entity.RoomMember) []modelsv1.Member {
	members := make([]modelsv1.Member, 0)
	for _, m := range roomMembers {
		if m == nil {
			continue
		}

		member := modelsv1.Member{
			ID:       m.ID,
			Nickname: m.Nickname,
		}
		if m.AvatarURL != "" {
			member.Avatar = &m.AvatarURL
		}

		members = append(members, member)
	}

	sort.Slice(members, func(i int, j int) bool {
		return members[i].ID < members[j].ID
	})

	return members
}

// --- errors ---

func errNotAdminOrApiToken(ctx context.Context, uuid string, name string) error {
	subject := common.GetSubject(ctx)
	aulogging.Warnf(ctx, "unauthorized attempt to access admin-only room %s (%s) by %s", uuid, url.PathEscape(name), subject)
	return common.NewForbidden(ctx, common.AuthForbidden, common.Details("you are not authorized for this operation - the attempt has been logged"))
}

func errRoomNotFound(ctx context.Context) error {
	return common.NewNotFound(ctx, common.RoomIDNotFound, common.Details("room does not exist"))
}

func errRoomNotEmpty(ctx context.Context) error {
	return common.NewConflict(ctx, common.RoomNotEmpty, common.Details("room is not empty and room deletion is a dangerous operation - please remove all occupants first to ensure you really mean this (also prevents possible problems with concurrent updates)"))
}

func errCouldNotGetValidator(ctx context.Context) error {
	return common.NewInternalServerError(ctx, common.InternalErrorMessage, common.Details("unexpected error when parsing user claims"))
}

func errRoomRead(ctx context.Context, details string) error {
	return common.NewInternalServerError(ctx, common.RoomReadError, common.Details(details))
}

func errRoomWrite(ctx context.Context, details string) error {
	return common.NewInternalServerError(ctx, common.RoomWriteError, common.Details(details))
}

func errInternal(ctx context.Context, details string) error {
	return common.NewInternalServerError(ctx, common.InternalErrorMessage, common.Details(details))
}
