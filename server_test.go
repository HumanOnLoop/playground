package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadMotd(t *testing.T) {
	path := filepath.Join(t.TempDir(), "motd")
	if err := os.WriteFile(path, []byte("  Ship it.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got, err := loadMotd(path); err != nil || got != "Ship it." {
		t.Errorf("loadMotd = %q, %v", got, err)
	}
	if _, err := loadMotd(filepath.Join(t.TempDir(), "missing")); err == nil {
		t.Error("no error for a missing file")
	}
}
