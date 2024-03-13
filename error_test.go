package xmate

import (
	"errors"
	"net/http"
	"testing"
)

func TestHTTPError_Is(t *testing.T) {
	var err1, err2 error

	err1 = Errorf(http.StatusConflict, "user already exists")
	err2 = Errorf(http.StatusConflict, "user already exists")
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
