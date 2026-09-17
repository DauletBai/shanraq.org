package web

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBookSamplePublicReadingAndDownloads(t *testing.T) {
	prefix := "/shop/go-book-sample/0.21.0-sample/"
	for _, path := range []string{"", "read/cover/go-book-cover-v3.png", "read/toc.html", "read/12-templates.html", "go-book-preview.pdf", "go-book-preview.epub", "go-book-preview-html.zip", "go-book-preview-code.zip"} {
		response := httptest.NewRecorder()
		StaticHandler().ServeHTTP(response, httptest.NewRequest("GET", prefix+path, nil))
		if response.Code != 200 || response.Body.Len() == 0 {
			t.Fatalf("public sample %s: status %d, bytes %d", path, response.Code, response.Body.Len())
		}
	}
	response := httptest.NewRecorder()
	StaticHandler().ServeHTTP(response, httptest.NewRequest("GET", prefix+"read/13-forms.html", nil))
	if response.Code != 404 {
		t.Fatalf("paid chapter returned %d", response.Code)
	}
	data, err := staticFiles.ReadFile("static" + prefix + "sample.json")
	if err != nil {
		t.Fatal(err)
	}
	var sample struct {
		Chapters []string `json:"chapters"`
	}
	if err := json.Unmarshal(data, &sample); err != nil {
		t.Fatal(err)
	}
	if strings.Join(sample.Chapters, ",") != "00,01,02,03,04,05,06,07,08,09,10,11,12" {
		t.Fatal("unexpected sample chapter boundary")
	}
}
