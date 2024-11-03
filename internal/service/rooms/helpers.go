package roomservice

import (
	"context"
	"errors"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"github.com/eurofurence/reg-room-service/internal/application/common"
	"github.com/eurofurence/reg-room-service/internal/repository/downstreams"
	"github.com/eurofurence/reg-room-service/internal/repository/downstreams/attendeeservice"
)

func (r *roomService) loggedInUserValidRegistrationBadgeNo(ctx context.Context) (attendeeservice.Attendee, error) {
	myRegIDs, err := r.AttSrv.ListMyRegistrationIds(ctx)
	if err != nil {
		aulogging.WarnErrf(ctx, err, "failed to obtain registrations for currently logged in user: %s", err.Error())
		return attendeeservice.Attendee{}, common.NewBadGateway(ctx, common.DownstreamAttSrv, common.Details("downstream error when contacting attendee service"))
	}
	if len(myRegIDs) == 0 {
		aulogging.InfoErr(ctx, err, "currently logged in user has no registrations - cannot be in a group")
		return attendeeservice.Attendee{}, common.NewForbidden(ctx, common.NoSuchAttendee, common.Details("you do not have a valid registration"))
	}
	myID := myRegIDs[0]

	if err := r.checkAttending(ctx, myID); err != nil {
		return attendeeservice.Attendee{}, err
	}

	attendee, err := r.AttSrv.GetAttendee(ctx, myID)
	if err != nil {
		return attendeeservice.Attendee{}, err
	}
	// ensure ID set in Attendee
	attendee.ID = myID

	return attendee, nil
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
