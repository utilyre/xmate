package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/utilyre/xmate"
)

func main() {
	r := NewRouter()

	r.Use(func(next xmate.Handler) xmate.Handler {
		return xmate.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
			w.Header().Set("X-Powered-By", "xmate")
			return next.ServeHTTP(w, r)
		})
	})

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) error {
		return xmate.WriteText(w, http.StatusOK, "Hello World")
	})

	log.Fatal(http.ListenAndServe(":3000", r.router))
}

type MiddlewareFunc func(next xmate.Handler) xmate.Handler

type Router struct {
	router       *mux.Router
	errorHandler xmate.ErrorHandler
}

func NewRouter() Router {
	return Router{
		router: mux.NewRouter(),
		errorHandler: func(w http.ResponseWriter, r *http.Request) {
			err := r.Context().Value(xmate.ErrorKey{}).(error)

			httpErr := new(xmate.HTTPError)
			if !errors.As(err, &httpErr) {
				httpErr.Code = http.StatusInternalServerError
				httpErr.Message = "Internal Server Error"

				log.Printf("%s %s failed: %s\n", r.Method, r.URL.Path, err)
			}

			xmate.WriteText(w, httpErr.Code, httpErr.Message)
		},
	}
}

func (r Router) Subrouter(prefix string) Router {
	return Router{
		router:       r.router.PathPrefix(prefix).Subrouter(),
		errorHandler: r.errorHandler,
	}
}

func (r Router) Use(mwf MiddlewareFunc) {
	r.router.Use(func(next http.Handler) http.Handler {
		xNext := xmate.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
			next.ServeHTTP(w, r)
			return nil
		})

		return r.errorHandler.Handle(mwf(xNext))
	})
}

func (r Router) Handle(path string, handler xmate.Handler) *mux.Route {
	return r.router.Handle(path, r.errorHandler.Handle(handler))
}

func (r Router) HandleFunc(path string, handler xmate.HandlerFunc) *mux.Route {
	return r.router.HandleFunc(path, r.errorHandler.HandleFunc(handler))
}
