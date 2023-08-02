package groups

import (
	"context"
	"net/http"

	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
	"github.com/eurofurence/reg-room-service/internal/logging"
)

type UpdateGroupRequest struct{}

func (h *Handler) UpdateGroup(ctx context.Context, req *UpdateGroupRequest, logger logging.Logger) (*modelsv1.Group, error) {
	return nil, nil
}

func (h *Handler) UpdateGroupRequest(_ *http.Request) (*UpdateGroupRequest, error) {
	return new(UpdateGroupRequest), nil
}

func (h *Handler) UpdateGroupResponse(ctx context.Context, res *modelsv1.Group, w http.ResponseWriter) error {
	return nil
}
