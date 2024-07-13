package groupservice

import (
	"context"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
	"github.com/eurofurence/reg-room-service/internal/config"
	apierrors "github.com/eurofurence/reg-room-service/internal/errors"
	"github.com/eurofurence/reg-room-service/internal/repository/downstreams/attendeeservice"
	"github.com/eurofurence/reg-room-service/internal/web/common"
)

func (g *groupService) loggedInUserValidRegistrationBadgeNo(ctx context.Context) (int64, error) {
	myRegIDs, err := g.AttSrv.ListMyRegistrationIds(ctx)
	if err != nil {
		aulogging.WarnErrf(ctx, err, "failed to obtain registrations for currently logged in user: %s", err.Error())
		return 0, apierrors.NewBadGateway(common.DownstreamAttSrv, "downstream error when contacting attendee service")
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
		return apierrors.NewBadGateway(common.DownstreamAttSrv, "downstream error when contacting attendee service")
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

func publicInfo(grp *modelsv1.Group, myID int32) *modelsv1.Group {
	if grp == nil {
		return nil
	}

	return &modelsv1.Group{
		ID:          grp.ID,
		Name:        grp.Name,
		Flags:       grp.Flags,
		Comments:    nil,
		MaximumSize: grp.MaximumSize,
		Owner:       0,
		Members:     maskMembers(grp.Members, myID),
		Invites:     filterInvites(grp.Invites, myID),
	}
}

func maskMembers(members []modelsv1.Member, myID int32) []modelsv1.Member {
	result := make([]modelsv1.Member, 0)
	for _, member := range members {
		if member.ID == myID {
			result = append(result, member)
		} else {
			result = append(result, modelsv1.Member{
				ID: 0, // cannot see who other members are, just that there is one
			})
		}
	}
	return result
}

func filterInvites(members []modelsv1.Member, myID int32) []modelsv1.Member {
	result := make([]modelsv1.Member, 0)
	for _, member := range members {
		if member.ID == myID {
			// can only see myself if invited
			// TODO filter in result list of groups if banned instead
			result = append(result, member)
		}
	}
	return result
}

func groupContains(group *modelsv1.Group, memberID int32) bool {
	if group != nil {
		for _, member := range group.Members {
			if member.ID == memberID {
				return true
			}
		}
	}
	return false
}

func groupInvited(group *modelsv1.Group, invitedMemberID int32) bool {
	if group != nil {
		for _, member := range group.Invites {
			if member.ID == invitedMemberID {
				return true
			}
		}
	}
	return false
}

func groupHasFlag(group *modelsv1.Group, wantedFlag string) bool {
	if group != nil {
		for _, flag := range group.Flags {
			if flag == wantedFlag {
				return true
			}
		}
	}
	return false
}
