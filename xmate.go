// Package xmate provides missing convenient functionality for net/http.
package xmate

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

// WriteText writes body to w along with a proper header for text/plain.
func WriteText(w http.ResponseWriter, code int, body string) error {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	_, err := fmt.Fprintln(w, body)
	return err
}

// WriteHTML writes the template associated with t that has the given name to w
// along with a proper header for text/html.
func WriteHTML(w http.ResponseWriter, t *template.Template, code int, name string, data any) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	return t.ExecuteTemplate(w, name, data)
}

// A Map represents arbitrary JSON data.
type Map map[string]any

// WriteJSON writes body to w along with a proper header for application/json.
func WriteJSON(w http.ResponseWriter, code int, body any) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	return json.NewEncoder(w).Encode(body)
}
