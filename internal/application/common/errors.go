package common

import (
	"context"
	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
	"net/http"
	"net/url"
	"time"
)

// APIError allows lower layers of the service to provide detailed information about an error.
//
// While this breaks layer separation somewhat, it avoids having to map errors all over the place.
type APIError interface {
	error
	Status() int
	Response() modelsv1.Error
	InternalCauses() []error // not for sending to the client, but useful for logging
}

// ErrorMessageCode is a key to use for error messages in frontends or other automated systems interacting
// with our API. It avoids having to parse human-readable language for error classification beyond the
// http status.
type ErrorMessageCode string

// Note: the OpenAPI Spec is the leading document for error codes. This should directly follow the list and explanations
// in the Error schema.
const (
	DownstreamAttSrv ErrorMessageCode = "attendee.validation.error"     // attendee service downstream failure
	NoSuchAttendee   ErrorMessageCode = "attendee.notfound"             // no such attendee, probably invalid badge number or your user has no registration
	NotAttending     ErrorMessageCode = "attendee.status.not.attending" // attendee has a registration, but it is not in a status that allows being in a room, e.g. cancelled, waiting list

	AuthForbidden    ErrorMessageCode = "auth.forbidden"    // permissions missing or not a registered attendee
	AuthUnauthorized ErrorMessageCode = "auth.unauthorized" // token missing completely or invalid or expired

	GroupBanDuplicate      ErrorMessageCode = "group.ban.duplicate"       // an auto-decline entry with this badge number already exists - cannot add again
	GroupBanNotFound       ErrorMessageCode = "group.ban.notfound"        // an auto-decline entry with this badge number did not exist - removal failed
	GroupDataDuplicate     ErrorMessageCode = "group.data.duplicate"      // group with same name already exists, cannot create or rename
	GroupDataInvalid       ErrorMessageCode = "group.data.invalid"        // invalid field contents
	GroupIDInvalid         ErrorMessageCode = "group.id.invalid"          // invalid uuid id format
	GroupIDNotFound        ErrorMessageCode = "group.id.notfound"         // no such group id
	GroupInviteMismatch    ErrorMessageCode = "group.invite.mismatch"     // nickname or invitation code did not match - invite not sent or confirmed
	GroupMailError         ErrorMessageCode = "group.mail.error"          // mail service reported error when sending notification email - usually the operation will still have proceeded
	GroupMemberConflict    ErrorMessageCode = "group.member.conflict"     // attendee is already in or has been invited to another group
	GroupMemberDuplicate   ErrorMessageCode = "group.member.duplicate"    // attendee is already in or invited to this group
	GroupMemberNotFound    ErrorMessageCode = "group.member.notfound"     // attendee is not in or invited to this group
	GroupOwnerNotInGroup   ErrorMessageCode = "group.owner.notingroup"    // requested owner is not part of this group
	GroupOwnerCannotRemove ErrorMessageCode = "group.owner.cannot.remove" // this attendee is currently the owner of the group. Either change the owner first, or disband the group completely
	GroupReadError         ErrorMessageCode = "group.read.error"          // database error
	GroupSizeFull          ErrorMessageCode = "group.size.full"           // group has reached its maximum size
	GroupWriteError        ErrorMessageCode = "group.write.error"         // database error

	InternalErrorMessage ErrorMessageCode = "http.error.internal"  // Internal error
	RequestParseFailed   ErrorMessageCode = "request.parse.failed" // Request could not be parsed properly

	RoomDataDuplicate     ErrorMessageCode = "room.data.duplicate"     // room with same name already exists, cannot create or rename
	RoomDataInvalid       ErrorMessageCode = "room.data.invalid"       // invalid field contents
	RoomIDInvalid         ErrorMessageCode = "room.id.invalid"         // invalid uuid id format
	RoomIDNotFound        ErrorMessageCode = "room.id.notfound"        // no such room
	RoomOccupantConflict  ErrorMessageCode = "room.occupant.conflict"  // attendee is already in another room
	RoomOccupantDuplicate ErrorMessageCode = "room.occupant.duplicate" // attendee is already in this room
	RoomOccupantNotFound  ErrorMessageCode = "room.occupant.notfound"  // attendee is not in any/this room
	RoomNotEmpty          ErrorMessageCode = "room.not.empty"          // cannot delete a room that isn't empty
	RoomReadError         ErrorMessageCode = "room.read.error"         // database error
	RoomSizeFull          ErrorMessageCode = "room.size.full"          // not enough space in room to add another member
	RoomSizeTooSmall      ErrorMessageCode = "room.size.too.small"     // too many occupants in room to allow reducing size
	RoomWriteError        ErrorMessageCode = "room.write.error"        // database error
)

