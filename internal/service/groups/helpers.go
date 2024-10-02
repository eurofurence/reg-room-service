package groupservice

import (
	"context"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
	"github.com/eurofurence/reg-room-service/internal/application/common"
	"github.com/eurofurence/reg-room-service/internal/repository/config"
	"github.com/eurofurence/reg-room-service/internal/repository/downstreams/attendeeservice"
)

func (g *groupService) loggedInUserValidRegistrationBadgeNo(ctx context.Context) (attendeeservice.Attendee, error) {
	myRegIDs, err := g.AttSrv.ListMyRegistrationIds(ctx)
	if err != nil {
		aulogging.WarnErrf(ctx, err, "failed to obtain registrations for currently logged in user: %s", err.Error())
		return attendeeservice.Attendee{}, common.NewBadGateway(ctx, common.DownstreamAttSrv, common.Details("downstream error when contacting attendee service"))
	}
	if len(myRegIDs) == 0 {
		aulogging.InfoErr(ctx, err, "currently logged in user has no registrations - cannot be in a group")
		return attendeeservice.Attendee{}, common.NewForbidden(ctx, common.NoSuchAttendee, common.Details("you do not have a valid registration"))
	}
	myID := myRegIDs[0]

	if err := g.checkAttending(ctx, myID); err != nil {
		return attendeeservice.Attendee{}, err
	}

	attendee, err := g.AttSrv.GetAttendee(ctx, myID)
	if err != nil {
		return attendeeservice.Attendee{}, err
	}
	// ensure ID set in Attendee
	attendee.ID = myID

	return attendee, nil
}

func (g *groupService) checkAttending(ctx context.Context, badgeNo int64) error {
	myStatus, err := g.AttSrv.GetStatus(ctx, badgeNo)
	if err != nil {
		aulogging.WarnErrf(ctx, err, "failed to obtain status for badge number %d: %s", badgeNo, err.Error())
		return common.NewBadGateway(ctx, common.DownstreamAttSrv, common.Details("downstream error when contacting attendee service"))
	}

	switch myStatus {
	case attendeeservice.StatusApproved, attendeeservice.StatusPartiallyPaid, attendeeservice.StatusPaid, attendeeservice.StatusCheckedIn:
		return nil
	default:
		return common.NewForbidden(ctx, common.NotAttending, common.Details("registration is not in attending status"))
	}
}

func maxGroupSize() uint {
	conf, err := config.GetApplicationConfig()
	if err != nil {
		panic("configuration not loaded before call to maxGroupSize() - this is a bug")
	}
	return conf.Service.MaxGroupSize
}

func allowedFlags() []string {
	conf, err := config.GetApplicationConfig()
	if err != nil {
		panic("configuration not loaded before call to allowedFlags() - this is a bug")
	}
	return conf.Service.GroupFlags
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
