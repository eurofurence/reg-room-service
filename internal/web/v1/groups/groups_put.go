package groups

import (
	"context"
	"net/http"

	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
)

type UpdateGroupRequest struct{}

func (h *Handler) UpdateGroup(ctx context.Context, req *UpdateGroupRequest) (*modelsv1.Group, error) {
	return nil, nil
}

func (h *Handler) UpdateGroupRequest(_ *http.Request) (*UpdateGroupRequest, error) {
	return new(UpdateGroupRequest), nil
}

func (h *Handler) UpdateGroupResponse(ctx context.Context, res *modelsv1.Group, w http.ResponseWriter) error {
	return nil
}
