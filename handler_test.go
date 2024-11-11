package xmate_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/utilyre/xmate/v3"
)

func handleError(w http.ResponseWriter, r *http.Request, err error) {
	_ = xmate.WriteText(w, http.StatusInternalServerError, err.Error())
}

func handleEcho(w http.ResponseWriter, r *http.Request) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if len(body) == 0 {
		return xmate.WriteText(w, http.StatusBadRequest, "missing request body")
	}

	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(body)
	return err
}

func TestErrorHandler(t *testing.T) {
	eh := xmate.ErrorHandler(handleError)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	eh.HandleFunc(handleEcho).ServeHTTP(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("resp.StatusCode = %d; want %d", resp.StatusCode, http.StatusBadRequest)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != "missing request body\n" {
		t.Fatalf("body = '%s'; want 'missing request body\n'", body)
	}
}
