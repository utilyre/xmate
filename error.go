package xmate

import "fmt"

// An HTTPError represents a custom HTTP error.
type HTTPError struct {
	Code    int    // status code
	Message string // response message
}

// Errorf formats according to a format specifier and returns the string along
// with code as an HTTPError.
func Errorf(code int, format string, a ...any) HTTPError {
	return HTTPError{
		Code:    code,
		Message: fmt.Sprintf(format, a...),
	}
}

// Error returns e.Message.
func (e HTTPError) Error() string {
	return e.Message
}

// Is checks whether e and target have the same code and message.
func (e HTTPError) Is(target error) bool {
	httpErr, ok := target.(HTTPError)
	return ok && e.Code == httpErr.Code && e.Message == httpErr.Message
}
