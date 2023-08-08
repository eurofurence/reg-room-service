package common

import (
	"net/url"
	"time"
)

// ServiceError contains information
// which is required to let the application know which status code we want to send
// type ServiceError struct {
// 	Status int
// }

// APIErrorMessage type holds predefined error message constructs for the clients
type APIErrorMessage string

type serviceError struct {
	errorMessage APIErrorMessage
}

// ErrorFromMessage will construct a new error that can hold
// a predefined error message.
func ErrorFromMessage(message APIErrorMessage) error {
	return &serviceError{message}
}

// Error implements the `error` interface
func (s *serviceError) Error() string {
	return string(s.errorMessage)
}

// APIError is the generic return type for any Failure
// during endpoint operations
type APIError struct {
	RequestID string          `json:"requestid"`
	Message   APIErrorMessage `json:"message"`
	Timestamp int64           `json:"timestamp"`
	Details   url.Values      `json:"details"`
}

// NewAPIError creates a new instance of the `APIError` which will be returned
// to the client if an operation fails
func NewAPIError(reqID string, message APIErrorMessage, details url.Values) *APIError {
	return &APIError{
		RequestID: reqID,
		Message:   message,
		Timestamp: time.Now().Unix(),
		Details:   details,
	}
}
