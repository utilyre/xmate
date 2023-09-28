# xmate

Package xmate provides missing convenient functionality for net/http.

## Usage

You'll probably want to make a custom router type which takes care of wrapping
your handlers with type `xmate.ErrorHandler`, for which see the
[examples][examples] folder.

However here is a basic example

```go
package main

import (
	"errors"
	"net/http"
	"log"

	"github.com/utilyre/xmate"
)

func main() {
	mux := http.NewServeMux()
	eh := xmate.ErrorHandler(func(w http.ResponseWriter, r *http.Request) {
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
	})

	mux.HandleFunc("/", eh.HandleFunc(func(w http.ResponseWriter, r *http.Request) error {
		return xmate.WriteText(w, http.StatusOK, "Hello world!")
	}))

	http.ListenAndServe(":3000", mux)
}
```

[examples]: https://github.com/utilyre/xmate/tree/main/examples
