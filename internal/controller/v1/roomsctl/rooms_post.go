package roomsctl

import (
	"context"
	"errors"
	"github.com/eurofurence/reg-room-service/internal/application/common"
	"github.com/eurofurence/reg-room-service/internal/controller/v1/util"
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/url"
	"path"

	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
)

type CreateRoomRequest struct {
	// Room is the expected representation for the request body
	Room modelsv1.RoomCreate
}

// CreateRoom creates a new room without assignment.
//
// Endpoint access only for admin users or api token.
//
// Successful operations return status 201 with a location header that points to the created resource.
func (h *Controller) CreateRoom(ctx context.Context, req *CreateRoomRequest, w http.ResponseWriter) (*modelsv1.Empty, error) {
	newGroupUUID, err := h.svc.CreateRoom(ctx, &(req.Room))
	if err != nil {
		return nil, err
	}

	requestURL, ok := ctx.Value(common.CtxKeyRequestURL{}).(*url.URL)
	if !ok {
		return nil, errors.New("could not retrieve base URL from context - this is an implementation error")
	}

	w.Header().Set("Location", path.Join(requestURL.Path, newGroupUUID))
	return nil, nil
}

func (h *Controller) CreateRoomRequest(r *http.Request, w http.ResponseWriter) (*CreateRoomRequest, error) {
	var room modelsv1.RoomCreate

	if err := util.NewStrictJSONDecoder(r.Body).Decode(&room); err != nil {
		return nil, common.NewBadRequest(r.Context(), common.RoomDataInvalid, common.Details("invalid json provided"))
	}

	crr := &CreateRoomRequest{
		Room: room,
	}

	return crr, nil
}

func (h *Controller) CreateRoomResponse(ctx context.Context, _ *modelsv1.Empty, w http.ResponseWriter) error {
	w.WriteHeader(http.StatusCreated)
	return nil
}

type AddToRoomRequest struct {
	// RoomID is the uuid of the room
	RoomID string
	// BadgeNumber is the registration number of an attendee
	BadgeNumber int64
}

// AddToRoom adds an attendee to a room.
//
// See OpenAPI Spec for further details.
func (h *Controller) AddToRoom(ctx context.Context, req *AddToRoomRequest, w http.ResponseWriter) (*modelsv1.Empty, error) {
	err := h.svc.AddOccupantToRoom(ctx, req.RoomID, req.BadgeNumber)
	return &modelsv1.Empty{}, err
}

func (h *Controller) AddToRoomRequest(r *http.Request, w http.ResponseWriter) (*AddToRoomRequest, error) {
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

	return &AddToRoomRequest{
		RoomID:      roomID,
		BadgeNumber: badgeNumber,
	}, nil
}

func (h *Controller) AddToRoomResponse(ctx context.Context, _ *modelsv1.Empty, w http.ResponseWriter) error {
	w.WriteHeader(http.StatusNoContent)
	return nil
}
