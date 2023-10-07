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
// Endpoint access only for admin users.
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
func (h *Handler) CreateRoomResponse(ctx context.Context, _ *modelsv1.Empty, w http.ResponseWriter) error {
	return nil
}

// AddRoomMemberRequest is the request type for the AddRoomMember operation.
type AddRoomMemberRequest struct {
}

// AddRoomMember adds an attendee to a room as an individual.
// The attendee must not be member of a group.
// Admin only.
func (h *Handler) AddRoomMember(ctx context.Context, req *AddRoomMemberRequest, w http.ResponseWriter) (*modelsv1.Empty, error) {
	return nil, nil
}

// AddRoomMemberRequest validates and creates a request for the AddRoomMember operation.
func (h *Handler) AddRoomMemberRequest(r *http.Request, w http.ResponseWriter) (*AddRoomMemberRequest, error) {
	return nil, nil
}

// Add AddRoomMemberResponse writes out the response for the AddRoomMember operation.
func (h *Handler) AddRoomMemberResponse(ctx context.Context, _ *modelsv1.Empty, w http.ResponseWriter) error {
	return nil
}

// AddGroupRequest is the request type for the AddGroup operation.
type AddGroupRequest struct{}

// AddGroup adds a group to a room.
// This locks the group against membership changes done by regular attendees. Admins may still change group membership.
// Admin only.
func (h *Handler) AddGroup(ctx context.Context, req *AddGroupRequest, w http.ResponseWriter) (*modelsv1.Empty, error) {
	return nil, nil
}

// AddGroupRequest validates and creates the request for the AddGroup operation.
func (h *Handler) AddGroupRequest(r *http.Request, w http.ResponseWriter) (*AddGroupRequest, error) {
	return nil, nil
}

// AddGroupResponse writes out the response for the AddGroup operation.
func (h *Handler) AddGroupResponse(ctx context.Context, _ *modelsv1.Empty, w http.ResponseWriter) error {
	return nil
}
