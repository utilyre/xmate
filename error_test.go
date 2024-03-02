package xmate

import (
	"errors"
	"net/http"
	"testing"
)

func TestErrorf(t *testing.T) {
	httpErr := Errorf(http.StatusNotFound, "user with email not found").(HTTPError)
	if httpErr.Code != http.StatusNotFound {
		t.Errorf("httpErr.Code = %d; want %d", httpErr.Code, http.StatusNotFound)
	}
	if httpErr.Message != "user with email not found" {
		t.Errorf("httpErr.Message = '%s'; want 'user with email not found'", httpErr.Message)
	}
}

func TestHTTPError_Is(t *testing.T) {
	err1 := Errorf(http.StatusConflict, "user already exists")
	err2 := Errorf(http.StatusConflict, "user already exists")
	if !errors.Is(err1, err2) {
		t.Errorf("%#[1]v != %#[2]v; want %#[1]v = %#[2]v", err1, err2)
	}

	err1 = Errorf(http.StatusForbidden, "you do not have the privileges to access this endpoint")
	err2 = Errorf(http.StatusForbidden, "you are not allowed")
	if errors.Is(err1, err2) {
		t.Errorf("%#[1]v = %#[2]v; want %#[1]v != %#[2]v", err1, err2)
	}

	err1 = Errorf(http.StatusUnauthorized, "missing or malformed JWT")
	err2 = Errorf(http.StatusBadRequest, "missing or malformed JWT")
	if errors.Is(err1, err2) {
		t.Errorf("%#[1]v = %#[2]v; want %#[1]v != %#[2]v", err1, err2)
	}

	err1 = Errorf(http.StatusInternalServerError, "something went wrong")
	err2 = Errorf(http.StatusNotFound, "page not found")
	if errors.Is(err1, err2) {
		t.Errorf("%#[1]v = %#[2]v; want %#[1]v != %#[2]v", err1, err2)
	}

	err1 = Errorf(http.StatusMethodNotAllowed, "method POST is not allowed")
	err2 = errors.New("method POST is not allowed")
	if errors.Is(err1, err2) {
		t.Errorf("%#[1]v = %#[2]v; want %#[1]v != %#[2]v", err1, err2)
	}
}
