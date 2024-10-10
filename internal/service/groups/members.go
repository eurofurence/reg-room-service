package groupservice

import (
	"context"
	"errors"
	"github.com/eurofurence/reg-room-service/internal/application/common"
	"gorm.io/gorm"
)

// AddGroupMemberParams is the request type for the AddMemberToGroup operation.
//
// See OpenAPI spec for more details.
type AddGroupMemberParams struct {
	// GroupID is the ID of the group where a user should be added
	GroupID string
	// BadgeNumber is the registration number of a user
	BadgeNumber int64
	// Nickname is the nickname of a registered user that should receive
	// an invitation Email.
	Nickname string
	// Code is the invite code that can be used to join a group.
	Code string
	// Force is an admin only flag that allows to bypass the
	// validations.
	Force bool
}

func (g *groupService) AddMemberToGroup(ctx context.Context, req *AddGroupMemberParams) error {
	// TODO also remove ban if appropriate

	grp, err := g.DB.GetGroupByID(ctx, req.GroupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.NewNotFound(ctx, common.GroupIDNotFound, common.Details("this group does not exist"))
		} else {
			return errGroupRead(ctx, err.Error())
		}
	}

	gm, err := g.DB.GetGroupMembershipByAttendeeID(ctx, req.BadgeNumber)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// TODO handle new membership - this is very incomplete - auth etc.

			gm := g.DB.NewEmptyGroupMembership(ctx, req.GroupID, req.BadgeNumber, req.Nickname)

			// TODO new membership, decide whether invite or direct add

			err := g.DB.AddGroupMembership(ctx, gm)
			if err != nil {
				return errGroupWrite(ctx, err.Error())
			}

			return nil
		} else {
			return errGroupRead(ctx, err.Error())
		}
	}

	if grp.ID != gm.GroupID {
		return common.NewConflict(ctx, common.GroupMemberConflict, common.Details("this attendee is already invited to another group or in another group"))
	}

	if !gm.IsInvite {
		return common.NewConflict(ctx, common.GroupMemberDuplicate, common.Details("this attendee is already a member of this group"))
	}

	// TODO handle confirming an invitation - very incomplete - auth etc.

	gm.IsInvite = false

	err = g.DB.UpdateGroupMembership(ctx, gm)
	if err != nil {
		return errGroupWrite(ctx, err.Error())
	}

	return nil
}

// RemoveGroupMemberParams is the request type for the RemoveMemberFromGroup operation.
//
// See OpenAPI spec for more details.
type RemoveGroupMemberParams struct {
	// GroupID is the ID of the group where a user should be added
	GroupID string
	// BadgeNumber is the registration number of a user
	BadgeNumber int64
	// AutoDeny future invitations (effectively creates or removes a ban)
	AutoDeny bool
}

func (g *groupService) RemoveMemberFromGroup(ctx context.Context, req *RemoveGroupMemberParams) error {
	// TODO this is very incomplete - auth etc.

	grp, err := g.DB.GetGroupByID(ctx, req.GroupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.NewNotFound(ctx, common.GroupIDNotFound, common.Details("this group does not exist"))
		} else {
			return errGroupRead(ctx, err.Error())
		}
	}

	gm, err := g.DB.GetGroupMembershipByAttendeeID(ctx, req.BadgeNumber)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.NewNotFound(ctx, common.GroupMemberNotFound, common.Details("this attendee is not in any group"))
		} else {
			return errInternal(ctx, err.Error())
		}
	}

	if grp.ID != gm.GroupID {
		return common.NewConflict(ctx, common.GroupMemberConflict, common.Details("this attendee is invited to a different group or in a different group"))
	}

	// TODO also add ban if requested

	err = g.DB.DeleteGroupMembership(ctx, req.BadgeNumber)
	if err != nil {
		return errGroupWrite(ctx, err.Error())
	}

	return nil
}
