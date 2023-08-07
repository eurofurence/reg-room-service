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

const (
	// TransactionParseErrorMessage indicates a json body parse error
	TransactionParseErrorMessage APIErrorMessage = "transaction.parse.error"
	// field data failed to validate, see details field for more information
	TransactionDataInvalidMessage APIErrorMessage = "transaction.data.invalid"
	// duplicate referenceId
	TransactionDataDuplicateMessage APIErrorMessage = "transaction.data.duplicate"
	// database error
	TransactionWriteErrorMessage APIErrorMessage = "transaction.write.error"
	// database error
	TransactionReadErrorMessage APIErrorMessage = "transaction.read.error"
	// adapter failure while creating payment link
	TransactionPaylingErrorMessage APIErrorMessage = "transaction.paylink.error"
	// no such transaction in the database
	TransactionIDNotFoundMessage APIErrorMessage = "transaction.id.notfound"
	// syntactically invalid transaction id, must be positive integer
	TransactionIDInvalidMessage APIErrorMessage = "transaction.id.invalid"
	// deletion is not possible, e.g. because the grace period has expired for a valid payment
	TransactionCannotDeleteMessage APIErrorMessage = "transaction.cannot.delete"
	// token missing completely or invalid or expired
	AuthUnauthorizedMessage APIErrorMessage = "auth.unauthorized"
	// permissions missing
	AuthForbiddenMessage APIErrorMessage = "auth.forbidden"
	// Request could not be parsed properly
	RequestParseErrorMessage APIErrorMessage = "request.parse.failed"
	// Request created a conflict
	RequestConflictMessage APIErrorMessage = "request.conflict"
	// Internal error
	InternalErrorMessage APIErrorMessage = "http.error.internal"
	// Unknown error
	UnknownErrorMessage APIErrorMessage = "http.error.unkonwn"
)

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
