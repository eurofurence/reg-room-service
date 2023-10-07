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

type CreateGroupRequest struct {
	Group modelsv1.Group
}

func (h *Handler) CreateGroup(ctx context.Context, req *CreateGroupRequest, w http.ResponseWriter) (*modelsv1.Empty, error) {
	// TODO
	return nil, nil
}

func (h *Handler) CreateGroupRequest(r *http.Request, w http.ResponseWriter) (*CreateGroupRequest, error) {
	var group modelsv1.Group

	if err := util.NewStrictJSONDecoder(r.Body).Decode(&group); err != nil {
		common.SendHTTPStatusErrorResponse(r.Context(), w, apierrors.NewBadRequest(
			"group.data.invalid", "please check if your provided JSON is valid",
		))
		return nil, err
	}

	cgr := &CreateGroupRequest{
		Group: group,
	}

	return cgr, nil
}

func (h *Handler) CreateGroupResponse(ctx context.Context, _ *modelsv1.Empty, w http.ResponseWriter) error {
	// TODO
	return nil
}

type AddMemberToGroupRequest struct {
	GroupID     string
	BadgeNumber uint
}

// AddMemberToGroup
//
// Route: /groups/{uuid}/members/{badgenumber}
func (h *Handler) AddMemberToGroup(ctx context.Context, req *AddMemberToGroupRequest, w http.ResponseWriter) (*modelsv1.Empty, error) {
	return nil, nil
}

func (h *Handler) AddMemberToGroupRequest(r *http.Request, w http.ResponseWriter) (*AddMemberToGroupRequest, error) {
	groupID := chi.URLParam(r, "uuid")
	if _, err := uuid.Parse(groupID); err != nil {
		common.SendHTTPStatusErrorResponse(r.Context(), w, apierrors.NewBadRequest("group.id.invalid", ""))
		return nil, err
	}

	badge := chi.URLParam(r, "badgenumber")
	badgeNumber, err := util.ParseUInt[uint](badge)
	if err != nil {
		common.SendHTTPStatusErrorResponse(r.Context(), w, apierrors.NewBadRequest("group.data.invalid", "invalid type for badge number"))
		return nil, err
	}

	query := r.URL.Query()
	query.Get("nickname")
	query.Get("code")
	query.Get("force")

	req := &AddMemberToGroupRequest{
		GroupID:     groupID,
		BadgeNumber: badgeNumber,
	}

	return req, nil
}

func (h *Handler) AddMemberToGroupResponse(ctx context.Context, _ *modelsv1.Empty, w http.ResponseWriter) error {
	return nil
}
