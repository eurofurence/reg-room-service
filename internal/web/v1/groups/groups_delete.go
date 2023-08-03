package groups

import (
	"context"
	"net/http"
)

type DeleteGroupRequest struct{}

func (h *Handler) DeleteGroup(ctx context.Context, req *DeleteGroupRequest) (*Empty, error) {
	return nil, nil
}

func (h *Handler) DeleteGroupRequest(_ *http.Request) (*DeleteGroupRequest, error) {
	return nil, nil
}

func (h *Handler) DeleteGroupResponse(ctx context.Context, _ *Empty, w http.ResponseWriter) error {
	return nil
}

type RemoveGroupMemberRequest struct{}

func (h *Handler) RemoveGroupMember(ctx context.Context, req *RemoveGroupMemberRequest) (*Empty, error) {
	return nil, nil
}

func (h *Handler) RemoveGroupMemberRequest(r *http.Request) (*RemoveGroupMemberRequest, error) {
	return nil, nil
}

func (h *Handler) RemoveGroupMemberResponse(ctx context.Context, _ *Empty, w http.ResponseWriter) error {
	return nil
}