// construct specific API errors

func NewBadRequest(ctx context.Context, message ErrorMessageCode, details url.Values, internalCauses ...error) error {
	return NewAPIError(ctx, http.StatusBadRequest, message, details, internalCauses...)
}

func NewUnauthorized(ctx context.Context, message ErrorMessageCode, details url.Values, internalCauses ...error) error {
	return NewAPIError(ctx, http.StatusUnauthorized, message, details, internalCauses...)
}

func NewForbidden(ctx context.Context, message ErrorMessageCode, details url.Values, internalCauses ...error) error {
	return NewAPIError(ctx, http.StatusForbidden, message, details, internalCauses...)
}

func NewNotFound(ctx context.Context, message ErrorMessageCode, details url.Values, internalCauses ...error) error {
	return NewAPIError(ctx, http.StatusNotFound, message, details, internalCauses...)
}

func NewConflict(ctx context.Context, message ErrorMessageCode, details url.Values, internalCauses ...error) error {
	return NewAPIError(ctx, http.StatusConflict, message, details, internalCauses...)
}

func NewInternalServerError(ctx context.Context, message ErrorMessageCode, details url.Values, internalCauses ...error) error {
	return NewAPIError(ctx, http.StatusInternalServerError, message, details, internalCauses...)
}

func NewBadGateway(ctx context.Context, message ErrorMessageCode, details url.Values, internalCauses ...error) error {
	return NewAPIError(ctx, http.StatusBadGateway, message, details, internalCauses...)
}

// check for API errors

func IsBadRequestError(err error) bool {
	return isAPIErrorWithStatus(http.StatusBadRequest, err)
}

func IsUnauthorizedError(err error) bool {
	return isAPIErrorWithStatus(http.StatusUnauthorized, err)
}

func IsForbiddenError(err error) bool {
	return isAPIErrorWithStatus(http.StatusForbidden, err)
}

func IsNotFoundError(err error) bool {
	return isAPIErrorWithStatus(http.StatusNotFound, err)
}

func IsConflictError(err error) bool {
	return isAPIErrorWithStatus(http.StatusConflict, err)
}

func IsBadGatewayError(err error) bool {
	return isAPIErrorWithStatus(http.StatusBadGateway, err)
}

func IsInternalServerError(err error) bool {
	return isAPIErrorWithStatus(http.StatusInternalServerError, err)
}

func IsAPIError(err error) bool {
	_, ok := err.(APIError)
	return ok
}

const isoDateTimeFormat = "2006-01-02T15:04:05-07:00"

// NewAPIError creates a generic API error from directly provided information.
func NewAPIError(ctx context.Context, status int, message ErrorMessageCode, details url.Values, internalCauses ...error) APIError {

	return &StatusError{
		errStatus: status,
		response: modelsv1.Error{
			Timestamp: time.Now().Format(isoDateTimeFormat),
			Requestid: GetRequestID(ctx),
			Message:   string(message),
			Details:   details,
		},
		internalCauses: internalCauses,
	}
}

var _ error = (*StatusError)(nil)

type StatusError struct {
	errStatus      int
	response       modelsv1.Error
	internalCauses []error
}

func (se *StatusError) Error() string {
	return se.response.Message
}

func (se *StatusError) Status() int {
	return se.errStatus
}

func (se *StatusError) Response() modelsv1.Error {
	return se.response
}

func (se *StatusError) InternalCauses() []error {
	return se.internalCauses
}

func isAPIErrorWithStatus(status int, err error) bool {
	apiError, ok := err.(APIError)
	return ok && status == apiError.Status()
}

func Details(details string) url.Values {
	return url.Values{"details": []string{details}}
}
