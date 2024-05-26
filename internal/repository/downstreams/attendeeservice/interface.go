package attendeeservice

import (
	"context"
)

type Status string

var (
	StatusNew           Status = "new"
	StatusApproved      Status = "approved"
	StatusPartiallyPaid Status = "partially paid"
	StatusPaid          Status = "paid"
	StatusCheckedIn     Status = "checked in"
	StatusCancelled     Status = "cancelled"
	StatusWaiting       Status = "waiting"
	StatusDeleted       Status = "deleted"
)

type AttendeeService interface {
	// ListMyRegistrationIds which attendee ids belong to the current user?
	//
	// If your request was made with an api token, this will fail and should not be called.
	//
	// Admins are treated like normal users for this request, and will also only receive badge numbers
	// they have personally registered.
	//
	// Forwards the jwt from the request.
	ListMyRegistrationIds(ctx context.Context) ([]int64, error)

	// GetStatus obtains the status for a given attendee id.
	//
	// Nonexistent registrations will return StatusDeleted because the distinction isn't important.
	//
	// If your request was made by an admin, you can read everyone's status. A user can only read their own status.
	//
	// Forwards the jwt from the request.
	GetStatus(ctx context.Context, id int64) (Status, error)
}
