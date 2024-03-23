package common

import (
	"net/url"
	"time"
)

const (
	AuthUnauthorizedMessage  string = "auth.unauthorized"    // token missing completely or invalid or expired
	AuthForbiddenMessage     string = "auth.forbidden"       // permissions missing
	RequestParseErrorMessage string = "request.parse.failed" // Request could not be parsed properly
	InternalErrorMessage     string = "http.error.internal"  // Internal error
	UnknownErrorMessage      string = "http.error.unknown"   // Unknown error

	GroupIDInvalidMessage   string = "group.id.invalid"
	GroupDataInvalidMessage string = "group.data.invalid"
	GroupIDNotFoundMessage  string = "group.id.notfound"

	GroupMemberNotFound string = "group.member.notfound"
)

// ServiceError contains information
// which is required to let the application know which status code we want to send
// type ServiceError struct {
// 	Status int
// }

type serviceError struct {
	errorMessage string
}

// ErrorFromMessage will construct a new error that can hold
// a predefined error message.
func ErrorFromMessage(message string) error {
	return &serviceError{message}
}

// Error implements the `error` interface.
func (s *serviceError) Error() string {
	return string(s.errorMessage)
}

// APIError is the generic return type for any Failure
// during endpoint operations.
type APIError struct {
	RequestID string     `json:"requestid"`
	Message   string     `json:"message"`
	Timestamp string     `json:"timestamp"`
	Details   url.Values `json:"details"`
}

// NewAPIError creates a new instance of the `APIError` which will be returned
// to the client if an operation fails.
func NewAPIError(reqID string, message string, details url.Values) *APIError {
	return &APIError{
		RequestID: reqID,
		Message:   message,
		Timestamp: time.Now().Format(time.RFC3339),
		Details:   details,
	}
}
