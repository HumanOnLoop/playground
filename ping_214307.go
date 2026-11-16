package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// handlePing214307 answers GET /ping-214307 with {"ok":true} and any other method with 405.
func handlePing214307(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]bool{"ok": true}); err != nil {
		log.Printf("write /ping-214307 response: %v", err)
	}
}
