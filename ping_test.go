package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPingAnswersPong(t *testing.T) {
	srv := httptest.NewServer(newHandler())
	defer srv.Close()
	resp, err := http.Get(srv.URL + "/ping")
	if err != nil {
		t.Fatalf("GET /ping: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
	if ctype := resp.Header.Get("Content-Type"); ctype != "text/plain" {
		t.Errorf("Content-Type = %q, want text/plain", ctype)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	if got := string(body); got != "pong" {
		t.Errorf("body = %q, want %q", got, "pong")
	}
}

func TestPingHandlerDirect(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	handlePing(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if ctype := rec.Header().Get("Content-Type"); ctype != "text/plain" {
		t.Errorf("Content-Type = %q, want text/plain", ctype)
	}
	if got := rec.Body.String(); got != "pong" {
		t.Errorf("body = %q, want %q", got, "pong")
	}
}

func TestPingRejectsPost(t *testing.T) {
	rec := httptest.NewRecorder()
	newHandler().ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/ping", nil))
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("POST /ping: status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
	if got := rec.Header().Get("Allow"); got != http.MethodGet {
		t.Errorf("POST /ping: Allow = %q, want GET", got)
	}
}
