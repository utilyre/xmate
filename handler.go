package xmate

import (
	"context"
	"errors"
	"net/http"
	"sync/atomic"
)

var defaultErrorHandler atomic.Value

func init() {
	defaultErrorHandler.Store(ErrorHandler(func(w http.ResponseWriter, r *http.Request) {
		err := r.Context().Value(KeyError).(error)

		if httpErr := (&HTTPError{}); errors.As(err, &httpErr) {
			_ = WriteText(w, httpErr.Code, httpErr.Message)
			return
		}

		_ = WriteText(w, http.StatusInternalServerError, "Internal Server Error")
	}))
}

// Default returns the default error handler.
func Default() ErrorHandler {
	return defaultErrorHandler.Load().(ErrorHandler)
}

// SetDefault sets the default error handler, which is used by top level
// functions Handle and HandleFunc.
func SetDefault(eh ErrorHandler) {
	defaultErrorHandler.Store(eh)
}

// Handle calls Handle on the default error handler.
func Handle(next Handler) http.Handler {
	return Default().Handle(next)
}

// HandleFunc calls HandleFunc on the default error handler.
func HandleFunc(next HandlerFunc) http.HandlerFunc {
	return Default().HandleFunc(next)
}

// Key represents a package level context key.
type Key int

const (
	// KeyError is used to associate error values in a request context.
	KeyError Key = iota + 1
)

// Handler responds to an HTTP request.
type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request) error
}

// HandlerFunc is an adapter to allow the use of ordinary functions as HTTP
// handlers.
type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

// ServeHTTP calls f(w, r).
func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	return f(w, r)
}

// ErrorHandler responds to an HTTP request when an error occurs in the chain.
type ErrorHandler func(w http.ResponseWriter, r *http.Request)

// Handle converts a Handler to http.Handler by handling its error.
func (eh ErrorHandler) Handle(next Handler) http.Handler {
	return eh.handle(next)
}

// HandleFunc converts a HandlerFunc to http.HandlerFunc by handling its error.
func (eh ErrorHandler) HandleFunc(next HandlerFunc) http.HandlerFunc {
	return eh.handle(next)
}

func (eh ErrorHandler) handle(next Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := next.ServeHTTP(w, r); err != nil {
			r2 := r.WithContext(context.WithValue(r.Context(), KeyError, err))
			eh(w, r2)
		}
	}
}
