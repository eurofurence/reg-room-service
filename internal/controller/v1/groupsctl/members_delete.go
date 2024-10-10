package groupsctl

import (
	"context"
	"github.com/eurofurence/reg-room-service/internal/controller/v1/util"
	groupservice "github.com/eurofurence/reg-room-service/internal/service/groups"
	"github.com/go-chi/chi/v5"
	"net/http"

	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
	"github.com/eurofurence/reg-room-service/internal/application/common"
)

// RemoveGroupMember removes a group member or revokes an invitation.
//
// Details see OpenAPI spec.
func (h *Controller) RemoveGroupMember(ctx context.Context, req *groupservice.RemoveGroupMemberParams, w http.ResponseWriter) (*modelsv1.Empty, error) {
	return &modelsv1.Empty{}, h.svc.RemoveMemberFromGroup(ctx, req)
}

// RemoveGroupMemberRequest validates and creates the request for the RemoveGroupMember operation.
func (h *Controller) RemoveGroupMemberRequest(r *http.Request, w http.ResponseWriter) (*groupservice.RemoveGroupMemberParams, error) {
	ctx := r.Context()
	query := r.URL.Query()

	groupID := chi.URLParam(r, "uuid")
	if err := validateGroupID(ctx, groupID); err != nil {
		return nil, err
	}

	badge := chi.URLParam(r, "badgenumber")
	badgeNumber, err := util.ParseInt[int64](badge)
	if err != nil {
		return nil, common.NewBadRequest(ctx, common.RequestParseFailed, common.Details("invalid badge number - must be positive integer"), err)
	}
	if badgeNumber < 1 {
		return nil, common.NewBadRequest(ctx, common.GroupDataInvalid, common.Details("invalid badge number - must be positive integer"))
	}

	autodeny, err := util.ParseOptionalBool(query.Get("autodeny"))
	if err != nil {
		return nil, common.NewBadRequest(ctx, common.RequestParseFailed, common.Details("invalid autodeny parameter, try true, 1, false, 0 or omit"), err)
	}

	return &groupservice.RemoveGroupMemberParams{
		GroupID:     groupID,
		BadgeNumber: badgeNumber,
		AutoDeny:    autodeny,
	}, nil
}

// RemoveGroupMemberResponse writes out a `No Content` status.
func (h *Controller) RemoveGroupMemberResponse(_ context.Context, _ *modelsv1.Empty, w http.ResponseWriter) error {
	w.WriteHeader(http.StatusNoContent)
	return nil
}
