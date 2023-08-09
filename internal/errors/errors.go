package errors

import (
	"errors"
	"fmt"
	"net/http"
)

var _ error = (*StatusError)(nil)

type KnownReason string

const (
	KnownReasonBadRequest          KnownReason = "BadRequest"
	KnownReasonUnauthorized        KnownReason = "Unauthorized"
	KnownReasonForbidden           KnownReason = "Forbidden"
	KnownReasonNotFound            KnownReason = "NotFound"
	KnownReasonConflict            KnownReason = "Conflict"
	KnownReasonInternalServerError KnownReason = "InternalServerError"
	KnownReasonUnknown             KnownReason = "Unknown"
)

type Status struct {
	Reason  KnownReason
	Code    int
	Message string
	Details string
}

type StatusError struct {
	ErrStatus Status
}

type APIStatus interface {
	Status() Status
}

func (se *StatusError) Error() string {
	return fmt.Sprintf("%s - %s", se.Status().Message, se.Status().Details)
}

func (se *StatusError) Status() Status {
	return se.ErrStatus
}

// NewBadRequest creates a new StatusError with error code 400
func NewBadRequest(message, details string) APIStatus {
	return &StatusError{
		ErrStatus: Status{
			Reason:  KnownReasonBadRequest,
			Code:    http.StatusBadRequest,
			Message: message,
			Details: details,
		},
	}
}

// NewUnauthorized creates a new StatusError with error code 401
func NewUnauthorized(message, details string) APIStatus {
	return &StatusError{
		ErrStatus: Status{
			Reason:  KnownReasonUnauthorized,
			Code:    http.StatusUnauthorized,
			Message: message,
			Details: details,
		},
	}
}

// NewForbidden creates a new StatusError with error code 403
func NewForbidden(message, details string) APIStatus {
	return &StatusError{
		ErrStatus: Status{
			Reason:  KnownReasonForbidden,
			Code:    http.StatusForbidden,
			Message: message,
			Details: details,
		},
	}
}

// NewNotFound creates a new StatusError with error code 404
func NewNotFound(message, details string) APIStatus {
	return &StatusError{
		ErrStatus: Status{
			Reason:  KnownReasonNotFound,
			Code:    http.StatusNotFound,
			Message: message,
			Details: details,
		},
	}
}

// NewConflict creates a new StatusError with error code 409
func NewConflict(message, details string) APIStatus {
	return &StatusError{
		ErrStatus: Status{
			Reason:  KnownReasonConflict,
			Code:    http.StatusConflict,
			Message: message,
			Details: details,
		},
	}
}

// NewInternalServerError creates a new StatusError with error code 500
func NewInternalServerError(message, details string) APIStatus {
	return &StatusError{
		ErrStatus: Status{
			Reason:  KnownReasonInternalServerError,
			Code:    http.StatusInternalServerError,
			Message: message,
			Details: details,
		},
	}
}

func isReasonOrCodeForError(expectedReason KnownReason, status int, err error) bool {
	errReason, code := reasonAndStatusCode(err)

	if errReason == expectedReason {
		return true
	}

	if code == status {
		return true
	}

	return false
}

// IsBadRequestError checks if error is of type `bad request`
func IsBadRequestError(err error) bool {
	return isReasonOrCodeForError(KnownReasonBadRequest, http.StatusBadRequest, err)
}

// IsUnauthorizedError checks if error is of type `unauthorized`
func IsUnauthorizedError(err error) bool {
	return isReasonOrCodeForError(KnownReasonUnauthorized, http.StatusUnauthorized, err)
}

// IsForbiddenError checks if error is of type `forbidden`
func IsForbiddenError(err error) bool {
	return isReasonOrCodeForError(KnownReasonForbidden, http.StatusForbidden, err)
}

// IsNotFoundError checks if error is of type `not found`
func IsNotFoundError(err error) bool {
	return isReasonOrCodeForError(KnownReasonNotFound, http.StatusNotFound, err)
}

// IsConflictError checks if error is of type `conflict`
func IsConflictError(err error) bool {
	return isReasonOrCodeForError(KnownReasonConflict, http.StatusConflict, err)
}

// IsInternalServerError checks if error is of type `internal server error`
func IsInternalServerError(err error) bool {
	return isReasonOrCodeForError(KnownReasonInternalServerError, http.StatusInternalServerError, err)
}

// IsUnknownError checks if error is of type `unknown`
func IsUnknownError(err error) bool {
	return isReasonOrCodeForError(KnownReasonUnknown, 0, err)
}

// AsAPIStatus checks if the error is of type `APIStatus`
// and returns nil if not
func AsAPIStatus(err error) APIStatus {
	if status, ok := err.(APIStatus); ok || errors.As(err, &status) {
		return status
	}

	return nil
}

func reasonAndStatusCode(err error) (KnownReason, int) {
	if status := AsAPIStatus(err); status != nil {
		return status.Status().Reason, status.Status().Code
	}

	return KnownReasonUnknown, 0
}
