package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// greetings9d0d00 maps a language code to the greeting /greet-9d0d00 answers with.
var greetings9d0d00 = map[string]string{
	"en": "Hello",
	"sv": "Hej",
}

// handleGreet9d0d00 answers GET /greet-9d0d00?lang=<code> with {"text": "<greeting>"} and any other method with
// 405. A missing or unknown lang gets the English greeting.
func handleGreet9d0d00(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	text, ok := greetings9d0d00[r.URL.Query().Get("lang")]
	if !ok {
		text = greetings9d0d00["en"]
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"text": text}); err != nil {
		log.Printf("write /greet-9d0d00 response: %v", err)
	}
}
