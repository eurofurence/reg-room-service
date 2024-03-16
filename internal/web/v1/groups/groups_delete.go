package groups

import (
	"context"
	"net/http"

	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
)

type DeleteGroupRequest struct{}

func (h *Controller) DeleteGroup(ctx context.Context, req *DeleteGroupRequest, w http.ResponseWriter) (*modelsv1.Empty, error) {
	return nil, nil
}

func (h *Controller) DeleteGroupRequest(_ *http.Request, w http.ResponseWriter) (*DeleteGroupRequest, error) {
	return nil, nil
}

func (h *Controller) DeleteGroupResponse(ctx context.Context, _ *modelsv1.Empty, w http.ResponseWriter) error {
	return nil
}

type RemoveGroupMemberRequest struct{}

func (h *Controller) RemoveGroupMember(ctx context.Context, req *RemoveGroupMemberRequest, w http.ResponseWriter) (*modelsv1.Empty, error) {
	return nil, nil
}

func (h *Controller) RemoveGroupMemberRequest(r *http.Request, w http.ResponseWriter) (*RemoveGroupMemberRequest, error) {
	return nil, nil
}

func (h *Controller) RemoveGroupMemberResponse(ctx context.Context, _ *modelsv1.Empty, w http.ResponseWriter) error {
	return nil
}
