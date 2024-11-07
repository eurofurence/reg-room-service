package groupservice

import (
	"context"
	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
	"github.com/eurofurence/reg-room-service/internal/repository/database"
	"github.com/eurofurence/reg-room-service/internal/repository/downstreams/attendeeservice"
	"github.com/eurofurence/reg-room-service/internal/repository/downstreams/mailservice"
)

// Service defines the interface for the service function implementations for the group endpoints.
type Service interface {
	GetGroupByID(ctx context.Context, groupID string) (*modelsv1.Group, error)
	CreateGroup(ctx context.Context, group *modelsv1.GroupCreate) (string, error)
	UpdateGroup(ctx context.Context, group *modelsv1.Group) error
	DeleteGroup(ctx context.Context, groupID string) error
	// AddMemberToGroup adds the member to the group.
	//
	// Returns a possibly empty url extension to be appended to the Location for accepting an invitation,
	// if applicable. Unless empty, the url extension is something like "?code=<join code>".
	//
	// The same link will also be included in the email sent to the invited attendee.
	AddMemberToGroup(ctx context.Context, req *AddGroupMemberParams) (string, error)
	RemoveMemberFromGroup(ctx context.Context, req *RemoveGroupMemberParams) error
	FindGroups(ctx context.Context, minSize uint, maxSize int, memberIDs []int64, public bool) ([]*modelsv1.Group, error)
	FindMyGroup(ctx context.Context) (*modelsv1.Group, error)
}

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

func New(db database.Repository, attsrv attendeeservice.AttendeeService, mailsrv mailservice.MailService) Service {
	return &groupService{
		DB:      db,
		AttSrv:  attsrv,
		MailSrv: mailsrv,
	}
}

type groupService struct {
	DB      database.Repository
	AttSrv  attendeeservice.AttendeeService
	MailSrv mailservice.MailService
}
