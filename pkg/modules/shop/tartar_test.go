package shop

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"regexp"
	"shanraq.org/web"
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
		tartarPath, tartarPath + "/preface", tartarPath + "/prologue",
		tartarPath + "/chapter-1", tartarPath + "/chapter-2",
		tartarPath + "/chapter-3", tartarPath + "/chapter-4", tartarPath + "/chapter-5",
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
	for path, marker := range map[string]string{
		tartarPath:                "Часть I. Имя народа и память речи",
		tartarPath + "/prologue":  "5 февраля 1936 года",
		tartarPath + "/chapter-1": "Ру — Русь — Русский",
		tartarPath + "/chapter-4": "atkinson-bashkirs-1803-detail.jpg",
		tartarPath + "/chapter-5": "Источники и пояснения",
	} {
		response := httptest.NewRecorder()
		router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, path+"?lang=ru", nil))
		if !strings.Contains(response.Body.String(), marker) {
			t.Errorf("%s: missing %q", path, marker)
		}
	}

	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, tartarPath+"/missing?lang=ru", nil))
	if response.Code != http.StatusNotFound {
		t.Errorf("unknown section: status %d, want 404", response.Code)
	}
}

// Every footnote in the excerpt must resolve on its own chapter page.
func TestTartarExcerptFootnotesResolve(t *testing.T) {
	content, err := loadTartarExcerpt()
	if err != nil {
		t.Fatal(err)
	}
	ref := regexp.MustCompile(`href="#(source-[0-9]+)"`)
	for slug, section := range content {
		body := string(section)
		for _, match := range ref.FindAllStringSubmatch(body, -1) {
			if !strings.Contains(body, `id="`+match[1]+`"`) {
				t.Errorf("%s: missing %s", slug, match[1])
			}
		}
		if strings.Contains(body, "../assets/") {
			t.Errorf("%s: print-only image path", slug)
		}
	}
}

// The deployed excerpt must be exactly the seven sections verified at import.
func TestTartarExcerptEditionMatchesEmbeddedContent(t *testing.T) {
	data, err := fs.ReadFile(web.StaticFS(), "shop/tartar-i-ya-sample/edition.json")
	if err != nil {
		t.Fatal(err)
	}
	var edition struct {
		Edition  string `json:"edition"`
		Date     string `json:"edition_date"`
		Sections []struct {
			Slug string `json:"slug"`
			Hash string `json:"sha256"`
		} `json:"sections"`
	}
	if err := json.Unmarshal(data, &edition); err != nil {
		t.Fatal(err)
	}
	if edition.Edition != "0.88" || edition.Date != "2026-11-04" || len(edition.Sections) != len(tartarSections) {
		t.Fatalf("unexpected excerpt manifest: %+v", edition)
	}
	content, err := loadTartarExcerpt()
	if err != nil {
		t.Fatal(err)
	}
	for i, section := range edition.Sections {
		if section.Slug != tartarSections[i].Slug {
			t.Fatalf("unexpected section order: %s", section.Slug)
		}
		sum := sha256.Sum256([]byte(content[section.Slug]))
		if hex.EncodeToString(sum[:]) != section.Hash {
			t.Errorf("%s differs from verified edition", section.Slug)
		}
	}
	if !strings.Contains(string(content["preface"]), "Рудненский индустриальный институт") {
		t.Error("updated author biography is missing")
	}
}
