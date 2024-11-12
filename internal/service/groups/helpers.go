package groupservice

import (
	"context"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
	"github.com/eurofurence/reg-room-service/internal/application/common"
	"github.com/eurofurence/reg-room-service/internal/repository/config"
	"github.com/eurofurence/reg-room-service/internal/repository/downstreams/attendeeservice"
)

// filterGroupAndFieldVisibilityForAttendee takes a fully populated group model, and filters it using
// the visibility rules for a regular attendee (non-admin).
//
// Warning: **may return nil** if the attendee is not allowed to see the group at all. In this case, the
// calling function should return a proper error, or just omit the group from the result listing.
func (g *groupService) filterGroupAndFieldVisibilityForAttendee(group *modelsv1.Group, attendee attendeeservice.Attendee) *modelsv1.Group {
	if group == nil || attendee.ID <= 0 {
		return nil
	}

	if group.Owner == attendee.ID {
		// owner can see all group info
		return group
	} else if groupContains(group, attendee.ID) {
		// group members can see all group info, but no invites
		return &modelsv1.Group{
			ID:          group.ID,
			Name:        group.Name,
			Flags:       group.Flags,
			Comments:    group.Comments,
			MaximumSize: group.MaximumSize,
			Owner:       group.Owner,
			Members:     group.Members,
			Invites:     nil,
		}
	} else if groupInvited(group, attendee.ID) {
		// group invitees can see masked members and only their own invite
		return &modelsv1.Group{
			ID:          group.ID,
			Name:        group.Name,
			Flags:       group.Flags,
			Comments:    nil, // hide comment
			MaximumSize: group.MaximumSize,
			Owner:       group.Owner,
			Members:     maskMembers(group.Members, attendee.ID),
			Invites:     filterInvites(group.Invites, attendee.ID),
		}
	} else if groupHasFlag(group, "public") {
		// non-members get even less information
		return &modelsv1.Group{
			ID:          group.ID,
			Name:        group.Name,
			Flags:       group.Flags,
			Comments:    nil, // hide comment
			MaximumSize: group.MaximumSize,
			Owner:       group.Owner,
			Members:     maskMembers(group.Members, attendee.ID),
			Invites:     nil,
		}
	} else {
		return nil
	}
}

// loggedInUserValidRegistration obtains the attendee record for the currently logged-in user,
// but only if they have a valid registration with attending status.
//
// It will return a suitable common.APIError if no registration or not in attending status
// (or if attendee service fails to respond).
//
// Here, admins are treated exactly the same as normal users.
func (g *groupService) loggedInUserValidRegistration(ctx context.Context) (attendeeservice.Attendee, error) {
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
	status, err := g.AttSrv.GetStatus(ctx, badgeNo)
	if err != nil {
		aulogging.WarnErrf(ctx, err, "failed to obtain status for badge number %d: %s", badgeNo, err.Error())
		return common.NewBadGateway(ctx, common.DownstreamAttSrv, common.Details("downstream error when contacting attendee service"))
	}

	switch status {
	case attendeeservice.StatusApproved, attendeeservice.StatusPartiallyPaid, attendeeservice.StatusPaid, attendeeservice.StatusCheckedIn:
		return nil
	default:
		return common.NewForbidden(ctx, common.NotAttending, common.Details("registration is not in attending status"))
	}
}

func maxGroupSize() int64 {
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

func maskMembers(members []modelsv1.Member, myID int64) []modelsv1.Member {
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

func filterInvites(members []modelsv1.Member, myID int64) []modelsv1.Member {
	result := make([]modelsv1.Member, 0)
	for _, member := range members {
		if member.ID == myID {
			// can only see myself if invited
			result = append(result, member)
		}
	}
	return result
}

func groupContains(group *modelsv1.Group, memberID int64) bool {
	if group != nil {
		for _, member := range group.Members {
			if member.ID == memberID {
				return true
			}
		}
	}
	return false
}

func groupInvited(group *modelsv1.Group, invitedMemberID int64) bool {
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

func groupLessByName(left *modelsv1.Group, right *modelsv1.Group) bool {
	if left == nil || right == nil {
		return left == nil && right != nil
	}
	return left.Name < right.Name
}
