package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGreeting(t *testing.T) {
	for name, want := range map[string]string{"": "Hello, world!", "HOTL": "Hello, HOTL!"} {
		if got := Greeting(name); got != want {
			t.Errorf("Greeting(%q) = %q, want %q", name, got, want)
		}
	}
}

func getVersion(t *testing.T) map[string]string {
	t.Helper()
	rec := httptest.NewRecorder()
	newHandler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/version", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode body %q: %v", rec.Body.String(), err)
	}
	return body
}

func TestVersionDefaultsToDev(t *testing.T) {
	body := getVersion(t)
	if len(body) != 1 || body["version"] != "dev" {
		t.Errorf("body = %v, want {\"version\": \"dev\"}", body)
	}
}

func TestVersionUsesBuildVersion(t *testing.T) {
	old := version
	version = "1.2.3"
	t.Cleanup(func() { version = old })
	if got := getVersion(t)["version"]; got != "1.2.3" {
		t.Errorf("version = %q, want 1.2.3", got)
	}
}

func TestVersionOtherMethods(t *testing.T) {
	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch} {
		rec := httptest.NewRecorder()
		newHandler().ServeHTTP(rec, httptest.NewRequest(method, "/version", nil))
		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("%s /version status = %d, want %d", method, rec.Code, http.StatusMethodNotAllowed)
		}
	}
}
