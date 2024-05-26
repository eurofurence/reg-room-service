package groupservice

import (
	"context"
	"errors"
	"fmt"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"slices"
	"strings"

	"gorm.io/gorm"

	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
	"github.com/eurofurence/reg-room-service/internal/entity"
	apierrors "github.com/eurofurence/reg-room-service/internal/errors"
	"github.com/eurofurence/reg-room-service/internal/repository/database"
	"github.com/eurofurence/reg-room-service/internal/service/rbac"
	"github.com/eurofurence/reg-room-service/internal/util/ptr"
	"github.com/eurofurence/reg-room-service/internal/web/common"
	"github.com/eurofurence/reg-room-service/internal/web/v1/util"
)

var (
	errGroupIDNotFound      = apierrors.NewNotFound(common.GroupIDNotFoundMessage, "unable to find group in database")
	errGroupHasNoMembers    = apierrors.NewInternalServerError(common.GroupMemberNotFound, "unable to find members in group")
	errCouldNotGetValidator = apierrors.NewInternalServerError(common.InternalErrorMessage, "unexpected error when parsing user claims")
)

// Service defines the interface for the service function implementations for the group endpoints.
//
// TODO ListGroups
// TODO GetMyGroup
// TODO Remove member from group
type Service interface {
	GetGroupByID(ctx context.Context, groupID string) (*modelsv1.Group, error)
	CreateGroup(ctx context.Context, group modelsv1.GroupCreate) (string, error)
	UpdateGroup(ctx context.Context, group modelsv1.Group) error
	DeleteGroup(ctx context.Context, groupID string) error
	AddMemberToGroup(ctx context.Context, req AddGroupMemberParams) error
}

func NewService(repository database.Repository) Service {
	return &groupService{DB: repository}
}

type groupService struct {
	DB database.Repository
}

// GetGroupByID attempts to retrieve a group and its members from the database by a given ID.
func (g *groupService) GetGroupByID(ctx context.Context, groupID string) (*modelsv1.Group, error) {
	grp, err := g.DB.GetGroupByID(ctx, groupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errGroupIDNotFound
		}

		return nil, apierrors.NewInternalServerError(common.InternalErrorMessage, err.Error())
	}

	groupMembers, err := g.DB.GetGroupMembersByGroupID(ctx, groupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errGroupHasNoMembers
		}
	}

	return &modelsv1.Group{
		ID:          grp.ID,
		Name:        grp.Name,
		Flags:       aggregateFlags(grp.Flags),
		Comments:    &grp.Comments,
		MaximumSize: ptr.To(int32(grp.MaximumSize)),
		Owner:       int32(grp.Owner),
		Members:     toMembers(groupMembers),
		Invites:     nil,
	}, nil
}

// CreateGroup creates a new group in the database.
// Additionally, the group will add the owner as the initial group member.
//
// Admins can specify a specific group owner.
func (g *groupService) CreateGroup(ctx context.Context, group modelsv1.GroupCreate) (string, error) {
	validator, err := rbac.NewValidator(ctx)
	if err != nil {
		aulogging.ErrorErrf(ctx, err, "Could not retrieve RBAC validator from context. [error]: %v", err)
		return "", errCouldNotGetValidator
	}

	ownerID, err := util.ParseUInt[uint](validator.Subject())
	if err != nil {
		aulogging.WarnErrf(ctx, err, "subject has an unexpected value %q", validator.Subject())
		return "", apierrors.NewInternalServerError(common.InternalErrorMessage, "subject should have a valid numerical value")
	}
	if validator.IsAdmin() {
		ownerID = uint(group.Owner)
	}

	// Create a new group in the database
	groupID, err := g.DB.AddGroup(ctx, &entity.Group{
		Name:        group.Name,
		Flags:       fmt.Sprintf(",%s,", strings.Join(group.Flags, ",")),
		Comments:    ptr.Deref(group.Comments),
		MaximumSize: 6, // TODO add from config
		Owner:       ownerID,
	})

	if err != nil {
		return "", err
	}

	gm := g.DB.NewEmptyGroupMembership(ctx, groupID, ownerID)
	return groupID, g.DB.AddGroupMembership(ctx, gm)
}

// AddGroupMemberParams is the request type for the AddMemberToGroup operation.
type AddGroupMemberParams struct {
	// GroupID is the ID of the group where a user should be added
	GroupID string
	// BadgeNumber is the registration number of a user
	BadgeNumber uint
	// Nickname is the nickname of a registered user that should receive
	// an invitation Email.
	Nickname string
	// Code is the invite code that can be used to join a group.
	Code string
	// Force is an admin only flag that allows to bypass the
	// validations.
	Force bool
}

// AddMemberToGroup TODO...
func (g *groupService) AddMemberToGroup(ctx context.Context, req AddGroupMemberParams) error {
	gm := g.DB.NewEmptyGroupMembership(ctx, req.GroupID, req.BadgeNumber)

	err := g.DB.AddGroupMembership(ctx, gm)
	if err != nil {
		return apierrors.NewInternalServerError(common.InternalErrorMessage, err.Error())
	}

	return nil
}

