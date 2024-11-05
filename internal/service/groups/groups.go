package groupservice

import (
	"context"
	"errors"
	"fmt"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"github.com/eurofurence/reg-room-service/internal/controller/v1/util"
	"net/url"
	"slices"
	"sort"
	"strings"

	"gorm.io/gorm"

	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
	"github.com/eurofurence/reg-room-service/internal/application/common"
	"github.com/eurofurence/reg-room-service/internal/entity"
	"github.com/eurofurence/reg-room-service/internal/service/rbac"
)

// FindMyGroup finds the group containing the currently logged in attendee.
//
// This even works for admins, who are treated exactly as if they were a regular user.
//
// Uses the attendee service to look up the badge number and the registration status.
func (g *groupService) FindMyGroup(ctx context.Context) (*modelsv1.Group, error) {
	attendee, err := g.loggedInUserValidRegistration(ctx)
	if err != nil {
		return nil, err
	}

	groups, err := g.findGroupsFullAccess(ctx, 0, -1, []int64{attendee.ID})
	if err != nil {
		return nil, err
	}

	if len(groups) == 0 {
		return nil, errNoGroup(ctx)
	}
	if len(groups) > 1 {
		return nil, errInternal(ctx, "multiple group memberships found - this is a bug")
	}

	filtered := g.filterGroupAndFieldVisibilityForAttendee(groups[0], attendee)
	if filtered == nil {
		// this should never happen because then the group would not have been found by findGroupsFullAccess.
		aulogging.Warnf(ctx, "group %s found as containing attendee %d but got filtered - possible bug", groups[0].ID, attendee.ID)
		return nil, errNoAccess(ctx)
	}

	return filtered, nil
}

// FindGroups finds groups by size (number of members) and member badge numbers.
//
// A group matches if its size is in the range (maxSize -1 means no limit), and if it
// contains at least one of the specified badge numbers (if memberIDs is not empty).
//
// Admin or Api Key authorization: can see all groups.
//
// Normal users: can only see groups visible to them. If public groups are enabled in configuration,
// this means all groups that are public and from which the user wasn't banned. Not all fields
// will be filled in the results to protect the privacy of group members.
func (g *groupService) FindGroups(ctx context.Context, minSize uint, maxSize int, memberIDs []int64) ([]*modelsv1.Group, error) {
	validator, err := rbac.NewValidator(ctx)
	if err != nil {
		aulogging.ErrorErrf(ctx, err, "Could not retrieve RBAC validator from context. [error]: %v", err)
		return make([]*modelsv1.Group, 0), errCouldNotGetValidator(ctx)
	}

	if validator.IsAdmin() || validator.IsAPITokenCall() {
		return g.findGroupsFullAccess(ctx, minSize, maxSize, memberIDs)
	} else if validator.IsUser() {
		result := make([]*modelsv1.Group, 0)

		// ensure attending registration
		attendee, err := g.loggedInUserValidRegistration(ctx)
		if err != nil {
			return result, err
		}

		// normal users cannot specify memberIDs to filter for - ignore if set
		unchecked, err := g.findGroupsFullAccess(ctx, minSize, maxSize, nil)
		if err != nil {
			return result, err
		}

		for _, group := range unchecked {
			filtered := g.filterGroupAndFieldVisibilityForAttendee(group, attendee)
			if filtered != nil {
				result = append(result, filtered)
			}
		}

		return result, nil
	} else {
		return make([]*modelsv1.Group, 0), errNotAttending(ctx) // shouldn't ever happen, just in case
	}
}

// findGroupsFullAccess searches for groups without permission checks.
//
// It returns all matching groups unfiltered, mapped to the API model with all fields visible.
func (g *groupService) findGroupsFullAccess(ctx context.Context, minSize uint, maxSize int, memberIDs []int64) ([]*modelsv1.Group, error) {
	result := make([]*modelsv1.Group, 0)

	groupIDs, err := g.DB.FindGroups(ctx, minSize, maxSize, memberIDs)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return result, nil
		}

		aulogging.ErrorErrf(ctx, err, "find groups failed: %s", err.Error())
		return result, errInternal(ctx, "database error while finding groups - see logs for details")
	}

	for _, id := range groupIDs {
		group, err := g.getGroupByIDFullAccess(ctx, id)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				aulogging.WarnErrf(ctx, err, "find groups failed to read group %s - maybe intermittent change: %s", id, err.Error())
				return make([]*modelsv1.Group, 0), errInternal(ctx, "database error while finding groups - see logs for details")
			}
		}

		result = append(result, group)
	}

	return result, nil
}

