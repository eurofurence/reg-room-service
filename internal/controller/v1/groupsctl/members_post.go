package groupsctl

import (
	"context"
	"errors"
	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
	"github.com/eurofurence/reg-room-service/internal/application/common"
	"github.com/eurofurence/reg-room-service/internal/controller/v1/util"
	groupservice "github.com/eurofurence/reg-room-service/internal/service/groups"
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/url"
)

// AddMemberToGroup adds an attendee to a group.
//
// Details see OpenAPI spec.
func (h *Controller) AddMemberToGroup(ctx context.Context, req *groupservice.AddGroupMemberParams, w http.ResponseWriter) (*modelsv1.Empty, error) {
	requestURL, ok := ctx.Value(common.CtxKeyRequestURL{}).(*url.URL)
	if !ok {
		return nil, errors.New("could not retrieve base URL from context - this is an implementation error")
	}

	urlExtension, err := h.svc.AddMemberToGroup(ctx, req)
	if err != nil {
		return &modelsv1.Empty{}, err
	}

	w.Header().Set("Location", requestURL.Path+urlExtension)

	return &modelsv1.Empty{}, nil
}

// AddMemberToGroupRequest validates and creates the request for the AddMemberToGroup operation.
func (h *Controller) AddMemberToGroupRequest(r *http.Request, w http.ResponseWriter) (*groupservice.AddGroupMemberParams, error) {
	ctx := r.Context()
	query := r.URL.Query()

	groupID := chi.URLParam(r, "uuid")
	if err := validateGroupID(ctx, groupID); err != nil {
		return nil, err
	}

	badge := chi.URLParam(r, "badgenumber")
	badgeNumber, err := util.ParseInt[int64](badge)
	if err != nil {
		return nil, common.NewBadRequest(ctx, common.RequestParseFailed, common.Details("invalid badge number - must be positive integer"), err)
	}
	if badgeNumber < 1 {
		return nil, common.NewBadRequest(ctx, common.GroupDataInvalid, common.Details("invalid badge number - must be positive integer"))
	}

	force, err := util.ParseOptionalBool(query.Get("force"))
	if err != nil {
		return nil, common.NewBadRequest(ctx, common.RequestParseFailed, common.Details("invalid force parameter, try true, 1, false, 0 or omit"), err)
	}

	return &groupservice.AddGroupMemberParams{
		GroupID:     groupID,
		BadgeNumber: badgeNumber,
		Nickname:    query.Get("nickname"),
		Code:        query.Get("code"),
		Force:       force,
	}, nil
}

// AddMemberToGroupResponse writes out the response for the AddMemberToGroup operation.
func (h *Controller) AddMemberToGroupResponse(_ context.Context, _ *modelsv1.Empty, w http.ResponseWriter) error {
	w.WriteHeader(http.StatusNoContent)
	return nil
}
