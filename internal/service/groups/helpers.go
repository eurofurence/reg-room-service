package groupservice

import (
	"context"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"github.com/eurofurence/reg-room-service/internal/config"
	apierrors "github.com/eurofurence/reg-room-service/internal/errors"
	"github.com/eurofurence/reg-room-service/internal/repository/downstreams/attendeeservice"
	"github.com/eurofurence/reg-room-service/internal/web/common"
)

func (g *groupService) loggedInUserValidRegistrationBadgeNo(ctx context.Context) (int64, error) {
	myRegIDs, err := g.AttSrv.ListMyRegistrationIds(ctx)
	if err != nil {
		aulogging.WarnErrf(ctx, err, "failed to obtain registrations for currently logged in user: %s", err.Error())
		return 0, apierrors.NewInternalServerError(common.InternalErrorMessage, "downstream error when contacting attendee service")
	}
	if len(myRegIDs) == 0 {
		aulogging.InfoErr(ctx, err, "currently logged in user has no registrations - cannot create a group")
		return 0, apierrors.NewNotFound(common.NoSuchAttendee, "you do not have a valid registration")
	}
	myID := myRegIDs[0]

	if err := g.checkAttending(ctx, myID); err != nil {
		return 0, err
	}
	return myID, nil
}

func (g *groupService) checkAttending(ctx context.Context, badgeNo int64) error {
	myStatus, err := g.AttSrv.GetStatus(ctx, badgeNo)
	if err != nil {
		aulogging.WarnErrf(ctx, err, "failed to obtain status for badge number %d: %s", badgeNo, err.Error())
		return apierrors.NewInternalServerError(common.InternalErrorMessage, "downstream error when contacting attendee service")
	}

	switch myStatus {
	case attendeeservice.StatusApproved, attendeeservice.StatusPartiallyPaid, attendeeservice.StatusPaid, attendeeservice.StatusCheckedIn:
		return nil
	default:
		return apierrors.NewForbidden(common.NotAttending, "registration is not in attending status")
	}
}

func maxGroupSize() uint {
	conf, err := config.GetApplicationConfig()
	if err != nil {
		panic("configuration not loaded before call to maxGroupSize() - this is a bug")
	}
	return conf.Service.MaxGroupSize
}
