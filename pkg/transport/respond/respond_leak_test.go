package respond

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// A 5xx must never carry the underlying error to the client: at that level the
// text is whatever the database or a driver produced, and it has a habit of
// containing SQL, column names and connection strings.
func TestErrorMasksServerFaults(t *testing.T) {
	secret := `pq: relation "auth_users" does not exist (host=10.0.0.5 user=shanraq)`

	var logged error
	old := LogInternal
	LogInternal = func(status int, err error) { logged = err }
	defer func() { LogInternal = old }()

	w := httptest.NewRecorder()
	Error(w, http.StatusInternalServerError, errors.New(secret))

	body := w.Body.String()
	if strings.Contains(body, "auth_users") || strings.Contains(body, "10.0.0.5") {
		t.Errorf("server fault leaked to the client: %s", body)
	}
	var got map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("response is not JSON: %v", err)
	}
	if got["error"] != internalMessage {
		t.Errorf("error = %q, want %q", got["error"], internalMessage)
	}
	// Masking must not also silence the fault.
	if logged == nil || !strings.Contains(logged.Error(), "auth_users") {
		t.Errorf("the detail did not reach the log: %v", logged)
	}
}

// A 4xx exists to tell the caller what was wrong with their request, so its
// text must survive.
func TestErrorKeepsClientMessages(t *testing.T) {
	for _, status := range []int{http.StatusBadRequest, http.StatusForbidden, http.StatusConflict} {
		w := httptest.NewRecorder()
		Error(w, status, errors.New("email already registered"))
		if !strings.Contains(w.Body.String(), "email already registered") {
			t.Errorf("status %d dropped its message: %s", status, w.Body.String())
		}
	}
}
