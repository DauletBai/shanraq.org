package site

import "testing"

func TestTranslationsCoverAllLangs(t *testing.T) {
	for key, m := range messages {
		for _, lang := range Langs {
			if v, ok := m[lang]; !ok || v == "" {
				t.Errorf("translation key %q missing %s", key, lang)
			}
		}
	}
}
