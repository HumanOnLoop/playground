// Command playground is the target repo for HOTL's first real agent runs: small, safe changes an agent can plan,
// implement, test and merge without touching anything that matters.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
)

// version is the build version, set with -ldflags "-X main.version=...".
var version = "dev"

// Greeting returns the greeting for name.
func Greeting(name string) string {
	if name == "" {
		name = "world"
	}
	return fmt.Sprintf("Hello, %s!", name)
}

// newHandler returns the server's HTTP handler. Requests to /version with a method other than GET get 405.
func newHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /version", handleVersion)
	return mux
}

func handleVersion(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"version": version}); err != nil {
		log.Printf("write /version response: %v", err)
	}
}

func main() {
	addr := flag.String("addr", ":8080", "address to listen on")
	flag.Parse()
	fmt.Println(Greeting(""))
	log.Printf("playground %s listening on %s", version, *addr)
	log.Fatal(http.ListenAndServe(*addr, newHandler()))
}
