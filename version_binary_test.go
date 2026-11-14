package main

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

// buildAndServe builds the playground binary with the given extra go build args, starts it on a free port and
// returns the server's base URL. The process is stopped when the test ends.
func buildAndServe(t *testing.T, buildArgs ...string) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "playground")
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}
	goTool := filepath.Join(runtime.GOROOT(), "bin", "go")
	args := append([]string{"build", "-o", bin}, buildArgs...)
	args = append(args, ".")
	if out, err := exec.Command(goTool, args...).CombinedOutput(); err != nil {
		t.Fatalf("go build: %v\n%s", err, out)
	}

	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("find free port: %v", err)
	}
	addr := l.Addr().String()
	l.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, bin, "-addr", addr)
	if err := cmd.Start(); err != nil {
		cancel()
		t.Fatalf("start binary: %v", err)
	}
	t.Cleanup(func() {
		cancel()
		_ = cmd.Wait()
	})

	base := "http://" + addr
	deadline := time.Now().Add(10 * time.Second)
	for {
		resp, err := http.Get(base + "/version")
		if err == nil {
			resp.Body.Close()
			return base
		}
		if time.Now().After(deadline) {
			t.Fatalf("server at %s did not come up: %v", addr, err)
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func fetchVersion(t *testing.T, base string) (int, map[string]string) {
	t.Helper()
	resp, err := http.Get(base + "/version")
	if err != nil {
		t.Fatalf("GET /version: %v", err)
	}
	defer resp.Body.Close()
	var body map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	return resp.StatusCode, body
}

// Criterion 1 and 2: a binary built without a version set answers GET /version with 200 {"version": "dev"}.
func TestVersion_BinaryBuiltWithoutVersionReportsDev(t *testing.T) {
	base := buildAndServe(t)
	code, body := fetchVersion(t, base)
	if code != http.StatusOK {
		t.Errorf("status = %d, want 200", code)
	}
	if len(body) != 1 || body["version"] != "dev" {
		t.Errorf("body = %v, want {\"version\": \"dev\"}", body)
	}
}

// Criterion 1 and the ldflags constraint: -ldflags "-X main.version=..." sets the reported version.
func TestVersion_BinaryBuiltWithLdflagsReportsThatVersion(t *testing.T) {
	base := buildAndServe(t, "-ldflags", "-X main.version=v9.8.7")
	code, body := fetchVersion(t, base)
	if code != http.StatusOK {
		t.Errorf("status = %d, want 200", code)
	}
	if len(body) != 1 || body["version"] != "v9.8.7" {
		t.Errorf("body = %v, want {\"version\": \"v9.8.7\"}", body)
	}
}

// Criterion 3: POST /version answers 405, both in-process and against the running binary.
func TestVersion_PostIsNotAllowed(t *testing.T) {
	t.Run("handler", func(t *testing.T) {
		rec := httptest.NewRecorder()
		newHandler().ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/version", nil))
		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("status = %d, want 405", rec.Code)
		}
	})
	t.Run("binary", func(t *testing.T) {
		base := buildAndServe(t)
		resp, err := http.Post(base+"/version", "application/json", nil)
		if err != nil {
			t.Fatalf("POST /version: %v", err)
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusMethodNotAllowed {
			t.Errorf("status = %d, want 405", resp.StatusCode)
		}
	})
}
