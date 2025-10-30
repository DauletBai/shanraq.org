package respond

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDecode(t *testing.T) {
	body := bytes.NewBufferString(`{"value":"ok"}`)
	req, _ := http.NewRequest(http.MethodPost, "/", body)
	var out map[string]string
	if err := Decode(req, &out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if out["value"] != "ok" {
		t.Fatalf("unexpected decode output: %+v", out)
	}
}

func TestJSON(t *testing.T) {
	rec := httptest.NewRecorder()
	JSON(rec, http.StatusAccepted, map[string]string{"status": "ok"})
	if rec.Code != http.StatusAccepted {
		t.Fatalf("expected %d, got %d", http.StatusAccepted, rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("unexpected content type %s", ct)
	}
}

func TestError(t *testing.T) {
	rec := httptest.NewRecorder()
	Error(rec, http.StatusBadRequest, assertErr("boom"))
	if rec.Body.Len() == 0 {
		t.Fatalf("expected error body")
	}
}

type assertErr string

func (e assertErr) Error() string { return string(e) }
