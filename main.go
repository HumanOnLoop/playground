// Command playground is the target repo for HOTL's first real agent runs: small, safe changes an agent can plan,
// implement, test and merge without touching anything that matters.
package main

import "fmt"

// Greeting returns the greeting for name.
func Greeting(name string) string {
	if name == "" {
		name = "world"
	}
	return fmt.Sprintf("Hello, %s!", name)
}

func main() { fmt.Println(Greeting("")) }
