package groupservice

import (
	"context"
	"errors"
	"fmt"
	"github.com/eurofurence/reg-room-service/internal/logging"
	"slices"
	"strings"

	"github.com/eurofurence/reg-room-service/internal/web/common"

	"gorm.io/gorm"

	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
	"github.com/eurofurence/reg-room-service/internal/entity"
	apierrors "github.com/eurofurence/reg-room-service/internal/errors"
	"github.com/eurofurence/reg-room-service/internal/repository/database"
	"github.com/eurofurence/reg-room-service/internal/util/ptr"
)

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
			return nil, apierrors.NewNotFound(common.GroupIDNotFoundMessage, fmt.Sprintf("no record found for id %q", groupID))
		}

		return nil, apierrors.NewInternalServerError(common.InternalErrorMessage, err.Error())
	}

	groupMembers, err := g.DB.GetGroupMembersByGroupID(ctx, groupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apierrors.NewInternalServerError(common.GroupMemberNotFound, fmt.Sprintf("unable to find members for group %s", groupID))
		}
	}

	return &modelsv1.Group{
		ID:          grp.ID,
		Name:        grp.Name,
		Flags:       aggregateFlags(grp.Flags),
		Comments:    &grp.Comments,
		MaximumSize: ptr.To(int32(grp.MaximumSize)),
		Owner:       int32(grp.Owner),
		Members:     ToMembers(groupMembers),
		Invites:     nil,
	}, nil
}

// CreateGroup creates a new group in the database.
// Additionally, the group will add the owner as the initial group member.
//
// Admins can specify a specific group owner.
func (g *groupService) CreateGroup(ctx context.Context, group modelsv1.GroupCreate) (string, error) {
	// TODO check the token if the function was invoked by an admin
	// 	if an admin authored the creation and provided a custom owner ID, the group owner has to be set to the
	// 	provided owner ID

	ownerID := uint(42)
	isAdmin := false
	if isAdmin {
		ownerID = uint(group.Owner)
	}

	// Create a new group in the database
	groupID, err := g.DB.AddGroup(ctx, &entity.Group{
		Name:        group.Name,
		Flags:       fmt.Sprintf(",%s,", strings.Join(group.Flags, ",")),
		Comments:    ptr.Deref(group.Comments),
		MaximumSize: 6,       // TODO add from config
		Owner:       ownerID, // TODO read from attendee service (or passed in by admin)
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

func (g *groupService) AddMemberToGroup(ctx context.Context, req AddGroupMemberParams) error {
	gm := g.DB.NewEmptyGroupMembership(ctx, req.GroupID, req.BadgeNumber)

	err := g.DB.AddGroupMembership(ctx, gm)
	if err != nil {
		return apierrors.NewInternalServerError(common.InternalErrorMessage, err.Error())
	}

	return nil
}

func (g *groupService) UpdateGroup(ctx context.Context, group modelsv1.Group) error {
	// TODO retrieve badge number from context

	updateGroup := &entity.Group{
		Base:        entity.Base{ID: group.ID},
		Name:        group.Name,
		Flags:       fmt.Sprintf(",%s,", strings.Join(group.Flags, ",")),
		Comments:    ptr.Deref(group.Comments),
		MaximumSize: uint(ptr.Deref(group.MaximumSize)),
		Owner:       uint(group.Owner),
	}

	return g.DB.UpdateGroup(ctx, updateGroup)
}

// DeleteGroup removes all members from the group and sets a deletion timestamp.
func (g *groupService) DeleteGroup(ctx context.Context, groupID string) error {
	logger := logging.LoggerFromContext(ctx)
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
		logger.Warn("group %s was already marked for deletion", groupID)
		return nil
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
			logger.Error("error occurred when trying to remove member with ID %d from group %s. [error]: %s", member.ID, groupID, err.Error())
			return apierrors.NewInternalServerError(
				common.InternalErrorMessage,
				fmt.Sprintf("could not remove member %d from group %s", member.ID, groupID))
		}
	}

	if err := g.DB.SoftDeleteGroupByID(ctx, groupID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apierrors.NewNotFound(common.GroupIDNotFoundMessage, fmt.Sprintf("couldn't find group with ID %s", groupID))
		}

		logger.Error("unexpected error. [error]: %s", err.Error())
		return apierrors.NewInternalServerError(
			common.InternalErrorMessage, "unexpected error occurred during deletion of group")
	}

	return nil
}

func ToMembers(groupMembers []*entity.GroupMember) []modelsv1.Member {
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
