package main

import "testing"

func TestFarewell360095(t *testing.T) {
	if got, want := Farewell360095(), "goodbye"; got != want {
		t.Errorf("Farewell360095() = %q, want %q", got, want)
	}
}
