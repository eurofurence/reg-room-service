package groupsctl

import (
	"context"
	"github.com/go-chi/chi/v5"
	"net/http"

	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
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
	if err := validateGroupID(r.Context(), groupID); err != nil {
		return nil, err
	}

	return &DeleteGroupRequest{groupID: groupID}, nil
}

// DeleteGroupResponse writes out a `No Content` status, if the operation was successful.
func (h *Controller) DeleteGroupResponse(_ context.Context, _ *modelsv1.Empty, w http.ResponseWriter) error {
	w.WriteHeader(http.StatusNoContent)
	return nil
}
