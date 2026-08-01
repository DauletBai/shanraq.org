package sms

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewGate(t *testing.T) {
	if _, ok, err := New(Config{Provider: ""}); ok || err != nil {
		t.Fatalf("empty provider: want disabled, got ok=%v err=%v", ok, err)
	}
	if _, _, err := New(Config{Provider: "mobizon"}); err == nil {
		t.Fatal("mobizon without api key: want error")
	}
	if _, _, err := New(Config{Provider: "smsc", Login: "x"}); err == nil {
		t.Fatal("smsc without password: want error")
	}
	if _, _, err := New(Config{Provider: "carrierpigeon", APIKey: "x"}); err == nil {
		t.Fatal("unknown provider: want error")
	}
	if _, ok, err := New(Config{Provider: "mobizon", APIKey: "k"}); !ok || err != nil {
		t.Fatalf("valid mobizon: want ok, got ok=%v err=%v", ok, err)
	}
}

func TestSendMobizon(t *testing.T) {
	var gotRecipient, gotText, gotKey string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotRecipient = r.URL.Query().Get("recipient")
		gotText = r.URL.Query().Get("text")
		gotKey = r.URL.Query().Get("apiKey")
		w.Write([]byte(`{"code":0,"message":""}`))
	}))
	defer srv.Close()

	c, _, err := New(Config{Provider: "mobizon", APIKey: "secret", BaseURL: srv.URL})
	if err != nil {
		t.Fatal(err)
	}
	if err := c.SendSMS(context.Background(), "+7 (707) 915-22-06", "code 123"); err != nil {
		t.Fatalf("send: %v", err)
	}
	if gotRecipient != "77079152206" {
		t.Errorf("recipient = %q, want bare digits 77079152206", gotRecipient)
	}
	if gotText != "code 123" || gotKey != "secret" {
		t.Errorf("text=%q key=%q", gotText, gotKey)
	}
}

func TestSendMobizonError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"code":1,"message":"invalid recipient"}`))
	}))
	defer srv.Close()
	c, _, _ := New(Config{Provider: "mobizon", APIKey: "k", BaseURL: srv.URL})
	if err := c.SendSMS(context.Background(), "77079152206", "x"); err == nil {
		t.Fatal("nonzero code: want error")
	}
}

func TestSendSMSC(t *testing.T) {
	var gotPhones, gotLogin string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPhones = r.URL.Query().Get("phones")
		gotLogin = r.URL.Query().Get("login")
		w.Write([]byte(`{"id":42,"cnt":1}`))
	}))
	defer srv.Close()
	c, _, _ := New(Config{Provider: "smsc", Login: "me", Password: "pw", BaseURL: srv.URL})
	if err := c.SendSMS(context.Background(), "77079152206", "x"); err != nil {
		t.Fatalf("send: %v", err)
	}
	if gotPhones != "77079152206" || gotLogin != "me" {
		t.Errorf("phones=%q login=%q", gotPhones, gotLogin)
	}
}

func TestSendSMSCError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"error":"authorization failed","error_code":2}`))
	}))
	defer srv.Close()
	c, _, _ := New(Config{Provider: "smsc", Login: "me", Password: "bad", BaseURL: srv.URL})
	if err := c.SendSMS(context.Background(), "77079152206", "x"); err == nil {
		t.Fatal("smsc error payload: want error")
	}
}