// GetGroupByID retrieves a group and its members and invites from the database by a given ID.
func (g *groupService) GetGroupByID(ctx context.Context, groupID string) (*modelsv1.Group, error) {
	validator, err := rbac.NewValidator(ctx)
	if err != nil {
		aulogging.ErrorErrf(ctx, err, "Could not retrieve RBAC validator from context. [error]: %v", err)
		return nil, errCouldNotGetValidator(ctx)
	}

	if validator.IsAdmin() {
		// admins are allowed access
		return g.getGroupByIDFullAccess(ctx, groupID)
	} else if validator.IsUser() {
		// ensure attending registration
		attendee, err := g.loggedInUserValidRegistration(ctx)
		if err != nil {
			return nil, err
		}

		group, err := g.getGroupByIDFullAccess(ctx, groupID)
		if err != nil {
			return nil, err
		}

		filtered := g.filterGroupAndFieldVisibilityForAttendee(group, attendee)
		if filtered == nil {
			return nil, errNoAccess(ctx)
		}

		return filtered, nil
	} else {
		return nil, errNotAttending(ctx) // shouldn't ever happen, just in case
	}
}

// getGroupByIDFullAccess retrieves a group and its members and invites from the database by a given ID.
//
// It does so without permission checks and returns all group fields unfiltered.
func (g *groupService) getGroupByIDFullAccess(ctx context.Context, groupID string) (*modelsv1.Group, error) {
	grp, err := g.DB.GetGroupByID(ctx, groupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errGroupIDNotFound(ctx)
		}

		return nil, errGroupRead(ctx, err.Error())
	}

	groupMembers, err := g.DB.GetGroupMembersByGroupID(ctx, groupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// this is an error because a group should always contain its owner, or should have been deleted
			return nil, errGroupHasNoMembers(ctx)
		}

		return nil, errGroupRead(ctx, err.Error())
	}

	return &modelsv1.Group{
		ID:          grp.ID,
		Name:        grp.Name,
		Flags:       aggregateFlags(grp.Flags),
		Comments:    common.ToOmitEmpty(grp.Comments),
		MaximumSize: grp.MaximumSize,
		Owner:       grp.Owner,
		Members:     toMembers(groupMembers),
		Invites:     toInvites(groupMembers),
	}, nil
}

// CreateGroup creates a new group in the database.
// Additionally, the group will add the owner as the initial group member.
//
// Admins can specify a specific group owner.
func (g *groupService) CreateGroup(ctx context.Context, group *modelsv1.GroupCreate) (string, error) {
	validator, err := rbac.NewValidator(ctx)
	if err != nil {
		aulogging.ErrorErrf(ctx, err, "Could not retrieve RBAC validator from context. [error]: %v", err)
		return "", errCouldNotGetValidator(ctx)
	}

	var ownerID int64
	var nickname string
	if validator.IsAdmin() {
		ownerID = group.Owner
		if ownerID > 0 {
			attendee, err := g.AttSrv.GetAttendee(ctx, int64(ownerID))
			if err != nil {
				return "", err
			}
			nickname = attendee.Nickname
		}
	}
	if ownerID == 0 {
		attendee, err := g.loggedInUserValidRegistration(ctx)
		if err != nil {
			return "", err
		}
		ownerID = attendee.ID
		nickname = attendee.Nickname
	}

	if group.MaximumSize == 0 {
		group.MaximumSize = maxGroupSize()
	}

	validation := validateGroupCreate(group)
	if len(validation) > 0 {
		return "", common.NewBadRequest(ctx, common.GroupDataInvalid, validation)
	}

	// Create a new group in the database
	groupID, err := g.DB.AddGroup(ctx, &entity.Group{
		Name:        group.Name,
		Flags:       fmt.Sprintf(",%s,", strings.Join(group.Flags, ",")),
		Comments:    common.Deref(group.Comments),
		MaximumSize: maxGroupSize(),
		Owner:       ownerID,
	})

	if err != nil {
		return "", err
	}

	gm := g.DB.NewEmptyGroupMembership(ctx, groupID, ownerID, nickname)
	gm.IsInvite = false
	return groupID, g.DB.AddGroupMembership(ctx, gm)
}

