package groupsctl

import (
	"context"
	"github.com/eurofurence/reg-room-service/internal/controller/v1/util"
	"github.com/go-chi/chi/v5"
	"net/http"

	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
	"github.com/eurofurence/reg-room-service/internal/application/common"
)

// RemoveGroupMemberRequest holds information, which is required to call the RemoveGroupMember operation.
type RemoveGroupMemberRequest struct {
	groupID     string
	badgeNumber uint
}

// RemoveGroupMember removes a group member or revokes an invitation.
//
// Details see OpenAPI spec.
func (h *Controller) RemoveGroupMember(ctx context.Context, req *RemoveGroupMemberRequest, w http.ResponseWriter) (*modelsv1.Empty, error) {
	// TODO

	return nil, nil
}

// RemoveGroupMemberRequest validates and creates the request for the RemoveGroupMember operation.
func (h *Controller) RemoveGroupMemberRequest(r *http.Request, w http.ResponseWriter) (*RemoveGroupMemberRequest, error) {
	const uuidParam, badeNumberParam = "uuid", "badgenumber"

	ctx := r.Context()

	groupID := chi.URLParam(r, uuidParam)
	if err := validateGroupID(ctx, groupID); err != nil {
		return nil, err
	}

	badgeNumber, err := util.ParseUInt[uint](chi.URLParam(r, badeNumberParam))
	if err != nil {
		return nil, common.NewBadRequest(ctx, common.GroupDataInvalid, common.Details("invalid type for badge number"))
	}

	return &RemoveGroupMemberRequest{groupID, badgeNumber}, nil
}

// RemoveGroupMemberResponse writes out a `No Content` status.
func (h *Controller) RemoveGroupMemberResponse(_ context.Context, _ *modelsv1.Empty, w http.ResponseWriter) error {
	w.WriteHeader(http.StatusNoContent)
	return nil
}
