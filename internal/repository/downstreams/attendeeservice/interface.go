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

type Attendee struct {
	ID       int64  `json:"id"`       // badge number
	Nickname string `json:"nickname"` // fan name

	Email string `json:"email"`

	SpokenLanguages      string `json:"spoken_languages"`      // configurable subset of configured language codes, comma separated (de,en)
	RegistrationLanguage string `json:"registration_language"` // one out of configurable subset of RFC 5646 locales (default en-US)

	// comma separated lists, allowed choices are convention dependent
	Flags    string `json:"flags"`    // hc,anon,ev
	Packages string `json:"packages"` // room-none,attendance,stage,sponsor,sponsor2
	Options  string `json:"options"`  // art,anim,music,suit
}

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

	// GetAttendee obtains part of the registration information for given attendee id.
	//
	// Used for internal nickname lookups, etc.
	//
	// Uses the api token for full access, so access control must be performed in the implementation.
	GetAttendee(ctx context.Context, id int64) (Attendee, error)
}
