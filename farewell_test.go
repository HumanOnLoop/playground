package main

import "testing"

func TestFarewell(t *testing.T) {
	if got, want := Farewell(), "goodbye"; got != want {
		t.Errorf("Farewell() = %q, want %q", got, want)
	}
}
