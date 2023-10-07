package groups

import (
	"context"
	"net/http"

	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
	apierrors "github.com/eurofurence/reg-room-service/internal/errors"
	"github.com/eurofurence/reg-room-service/internal/web/common"
	"github.com/eurofurence/reg-room-service/internal/web/v1/util"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ListGroupsRequest struct {
	MemberIDs []string
	MinSize   uint
	MaxSize   uint
}

func (h *Handler) ListGroups(ctx context.Context, req *ListGroupsRequest, w http.ResponseWriter) (*modelsv1.GroupList, error) {
	// TODO implement

	return nil, nil
}

func (h *Handler) ListGroupsRequest(r *http.Request, w http.ResponseWriter) (*ListGroupsRequest, error) {
	var req ListGroupsRequest

	ctx := r.Context()
	query := r.URL.Query()

	queryIDs := query.Get("member_ids")
	memberIDs, err := util.ParseMemberIDs(queryIDs)
	if err != nil {
		common.SendHTTPStatusErrorResponse(ctx, w, apierrors.NewBadRequest("group.data.invalid", err.Error()))
		return nil, err
	}

	req.MemberIDs = memberIDs
	if minSize := query.Get("min_size"); minSize != "" {
		val, err := util.ParseUInt[uint](minSize)
		if err != nil {
			common.SendHTTPStatusErrorResponse(ctx, w, apierrors.NewBadRequest("group.data.invalid", err.Error()))
			return nil, err
		}

		req.MinSize = val
	}

	if maxSize := query.Get("max_size"); maxSize != "" {
		val, err := util.ParseUInt[uint](maxSize)
		if err != nil {
			common.SendHTTPStatusErrorResponse(ctx, w, apierrors.NewBadRequest("group.data.invalid", err.Error()))
			return nil, err
		}

		req.MaxSize = val
	}

	return &req, nil
}

// ListGroupsResponse writes out the result from the ListGroups operation.
func (h *Handler) ListGroupsResponse(ctx context.Context, res *modelsv1.GroupList, w http.ResponseWriter) error {
	return nil
}

type FindMyGroupRequest struct{}

// FindMyGroup TODO
func (h *Handler) FindMyGroup(ctx context.Context, req *FindMyGroupRequest, w http.ResponseWriter) (*modelsv1.Group, error) {
	return nil, nil
}

func (h *Handler) FindMyGroupRequest(r *http.Request, w http.ResponseWriter) (*FindMyGroupRequest, error) {
	// Endpoint only requires JWT token for now.
	return nil, nil
}

func (h *Handler) FindMyGroupResponse(ctx context.Context, res *modelsv1.Group, w http.ResponseWriter) error {
	return nil
}

type FindGroupByIDRequest struct {
	GroupID string
}

func (h *Handler) FindGroupByID(ctx context.Context, req *FindGroupByIDRequest, w http.ResponseWriter) (*modelsv1.Group, error) {
	return nil, nil
}

func (h *Handler) FindGroupByIDRequest(r *http.Request, w http.ResponseWriter) (*FindGroupByIDRequest, error) {
	groupID := chi.URLParam(r, "uuid")
	if _, err := uuid.Parse(groupID); err != nil {
		common.SendHTTPStatusErrorResponse(r.Context(), w, apierrors.NewBadRequest("group.id.invalid", ""))
		return nil, err
	}

	req := &FindGroupByIDRequest{
		GroupID: groupID,
	}

	return req, nil
}

func (h *Handler) FindGroupByIDResponse(ctx context.Context, res *modelsv1.Group, w http.ResponseWriter) error {
	return nil
}