func validateGroupCreate(group *modelsv1.GroupCreate) url.Values {
	return validate(group.Name, group.Flags, group.MaximumSize)
}

func validateGroup(group *modelsv1.Group) url.Values {
	return validate(group.Name, group.Flags, group.MaximumSize)
}

// validate checks group fields for validity using the service configuration.
//
// Note that dynamic checks that require loading attendees or group members are implemented in the
// calling functions because those are authorization dependent (affected fields: owner).
func validate(name string, flags []string, maximumSize int64) url.Values {
	result := url.Values{}
	if len(name) == 0 {
		result.Set("name", "group name cannot be empty")
	}
	if len(name) > 50 {
		result.Set("name", "group name too long, max 50 characters")
	}
	allowed := allowedFlags()
	for _, flag := range flags {
		if !util.SliceContains(flag, allowed) {
			result.Set("flags", fmt.Sprintf("no such flag '%s'", url.PathEscape(flag)))
		}
	}

	if maximumSize <= 0 || maximumSize > maxGroupSize() {
		result.Set("maximum_size", "group maximum_size out of bounds")
	}

	return result
}

// UpdateGroup updates an existing group by uuid. Note that you cannot use this to change the group members!
//
// Admins or the current group owner can change the group owner to any member of the group.
func (g *groupService) UpdateGroup(ctx context.Context, group *modelsv1.Group) error {
	validator, err := rbac.NewValidator(ctx)
	if err != nil {
		aulogging.ErrorErrf(ctx, err, "Could not retrieve RBAC validator from context. [error]: %v", err)
		return errInternal(ctx, "unexpected error when parsing user claims")
	}

	dbGroup, err := g.DB.GetGroupByID(ctx, group.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errGroupIDNotFound(ctx)
		} else {
			return errGroupRead(ctx, err.Error())
		}
	}

	if validator.IsAdmin() || validator.IsAPITokenCall() {
		// admins and api token are allowed to make changes to any group

		// owner validation for admins
		if group.Owner > 0 {
			if _, err := g.AttSrv.GetAttendee(ctx, group.Owner); err != nil {
				return err
			}

			if err := g.checkAttending(ctx, group.Owner); err != nil {
				return err
			}
		} else {
			// keep owner unchanged to avoid accidental assignments
			group.Owner = dbGroup.Owner
		}
	} else if validator.IsUser() {
		attendee, err := g.loggedInUserValidRegistration(ctx)
		if err != nil {
			return err
		}

		if dbGroup.Owner != attendee.ID {
			return common.NewForbidden(ctx, common.AuthForbidden, common.Details("only the group owner or an admin can change a group"))
		}
	} else {
		return errNotAttending(ctx) // shouldn't ever happen, just in case
	}

	if group.MaximumSize == 0 {
		group.MaximumSize = maxGroupSize()
	}

	validation := validateGroup(group)
	if len(validation) > 0 {
		return common.NewBadRequest(ctx, common.GroupDataInvalid, validation)
	}

	// TODO check that new group size not too small (counting invitations and members)

	if dbGroup.Owner != group.Owner {
		err := g.canChangeGroupOwner(ctx, group)
		if err != nil {
			return err
		}
	}

	// do not touch fields that we do not wish to change, like createdAt or referenced members
	dbGroup.Name = group.Name
	dbGroup.Flags = fmt.Sprintf(",%s,", strings.Join(group.Flags, ","))
	dbGroup.Comments = common.Deref(group.Comments)
	dbGroup.MaximumSize = group.MaximumSize
	dbGroup.Owner = group.Owner

	return g.DB.UpdateGroup(ctx, dbGroup)
}

func (g *groupService) canChangeGroupOwner(ctx context.Context, group *modelsv1.Group) error {
	members, err := g.DB.GetGroupMembersByGroupID(ctx, group.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errGroupHasNoMembers(ctx)
		}
		aulogging.ErrorErrf(ctx, err, "unexpected error %v", err)
		return errInternal(ctx, "unexpected error occurrec")
	}
	found := false
	for _, member := range members {
		if member.ID == group.Owner {
			found = true
			break
		}
	}
	if !found {
		return errNewOwnerNotMember(ctx)
	}
	return nil
}

