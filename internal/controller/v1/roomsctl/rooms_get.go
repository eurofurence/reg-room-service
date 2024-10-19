package roomsctl

import (
	"context"
	"github.com/eurofurence/reg-room-service/internal/application/common"
	"github.com/eurofurence/reg-room-service/internal/application/web"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
	"net/url"

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
	room, err := h.svc.GetRoomByID(ctx, req.UUID)
	if err != nil {
		return nil, err
	}

	return room, nil
}

func (h *Controller) GetRoomByIDRequest(r *http.Request, w http.ResponseWriter) (*GetRoomByIDRequest, error) {
	roomID := chi.URLParam(r, "uuid")
	if _, err := uuid.Parse(roomID); err != nil {
		return nil, common.NewBadRequest(r.Context(), common.RoomIDInvalid, url.Values{"details": []string{"you must specify a valid uuid"}})
	}

	req := &GetRoomByIDRequest{
		UUID: roomID,
	}

	return req, nil
}

func (h *Controller) GetRoomByIDResponse(ctx context.Context, res *modelsv1.Room, w http.ResponseWriter) error {
	return web.EncodeWithStatus(http.StatusOK, res, w)
}
