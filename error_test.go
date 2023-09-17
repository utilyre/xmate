package xmate

import (
	"errors"
	"net/http"
	"testing"
)

func TestNewHTTPError(t *testing.T) {
	httpErr := NewHTTPError(http.StatusNotFound, "user with email not found").(*HTTPError)
	if httpErr.Code != http.StatusNotFound {
		t.Errorf("httpErr.Code = %d; want %d", httpErr.Code, http.StatusNotFound)
	}
	if httpErr.Message != "user with email not found" {
		t.Errorf("httpErr.Message = '%s'; want 'user with email not found'", httpErr.Error())
	}

	httpErr = NewHTTPError(http.StatusInternalServerError).(*HTTPError)
	if httpErr.Code != http.StatusInternalServerError {
		t.Errorf("httpErr.Code = %d; want %d", httpErr.Code, http.StatusInternalServerError)
	}
	if httpErr.Message != "Internal Server Error" {
		t.Errorf("httpErr.Message = '%s'; want 'Internal Server Error'", httpErr.Error())
	}
}

func TestHTTPErrorIs(t *testing.T) {
	err1 := NewHTTPError(http.StatusForbidden, "you do not have the privileges to access this endpoint")
	err2 := NewHTTPError(http.StatusForbidden, "you are not allowed")

	if !errors.Is(err1, err2) {
		t.Error("err1 != err2; want err1 == err2")
	}
}
