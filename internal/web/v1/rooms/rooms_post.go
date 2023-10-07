package rooms

import (
	"context"
	"net/http"

	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
)

type CreateRoomRequest struct {
	Room modelsv1.Room
}

// CreateRoom creates a new room without assignment.
// Required Permissions: [admin]
//
// Successful operations return status 201 with a location header that points to the created resource.
func (h *Handler) CreateRoom(ctx context.Context, req *CreateRoomRequest, w http.ResponseWriter) (*modelsv1.Empty, error) {
	return nil, nil
}

// CreateRoomRequest validates and converts the request parameters into a `CreateRoomRequest` type.
func (h *Handler) CreateRoomRequest(r *http.Request, w http.ResponseWriter) (*CreateRoomRequest, error) {
	return nil, nil
}

// CreateRoomResponse will write out the response of the create room request.
func (h *Handler) CreateRoomResponse(ctx context.Context, e *modelsv1.Empty, w http.ResponseWriter) error {
	return nil
}
