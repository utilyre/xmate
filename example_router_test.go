package xmate_test

import (
	"errors"
	"log"
	"net/http"

	"github.com/utilyre/xmate"
)

type serveMux = http.ServeMux

type Router struct {
	*serveMux
	handler xmate.ErrorHandler
}

func NewRouter(handler xmate.ErrorHandler) *Router {
	return &Router{
		serveMux: http.NewServeMux(),
		handler:  handler,
	}
}

func (r *Router) Handle(pattern string, handler xmate.Handler) {
	r.serveMux.Handle(pattern, r.handler.Handle(handler))
}

func (r *Router) HandleFunc(pattern string, handler xmate.HandlerFunc) {
	r.serveMux.HandleFunc(pattern, r.handler.HandleFunc(handler))
}

func Example() {
	r := NewRouter(handleError)
	r.HandleFunc("/", handleIndex)
	log.Fatal(http.ListenAndServe(":8080", r))
}

func handleError(w http.ResponseWriter, r *http.Request) {
	err := r.Context().Value(xmate.KeyError).(error)

	if httpErr := (xmate.HTTPError{}); errors.As(err, &httpErr) {
		_ = xmate.WriteText(w, httpErr.Code, httpErr.Message)
		return
	}

	log.Printf("%s %s failed: %v\n", r.Method, r.URL.Path, err)
	_ = xmate.WriteText(w, http.StatusInternalServerError, "Internal Server Error")
}

func handleIndex(w http.ResponseWriter, r *http.Request) error {
	return xmate.WriteText(w, http.StatusOK, "Hello world!")
}
