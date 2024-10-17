package roomsctl

import (
	"context"
	"net/http"

	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
)

type ListRoomsRequest struct {
	OccupantIDs  []int64
	MinSize      uint
	MaxSize      uint
	MinOccupants uint
	MaxOccupants uint
}

func (h *Controller) ListRooms(ctx context.Context, req *ListRoomsRequest, w http.ResponseWriter) (*modelsv1.RoomList, error) {
	return nil, nil
}

func (h *Controller) ListRoomsRequest(r *http.Request, w http.ResponseWriter) (*ListRoomsRequest, error) {
	return nil, nil
}

func (h *Controller) ListRoomsResponse(ctx context.Context, res *modelsv1.RoomList, w http.ResponseWriter) error {
	return nil
}

type FindMyRoomRequest struct{}

// FindMyRoom gets the room you are in. Must have a valid registration.
//
// See OpenAPI Spec for further details.
func (h *Controller) FindMyRoom(ctx context.Context, req *FindMyRoomRequest, w http.ResponseWriter) (*modelsv1.Room, error) {
	return nil, nil
}

func (h *Controller) FindMyRoomRequest(r *http.Request, w http.ResponseWriter) (*FindMyRoomRequest, error) {
	return nil, nil
}

func (h *Controller) FindMyRoomResponse(ctx context.Context, res *modelsv1.Room, w http.ResponseWriter) error {
	return nil
}

type GetRoomByIDRequest struct {
	UUID string
}

// GetRoomByID returns a single room. Admin/API key only.
//
// See OpenAPI Spec for further details.
func (h *Controller) GetRoomByID(ctx context.Context, req *GetRoomByIDRequest, w http.ResponseWriter) (*modelsv1.Room, error) {
	return nil, nil
}

func (h *Controller) GetRoomByIDRequest(r *http.Request, w http.ResponseWriter) (*GetRoomByIDRequest, error) {
	return nil, nil
}

func (h *Controller) GetRoomByIDResponse(ctx context.Context, res *modelsv1.Room, w http.ResponseWriter) error {
	return nil
}
