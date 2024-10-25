package roomsctl

import (
	"context"
	"errors"
	"github.com/eurofurence/reg-room-service/internal/application/common"
	"github.com/eurofurence/reg-room-service/internal/controller/v1/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-http-utils/headers"
	"net/http"
	"net/url"

	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
)

// UpdateRoomRequest is the request type for the UpdateRoom operation.
type UpdateRoomRequest struct {
	Room modelsv1.Room
}

// UpdateRoom updates an existing room by uuid. Note that you cannot use this to change the room members!
//
// See OpenAPI Spec for further details.
func (h *Controller) UpdateRoom(ctx context.Context, req *UpdateRoomRequest, w http.ResponseWriter) (*modelsv1.Empty, error) {
	if err := h.svc.UpdateRoom(ctx, &(req.Room)); err != nil {
		return nil, err
	}

	reqURL, ok := ctx.Value(common.CtxKeyRequestURL{}).(*url.URL)
	if !ok {
		return nil, errors.New("unable to retrieve URL from context - this is an implementation error")
	}

	w.Header().Set(headers.Location, reqURL.Path)

	return nil, nil
}

func (h *Controller) UpdateRoomRequest(r *http.Request, w http.ResponseWriter) (*UpdateRoomRequest, error) {
	ctx := r.Context()

	roomID := chi.URLParam(r, "uuid")
	if err := validateRoomID(ctx, roomID); err != nil {
		return nil, err
	}

	var room modelsv1.Room

	if err := util.NewStrictJSONDecoder(r.Body).Decode(&room); err != nil {
		return nil, common.NewBadRequest(ctx, common.RoomDataInvalid, common.Details("invalid json provided"))
	}

	room.ID = roomID
	return &UpdateRoomRequest{Room: room}, nil
}

func (h *Controller) UpdateRoomResponse(ctx context.Context, res *modelsv1.Empty, w http.ResponseWriter) error {
	w.WriteHeader(http.StatusNoContent)
	return nil
}
