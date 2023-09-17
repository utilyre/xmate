package xmate

import "net/http"

type HTTPError struct {
	Code    int
	Message string
}

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

func (e *HTTPError) Error() string {
	return e.Message
}

func (e *HTTPError) Is(target error) bool {
	httpErr, ok := target.(*HTTPError)
	if !ok {
		return false
	}

	return e.Code == httpErr.Code
}
