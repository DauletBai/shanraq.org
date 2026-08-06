package respond

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
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

// internalMessage is what a caller is told when something broke on our side.
// Deliberately fixed and uninformative: the caller can do nothing with the
// detail, and an attacker can do quite a lot.
const internalMessage = "internal error"

// LogInternal receives every error that Error hides from the client, so masking
// the response does not also mask the failure. The application points this at
// the real logger during boot; until then it writes to the standard logger, so
// a fault is never silently swallowed. Set once at startup, never concurrently.
var LogInternal = func(status int, err error) {
	log.Printf("respond: %d %v", status, err)
}

// Error sends {"error": "..."} as JSON.
//
// A 4xx carries its text through: those errors exist to tell the caller what
// was wrong with their request. A 5xx does not. err.Error() at that level is
// whatever the database, the driver or a third-party client produced, and it
// has a habit of containing SQL fragments, column names and connection strings
// — handed to anyone who can make a request fail. The detail goes to the log;
// the caller gets one fixed sentence.
func Error(w http.ResponseWriter, status int, err error) {
	if status >= http.StatusInternalServerError {
		if LogInternal != nil {
			LogInternal(status, err)
		}
		JSON(w, status, map[string]string{"error": internalMessage})
		return
	}
	JSON(w, status, map[string]string{"error": err.Error()})
}

// Validation writes validation errors as JSON.
func Validation(w http.ResponseWriter, fields map[string]string) {
	JSON(w, http.StatusBadRequest, map[string]any{
		"error":  "validation failed",
		"fields": fields,
	})
}
