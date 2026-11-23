package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

// hashSuffixed matches file base names ending in a commit-hash-like suffix, such as farewell_360095.go or
// greet_9d0d00_test.go: an underscore or hyphen, six lowercase hex digits, an optional _test, then the extension.
var hashSuffixed = regexp.MustCompile(`[_-][0-9a-f]{6}(_test)?\.[a-z]+$`)

func TestHashSuffixedNamePattern(t *testing.T) {
	for name, want := range map[string]bool{
		"farewell_360095.go":      true,
		"farewell_360095_test.go": true,
		"ping-31d29b.go":          true,
		"main.go":                 false,
		"farewell.go":             false,
		"ping_test.go":            false,
		"version_binary_test.go":  false,
	} {
		if got := hashSuffixed.MatchString(name); got != want {
			t.Errorf("hashSuffixed.MatchString(%q) = %v, want %v", name, got, want)
		}
	}
}

func TestNoCommitHashSuffixedFiles(t *testing.T) {
	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && d.Name() == ".git" {
			return filepath.SkipDir
		}
		if !d.IsDir() && hashSuffixed.MatchString(d.Name()) {
			t.Errorf("file %s has a commit-hash-like suffix in its name", path)
		}
		return nil
	})
	if err != nil && !os.IsNotExist(err) {
		t.Fatalf("walk repository: %v", err)
	}
}
