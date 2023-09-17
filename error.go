package xmate

import "net/http"

// A HTTPError represents a custom HTTP error.
type HTTPError struct {
	Code    int
	Message string
}

// NewHTTPError returns a new HTTP error instance.
func NewHTTPError(code int, message ...string) error {
	msg := ""
	if len(message) > 0 {
		msg = message[0]
	} else {
		msg = http.StatusText(code)
	}

	return &HTTPError{
		Code:    code,
		Message: msg,
	}
}

// Error returns e.Message.
func (e *HTTPError) Error() string {
	return e.Message
}

// Is checks whether e and target have the same code.
func (e *HTTPError) Is(target error) bool {
	httpErr, ok := target.(*HTTPError)
	if !ok {
		return false
	}

	return e.Code == httpErr.Code
}
