package roomsctl

import (
	"context"
	"github.com/eurofurence/reg-room-service/internal/application/common"
	"github.com/eurofurence/reg-room-service/internal/application/web"
	"github.com/eurofurence/reg-room-service/internal/controller/v1/util"
	roomservice "github.com/eurofurence/reg-room-service/internal/service/rooms"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
	"net/url"

	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
)

func (h *Controller) ListRooms(ctx context.Context, req *roomservice.FindRoomParams, w http.ResponseWriter) (*modelsv1.RoomList, error) {
	rooms, err := h.svc.FindRooms(ctx, req)
	if err != nil {
		return nil, err
	}

	return &modelsv1.RoomList{
		Rooms: rooms,
	}, nil
}

func (h *Controller) ListRoomsRequest(r *http.Request, w http.ResponseWriter) (*roomservice.FindRoomParams, error) {
	var req roomservice.FindRoomParams

	ctx := r.Context()
	query := r.URL.Query()

	queryIDs := query.Get("occupant_ids")
	memberIDs, err := util.ParseMemberIDs(queryIDs)
	if err != nil {
		return nil, common.NewBadRequest(ctx, common.RequestParseFailed, common.Details(err.Error()))
	}
	req.MemberIDs = memberIDs

	if minSize := query.Get("min_size"); minSize != "" {
		val, err := util.ParseUInt[uint](minSize)
		if err != nil {
			return nil, common.NewBadRequest(ctx, common.RequestParseFailed, common.Details(err.Error()))
		}

		req.MinSize = val
	}

	if maxSize := query.Get("max_size"); maxSize != "" {
		val, err := util.ParseUInt[uint](maxSize)
		if err != nil {
			return nil, common.NewBadRequest(ctx, common.RequestParseFailed, common.Details(err.Error()))
		}

		req.MaxSize = val
	}

	if minOccupants := query.Get("min_occupants"); minOccupants != "" {
		val, err := util.ParseUInt[uint](minOccupants)
		if err != nil {
			return nil, common.NewBadRequest(ctx, common.RequestParseFailed, common.Details(err.Error()))
		}

		req.MinOccupants = val
	}

	if maxOccupants := query.Get("max_occupants"); maxOccupants != "" {
		val, err := util.ParseInt[int](maxOccupants)
		if err != nil {
			return nil, common.NewBadRequest(ctx, common.RequestParseFailed, common.Details(err.Error()))
		}
		if val < -1 {
			return nil, common.NewBadRequest(ctx, common.RequestParseFailed, common.Details("max_occupants cannot be less than -1"))
		}

		req.MaxOccupants = val
	} else {
		req.MaxOccupants = -1
	}

	return &req, nil
}

func (h *Controller) ListRoomsResponse(ctx context.Context, res *modelsv1.RoomList, w http.ResponseWriter) error {
	return web.EncodeWithStatus(http.StatusOK, res, w)
}

type FindMyRoomRequest struct{}

// FindMyRoom gets the room you are in. Must have a valid registration.
//
// See OpenAPI Spec for further details.
func (h *Controller) FindMyRoom(ctx context.Context, req *FindMyRoomRequest, w http.ResponseWriter) (*modelsv1.Room, error) {
	room, err := h.svc.FindMyRoom(ctx)
	if err != nil {
		return nil, err
	}

	return room, nil
}

func (h *Controller) FindMyRoomRequest(r *http.Request, w http.ResponseWriter) (*FindMyRoomRequest, error) {
	// Endpoint only requires logged-in user
	return &FindMyRoomRequest{}, nil
}

func (h *Controller) FindMyRoomResponse(ctx context.Context, res *modelsv1.Room, w http.ResponseWriter) error {
	return web.EncodeWithStatus(http.StatusOK, res, w)
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
