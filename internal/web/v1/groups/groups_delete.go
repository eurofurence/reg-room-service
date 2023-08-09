package groups

import (
	"context"
	"net/http"
)

type DeleteGroupRequest struct{}

func (h *Handler) DeleteGroup(ctx context.Context, req *DeleteGroupRequest, w http.ResponseWriter) (*Empty, error) {
	return nil, nil
}

func (h *Handler) DeleteGroupRequest(_ *http.Request, w http.ResponseWriter) (*DeleteGroupRequest, error) {
	return nil, nil
}

func (h *Handler) DeleteGroupResponse(ctx context.Context, _ *Empty, w http.ResponseWriter) error {
	return nil
}

type RemoveGroupMemberRequest struct{}

func (h *Handler) RemoveGroupMember(ctx context.Context, req *RemoveGroupMemberRequest, w http.ResponseWriter) (*Empty, error) {
	return nil, nil
}

func (h *Handler) RemoveGroupMemberRequest(r *http.Request, w http.ResponseWriter) (*RemoveGroupMemberRequest, error) {
	return nil, nil
}

func (h *Handler) RemoveGroupMemberResponse(ctx context.Context, _ *Empty, w http.ResponseWriter) error {
	return nil
}
