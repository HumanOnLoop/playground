package main

import (
	"encoding/json"
	"io"
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

// getGreet9d0d00 sends GET /greet-9d0d00<query> to a running test server and returns the status, Content-Type
// and raw body.
func getGreet9d0d00(t *testing.T, base, query string) (int, string, string) {
	t.Helper()
	resp, err := http.Get(base + "/greet-9d0d00" + query)
	if err != nil {
		t.Fatalf("GET /greet-9d0d00%s: %v", query, err)
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	return resp.StatusCode, resp.Header.Get("Content-Type"), string(raw)
}

func assertGreetJSON(t *testing.T, code int, ctype, raw, want string) {
	t.Helper()
	if code != http.StatusOK {
		t.Errorf("status = %d, want 200", code)
	}
	if ctype != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", ctype)
	}
	var body map[string]string
	if err := json.Unmarshal([]byte(raw), &body); err != nil {
		t.Fatalf("body %q is not a JSON object of strings: %v", raw, err)
	}
	if len(body) != 1 || body["text"] != want {
		t.Errorf("body = %s, want {\"text\": %q}", raw, want)
	}
}

// Criterion 1: GET /greet-9d0d00?lang=en answers 200 {"text": "Hello"}.
func TestGreet9d0d00_LangEnAnswersHello(t *testing.T) {
	srv := httptest.NewServer(newHandler())
	defer srv.Close()
	code, ctype, raw := getGreet9d0d00(t, srv.URL, "?lang=en")
	assertGreetJSON(t, code, ctype, raw, "Hello")
}

// Criterion 2: GET /greet-9d0d00?lang=sv answers 200 {"text": "Hej"}.
func TestGreet9d0d00_LangSvAnswersHej(t *testing.T) {
	srv := httptest.NewServer(newHandler())
	defer srv.Close()
	code, ctype, raw := getGreet9d0d00(t, srv.URL, "?lang=sv")
	assertGreetJSON(t, code, ctype, raw, "Hej")
}

// Criterion 3: GET /greet-9d0d00 without lang answers exactly as for lang=en.
func TestGreet9d0d00_NoLangAnswersAsEnglish(t *testing.T) {
	srv := httptest.NewServer(newHandler())
	defer srv.Close()
	enCode, enType, enRaw := getGreet9d0d00(t, srv.URL, "?lang=en")
	code, ctype, raw := getGreet9d0d00(t, srv.URL, "")
	assertGreetJSON(t, code, ctype, raw, "Hello")
	if code != enCode || ctype != enType || raw != enRaw {
		t.Errorf("without lang: %d %q %q; with lang=en: %d %q %q", code, ctype, raw, enCode, enType, enRaw)
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
