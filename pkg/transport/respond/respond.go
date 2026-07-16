package respond

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// maxJSONBody caps request bodies to guard against decompression/oversize abuse.
const maxJSONBody = 1 << 20 // 1 MiB

// Decode parses a single JSON value from the request body into dest. It caps the
// body size, rejects unknown fields, and rejects trailing garbage after the
// value (a valid object followed by junk is an error, not silently accepted).
func Decode(r *http.Request, dest any) error {
	dec := json.NewDecoder(io.LimitReader(r.Body, maxJSONBody))
	dec.DisallowUnknownFields()
	if err := dec.Decode(dest); err != nil {
		return fmt.Errorf("decode json: %w", err)
	}
	if dec.More() {
		return errors.New("decode json: unexpected trailing data after JSON value")
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

// Validation writes validation errors as JSON.
func Validation(w http.ResponseWriter, fields map[string]string) {
	JSON(w, http.StatusBadRequest, map[string]any{
		"error":  "validation failed",
		"fields": fields,
	})
}
