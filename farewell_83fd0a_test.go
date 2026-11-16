package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFarewell83fd0aAnswersInLanguage(t *testing.T) {
	srv := httptest.NewServer(newHandler())
	defer srv.Close()
	tests := []struct {
		name  string
		query string
		want  string
	}{
		{"english", "?lang=en", "Goodbye"},
		{"swedish", "?lang=sv", "Hej då"},
		{"no lang", "", "Goodbye"},
		{"empty lang", "?lang=", "Goodbye"},
		{"unknown lang", "?lang=fr", "Goodbye"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.Get(srv.URL + "/farewell-83fd0a" + tt.query)
			if err != nil {
				t.Fatalf("GET /farewell-83fd0a%s: %v", tt.query, err)
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

func TestFarewell83fd0aRejectsPost(t *testing.T) {
	rec := httptest.NewRecorder()
	newHandler().ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/farewell-83fd0a", nil))
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("POST /farewell-83fd0a: status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
	if got := rec.Header().Get("Allow"); got != http.MethodGet {
		t.Errorf("POST /farewell-83fd0a: Allow = %q, want GET", got)
	}
}
