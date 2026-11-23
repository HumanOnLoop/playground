package main

import "strings"

// Shout returns the phrase in upper case followed by "!".
func Shout(phrase string) string {
	return strings.ToUpper(phrase) + "!"
}
