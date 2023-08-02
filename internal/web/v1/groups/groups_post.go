package groups

import (
	"context"
	"net/http"

	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
	"github.com/eurofurence/reg-room-service/internal/logging"
)

type CreateGroupRequest struct {
	Group modelsv1.Group
}

func (h *Handler) CreateGroup(ctx context.Context, req *CreateGroupRequest, logger logging.Logger) (*Empty, error) {
	// TODO
	return nil, nil
}

func (h *Handler) CreateGroupRequest(r *http.Request) (*CreateGroupRequest, error) {
	// TODO
	return nil, nil
}

func (h *Handler) CreateGroupResponse(ctx context.Context, _ *Empty, w http.ResponseWriter) error {
	// TODO
	return nil
}

type AddMemberToGroupRequest struct{}

func (h *Handler) AddMemberToGroup(ctx context.Context, req *AddMemberToGroupRequest, logger logging.Logger) (*Empty, error) {
	return nil, nil
}

func (h *Handler) AddMemberToGroupRequest(r *http.Request) (*AddMemberToGroupRequest, error) {
	return nil, nil
}

func (h *Handler) AddMemberToGroupResponse(ctx context.Context, _ *Empty, w http.ResponseWriter) error {
	return nil
}
