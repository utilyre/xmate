package xmate

import "net/http"

// An HTTPError represents a custom HTTP error.
type HTTPError struct {
	Code    int    // status code
	Message string // response message
}

// NewHTTPError returns a new HTTP error instance.
//
// If message is not provided, http.StatusText(code) will be used.
func NewHTTPError(code int, message ...string) error {
	if len(message) == 0 {
		return &HTTPError{
			Code:    code,
			Message: http.StatusText(code),
		}
	}

	return &HTTPError{
		Code:    code,
		Message: message[0],
	}
}

// Error returns e.Message.
func (e *HTTPError) Error() string {
	return e.Message
}

// Is checks whether e and target have the same code and message.
func (e *HTTPError) Is(target error) bool {
	httpErr, ok := target.(*HTTPError)
	return ok && (e.Code == httpErr.Code && e.Message == httpErr.Message)
}
