package web

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

// The palette is the one accessibility property that can be checked without a
// browser, and it is the one that drifts: a colour is nudged for looks, and
// nobody notices that the caption under every photo has fallen to 3:1. Axe
// found exactly that — "многочисленные нарушения контрастности" — across both
// themes, and a report that has to be re-run by hand is a report that stops
// being run.
//
// This asserts WCAG 2.2 AA (1.4.3): 4.5:1 for body text against every surface
// it can sit on. It reads the same stylesheet the site ships.

var tokenRe = regexp.MustCompile(`(--[\w-]+):\s*(#[0-9a-fA-F]{3,8})\s*;`)

// themeTokens returns the colour tokens of the light theme and of the dark
// theme, the latter falling through to light for anything it does not redefine
// — which is how the cascade actually resolves them.
func themeTokens(t *testing.T) (light, dark map[string]string) {
	t.Helper()
	b, err := staticFiles.ReadFile("static/css/shanraq.css")
	if err != nil {
		t.Fatalf("read stylesheet: %v", err)
	}
	css := string(b)

	block := func(start string) string {
		i := strings.Index(css, start)
		if i < 0 {
			t.Fatalf("stylesheet has no %q block", start)
		}
		j := strings.Index(css[i:], "}")
		if j < 0 {
			t.Fatalf("unterminated %q block", start)
		}
		return css[i : i+j]
	}

	parse := func(s string) map[string]string {
		out := map[string]string{}
		for _, m := range tokenRe.FindAllStringSubmatch(s, -1) {
			out[m[1]] = strings.ToLower(m[2])
		}
		return out
	}

	light = parse(block(":root {"))
	dark = map[string]string{}
	for k, v := range light {
		dark[k] = v
	}
	for k, v := range parse(block(`:root[data-theme="dark"] {`)) {
		dark[k] = v
	}
	return light, dark
}

func channel(c float64) float64 {
	c /= 255
	if c <= 0.03928 {
		return c / 12.92
	}
	return math.Pow((c+0.055)/1.055, 2.4)
}

func luminance(t *testing.T, hex string) float64 {
	t.Helper()
	h := strings.TrimPrefix(hex, "#")
	if len(h) == 3 {
		h = string([]byte{h[0], h[0], h[1], h[1], h[2], h[2]})
	}
	if len(h) < 6 {
		t.Fatalf("cannot read colour %q", hex)
	}
	v := func(i int) float64 {
		n, err := strconv.ParseUint(h[i:i+2], 16, 8)
		if err != nil {
			t.Fatalf("cannot read colour %q: %v", hex, err)
		}
		return channel(float64(n))
	}
	return 0.2126*v(0) + 0.7152*v(2) + 0.0722*v(4)
}

func contrast(t *testing.T, fg, bg string) float64 {
	t.Helper()
	a, b := luminance(t, fg), luminance(t, bg)
	if a < b {
		a, b = b, a
	}
	return (a + 0.05) / (b + 0.05)
}

// textTokens are the colours that carry words; surfaces are what they sit on.
var (
	textTokens = []string{"--ink", "--ink-soft", "--muted", "--gold", "--gold-strong",
		"--teal", "--clay", "--danger", "--ok", "--st-ok", "--st-warn", "--st-info", "--st-off"}
	surfaceTokens = []string{"--paper", "--surface", "--surface-2"}
)

const aaText = 4.5

func TestPaletteMeetsAAInBothThemes(t *testing.T) {
	light, dark := themeTokens(t)
	for _, theme := range []struct {
		name   string
		tokens map[string]string
	}{{"light", light}, {"dark", dark}} {
		for _, bg := range surfaceTokens {
			bgHex, ok := theme.tokens[bg]
			if !ok {
				t.Fatalf("%s theme has no %s", theme.name, bg)
			}
			for _, fg := range textTokens {
				fgHex, ok := theme.tokens[fg]
				if !ok {
					t.Fatalf("%s theme has no %s", theme.name, fg)
				}
				if r := contrast(t, fgHex, bgHex); r < aaText {
					t.Errorf("%s theme: %s (%s) on %s (%s) is %s — WCAG AA needs %.1f:1",
						theme.name, fg, fgHex, bg, bgHex,
						fmt.Sprintf("%.2f:1", r), aaText)
				}
			}
		}
	}
}

// A solid accent fill takes the label colour --on-accent. The whole point of
// that token is that the answer flips with the theme, so both answers are
// checked against every accent that is ever filled.
func TestAccentFillsCarryALegibleLabel(t *testing.T) {
	light, dark := themeTokens(t)
	fills := []string{"--gold", "--gold-strong", "--teal", "--danger"}
	for _, theme := range []struct {
		name   string
		tokens map[string]string
	}{{"light", light}, {"dark", dark}} {
		label, ok := theme.tokens["--on-accent"]
		if !ok {
			t.Fatalf("%s theme defines no --on-accent", theme.name)
		}
		for _, fill := range fills {
			fillHex, ok := theme.tokens[fill]
			if !ok {
				continue
			}
			if r := contrast(t, label, fillHex); r < aaText {
				t.Errorf("%s theme: --on-accent (%s) on %s (%s) is %.2f:1 — WCAG AA needs %.1f:1",
					theme.name, label, fill, fillHex, r, aaText)
			}
		}
	}
}
