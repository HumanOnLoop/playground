package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// farewells413f24 maps a language code to the farewell /farewell-413f24 answers with.
var farewells413f24 = map[string]string{
	"en": "Goodbye",
	"sv": "Hej då",
}

// handleFarewell413f24 answers GET /farewell-413f24?lang=<code> with {"text": "<farewell>"} and any other method
// with 405. A missing or unknown lang gets the English farewell.
func handleFarewell413f24(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	text, ok := farewells413f24[r.URL.Query().Get("lang")]
	if !ok {
		text = farewells413f24["en"]
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"text": text}); err != nil {
		log.Printf("write /farewell-413f24 response: %v", err)
	}
}
