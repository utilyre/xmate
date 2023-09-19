package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/utilyre/xmate"
)

func main() {
	r := NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) error {
		return xmate.WriteText(w, http.StatusOK, "Hello World")
	})

	log.Fatal(http.ListenAndServe(":3000", r.mux))
}

type Router struct {
	mux          *http.ServeMux
	errorHandler xmate.ErrorHandler
}

func NewRouter() Router {
	return Router{
		mux: http.NewServeMux(),
		errorHandler: func(w http.ResponseWriter, r *http.Request) {
			err := r.Context().Value(xmate.ErrorKey{}).(error)

			httpErr := new(xmate.HTTPError)
			if !errors.As(err, &httpErr) {
				httpErr.Code = http.StatusInternalServerError
				httpErr.Message = http.StatusText(httpErr.Code)

				log.Printf("%s %s failed: %s\n", r.Method, r.URL.Path, err)
			}

			http.Error(w, httpErr.Message, httpErr.Code)
		},
	}
}

func (r Router) Handle(pattern string, handler xmate.Handler) {
	r.mux.Handle(pattern, r.errorHandler.Handle(handler))
}

func (r Router) HandleFunc(pattern string, handler xmate.HandlerFunc) {
	r.mux.HandleFunc(pattern, r.errorHandler.HandleFunc(handler))
}
