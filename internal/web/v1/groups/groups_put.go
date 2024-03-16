package groups

import (
	"context"
	"net/http"

	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
)

type UpdateGroupRequest struct{}

func (h *Controller) UpdateGroup(ctx context.Context, req *UpdateGroupRequest, w http.ResponseWriter) (*modelsv1.Group, error) {
	return nil, nil
}

func (h *Controller) UpdateGroupRequest(_ *http.Request, w http.ResponseWriter) (*UpdateGroupRequest, error) {
	return new(UpdateGroupRequest), nil
}

func (h *Controller) UpdateGroupResponse(ctx context.Context, res *modelsv1.Group, w http.ResponseWriter) error {
	return nil
}
