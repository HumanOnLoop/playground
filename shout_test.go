package main

import "testing"

func TestShout(t *testing.T) {
	tests := []struct {
		name   string
		phrase string
		want   string
	}{
		{name: "hello", phrase: "hello", want: "HELLO!"},
		{name: "empty phrase", phrase: "", want: "!"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := Shout(tc.phrase); got != tc.want {
				t.Errorf("Shout(%q) = %q, want %q", tc.phrase, got, tc.want)
			}
		})
	}
}
