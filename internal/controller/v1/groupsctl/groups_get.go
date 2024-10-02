package groupsctl

import (
	"context"
	"github.com/eurofurence/reg-room-service/internal/application/web"
	"github.com/eurofurence/reg-room-service/internal/controller/v1/util"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
	"github.com/eurofurence/reg-room-service/internal/application/common"
)

type ListGroupsRequest struct {
	MemberIDs []int64
	MinSize   uint
	MaxSize   int
}

func (h *Controller) ListGroups(ctx context.Context, req *ListGroupsRequest, w http.ResponseWriter) (*modelsv1.GroupList, error) {
	groups, err := h.svc.FindGroups(ctx, req.MinSize, req.MaxSize, req.MemberIDs)
	if err != nil {
		return nil, err
	}

	return &modelsv1.GroupList{
		Groups: groups,
	}, nil
}

func (h *Controller) ListGroupsRequest(r *http.Request, w http.ResponseWriter) (*ListGroupsRequest, error) {
	var req ListGroupsRequest

	ctx := r.Context()
	query := r.URL.Query()

	queryIDs := query.Get("member_ids")
	memberIDs, err := util.ParseMemberIDs(queryIDs)
	if err != nil {
		return nil, common.NewBadRequest(ctx, common.RequestParseFailed, common.Details(err.Error()))
	}

	req.MemberIDs = memberIDs
	if minSize := query.Get("min_size"); minSize != "" {
		val, err := util.ParseUInt[uint](minSize)
		if err != nil {
			return nil, common.NewBadRequest(ctx, common.RequestParseFailed, common.Details(err.Error()))
		}

		req.MinSize = val
	}

	if maxSize := query.Get("max_size"); maxSize != "" {
		val, err := util.ParseInt[int](maxSize)
		if err != nil {
			return nil, common.NewBadRequest(ctx, common.RequestParseFailed, common.Details(err.Error()))
		}
		if val < -1 {
			return nil, common.NewBadRequest(ctx, common.RequestParseFailed, common.Details("maxSize cannot be less than -1"))
		}

		req.MaxSize = val
	} else {
		req.MaxSize = -1
	}

	return &req, nil
}

func (h *Controller) ListGroupsResponse(ctx context.Context, res *modelsv1.GroupList, w http.ResponseWriter) error {
	return web.EncodeWithStatus(http.StatusOK, res, w)
}

type FindMyGroupRequest struct{}

func (h *Controller) FindMyGroup(ctx context.Context, req *FindMyGroupRequest, w http.ResponseWriter) (*modelsv1.Group, error) {
	group, err := h.svc.FindMyGroup(ctx)
	if err != nil {
		return nil, err
	}

	return group, nil
}

func (h *Controller) FindMyGroupRequest(r *http.Request, w http.ResponseWriter) (*FindMyGroupRequest, error) {
	// Endpoint only requires JWT token for now.
	return &FindMyGroupRequest{}, nil
}

func (h *Controller) FindMyGroupResponse(ctx context.Context, res *modelsv1.Group, w http.ResponseWriter) error {
	return web.EncodeWithStatus(http.StatusOK, res, w)
}

type FindGroupByIDRequest struct {
	GroupID string
}

func (h *Controller) FindGroupByID(ctx context.Context, req *FindGroupByIDRequest, w http.ResponseWriter) (*modelsv1.Group, error) {
	grp, err := h.svc.GetGroupByID(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}

	return grp, nil
}

func (h *Controller) FindGroupByIDRequest(r *http.Request, w http.ResponseWriter) (*FindGroupByIDRequest, error) {
	groupID := chi.URLParam(r, "uuid")
	if _, err := uuid.Parse(groupID); err != nil {
		return nil, common.NewBadRequest(r.Context(), common.GroupIDInvalid, url.Values{"details": []string{"you must specify a valid uuid"}})
	}

	req := &FindGroupByIDRequest{
		GroupID: groupID,
	}

	return req, nil
}

func (h *Controller) FindGroupByIDResponse(_ context.Context, res *modelsv1.Group, w http.ResponseWriter) error {
	return web.EncodeWithStatus(http.StatusOK, res, w)
}
