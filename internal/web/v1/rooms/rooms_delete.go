package rooms

import (
	"context"
	"net/http"

	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
)

// DeleteRoomRequest is the request type for the DeleteRoom operation.
type DeleteRoomRequest struct {
	UUID string
}

// DeleteRoom deletes an existing room by uuid.
//
// IMPORTANT: once an attendee has been billed for this room, this is a dangerous operation, as it may
// deprive them of a room reservation that you have confirmed! For this reason, you can only
// delete empty rooms.
//
// Admin only.
func (h *Handler) DeleteRoom(ctx context.Context, req *DeleteRoomRequest, w http.ResponseWriter) (*modelsv1.Empty, error) {
	return nil, nil
}

// DeleteRoomRequest validates and creates the request for the DeleteRoom operation.
func (h *Handler) DeleteRoomRequest(r *http.Request, w http.ResponseWriter) (*DeleteRoomRequest, error) {
	return nil, nil
}

// DeleteRoomResponse writes out the response for the DeleteRoom operation.
func (h *Handler) DeleteRoomResponse(ctx context.Context, _ *modelsv1.Empty, w http.ResponseWriter) error {
	return nil
}

// DeleteRoomMemberRequest is the request type for the DeleteRoomMember operation.
type DeleteRoomMemberRequest struct {
}

// DeleteRoomMember Removes the attendee with the given badge number from the room as an individual.
// You cannot change groups this way.
// Admin only.
func (h *Handler) DeleteRoomMember(ctx context.Context, req *DeleteRoomMemberRequest, w http.ResponseWriter) (*modelsv1.Empty, error) {
	return nil, nil
}

// DeleteRoomMemberRequest validates and creates a request for the DeleteRoomMember operation.
func (h *Handler) DeleteRoomMemberRequest(r *http.Request, w http.ResponseWriter) (*DeleteRoomMemberRequest, error) {
	return nil, nil
}

// Delete DeleteRoomMemberResponse writes out the response for the DeleteRoomMember operation.
func (h *Handler) DeleteRoomMemberResponse(ctx context.Context, _ *modelsv1.Empty, w http.ResponseWriter) error {
	return nil
}

// DeleteGroupRequest is the request type for the DeleteGroup operation.
type DeleteGroupRequest struct{}

// DeleteGroup removes the group from the room.
// Admin only.
func (h *Handler) DeleteGroup(ctx context.Context, req *DeleteGroupRequest, w http.ResponseWriter) (*modelsv1.Empty, error) {
	return nil, nil
}

// DeleteGroupRequest validates and creates the request for the DeleteGroup operation.
func (h *Handler) DeleteGroupRequest(r *http.Request, w http.ResponseWriter) (*DeleteGroupRequest, error) {
	return nil, nil
}

// DeleteGroupResponse writes out the response for the DeleteGroup operation.
func (h *Handler) DeleteGroupResponse(ctx context.Context, _ *modelsv1.Empty, w http.ResponseWriter) error {
	return nil
}