// DeleteGroup removes all members from the group and sets a deletion timestamp.
func (g *groupService) DeleteGroup(ctx context.Context, groupID string) error {
	validator, err := rbac.NewValidator(ctx)
	if err != nil {
		aulogging.ErrorErrf(ctx, err, "Could not retrieve RBAC validator from context. [error]: %v", err)
		return errInternal(ctx, "unexpected error when parsing user claims")
	}

	group, err := g.DB.GetGroupByID(ctx, groupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errGroupIDNotFound(ctx)
		}

		aulogging.Warnf(ctx, "failed to read group %s from db: %s", url.PathEscape(groupID), err.Error())
		return errGroupRead(ctx, "error retrieving group - see logs for details")
	}

	if validator.IsAdmin() || validator.IsAPITokenCall() {
		// admins and api token are allowed to make changes to any group
	} else if validator.IsUser() {
		attendee, err := g.loggedInUserValidRegistration(ctx)
		if err != nil {
			return err
		}

		if int64(group.Owner) != attendee.ID {
			return common.NewForbidden(ctx, common.AuthForbidden, common.Details("only the group owner or an admin can delete a group"))
		}
	} else {
		return errNotAttending(ctx) // shouldn't ever happen, just in case
	}

	members, err := g.DB.GetGroupMembersByGroupID(ctx, groupID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return errInternal(ctx, "failed to read group members during delete")
		}
		// empty group is ok
	}

	// first we have to remove all members, which have been part of the group and then
	for _, member := range members {
		if err := g.DB.DeleteGroupMembership(ctx, member.ID); err != nil {
			aulogging.ErrorErrf(ctx, err, "error occurred when trying to remove member with ID %d from group %s. [error]: %s", member.ID, groupID, err.Error())
			return errInternal(ctx,
				fmt.Sprintf("could not remove member %d from group %s", member.ID, groupID))
		}
	}

	if err := g.DB.DeleteGroupByID(ctx, groupID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errGroupIDNotFound(ctx)
		}

		aulogging.ErrorErrf(ctx, err, "unexpected error. [error]: %s", err.Error())
		return errInternal(ctx, "unexpected error occurred during deletion of group")
	}

	return nil
}

func toMembers(groupMembers []*entity.GroupMember) []modelsv1.Member {
	return toMembersFilteredSorted(groupMembers, false)
}

func toInvites(groupMembers []*entity.GroupMember) []modelsv1.Member {
	return toMembersFilteredSorted(groupMembers, true)
}

func toMembersFilteredSorted(groupMembers []*entity.GroupMember, invites bool) []modelsv1.Member {
	members := make([]modelsv1.Member, 0)
	for _, m := range groupMembers {
		if m == nil {
			continue
		}
		if m.IsInvite != invites {
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

func errNoGroup(ctx context.Context) error {
	return common.NewNotFound(ctx, common.GroupMemberNotFound, common.Details("not in a group"))
}

func errGroupIDNotFound(ctx context.Context) error {
	return common.NewNotFound(ctx, common.GroupIDNotFound, common.Details("group does not exist"))
}

func errGroupHasNoMembers(ctx context.Context) error {
	return common.NewInternalServerError(ctx, common.GroupMemberNotFound, common.Details("unable to find members in group"))
}

func errNewOwnerNotMember(ctx context.Context) error {
	return common.NewInternalServerError(ctx, common.GroupMemberNotFound, common.Details("new owner must be a member of the group"))
}

func errCouldNotGetValidator(ctx context.Context) error {
	return common.NewInternalServerError(ctx, common.InternalErrorMessage, common.Details("unexpected error when parsing user claims"))
}

func errNotAttending(ctx context.Context) error {
	return common.NewForbidden(ctx, common.NotAttending, common.Details("access denied - you must have a valid registration in status approved, (partially) paid, checked in"))
}

func errNoAccess(ctx context.Context) error {
	return common.NewForbidden(ctx, common.AuthForbidden, common.Details("access denied - you do not have access to this group"))
}

func errGroupRead(ctx context.Context, details string) error {
	return common.NewInternalServerError(ctx, common.GroupReadError, common.Details(details))
}

func errGroupWrite(ctx context.Context, details string) error {
	return common.NewInternalServerError(ctx, common.GroupWriteError, common.Details(details))
}

func errInternal(ctx context.Context, details string) error {
	return common.NewInternalServerError(ctx, common.InternalErrorMessage, common.Details(details))
}
