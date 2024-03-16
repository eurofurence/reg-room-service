package groups

import (
	"context"
	"fmt"
	apierrors "github.com/eurofurence/reg-room-service/internal/errors"
	"github.com/eurofurence/reg-room-service/internal/web/common"
	"github.com/eurofurence/reg-room-service/internal/web/v1/util"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"

	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
)

type UpdateGroupRequest struct {
	Group modelsv1.Group
}

func (h *Controller) UpdateGroup(ctx context.Context, req *UpdateGroupRequest, w http.ResponseWriter) (*modelsv1.Group, error) {
	return nil, nil
}

func (h *Controller) UpdateGroupRequest(r *http.Request, w http.ResponseWriter) (*UpdateGroupRequest, error) {
	groupID := chi.URLParam(r, "uuid")
	if err := uuid.Validate(groupID); err != nil {
		common.SendHTTPStatusErrorResponse(r.Context(), w, apierrors.NewBadRequest(common.GroupIDInvalidMessage, fmt.Sprintf("%q is not a vailid UUID", groupID)))
		return nil, err
	}

	var group modelsv1.Group

	if err := util.NewStrictJSONDecoder(r.Body).Decode(&group); err != nil {
		common.SendHTTPStatusErrorResponse(r.Context(), w, apierrors.NewBadRequest(common.GroupDataInvalidMessage, "invalid json provided"))
		return nil, err
	}

	return &UpdateGroupRequest{group}, nil
}

func (h *Controller) UpdateGroupResponse(_ context.Context, res *modelsv1.Group, w http.ResponseWriter) error {
	return common.EncodeWithStatus(http.StatusOK, res, w)
}
