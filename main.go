// Command playground is the target repo for HOTL's first real agent runs: small, safe changes an agent can plan,
// implement, test and merge without touching anything that matters.
//
// It runs a small HTTP server. GET /version reports the running version, which can be set at build time with
// -ldflags "-X main.version=...".
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
)

// version is the running version, set at build time with -ldflags "-X main.version=...".
var version = "dev"

// Greeting returns the greeting for name.
func Greeting(name string) string {
	if name == "" {
		name = "world"
	}
	return fmt.Sprintf("Hello, %s!", name)
}

// newHandler returns the server's HTTP handler.
func newHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/version", handleVersion)
	return mux
}

// handleVersion answers GET /version with {"version": "<version>"} and any other method with 405.
func handleVersion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"version": version}); err != nil {
		log.Printf("write /version response: %v", err)
	}
}

func main() {
	addr := flag.String("addr", ":8080", "address to listen on")
	flag.Parse()
	log.Printf("playground %s listening on %s", version, *addr)
	log.Fatal(http.ListenAndServe(*addr, newHandler()))
}
