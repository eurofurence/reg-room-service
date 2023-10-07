package rooms

import (
	"context"
	"net/http"

	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
)

type ListRoomsRequest struct {
	MemberIDs []string
	MinSize   uint
	MaxSize   uint
}

func (h *Handler) ListRooms(ctx context.Context, req *ListRoomsRequest, w http.ResponseWriter) (*modelsv1.RoomList, error) {
	return nil, nil
}

func (h *Handler) ListRoomsRequest(r *http.Request, w http.ResponseWriter) (*ListRoomsRequest, error) {
	return nil, nil
}

func (h *Handler) ListRoomsResponse(ctx context.Context, res *modelsv1.RoomList, w http.ResponseWriter) error {
	return nil
}

// FindMyRoomRequest ist the request type for the
// corresponding operation.
type FindMyRoomRequest struct{}

// FindMyRooom gets the room you are in. Must have a valid registration.
//
// Visibility of this information depends on the "final" flag that is set on the room, so admins can start planning
// room assignments without them becoming immediately visible to users.
//
// This endpoint works even for admins, giving them the room they are in.
// Because the user identity is taken from the logged in user, this does not work for Api Key authorization.
// Use the /rooms endpoint with member_id parameter instead.
func (h *Handler) FindMyRooom(ctx context.Context, req *FindMyRoomRequest, w http.ResponseWriter) (*modelsv1.Room, error) {
	return nil, nil
}

// FindMyRoomRequest validates and creates a `FindMyRoomRequest` from the http request.
func (h *Handler) FindMyRoomRequest(r *http.Request, w http.ResponseWriter) (*FindMyRoomRequest, error) {
	return nil, nil
}

// FindMyRoomResponse writes out the result from the FindMyRoom operation.
func (h *Handler) FindMyRoomResponse(ctx context.Context, res *modelsv1.Room, w http.ResponseWriter) error {
	return nil
}

// FindRoomByIDRequest is the request type for the
// FindRoomByUUID operation.
type FindRoomByIDRequest struct {
	UUID string
}

// FindRoomByUUID searched for a room by a given UUID and returns the room if it was found.
func (h *Handler) FindRoomByUUID(ctx context.Context, req *FindRoomByIDRequest, w http.ResponseWriter) (*modelsv1.Room, error) {
	return nil, nil
}

// FindRoomByUUIDRequest validates and create the request for the FindRoomByUUID operation.
func (h *Handler) FindRoomByUUIDRequest(r *http.Request, w http.ResponseWriter) (*FindRoomByIDRequest, error) {
	return nil, nil
}

// FindRoomByUUIDResponse writes out the response for the FindRoomByUUID operation.
func (h *Handler) FindRoomByUUIDResponse(ctx context.Context, res *modelsv1.Room, w http.ResponseWriter) error {
	return nil
}
