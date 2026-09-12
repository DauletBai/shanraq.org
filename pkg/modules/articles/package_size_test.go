package articles

import (
	"os"
	"strings"
	"testing"
)

// This package is being taken apart, not added to.
//
// It holds the site's first module and everything that was quickest to put
// beside it: articles, listings, ads, loans, courses, weather, exchange rates,
// analytics. That was cheap each time and expensive in total -- one Module
// struct now owns two dozen stores, and a reader looking for the weather page
// has to know it lives among the articles.
//
// The way out is one domain at a time into a module of its own, and this test
// is the ratchet that keeps the direction: the ceiling only ever moves down.
// A new domain does not go here at all -- it starts as its own module under
// pkg/modules, the way the shop does.
func TestArticlesPackageOnlyShrinks(t *testing.T) {
	// Lower this when a domain moves out. Never raise it.
	const ceiling = 162

	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("не прочитать каталог пакета: %v", err)
	}
	files := 0
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".go") {
			files++
		}
	}
	if files > ceiling {
		t.Fatalf("в пакете articles стало %d файлов при потолке %d.\n"+
			"Новому домену здесь не место: заведите модуль в pkg/modules,\n"+
			"как это сделано для shop. Потолок опускают, когда домен уезжает.", files, ceiling)
	}
	t.Logf("файлов в пакете: %d, потолок: %d", files, ceiling)
}
