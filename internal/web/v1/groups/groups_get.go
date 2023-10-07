package groups

import (
	"context"
	"net/http"

	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
	apierrors "github.com/eurofurence/reg-room-service/internal/errors"
	"github.com/eurofurence/reg-room-service/internal/web/common"
	"github.com/eurofurence/reg-room-service/internal/web/v1/util"
)

type ListGroupsRequest struct {
	MemberIDs []string
	MinSize   int
	MaxSize   int
}

func (h *Handler) ListGroups(ctx context.Context, req *ListGroupsRequest, w http.ResponseWriter) (*modelsv1.GroupList, error) {
	// TODO implement

	return nil, nil
}

func (h *Handler) ListGroupsRequest(r *http.Request, w http.ResponseWriter) (*ListGroupsRequest, error) {
	var req ListGroupsRequest

	queryIDs := r.URL.Query().Get("member_ids")
	memberIDs, err := util.ParseMemberIDs(queryIDs)
	if err != nil {
		common.SendHttpStatusErrorResponse(
			r.Context(),
			w,
			apierrors.NewBadRequest("test.test.test", err.Error()),
		)
		return nil, err
	}

	req.MemberIDs = memberIDs

	return &req, nil
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
