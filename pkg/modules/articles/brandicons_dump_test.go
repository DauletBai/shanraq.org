package articles

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Dumps every mark as a standalone SVG so it can be rasterised and inspected
// outside the test binary. Off unless BRANDICO_DUMP names a directory: this is
// a debugging aid for the one thing a unit test cannot check — what the mark
// actually looks like.
func TestDumpBrandMarks(t *testing.T) {
	dir := os.Getenv("BRANDICO_DUMP")
	if dir == "" {
		t.Skip("set BRANDICO_DUMP=<dir> to write the marks out")
	}
	for slug := range brandMarks {
		svg := string(brandIcon(slug))
		svg = strings.Replace(svg, `<svg class="brandico brandico--ink"`,
			`<svg xmlns="http://www.w3.org/2000/svg" width="96" height="96"`, 1)
		svg = strings.Replace(svg, `<svg class="brandico"`,
			`<svg xmlns="http://www.w3.org/2000/svg" width="96" height="96"`, 1)
		svg = strings.ReplaceAll(svg, "currentColor", "#111111")
		if err := os.WriteFile(filepath.Join(dir, slug+".svg"), []byte(svg), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}
