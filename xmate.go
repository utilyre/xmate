package xmate

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func WriteText(w http.ResponseWriter, code int, body string) error {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	_, err := fmt.Fprintln(w, body)
	return err
}

func WriteHTML(w http.ResponseWriter, code int, body string) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	_, err := fmt.Fprint(w, body)
	return err
}

type Map map[string]any

func WriteJSON(w http.ResponseWriter, code int, body any) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	return json.NewEncoder(w).Encode(body)
}
