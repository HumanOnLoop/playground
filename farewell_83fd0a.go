package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// farewells83fd0a maps a language code to the farewell /farewell-83fd0a answers with.
var farewells83fd0a = map[string]string{
	"en": "Goodbye",
	"sv": "Hej då",
}

// handleFarewell83fd0a answers GET /farewell-83fd0a?lang=<code> with {"text": "<farewell>"} and any other method
// with 405. A missing or unknown lang gets the English farewell.
func handleFarewell83fd0a(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	text, ok := farewells83fd0a[r.URL.Query().Get("lang")]
	if !ok {
		text = farewells83fd0a["en"]
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"text": text}); err != nil {
		log.Printf("write /farewell-83fd0a response: %v", err)
	}
}