// UpdateGroup updates an existing group by uuid. Note that you cannot use this to change the group members!
//
// Admins or the current group owner can change the group owner to any member of the group.
func (g *groupService) UpdateGroup(ctx context.Context, group modelsv1.Group) error {
	validator, err := rbac.NewValidator(ctx)
	if err != nil {
		aulogging.ErrorErrf(ctx, err, "Could not retrieve RBAC validator from context. [error]: %v", err)
		return apierrors.NewInternalServerError(common.InternalErrorMessage, "unexpected error when parsing user claims")
	}

	badgeNumber, err := util.ParseUInt[uint](validator.Subject())
	if err != nil {
		aulogging.WarnErrf(ctx, err, "subject has an unexpected value %q", validator.Subject())
		return apierrors.NewInternalServerError(common.InternalErrorMessage, "subject should have a valid numerical value")
	}

	getGroup, err := g.DB.GetGroupByID(ctx, group.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errGroupIDNotFound
		}
	}

	updateGroup := &entity.Group{
		Base:        entity.Base{ID: group.ID},
		Name:        group.Name,
		Flags:       fmt.Sprintf(",%s,", strings.Join(group.Flags, ",")),
		Comments:    ptr.Deref(group.Comments),
		MaximumSize: uint(ptr.Deref(group.MaximumSize)),
	}

	// Changes to the group owner can only be instigated by either the group owner
	// or forcefully by the admin.
	// In both cases a new owner can only be an already existing member in the group.
	switch {
	case validator.IsAdmin():
		fallthrough
	case getGroup.Owner == badgeNumber && group.Owner != int32(getGroup.Owner):
		if getGroup.Owner == uint(group.Owner) {
			// we are not changing the owner here
			break
		}

		err := g.changeGroupOwner(ctx, group, updateGroup)
		if err != nil {
			return err
		}
	}

	return g.DB.UpdateGroup(ctx, updateGroup)
}

func (g *groupService) changeGroupOwner(ctx context.Context, group modelsv1.Group, updateGroup *entity.Group) error {
	members, err := g.DB.GetGroupMembersByGroupID(ctx, group.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errGroupHasNoMembers
		}
		aulogging.ErrorErrf(ctx, err, "unexpected error %v", err)
		return apierrors.NewInternalServerError(common.InternalErrorMessage, "unexpected error occurrec")
	}
	found := false
	for _, member := range members {
		if member.ID == uint(group.Owner) {
			found = true
			break
		}
	}
	if !found {
		return errGroupHasNoMembers
	}
	updateGroup.Owner = uint(group.Owner)
	return nil
}

// DeleteGroup removes all members from the group and sets a deletion timestamp.
func (g *groupService) DeleteGroup(ctx context.Context, groupID string) error {
	validator, err := rbac.NewValidator(ctx)
	if err != nil {
		aulogging.ErrorErrf(ctx, err, "Could not retrieve RBAC validator from context. [error]: %v", err)
		return apierrors.NewInternalServerError(common.InternalErrorMessage, "unexpected error when parsing user claims")
	}

	group, err := g.DB.GetGroupByID(ctx, groupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apierrors.NewNotFound(common.GroupIDNotFoundMessage,
				fmt.Sprintf("couldn't find group for ID: %s", groupID))
		}

		return apierrors.NewInternalServerError(common.InternalErrorMessage,
			fmt.Sprintf("error when retrieving group with ID: %s", groupID))
	}

	if group.DeletedAt.Valid {
		// group is already deleted
		aulogging.Warnf(ctx, "group %s was already marked for deletion", groupID)
		return nil
	}

	badgeNumber, err := util.ParseUInt[uint](validator.Subject())
	if err != nil {
		aulogging.ErrorErrf(ctx, err, "subject has an unexpected value %q", validator.Subject())
		return apierrors.NewInternalServerError(common.InternalErrorMessage, "subject should have a valid numerical value")
	}

	if !validator.IsAdmin() || badgeNumber == group.Owner {
		return apierrors.NewForbidden(common.AuthForbiddenMessage, "only the group owner or an admin can delete a group")
	}

	members, err := g.DB.GetGroupMembersByGroupID(ctx, groupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apierrors.NewInternalServerError(common.GroupMemberNotFound, "at least one group member should be in this group")
		}
	}

	// first we have to remove all members, which have been part of the group and then
	for _, member := range members {
		if err := g.DB.DeleteGroupMembership(ctx, member.ID); err != nil {
			aulogging.ErrorErrf(ctx, err, "error occurred when trying to remove member with ID %d from group %s. [error]: %s", member.ID, groupID, err.Error())
			return apierrors.NewInternalServerError(
				common.InternalErrorMessage,
				fmt.Sprintf("could not remove member %d from group %s", member.ID, groupID))
		}
	}

	if err := g.DB.SoftDeleteGroupByID(ctx, groupID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apierrors.NewNotFound(common.GroupIDNotFoundMessage, fmt.Sprintf("couldn't find group with ID %s", groupID))
		}

		aulogging.ErrorErrf(ctx, err, "unexpected error. [error]: %s", err.Error())
		return apierrors.NewInternalServerError(
			common.InternalErrorMessage, "unexpected error occurred during deletion of group")
	}

	return nil
}

func toMembers(groupMembers []*entity.GroupMember) []modelsv1.Member {
	members := make([]modelsv1.Member, 0)
	for _, m := range groupMembers {
		if m == nil {
			continue
		}

		members = append(members, modelsv1.Member{
			ID:       int32(m.ID),
			Nickname: m.Nickname,
			Avatar:   &m.AvatarURL,
		})
	}

	return members
}

func aggregateFlags(input string) []string {
	if input == "" {
		return nil
	}

	tags := strings.Split(input, ",")
	tags = slices.DeleteFunc(tags, func(s string) bool {
		return s == ""
	})

	return tags
}
