# xmate

Package xmate provides missing convenient functionality for net/http.

## Usage

Here is a basic example

```go
package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/utilyre/xmate"
)

func main() {
	mux := http.NewServeMux()
	eh := xmate.ErrorHandler(handleError)

	mux.HandleFunc("/", eh.HandleFunc(handleIndex))

	http.ListenAndServe(":8080", mux)
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
```
