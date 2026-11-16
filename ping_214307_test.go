package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing214307AnswersOK(t *testing.T) {
	srv := httptest.NewServer(newHandler())
	defer srv.Close()
	resp, err := http.Get(srv.URL + "/ping-214307")
	if err != nil {
		t.Fatalf("GET /ping-214307: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
	if ctype := resp.Header.Get("Content-Type"); ctype != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", ctype)
	}
	var body map[string]bool
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if len(body) != 1 || !body["ok"] {
		t.Errorf("body = %v, want {\"ok\": true}", body)
	}
}

func TestPing214307RejectsPost(t *testing.T) {
	rec := httptest.NewRecorder()
	newHandler().ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/ping-214307", nil))
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("POST /ping-214307: status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
	if got := rec.Header().Get("Allow"); got != http.MethodGet {
		t.Errorf("POST /ping-214307: Allow = %q, want GET", got)
	}
}
