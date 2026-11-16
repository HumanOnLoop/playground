package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGreet9d0d00AnswersInLanguage(t *testing.T) {
	srv := httptest.NewServer(newHandler())
	defer srv.Close()
	tests := []struct {
		name  string
		query string
		want  string
	}{
		{"english", "?lang=en", "Hello"},
		{"swedish", "?lang=sv", "Hej"},
		{"no lang", "", "Hello"},
		{"empty lang", "?lang=", "Hello"},
		{"unknown lang", "?lang=fr", "Hello"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.Get(srv.URL + "/greet-9d0d00" + tt.query)
			if err != nil {
				t.Fatalf("GET /greet-9d0d00%s: %v", tt.query, err)
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
			if len(body) != 1 || body["text"] != tt.want {
				t.Errorf("body = %v, want {\"text\": %q}", body, tt.want)
			}
		})
	}
}

func TestGreet9d0d00RejectsPost(t *testing.T) {
	rec := httptest.NewRecorder()
	newHandler().ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/greet-9d0d00", nil))
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("POST /greet-9d0d00: status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
	if got := rec.Header().Get("Allow"); got != http.MethodGet {
		t.Errorf("POST /greet-9d0d00: Allow = %q, want GET", got)
	}
}
