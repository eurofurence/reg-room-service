package roomservice

import (
	"context"
	"errors"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"github.com/eurofurence/reg-room-service/internal/application/common"
	"github.com/eurofurence/reg-room-service/internal/entity"
	"github.com/eurofurence/reg-room-service/internal/repository/downstreams"
	"github.com/eurofurence/reg-room-service/internal/repository/downstreams/attendeeservice"
	"github.com/eurofurence/reg-room-service/internal/service/rbac"
	"gorm.io/gorm"
)

func (r *roomService) AddOccupantToRoom(ctx context.Context, roomID string, badgeNumber int64) error {
	validator, err := rbac.NewValidator(ctx)
	if err != nil {
		aulogging.ErrorErrf(ctx, err, "Could not retrieve RBAC validator from context. [error]: %v", err)
		return errCouldNotGetValidator(ctx)
	}

	if validator.IsAdmin() || validator.IsAPITokenCall() {
		room, existingMembership, err := r.roomMembershipExisting(ctx, roomID, badgeNumber) // existingMembership may be nil if not exists
		if err != nil {
			return err
		}

		occupant, err := r.validateRequestedAttendee(ctx, badgeNumber)
		if err != nil {
			return err
		}
		if err := r.checkAttending(ctx, badgeNumber); err != nil {
			return err
		}

		if existingMembership != nil {
			if room.ID != existingMembership.RoomID {
				return common.NewConflict(ctx, common.RoomOccupantConflict, common.Details("this attendee is already in another room"))
			} else {
				return common.NewConflict(ctx, common.RoomOccupantDuplicate, common.Details("this attendee is already in this room"))
			}
		}

		if err := r.checkRoomFull(ctx, roomID, room.Size); err != nil {
			return err
		}

		newMembership := r.DB.NewEmptyRoomMembership(ctx, roomID, badgeNumber)
		newMembership.Nickname = occupant.Nickname

		if err := r.DB.AddRoomMembership(ctx, newMembership); err != nil {
			return errRoomWrite(ctx, err.Error())
		}

		// TODO now check room size again, remove again if exceeded

		return nil
	} else {
		return errNotAdminOrApiToken(ctx, roomID, "(not loaded)")
	}
}

func (r *roomService) RemoveOccupantFromRoom(ctx context.Context, roomID string, badgeNumber int64) error {
	validator, err := rbac.NewValidator(ctx)
	if err != nil {
		aulogging.ErrorErrf(ctx, err, "Could not retrieve RBAC validator from context. [error]: %v", err)
		return errCouldNotGetValidator(ctx)
	}

	if validator.IsAdmin() || validator.IsAPITokenCall() {
		room, existingMembership, err := r.roomMembershipExisting(ctx, roomID, badgeNumber) // existingMembership may be nil if not exists
		if err != nil {
			return err
		}

		if _, err := r.validateRequestedAttendee(ctx, badgeNumber); err != nil {
			return err
		}
		// not attending is ok, we are removing the attendee

		if existingMembership == nil {
			return common.NewNotFound(ctx, common.RoomOccupantNotFound, common.Details("this attendee is not in any room"))
		}

		if room.ID != existingMembership.RoomID {
			return common.NewConflict(ctx, common.RoomOccupantConflict, common.Details("this attendee is in a different room"))
		}

		if err := r.DB.DeleteRoomMembership(ctx, badgeNumber); err != nil {
			return errRoomWrite(ctx, err.Error())
		}

		return nil
	} else {
		return errNotAdminOrApiToken(ctx, roomID, "(not loaded)")
	}
}

// --- helpers ---

func (r *roomService) checkRoomFull(ctx context.Context, roomID string, roomSize int64) error {
	memberIDs, err := r.DB.GetRoomMembersByRoomID(ctx, roomID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// empty room is acceptable
			return nil
		} else {
			return errRoomRead(ctx, err.Error())
		}
	}

	if len(memberIDs) >= int(roomSize) {
		return errRoomFull(ctx)
	}

	return nil
}

func (r *roomService) roomMembershipExisting(ctx context.Context, roomID string, badgeNumber int64) (*entity.Room, *entity.RoomMember, error) {
	room, err := r.DB.GetRoomByID(ctx, roomID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, errRoomNotFound(ctx)
		}

		return nil, nil, errRoomRead(ctx, err.Error())
	}

	member, err := r.DB.GetRoomMembershipByAttendeeID(ctx, badgeNumber)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// maybe ok, does not have a membership record
			return room, nil, nil
		} else {
			return room, nil, errInternal(ctx, err.Error())
		}
	}

	return room, member, nil
}

func (r *roomService) validateRequestedAttendee(ctx context.Context, badgeNo int64) (attendeeservice.Attendee, error) {
	if badgeNo <= 0 {
		return attendeeservice.Attendee{}, common.NewBadRequest(ctx, common.RoomDataInvalid, common.Details("attendee badge number must be positive integer"))
	}

	attendee, err := r.AttSrv.GetAttendee(ctx, badgeNo)
	if err != nil {
		if errors.Is(err, downstreams.ErrDownStreamNotFound) {
			return attendeeservice.Attendee{}, common.NewNotFound(ctx, common.NoSuchAttendee, common.Details("no such attendee"))
		} else {
			aulogging.WarnErrf(ctx, err, "failed to query for attendee with badge number %d: %s", badgeNo, err.Error())
			return attendeeservice.Attendee{}, common.NewBadGateway(ctx, common.DownstreamAttSrv, common.Details("failed to look up invited attendee - internal error, see logs for details"))
		}
	}

	return attendee, nil
}

func (r *roomService) checkAttending(ctx context.Context, badgeNo int64) error {
	status, err := r.AttSrv.GetStatus(ctx, badgeNo)
	if err != nil {
		aulogging.WarnErrf(ctx, err, "failed to obtain status for badge number %d: %s", badgeNo, err.Error())
		return common.NewBadGateway(ctx, common.DownstreamAttSrv, common.Details("downstream error when contacting attendee service"))
	}

	switch status {
	case attendeeservice.StatusApproved, attendeeservice.StatusPartiallyPaid, attendeeservice.StatusPaid, attendeeservice.StatusCheckedIn:
		return nil
	default:
		return common.NewConflict(ctx, common.NotAttending, common.Details("registration is not in attending status"))
	}
}
