# xmate

Package xmate provides missing convenient functionality for net/http.

## Usage

### Basic

In the simplest scenario, wrap the HTTP handlers and functional HTTP handlers
with `xmate.Handle` and `xmate.HandleFunc`, respectively.

> [!NOTE]
> These top-level functions convert the given HTTP handler into the standard
> `http.Handler` or `http.HandlerFunc` (depending on the function used) by
> handling its error.

```go
package main

import (
	"log"
	"net/http"

	"github.com/utilyre/xmate/v3"
)

func main() {
	mux := http.NewServeMux()

	sh := statusHandler{}

	mux.HandleFunc("GET /{$}", xmate.HandleFunc(handleHelloWorld))
	mux.Handle("GET /status", xmate.Handle(sh))

	log.Fatal(http.ListenAndServe(":8080", mux))
}

func handleHelloWorld(w http.ResponseWriter, r *http.Request) error {
	return xmate.WriteText(w, http.StatusOK, "hello world")
}

type statusHandler struct{}

func (sh statusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	return xmate.WriteJSON(w, http.StatusOK, map[string]any{
		"status": "healthy",
		"description": "service is ready to accept connections.",
	})
}
```

### In middleware

Instead of wrapping the middleware itself, wrap the function returned from it.

```go
package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/utilyre/xmate/v2"
)

func main() {
	r := chi.NewRouter()

	r.Use(AssertJSON)

	// ...

	log.Fatal(http.ListenAndServe(":8080", mux))
}

func AssertJSON(next http.Handler) http.Handler {
	return xmate.HandleFunc(func(w http.ResponseWriter, r *http.Request) error {
		if r.Header.Get("Content-Type") != "application/json" {
			return xmate.WriteText(w, http.StatusBadRequest, "Unsupported content type")
		}

		next.ServeHTTP(w, r)
		return nil
	})
}
```

### Multiple error handlers

TODO
