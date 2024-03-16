package groups

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
	"net/url"

	apierrors "github.com/eurofurence/reg-room-service/internal/errors"
	"github.com/eurofurence/reg-room-service/internal/logging"
	"github.com/eurofurence/reg-room-service/internal/web/common"
	"github.com/eurofurence/reg-room-service/internal/web/v1/util"

	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
)

type UpdateGroupRequest struct {
	Group modelsv1.Group
}

// UpdateGroup is used to update an existing group by uuid. Note that you cannot use this to change the group members!
//
//	Admins or the current group owner can change the group owner to any member of the group.
func (h *Controller) UpdateGroup(ctx context.Context, req *UpdateGroupRequest, w http.ResponseWriter) (*modelsv1.Empty, error) {
	logger := logging.LoggerFromContext(ctx)
	if err := h.ctrl.UpdateGroup(ctx, req.Group); err != nil {
		var statusErr apierrors.APIStatus
		if errors.As(err, &statusErr) {
			common.SendHTTPStatusErrorResponse(ctx, w, statusErr)
			return nil, err
		}

		common.SendHTTPStatusErrorResponse(ctx, w, apierrors.NewInternalServerError(common.InternalErrorMessage, "unexpected error when updating group"))
		return nil, err
	}

	reqURL, ok := ctx.Value(common.CtxKeyRequestURL{}).(*url.URL)
	if !ok {
		logger.Error("unable to retrieve URL from context")
		return nil, nil
	}

	w.Header().Set("Location", reqURL.Path)

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

	group.ID = groupID
	return &UpdateGroupRequest{group}, nil
}

func (h *Controller) UpdateGroupResponse(_ context.Context, _ *modelsv1.Empty, w http.ResponseWriter) error {
	w.WriteHeader(http.StatusOK)
	return nil
}
