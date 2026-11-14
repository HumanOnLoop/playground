package main

import (
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// Criterion 1: a running server answers GET /version with 200 and a JSON body {"version": "<version>"}.
func TestVersion_GetOverHTTPReturnsJSONVersion(t *testing.T) {
	srv := httptest.NewServer(newHandler())
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/version")
	if err != nil {
		t.Fatalf("GET /version: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	var body map[string]any
	if err := json.Unmarshal(raw, &body); err != nil {
		t.Fatalf("body %q is not JSON: %v", raw, err)
	}
	if len(body) != 1 {
		t.Errorf("body = %v, want exactly one key \"version\"", body)
	}
	if v, ok := body["version"].(string); !ok || v != version {
		t.Errorf("body[\"version\"] = %#v, want string %q", body["version"], version)
	}
}

// Criterion 3: a running server answers POST /version with 405.
func TestVersion_PostOverHTTPIsNotAllowed(t *testing.T) {
	srv := httptest.NewServer(newHandler())
	defer srv.Close()

	resp, err := http.Post(srv.URL+"/version", "application/json", strings.NewReader("{}"))
	if err != nil {
		t.Fatalf("POST /version: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want 405", resp.StatusCode)
	}
}

// Criteria 1 and 2 end to end: build the real binary, with and without -ldflags "-X main.version=...", run it and
// query GET /version.
func TestVersion_BuiltBinary(t *testing.T) {
	if testing.Short() {
		t.Skip("builds and runs the binary")
	}
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skip("go toolchain not on PATH")
	}
	for _, tc := range []struct {
		name, ldflags, want string
	}{
		{"WithoutVersionIsDev", "", "dev"},
		{"WithLdflagsVersion", "-X main.version=9.8.7-test", "9.8.7-test"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			bin := filepath.Join(t.TempDir(), "playground")
			args := []string{"build", "-o", bin}
			if tc.ldflags != "" {
				args = append(args, "-ldflags", tc.ldflags)
			}
			args = append(args, ".")
			build := exec.Command(goBin, args...)
			build.Env = os.Environ()
			if out, err := build.CombinedOutput(); err != nil {
				t.Fatalf("go build: %v\n%s", err, out)
			}

			addr := freeAddr(t)
			cmd := exec.Command(bin, "-addr", addr)
			if err := cmd.Start(); err != nil {
				t.Fatalf("start binary: %v", err)
			}
			t.Cleanup(func() { _ = cmd.Process.Kill(); _ = cmd.Wait() })

			var resp *http.Response
			deadline := time.Now().Add(10 * time.Second)
			for {
				resp, err = http.Get("http://" + addr + "/version")
				if err == nil {
					break
				}
				if time.Now().After(deadline) {
					t.Fatalf("server did not come up: %v", err)
				}
				time.Sleep(50 * time.Millisecond)
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				t.Fatalf("status = %d, want 200", resp.StatusCode)
			}
			var body map[string]string
			if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if body["version"] != tc.want {
				t.Errorf("version = %q, want %q", body["version"], tc.want)
			}
		})
	}
}

func freeAddr(t *testing.T) string {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	addr := l.Addr().String()
	l.Close()
	return addr
}
