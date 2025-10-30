package respond

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Decode parses JSON body into dest.
func Decode(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf("decode json: %w", err)
	}
	return nil
}

// JSON writes payload as JSON with provided status.
func JSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// Error sends {"error": "..."} JSON response.
func Error(w http.ResponseWriter, status int, err error) {
	JSON(w, status, map[string]string{"error": err.Error()})
}
