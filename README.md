# xmate

xmate is a Go library designed to enhance the standard `net/http` package by 
providing convenient functionality specifically focused on error handling in 
HTTP handlers.

The name "xmate" signifies "friend of mux," highlighting its role in 
simplifying the creation of HTTP handlers. While it doesn't enhance routing in 
the standard library, xmate makes writing handlers less verbose, allowing 
developers to focus on building applications more efficiently and with less 
boilerplate code.

## Features

- **Less Boilerplate Code**: Streamline your HTTP handler development with
	reduced boilerplate, allowing you to write cleaner and more maintainable code.

- **Error Handling**: Convert HTTP handlers that return errors into standard
	`net/http` handlers, simplifying error management and ensuring consistent
	responses.

- **Multiple Error Handlers**: Easily define different error handlers for
	various types of HTTP endpoints, providing tailored responses based on the
	specific context of each request.

- **Ease of Use**: Enjoy a user-friendly experience that simplifies the process
	of creating and managing HTTP handlers, making it accessible for developers
	of all skill levels.

## Usage

### Handlers

In the simplest scenario, wrap HTTP handlers and functional HTTP handlers with
`xmate.Handle` and `xmate.HandleFunc`, respectively.

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

### Middlewares

Instead of wrapping the middleware itself, wrap the function returned from it
with `xmate.HandleFunc`.

```go
package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/utilyre/xmate/v3"
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

Error handling may need to vary in different contexts. That's why you can
always instantiate a non-default error handler and use it on your HTTP handlers.

```go
package main

import (
	"log"
	"net/http"

	"github.com/utilyre/xmate/v3"
)

func main() {
	mux := http.NewServeMux()

	xmate.SetDefault(handleError) // changes the error handler used by top-level functions
	apiEH := xmate.ErrorHandler(handleAPIError) // custom handler for API routes

	mux.HandleFunc("GET /signup", xmate.HandleFunc(handleSignUpPage))
	mux.HandleFunc("GET /login", xmate.HandleFunc(handleLoginPage))

	mux.Handle("GET /api/v1/auth/signup", apiEH.HandleFunc(handleSignUp))
	mux.Handle("GET /api/v1/auth/login", apiEH.HandleFunc(handleLogin))

	log.Fatal(http.ListenAndServe(":8080", mux))
}

func handleError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Failed to %s %s: %v", r.Method, r.URL.Path, err)
	_ = xmate.WriteText(w, http.StatusInternalServerError, "Internal Server Error")
}

func handleAPIError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Failed to %s %s: %v", r.Method, r.URL.Path, err)
	_ = xmate.WriteJSON(w, http.StatusInternalServerError, map[string]any{
		"message": "Internal Server Error",
	})
}
```
