package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing31d29bAnswersPong(t *testing.T) {
	srv := httptest.NewServer(newHandler())
	defer srv.Close()
	resp, err := http.Get(srv.URL + "/ping-31d29b")
	if err != nil {
		t.Fatalf("GET /ping-31d29b: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
	if ctype := resp.Header.Get("Content-Type"); ctype != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", ctype)
	}
	var body map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if len(body) != 1 || body["pong"] != "31d29b" {
		t.Errorf("body = %v, want {\"pong\": \"31d29b\"}", body)
	}
}

func TestPing31d29bRejectsPost(t *testing.T) {
	rec := httptest.NewRecorder()
	newHandler().ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/ping-31d29b", nil))
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("POST /ping-31d29b: status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
	if got := rec.Header().Get("Allow"); got != http.MethodGet {
		t.Errorf("POST /ping-31d29b: Allow = %q, want GET", got)
	}
}
