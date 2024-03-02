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

	log.Fatal(http.ListenAndServe(":8080", mux))
}

func handleError(w http.ResponseWriter, r *http.Request) {
	err := r.Context().Value(xmate.KeyError).(error)

	if httpErr := new(xmate.HTTPError); errors.As(err, &httpErr) {
		_ = xmate.WriteText(w, httpErr.Code, httpErr.Message)
		return
	}

	log.Printf("%s %s failed: %v\n", r.Method, r.URL.Path, err)
	_ = xmate.WriteText(w, http.StatusInternalServerError, "Internal Server Error")
}

func handleIndex(w http.ResponseWriter, r *http.Request) error {
	return xmate.WriteText(w, http.StatusOK, "Hello world!")
}
```
