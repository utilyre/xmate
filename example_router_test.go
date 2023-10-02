package xmate_test

import (
	"errors"
	"log"
	"net/http"

	"github.com/utilyre/xmate"
)

type Router struct {
	mux *http.ServeMux
	eh  xmate.ErrorHandler
}

func NewRouter(eh xmate.ErrorHandler) *Router {
	return &Router{
		mux: http.NewServeMux(),
		eh:  eh,
	}
}

func (r *Router) Handle(pattern string, handler xmate.Handler) {
	r.mux.Handle(pattern, r.eh.Handle(handler))
}

func (r *Router) HandleFunc(pattern string, handler xmate.HandlerFunc) {
	r.mux.HandleFunc(pattern, r.eh.HandleFunc(handler))
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

func Example() {
	r := NewRouter(handleError)
	r.HandleFunc("/", handleIndex)
	log.Fatal(http.ListenAndServe(":8080", r))
}

func handleError(w http.ResponseWriter, r *http.Request) {
	err := r.Context().Value(xmate.ErrorKey{}).(error)

	httpErr := new(xmate.HTTPError)
	if !errors.As(err, &httpErr) {
		httpErr.Code = http.StatusInternalServerError
		httpErr.Message = "Internal Server Error"

		log.Printf("%s %s failed: %s\n", r.Method, r.URL.Path, err)
	}

	if err := xmate.WriteText(w, httpErr.Code, httpErr.Message); err != nil {
		log.Printf("%s %s failed to write error response: %s\n", r.Method, r.URL.Path, err)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) error {
	return xmate.WriteText(w, http.StatusOK, "Hello world!")
}
