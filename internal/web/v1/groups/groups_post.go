package groups

import (
	"context"
	"fmt"
	"github.com/eurofurence/reg-room-service/internal/logging"
	"net/http"
	"net/url"

	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
	apierrors "github.com/eurofurence/reg-room-service/internal/errors"
	"github.com/eurofurence/reg-room-service/internal/web/common"
	"github.com/eurofurence/reg-room-service/internal/web/v1/util"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// CreateGroupRequest is the request type for the AddMemberToGroup operation.
type CreateGroupRequest struct {
	// Group is the expected representation for the request body
	Group modelsv1.GroupCreate
}

// CreateGroup creates a new group, setting yourself as the owner.
//
// This also adds you as the first member of the group.
//
// You must have a valid registration.
//
// Note that the members and invites fields are ignored. The group is always created with just you as the owner
// (or for admins, if a different owner is specified via badge number, that owner).
func (h *Controller) CreateGroup(ctx context.Context, req *CreateGroupRequest, w http.ResponseWriter) (*modelsv1.Empty, error) {
	logger := logging.LoggerFromContext(ctx)

	newGroupUUID, err := h.ctrl.CreateGroup(ctx, req.Group)
	if err != nil || newGroupUUID == "" {
		return nil, err
	}

	requestURL, ok := ctx.Value(common.CtxKeyRequestURL{}).(*url.URL)
	if !ok {
		logger.Error("could not retrieve base URL from context")
		return nil, err
	}

	w.Header().Set("Location", fmt.Sprintf("%s/%s", requestURL.Path, newGroupUUID))
	return nil, nil
}

func (h *Controller) CreateGroupRequest(r *http.Request, w http.ResponseWriter) (*CreateGroupRequest, error) {
	var group modelsv1.GroupCreate

	if err := util.NewStrictJSONDecoder(r.Body).Decode(&group); err != nil {
		common.SendHTTPStatusErrorResponse(r.Context(), w, apierrors.NewBadRequest(
			"group.data.invalid", "please check if your provided JSON is valid",
		))
		return nil, err
	}

	cgr := &CreateGroupRequest{
		Group: group,
	}

	return cgr, nil
}

func (h *Controller) CreateGroupResponse(ctx context.Context, _ *modelsv1.Empty, w http.ResponseWriter) error {
	w.WriteHeader(http.StatusCreated)
	return nil
}

// AddMemberToGroupRequest is the request type for the AddMemberToGroup operation.
type AddMemberToGroupRequest struct {
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

// AddMemberToGroup adds an attendee to a group.
// Group owners may use this to send an invite email. The invite email will contain a link with a code which
// then allows the invited person to add themselves.
//
// Admins can add the force query parameter to just add. If they do not specify force=true, they are subject
// to the same limitations as every normal user.
//
// Users may only add themselves, and only if they have a valid invite code, and if they are registered for the convention.
//
// If an attendee is already in a group, or has already been individually assigned to a room, then they
// cannot be added to a group any more.
//
// If a group has already been assigned to a room, then only admins can change their members.
func (h *Controller) AddMemberToGroup(ctx context.Context, req *AddMemberToGroupRequest, w http.ResponseWriter) (*modelsv1.Empty, error) {
	return nil, nil
}

// AddMemberToGroupRequest validates and creates the request for the AddMemberToGroup operation.
func (h *Controller) AddMemberToGroupRequest(r *http.Request, w http.ResponseWriter) (*AddMemberToGroupRequest, error) {
	groupID := chi.URLParam(r, "uuid")
	if _, err := uuid.Parse(groupID); err != nil {
		common.SendHTTPStatusErrorResponse(r.Context(), w, apierrors.NewBadRequest("group.id.invalid", ""))
		return nil, err
	}

	badge := chi.URLParam(r, "badgenumber")
	badgeNumber, err := util.ParseUInt[uint](badge)
	if err != nil {
		common.SendHTTPStatusErrorResponse(r.Context(), w, apierrors.NewBadRequest("group.data.invalid", "invalid type for badge number"))
		return nil, err
	}

	query := r.URL.Query()
	req := &AddMemberToGroupRequest{
		GroupID:     groupID,
		BadgeNumber: badgeNumber,
		Nickname:    query.Get("nickname"),
		Code:        query.Get("code"),
	}

	force, err := util.ParseOptionalBool(query.Get("force"))
	if err != nil {
		common.SendHTTPStatusErrorResponse(r.Context(), w, apierrors.NewBadRequest("group.data.invalid", ""))
		return nil, err
	}

	req.Force = force
	return req, nil
}

// AddMemberToGroupResponse writes out the response for the AddMemberToGroup operation.
func (h *Controller) AddMemberToGroupResponse(ctx context.Context, _ *modelsv1.Empty, w http.ResponseWriter) error {
	return nil
}
