package db

import (
	"errors"
	"net"
	"testing"
)

func TestRetryable(t *testing.T) {
	// "connection refused" is a net.Error whose Timeout()/Temporary() are both
	// false — it must still be retryable (Postgres is often still starting).
	refused := &net.OpError{Op: "dial", Err: errors.New("connect: connection refused")}
	if !retryable(refused) {
		t.Error("connection refused should be retryable")
	}

	if !retryable(errors.New("failed: connection refused")) {
		t.Error("string 'connection refused' should be retryable")
	}
	if retryable(errors.New("password authentication failed")) {
		t.Error("auth failure should NOT be retryable")
	}
}

// timeoutErr is a net.Error that reports a timeout.
type timeoutErr struct{}

func (timeoutErr) Error() string   { return "i/o timeout" }
func (timeoutErr) Timeout() bool   { return true }
func (timeoutErr) Temporary() bool { return false }

func TestRetryableTimeout(t *testing.T) {
	if !retryable(timeoutErr{}) {
		t.Error("a net.Error timeout should be retryable")
	}
}
