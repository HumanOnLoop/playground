package main

import (
	"log"
	"net/http"
)

// handlePing answers GET /ping with text/plain "pong" and any other method with 405.
func handlePing(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	if _, err := w.Write([]byte("pong")); err != nil {
		log.Printf("write /ping response: %v", err)
	}
}
