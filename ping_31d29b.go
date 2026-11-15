package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// handlePing31d29b answers GET /ping-31d29b with {"pong": "31d29b"} and any other method with 405.
func handlePing31d29b(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"pong": "31d29b"}); err != nil {
		log.Printf("write /ping-31d29b response: %v", err)
	}
}
