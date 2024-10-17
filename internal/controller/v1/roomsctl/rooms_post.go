package roomsctl

import (
	"context"
	"net/http"

	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
)

type CreateRoomRequest struct {
	Room modelsv1.Room
}

// CreateRoom creates a new room without assignment.
// Endpoint access only for admin users.
//
// Successful operations return status 201 with a location header that points to the created resource.
func (h *Controller) CreateRoom(ctx context.Context, req *CreateRoomRequest, w http.ResponseWriter) (*modelsv1.Empty, error) {
	return nil, nil
}

func (h *Controller) CreateRoomRequest(r *http.Request, w http.ResponseWriter) (*CreateRoomRequest, error) {
	return nil, nil
}

func (h *Controller) CreateRoomResponse(ctx context.Context, _ *modelsv1.Empty, w http.ResponseWriter) error {
	return nil
}

type AddToRoomRequest struct {
}

// AddToRoom adds an attendee to a room.
//
// See OpenAPI Spec for further details.
func (h *Controller) AddToRoom(ctx context.Context, req *AddToRoomRequest, w http.ResponseWriter) (*modelsv1.Empty, error) {
	return nil, nil
}

func (h *Controller) AddToRoomRequest(r *http.Request, w http.ResponseWriter) (*AddToRoomRequest, error) {
	return nil, nil
}

func (h *Controller) AddToRoomResponse(ctx context.Context, _ *modelsv1.Empty, w http.ResponseWriter) error {
	return nil
}
