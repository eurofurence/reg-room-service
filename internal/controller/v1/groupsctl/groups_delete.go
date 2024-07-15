package groupsctl

import (
	"context"
	"github.com/eurofurence/reg-room-service/internal/application/web"
	"github.com/eurofurence/reg-room-service/internal/controller/v1/util"
	"github.com/go-chi/chi/v5"
	"net/http"

	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
	"github.com/eurofurence/reg-room-service/internal/application/common"
)

// DeleteGroupRequest holds information, which is required to call the DeleteGroup operation.
type DeleteGroupRequest struct {
	groupID string
}

// DeleteGroup disbands (and deletes) an existing group by uuid. Note that this will first kick everyone from the group!
//
// Only Admins or the current group owner can do this.
func (h *Controller) DeleteGroup(ctx context.Context, req *DeleteGroupRequest, w http.ResponseWriter) (*modelsv1.Empty, error) {
	err := h.svc.DeleteGroup(ctx, req.groupID)
	return nil, err
}

// DeleteGroupRequest parses and returns a request containing information to call the DeleteGroup function.
func (h *Controller) DeleteGroupRequest(r *http.Request, w http.ResponseWriter) (*DeleteGroupRequest, error) {
	groupID := chi.URLParam(r, "uuid")
	if err := validateGroupID(r.Context(), w, groupID); err != nil {
		return nil, err
	}

	return &DeleteGroupRequest{groupID: groupID}, nil
}

// DeleteGroupResponse writes out a `No Content` status, if the operation was successful.
func (h *Controller) DeleteGroupResponse(_ context.Context, _ *modelsv1.Empty, w http.ResponseWriter) error {
	w.WriteHeader(http.StatusNoContent)
	return nil
}

// RemoveGroupMemberRequest holds information, which is required to call the RemoveGroupMember operation.
type RemoveGroupMemberRequest struct {
	groupID     string
	badgeNumber uint
}

// RemoveGroupMember removes a group member or revokes an invitation
//
// Removes the attendee with the given badge number from the group (or its list of invitations).
// Possibly also add an entry to the group's auto-deny list.
//
// *Permissions*
//
// * Group owners can remove members/revoke invitations.
// * Members can remove themselves/decline invitations.
// * Admins can remove anyone/revoke their invitations.
//
// *Limitations*
//
// If a member is the current group owner, this fails with 409 conflict. First must reassign the group owner via
// an update to the group resource.
//
// *Auto-Deny*
//
// If the auto deny parameter is set to true, in addition to removing the group membership/invitation, the
// badge number is added to an auto-decline list. Further attempts to invite this attendee into the group
// are automatically declined.
func (h *Controller) RemoveGroupMember(ctx context.Context, req *RemoveGroupMemberRequest, w http.ResponseWriter) (*modelsv1.Empty, error) {
	return nil, nil
}

// RemoveGroupMemberRequest validates and creates the request for the RemoveGroupMember operation.
func (h *Controller) RemoveGroupMemberRequest(r *http.Request, w http.ResponseWriter) (*RemoveGroupMemberRequest, error) {
	const uuidParam, badeNumberParam = "uuid", "badgenumber"

	groupID := chi.URLParam(r, uuidParam)
	if err := validateGroupID(r.Context(), w, groupID); err != nil {
		return nil, err
	}

	badgeNumber, err := util.ParseUInt[uint](chi.URLParam(r, badeNumberParam))
	if err != nil {
		ctx := r.Context()
		web.SendErrorResponse(ctx, w,
			common.NewBadRequest(ctx, common.GroupDataInvalid, common.Details("invalid type for badge number")))
		return nil, err
	}

	return &RemoveGroupMemberRequest{groupID, badgeNumber}, nil
}

// RemoveGroupMemberResponse writes out a `No Content` status.
func (h *Controller) RemoveGroupMemberResponse(_ context.Context, _ *modelsv1.Empty, w http.ResponseWriter) error {
	w.WriteHeader(http.StatusNoContent)
	return nil
}
