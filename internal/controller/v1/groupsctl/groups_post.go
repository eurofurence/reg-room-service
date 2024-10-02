package groupsctl

import (
	"context"
	"errors"
	"github.com/eurofurence/reg-room-service/internal/controller/v1/util"
	"net/http"
	"net/url"
	"path"

	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
	"github.com/eurofurence/reg-room-service/internal/application/common"
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
		return nil, err
	}

	requestURL, ok := ctx.Value(common.CtxKeyRequestURL{}).(*url.URL)
	if !ok {
		return nil, errors.New("could not retrieve base URL from context - this is an implementation error")
	}

	w.Header().Set("Location", path.Join(requestURL.Path, newGroupUUID))
	return nil, nil
}

func (h *Controller) CreateGroupRequest(r *http.Request, w http.ResponseWriter) (*CreateGroupRequest, error) {
	var group modelsv1.GroupCreate

	if err := util.NewStrictJSONDecoder(r.Body).Decode(&group); err != nil {
		return nil, common.NewBadRequest(r.Context(), common.GroupDataInvalid, common.Details("invalid json provided"))
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
