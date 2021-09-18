package service_errors

import (
	"encoding/json"
)

// ErrorType is the list of allowed values for the error's type.
type ErrorType string

// ErrorCode is the list of allowed values for the error's code.
type ErrorCode string

// List of values that ErrorType can take.
const (
	ErrorTypeAPI            ErrorType = "api_error"
	ErrorTypeAPIConnection  ErrorType = "api_connection_error"
	ErrorTypeAuthentication ErrorType = "authentication_error"
	ErrorTypeDistributedTx  ErrorType = "distributed_tx_error"
	ErrorTypeInvalidRequest ErrorType = "invalid_request_error"
	ErrorTypePermission     ErrorType = "more_permissions_required"
	ErrorTypeRateLimit      ErrorType = "rate_limit_error"
)

// ServiceError is the response returned when a call is unsuccessful.
type ServiceError struct {
	Code ErrorCode `json:"code,omitempty"`

	// Err contains an internal error with an additional level of granularity
	// that can be used in some cases to get more detailed information about
	// what went wrong. For example, Err may hold a TxError that indicates
	// exactly what went wrong during distributed tx.
	Err error `json:"-"`

	HTTPStatusCode int       `json:"status,omitempty"`
	Msg            string    `json:"message"`
	Param          string    `json:"param,omitempty"`
	RequestID      string    `json:"request_id,omitempty"`
	Type           ErrorType `json:"type"`
}

// NewServiceError creates a new instance of the service error
func NewServiceError(code ErrorCode, err error, statusCode int, msg, param, requestId string, errType ErrorType) ServiceError {
	return ServiceError{
		Code:           code,
		Err:            err,
		HTTPStatusCode: statusCode,
		Msg:            msg,
		Param:          param,
		RequestID:      requestId,
		Type:           errType,
	}
}

// Error serializes the error object to JSON and returns it as a string.
func (e *ServiceError) Error() string {
	ret, _ := json.Marshal(e)
	return string(ret)
}

// APIConnectionError is a failure to connect to the Merchant Account Service API.
type APIConnectionError struct {
	err *ServiceError
}

// Error serializes the error object to JSON and returns it as a string.
func (e *APIConnectionError) Error() string {
	return e.err.Error()
}

// APIError is a catch all for any errors not covered by other types (and
// should be extremely uncommon).
type APIError struct {
	err *ServiceError
}

// Error serializes the error object to JSON and returns it as a string.
func (e *APIError) Error() string {
	return e.err.Error()
}

// AuthenticationError is a failure to properly authenticate during a request.
type AuthenticationError struct {
	err *ServiceError
}

// Error serializes the error object to JSON and returns it as a string.
func (e *AuthenticationError) Error() string {
	return e.err.Error()
}

// PermissionError results when you attempt to make an API request
// for which your account doesn't have the right permissions.
type PermissionError struct {
	err *ServiceError
}

// Error serializes the error object to JSON and returns it as a string.
func (e *PermissionError) Error() string {
	return e.err.Error()
}

// DistributedTxError results when an error occurs when performing a distributed tx.
type DistributedTxError struct {
	err *ServiceError
}

// Error serializes the error object to JSON and returns it as a string.
func (e *DistributedTxError) Error() string {
	return e.err.Error()
}

// InvalidRequestError results when an error occurs due to an invalid request being performed against this service
type InvalidRequestError struct {
	err *ServiceError
}

// Error serializes the error object to JSON and returns it as a string.
func (e *InvalidRequestError) Error() string {
	return e.err.Error()
}

// RatelimitError occurs in the face of errors due to rate limiting
type RatelimitError struct {
	err *ServiceError
}

// Error serializes the error object to JSON and returns it as a string.
func (e *RatelimitError) Error() string {
	return e.err.Error()
}
