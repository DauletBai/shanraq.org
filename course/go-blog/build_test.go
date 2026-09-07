// Package goblog holds no code of its own — only the guard below.
//
// The blog project is what the reader builds alongside the lessons, so an
// example that no longer compiles is worse than no example: it teaches the
// reader to distrust the course. Steps without a go.mod are covered by the
// repository's own `go build ./...`. A step that is its own module — the
// package split in step 6 — is invisible to that, so it is checked here.
//
// Building is not enough. Each of those modules carries its own tests, and a
// root `go test ./...` never reaches them: a nested module is invisible to the
// parent. So every step module is built, vetted and tested here, and the ones
// that start goroutines are tested again under the race detector — otherwise a
// green run would prove only that the code compiles.
package goblog

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
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
		ok := true
		for _, args := range [][]string{
			{"build", "./..."},
			{"vet", "./..."},
			{"test", "./..."},
		} {
			cmd := exec.Command("go", args...)
			cmd.Dir = dir
			if out, err := cmd.CombinedOutput(); err != nil {
				t.Errorf("%s: go %s не прошёл: %v\n%s", dir, args[0], err, out)
				ok = false
				break
			}
		}
		// The race detector costs a second per step, so it runs where there is
		// something to race: from the goroutines lesson onwards.
		if ok && racy(dir) {
			cmd := exec.Command("go", "test", "-race", "./...")
			cmd.Dir = dir
			if out, err := cmd.CombinedOutput(); err != nil {
				t.Errorf("%s: go test -race не прошёл: %v\n%s", dir, err, out)
				ok = false
			}
		}
		// go build leaves the binary behind; the repository stays clean.
		_ = os.Remove(filepath.Join(dir, dir))
		if !ok {
			continue
		}
		built++
	}
	if built == 0 {
		t.Error("ни один шаг не проверен — тест перестал что-либо доказывать")
	}
	t.Logf("шагов-модулей проверено: %d", built)
}

// racy reports whether the step has goroutines worth running the detector on.
// Shared state appears with the counter in step 27 and never leaves.
func racy(dir string) bool {
	n, err := strconv.Atoi(strings.TrimPrefix(dir, "step-"))
	return err == nil && n >= 27
}
