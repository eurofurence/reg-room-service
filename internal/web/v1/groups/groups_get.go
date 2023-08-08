package groups

import (
	"context"
	"net/http"

	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
)

type ListGroupsRequest struct {
	MemberIDs string
	MinSize   int
	MaxSize   int
}

func (h *Handler) ListGroups(ctx context.Context, req *ListGroupsRequest, w http.ResponseWriter) (*modelsv1.GroupList, error) {
	// TODO implement

	return nil, nil
}

func (h *Handler) ListGroupsRequest(r *http.Request, w http.ResponseWriter) (*ListGroupsRequest, error) {
	return nil, nil
}

func (h *Handler) ListGroupsResponse(ctx context.Context, res *modelsv1.GroupList, w http.ResponseWriter) error {
	return nil
}

type FindMyGroupRequest struct{}

func (h *Handler) FindMyGroup(ctx context.Context, req *FindMyGroupRequest, w http.ResponseWriter) (*modelsv1.Group, error) {
	return nil, nil
}

func (h *Handler) GetMyGroupRequest(r *http.Request, w http.ResponseWriter) (*FindMyGroupRequest, error) {
	return new(FindMyGroupRequest), nil
}

func (h *Handler) FindMyGroupResponse(ctx context.Context, res *modelsv1.Group, w http.ResponseWriter) error {
	return nil
}

type FindGroupByIDRequest struct{}

func (h *Handler) FindGroupByID(ctx context.Context, req *FindGroupByIDRequest, w http.ResponseWriter) (*modelsv1.Group, error) {
	return nil, nil
}

func (h *Handler) FindGroupByIDRequest(_ *http.Request, w http.ResponseWriter) (*FindGroupByIDRequest, error) {
	return new(FindGroupByIDRequest), nil
}

func (h *Handler) FindGroupByIDResponse(ctx context.Context, res *modelsv1.Group, w http.ResponseWriter) error {
	return nil
}
