package main

import "testing"

func TestGreeting(t *testing.T) {
	for name, want := range map[string]string{"": "Hello, world!", "HOTL": "Hello, HOTL!"} {
		if got := Greeting(name); got != want {
			t.Errorf("Greeting(%q) = %q, want %q", name, got, want)
		}
	}
}
