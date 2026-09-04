// Package goblog holds no code of its own — only the guard below.
//
// The blog project is what the reader builds alongside the lessons, so an
// example that no longer compiles is worse than no example: it teaches the
// reader to distrust the course. Steps without a go.mod are covered by the
// repository's own `go build ./...`. A step that is its own module — the
// package split in step 6 — is invisible to that, so it is built here.
package goblog

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestEveryStepModuleBuilds(t *testing.T) {
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	built := 0
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		dir := e.Name()
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err != nil {
			continue // part of the repository's own module, already built
		}
		cmd := exec.Command("go", "build", "./...")
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Errorf("%s не собирается: %v\n%s", dir, err, out)
			continue
		}
		// go build leaves the binary behind; the repository stays clean.
		_ = os.Remove(filepath.Join(dir, dir))
		built++
	}
	if built == 0 {
		t.Error("ни один шаг не проверен — тест перестал что-либо доказывать")
	}
	t.Logf("шагов-модулей собрано: %d", built)
}
