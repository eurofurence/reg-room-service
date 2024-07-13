package groups

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
	apierrors "github.com/eurofurence/reg-room-service/internal/errors"
	"github.com/eurofurence/reg-room-service/internal/web/common"
	"github.com/eurofurence/reg-room-service/internal/web/v1/util"
)

type ListGroupsRequest struct {
	MemberIDs []uint
	MinSize   uint
	MaxSize   int
}

func (h *Controller) ListGroups(ctx context.Context, req *ListGroupsRequest, w http.ResponseWriter) (*modelsv1.GroupList, error) {
	groups, err := h.svc.FindGroups(ctx, req.MinSize, req.MaxSize, req.MemberIDs)
	if err != nil {
		common.SendErrorResponse(ctx, w, err)
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
		val, err := util.ParseInt[int](maxSize)
		if err != nil {
			common.SendHTTPStatusErrorResponse(ctx, w, apierrors.NewBadRequest("group.data.invalid", err.Error()))
			return nil, err
		}
		if val < -1 {
			common.SendHTTPStatusErrorResponse(ctx, w, apierrors.NewBadRequest("group.data.invalid", "maxSize cannot be less than -1"))
			return nil, err
		}

		req.MaxSize = val
	} else {
		req.MaxSize = -1
	}

	return &req, nil
}

func (h *Controller) ListGroupsResponse(ctx context.Context, res *modelsv1.GroupList, w http.ResponseWriter) error {
	return common.EncodeWithStatus(http.StatusOK, res, w)
}

type FindMyGroupRequest struct{}

// FindMyGroup TODO.
func (h *Controller) FindMyGroup(ctx context.Context, req *FindMyGroupRequest, w http.ResponseWriter) (*modelsv1.Group, error) {
	return nil, nil
}

func (h *Controller) FindMyGroupRequest(r *http.Request, w http.ResponseWriter) (*FindMyGroupRequest, error) {
	// Endpoint only requires JWT token for now.
	return nil, nil
}

func (h *Controller) FindMyGroupResponse(ctx context.Context, res *modelsv1.Group, w http.ResponseWriter) error {
	return nil
}

type FindGroupByIDRequest struct {
	GroupID string
}

func (h *Controller) FindGroupByID(ctx context.Context, req *FindGroupByIDRequest, w http.ResponseWriter) (*modelsv1.Group, error) {
	grp, err := h.svc.GetGroupByID(ctx, req.GroupID)
	if err != nil {
		var statusErr apierrors.APIStatus
		if !errors.As(err, &statusErr) {
			common.SendHTTPStatusErrorResponse(ctx, w, apierrors.NewInternalServerError(common.InternalErrorMessage, err.Error()))
			return nil, err
		}

		common.SendHTTPStatusErrorResponse(ctx, w, statusErr)
		return nil, statusErr
	}

	return grp, nil
}

func (h *Controller) FindGroupByIDRequest(r *http.Request, w http.ResponseWriter) (*FindGroupByIDRequest, error) {
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

func (h *Controller) FindGroupByIDResponse(_ context.Context, res *modelsv1.Group, w http.ResponseWriter) error {
	return common.EncodeWithStatus(http.StatusOK, res, w)
}
