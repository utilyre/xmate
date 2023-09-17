package xmate

import (
	"context"
	"net/http"
)

type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request) error
}

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	return f(w, r)
}

type ErrorHandler func(w http.ResponseWriter, r *http.Request)

func (eh ErrorHandler) Handle(next Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := next.ServeHTTP(w, r); err != nil {
			r2 := r.WithContext(context.WithValue(r.Context(), "error", err))
			eh(w, r2)
		}
	})
}

func (eh ErrorHandler) HandleFunc(next HandlerFunc) http.Handler {
	return eh.Handle(next)
}
