package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

// loadMotd reads the message of the day from path.
func loadMotd(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(b)), nil
}

// motdHandler answers GET /motd with the message of the day as JSON.
func motdHandler(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		msg, err := loadMotd(path)
		if err != nil {
			slog.ErrorContext(r.Context(), "no message of the day", "error", err)
			http.Error(w, "no message today", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"message": msg})
	}
}
