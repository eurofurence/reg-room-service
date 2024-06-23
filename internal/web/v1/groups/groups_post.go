package groups

import (
	"context"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"net/http"
	"net/url"
	"path"

	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
	apierrors "github.com/eurofurence/reg-room-service/internal/errors"
	"github.com/eurofurence/reg-room-service/internal/web/common"
	"github.com/eurofurence/reg-room-service/internal/web/v1/util"
)

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
	newGroupUUID, err := h.svc.CreateGroup(ctx, req.Group)
	if err != nil {
		common.SendErrorResponse(ctx, w, err)
		return nil, err
	}

	requestURL, ok := ctx.Value(common.CtxKeyRequestURL{}).(*url.URL)
	if !ok {
		aulogging.Error(ctx, "could not retrieve base URL from context")
		common.SendErrorResponse(ctx, w, nil)
		return nil, nil
	}

	w.Header().Set("Location", path.Join(requestURL.Path, newGroupUUID))
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
