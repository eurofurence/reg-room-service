package groupsctl

import (
	"context"
	"github.com/eurofurence/reg-room-service/internal/controller/v1/util"
	groupservice "github.com/eurofurence/reg-room-service/internal/service/groups"
	"net/http"

	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
	"github.com/eurofurence/reg-room-service/internal/application/common"
	"github.com/go-chi/chi/v5"
)

// AddMemberToGroup adds an attendee to a group.
//
// Details see OpenAPI spec.
func (h *Controller) AddMemberToGroup(ctx context.Context, req *groupservice.AddGroupMemberParams, w http.ResponseWriter) (*modelsv1.Empty, error) {
	// TODO

	return &modelsv1.Empty{}, h.svc.AddMemberToGroup(ctx, *req)
}

// AddMemberToGroupRequest validates and creates the request for the AddMemberToGroup operation.
func (h *Controller) AddMemberToGroupRequest(r *http.Request, w http.ResponseWriter) (*groupservice.AddGroupMemberParams, error) {
	ctx := r.Context()

	groupID := chi.URLParam(r, "uuid")
	if err := validateGroupID(ctx, groupID); err != nil {
		return nil, err
	}

	badge := chi.URLParam(r, "badgenumber")
	badgeNumber, err := util.ParseInt[int64](badge)
	if err != nil {
		return nil, common.NewBadRequest(ctx, common.GroupDataInvalid, common.Details("invalid badge number - must be positive integer"), err)
	}
	if badgeNumber < 1 {
		return nil, common.NewBadRequest(ctx, common.GroupDataInvalid, common.Details("invalid badge number - must be positive integer"))
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
		return nil, common.NewBadRequest(ctx, common.GroupDataInvalid, nil, err)
	}

	req.Force = force
	return req, nil
}

// AddMemberToGroupResponse writes out the response for the AddMemberToGroup operation.
func (h *Controller) AddMemberToGroupResponse(_ context.Context, _ *modelsv1.Empty, w http.ResponseWriter) error {
	w.WriteHeader(http.StatusNoContent)
	return nil
}
