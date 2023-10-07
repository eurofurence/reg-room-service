package rooms

import (
	"context"
	"net/http"

	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
)

// UpdateRoomRequest is the request type for the UpdateRoom operation.
type UpdateRoomRequest struct {
	Room modelsv1.Room
}

// UpdateRoom updates an existing room by uuid. Note that you cannot use this to change the room members!
// Admin/Api Key only.
//
// Successful operations return status 201 with a location header that points to the updated resource.
func (h *Handler) UpdateRoom(ctx context.Context, req *UpdateRoomRequest, w http.ResponseWriter) (*modelsv1.Empty, error) {
	return nil, nil
}

// UpdateRoomRequest validates and creates the request for the UpdateRoom operation.
func (h *Handler) UpdateRoomRequest(r *http.Request, w http.ResponseWriter) (*UpdateRoomRequest, error) {
	return nil, nil
}

// UpdateRoomResponse writes out the response for the UpdateRoom operation.
func (h *Handler) UpdateRoomResponse(ctx context.Context, _ *modelsv1.Empty, w http.ResponseWriter) error {
	return nil
}
