package groups

import (
	"context"
	"net/http"

	"github.com/eurofurence/reg-room-service/internal/logging"
)

type DeleteGroupRequest struct{}

func (h *Handler) DeleteGroup(ctx context.Context, req *DeleteGroupRequest, logger logging.Logger) (*Empty, error) {
	return nil, nil
}

func (h *Handler) DeleteGroupRequest(_ *http.Request) (*DeleteGroupRequest, error) {
	return nil, nil
}

func (h *Handler) DeleteGroupResponse(ctx context.Context, _ *Empty, w http.ResponseWriter) error {
	return nil
}

type RemoveGroupMemberRequest struct{}

func (h *Handler) RemoveGroupMember(ctx context.Context, req *RemoveGroupMemberRequest, logger logging.Logger) (*Empty, error) {

	return nil, nil
}

func (h *Handler) RemoveGroupMemberRequest(r *http.Request) (*RemoveGroupMemberRequest, error) {
	return nil, nil
}

func (h *Handler) RemoveGroupMemberResponse(ctx context.Context, _ *Empty, w http.ResponseWriter) error {

	return nil
}
