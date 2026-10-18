package shop

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"shanraq.org/pkg/shanraq"
)

func TestTartarExcerptUsesSiteFrameAndKeepsAllSections(t *testing.T) {
	content, err := loadTartarExcerpt()
	if err != nil {
		t.Fatal(err)
	}
	module := &Module{rt: &shanraq.Runtime{Site: renderer(t), Logger: zap.NewNop()}, tartar: content}
	router := chi.NewRouter()
	router.Get(tartarPath, module.handleTartar)
	router.Get(tartarPath+"/{section}", module.handleTartar)

	for _, path := range []string{
		tartarPath, tartarPath + "/why", tartarPath + "/preface",
		tartarPath + "/chapter-1", tartarPath + "/chapter-2",
	} {
		response := httptest.NewRecorder()
		router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, path+"?lang=ru", nil))
		if response.Code != http.StatusOK {
			t.Errorf("%s: status %d", path, response.Code)
			continue
		}
		body := response.Body.String()
		for _, marker := range []string{"class=\"topstrip\"", "class=\"site-footer\"", "Тартар и Я"} {
			if !strings.Contains(body, marker) {
				t.Errorf("%s: missing %q", path, marker)
			}
		}
		if strings.Contains(body, "href=\"chapter-1.html\"") {
			t.Errorf("%s: old standalone navigation leaked into site page", path)
		}
	}

	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, tartarPath+"/missing?lang=ru", nil))
	if response.Code != http.StatusNotFound {
		t.Errorf("unknown section: status %d, want 404", response.Code)
	}
}
