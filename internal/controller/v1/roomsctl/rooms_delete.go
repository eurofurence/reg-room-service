package roomsctl

import (
	"context"
	"net/http"

	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
)

type DeleteRoomRequest struct {
	UUID string
}

// DeleteRoom deletes an existing room by uuid.
//
// See OpenAPI Spec for further details.
func (h *Controller) DeleteRoom(ctx context.Context, req *DeleteRoomRequest, w http.ResponseWriter) (*modelsv1.Empty, error) {
	return nil, nil
}

func (h *Controller) DeleteRoomRequest(r *http.Request, w http.ResponseWriter) (*DeleteRoomRequest, error) {
	return nil, nil
}

func (h *Controller) DeleteRoomResponse(ctx context.Context, _ *modelsv1.Empty, w http.ResponseWriter) error {
	return nil
}

type RemoveFromRoomRequest struct {
	UUID        string
	badgenumber int64
}

// RemoveFromRoom removes the attendee with the given badge number from the room.
//
// See OpenAPI Spec for further details.
func (h *Controller) RemoveFromRoom(ctx context.Context, req *RemoveFromRoomRequest, w http.ResponseWriter) (*modelsv1.Empty, error) {
	return nil, nil
}

func (h *Controller) RemoveFromRoomRequest(r *http.Request, w http.ResponseWriter) (*RemoveFromRoomRequest, error) {
	return nil, nil
}

func (h *Controller) RemoveFromRoomResponse(ctx context.Context, _ *modelsv1.Empty, w http.ResponseWriter) error {
	return nil
}
