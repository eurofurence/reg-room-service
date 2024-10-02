package groupsctl

import (
	"context"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"github.com/eurofurence/reg-room-service/internal/application/web"
	"github.com/eurofurence/reg-room-service/internal/controller/v1/util"
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/url"

	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
	"github.com/eurofurence/reg-room-service/internal/application/common"
)

type UpdateGroupRequest struct {
	Group modelsv1.Group
}

// UpdateGroup is used to update an existing group by uuid. Note that you cannot use this to change the group members!
//
// Admins or the current group owner can change the group owner to any member of the group.
func (h *Controller) UpdateGroup(ctx context.Context, req *UpdateGroupRequest, w http.ResponseWriter) (*modelsv1.Empty, error) {
	if err := h.svc.UpdateGroup(ctx, req.Group); err != nil {
		web.SendErrorResponse(ctx, w, err)
		return nil, err
	}

	reqURL, ok := ctx.Value(common.CtxKeyRequestURL{}).(*url.URL)
	if !ok {
		aulogging.Error(ctx, "unable to retrieve URL from context")
		return nil, nil
	}

	w.Header().Set("Location", reqURL.Path)

	return nil, nil
}

func (h *Controller) UpdateGroupRequest(r *http.Request, w http.ResponseWriter) (*UpdateGroupRequest, error) {
	ctx := r.Context()

	groupID := chi.URLParam(r, "uuid")
	if err := validateGroupID(ctx, w, groupID); err != nil {
		return nil, err
	}

	var group modelsv1.Group

	if err := util.NewStrictJSONDecoder(r.Body).Decode(&group); err != nil {
		web.SendErrorResponse(ctx, w, common.NewBadRequest(ctx, common.GroupDataInvalid, common.Details("invalid json provided")))
		return nil, err
	}

	group.ID = groupID
	return &UpdateGroupRequest{group}, nil
}

func (h *Controller) UpdateGroupResponse(_ context.Context, _ *modelsv1.Empty, w http.ResponseWriter) error {
	w.WriteHeader(http.StatusOK)
	return nil
}
