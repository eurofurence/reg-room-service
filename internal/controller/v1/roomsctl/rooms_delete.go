package roomsctl

import (
	"context"
	"github.com/eurofurence/reg-room-service/internal/application/common"
	"github.com/eurofurence/reg-room-service/internal/controller/v1/util"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
	"net/url"

	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
)

type DeleteRoomRequest struct {
	UUID string
}

// DeleteRoom deletes an existing room by uuid.
//
// See OpenAPI Spec for further details.
func (h *Controller) DeleteRoom(ctx context.Context, req *DeleteRoomRequest, w http.ResponseWriter) (*modelsv1.Empty, error) {
	err := h.svc.DeleteRoom(ctx, req.UUID)
	return nil, err
}

func (h *Controller) DeleteRoomRequest(r *http.Request, w http.ResponseWriter) (*DeleteRoomRequest, error) {
	roomID := chi.URLParam(r, "uuid")
	if _, err := uuid.Parse(roomID); err != nil {
		return nil, common.NewBadRequest(r.Context(), common.RoomIDInvalid, url.Values{"details": []string{"you must specify a valid uuid"}})
	}

	req := &DeleteRoomRequest{
		UUID: roomID,
	}

	return req, nil
}

func (h *Controller) DeleteRoomResponse(ctx context.Context, _ *modelsv1.Empty, w http.ResponseWriter) error {
	w.WriteHeader(http.StatusNoContent)
	return nil
}

type RemoveFromRoomRequest struct {
	// RoomID is the uuid of the room
	RoomID string
	// BadgeNumber is the registration number of an attendee
	BadgeNumber int64
}

// RemoveFromRoom removes the attendee with the given badge number from the room.
//
// See OpenAPI Spec for further details.
func (h *Controller) RemoveFromRoom(ctx context.Context, req *RemoveFromRoomRequest, _ http.ResponseWriter) (*modelsv1.Empty, error) {
	err := h.svc.RemoveOccupantFromRoom(ctx, req.RoomID, req.BadgeNumber)
	return &modelsv1.Empty{}, err
}

func (h *Controller) RemoveFromRoomRequest(r *http.Request, _ http.ResponseWriter) (*RemoveFromRoomRequest, error) {
	ctx := r.Context()

	roomID := chi.URLParam(r, "uuid")
	if err := validateRoomID(ctx, roomID); err != nil {
		return nil, err
	}

	badge := chi.URLParam(r, "badgenumber")
	badgeNumber, err := util.ParseInt[int64](badge)
	if err != nil {
		return nil, common.NewBadRequest(ctx, common.RequestParseFailed, common.Details("invalid badge number - must be positive integer"), err)
	}
	if badgeNumber < 1 {
		return nil, common.NewBadRequest(ctx, common.RoomDataInvalid, common.Details("invalid badge number - must be positive integer"))
	}

	return &RemoveFromRoomRequest{
		RoomID:      roomID,
		BadgeNumber: badgeNumber,
	}, nil
}

func (h *Controller) RemoveFromRoomResponse(_ context.Context, _ *modelsv1.Empty, w http.ResponseWriter) error {
	w.WriteHeader(http.StatusNoContent)
	return nil
}
