package groups

import (
	"context"
	groupservice "github.com/eurofurence/reg-room-service/internal/service/groups"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
	apierrors "github.com/eurofurence/reg-room-service/internal/errors"
	"github.com/eurofurence/reg-room-service/internal/web/common"
	"github.com/eurofurence/reg-room-service/internal/web/v1/util"
)

// AddMemberToGroup adds an attendee to a group.
// Group owners may use this to send an invitation email. The invite email will contain a link with a code which
// then allows the invited person to add themselves.
//
// Admins can add the force query parameter to just add. If they do not specify force=true, they are subject
// to the same limitations as every normal user.
//
// Users may only add themselves, and only if they have a valid invite code, and if they are registered for the convention.
//
// If an attendee is already in a group, or has already been individually assigned to a room, then they
// cannot be added to a group anymore.
//
// If a group has already been assigned to a room, then only admins can change their members.
func (h *Controller) AddMemberToGroup(ctx context.Context, req *groupservice.AddGroupMemberParams, w http.ResponseWriter) (*modelsv1.Empty, error) {
	return &modelsv1.Empty{}, h.svc.AddMemberToGroup(ctx, *req)
}

// AddMemberToGroupRequest validates and creates the request for the AddMemberToGroup operation.
func (h *Controller) AddMemberToGroupRequest(r *http.Request, w http.ResponseWriter) (*groupservice.AddGroupMemberParams, error) {
	groupID := chi.URLParam(r, "uuid")
	if _, err := uuid.Parse(groupID); err != nil {
		common.SendHTTPStatusErrorResponse(r.Context(), w, apierrors.NewBadRequest(common.GroupIDInvalidMessage, ""))
		return nil, err
	}

	badge := chi.URLParam(r, "badgenumber")
	badgeNumber, err := util.ParseUInt[uint](badge)
	if err != nil {
		common.SendHTTPStatusErrorResponse(r.Context(), w, apierrors.NewBadRequest(common.GroupDataInvalidMessage, "invalid type for badge number"))
		return nil, err
	}

	query := r.URL.Query()
	req := &groupservice.AddGroupMemberParams{
		GroupID:     groupID,
		BadgeNumber: badgeNumber,
		Nickname:    query.Get("nickname"),
		Code:        query.Get("code"),
	}

	force, err := util.ParseOptionalBool(query.Get("force"))
	if err != nil {
		common.SendHTTPStatusErrorResponse(r.Context(), w, apierrors.NewBadRequest(common.GroupDataInvalidMessage, ""))
		return nil, err
	}

	req.Force = force
	return req, nil
}

// AddMemberToGroupResponse writes out the response for the AddMemberToGroup operation.
func (h *Controller) AddMemberToGroupResponse(_ context.Context, _ *modelsv1.Empty, w http.ResponseWriter) error {
	w.WriteHeader(http.StatusNoContent)
	return nil
}
