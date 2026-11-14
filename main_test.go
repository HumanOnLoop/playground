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

// getVersion sends GET /version to a test server and returns the status, Content-Type and decoded body.
func getVersion(t *testing.T) (int, string, map[string]string) {
	t.Helper()
	srv := httptest.NewServer(newHandler())
	defer srv.Close()
	resp, err := http.Get(srv.URL + "/version")
	if err != nil {
		t.Fatalf("GET /version: %v", err)
	}
	defer resp.Body.Close()
	var body map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	return resp.StatusCode, resp.Header.Get("Content-Type"), body
}

func TestVersionDefaultsToDev(t *testing.T) {
	code, ctype, body := getVersion(t)
	if code != http.StatusOK {
		t.Errorf("status = %d, want %d", code, http.StatusOK)
	}
	if ctype != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", ctype)
	}
	if len(body) != 1 || body["version"] != "dev" {
		t.Errorf("body = %v, want {\"version\": \"dev\"}", body)
	}
}

func TestVersionReportsBuildVersion(t *testing.T) {
	old := version
	version = "1.2.3"
	t.Cleanup(func() { version = old })
	code, _, body := getVersion(t)
	if code != http.StatusOK || body["version"] != "1.2.3" {
		t.Errorf("got %d %v, want 200 {\"version\": \"1.2.3\"}", code, body)
	}
}

func TestVersionRejectsOtherMethods(t *testing.T) {
	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodDelete} {
		rec := httptest.NewRecorder()
		newHandler().ServeHTTP(rec, httptest.NewRequest(method, "/version", nil))
		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("%s /version: status = %d, want %d", method, rec.Code, http.StatusMethodNotAllowed)
		}
		if got := rec.Header().Get("Allow"); got != http.MethodGet {
			t.Errorf("%s /version: Allow = %q, want GET", method, got)
		}
	}
}
