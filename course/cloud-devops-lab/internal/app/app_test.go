package app

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestStorePersistsNotes(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "notes.json")
	store, err := loadStore(path)
	if err != nil {
		t.Fatal(err)
	}
	for _, note := range []string{"first deploy", "rollback checked"} {
		if err := store.add(note); err != nil {
			t.Fatal(err)
		}
	}
	reloaded, err := loadStore(path)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"first deploy", "rollback checked"}; !reflect.DeepEqual(reloaded.list(), want) {
		t.Fatalf("notes = %#v, want %#v", reloaded.list(), want)
	}
}

func TestConfigRejectsAddressWithoutColon(t *testing.T) {
	t.Setenv("CLOUDLAB_ADDR", "8080")
	if _, err := ConfigFromEnv("test"); err == nil {
		t.Fatal("expected invalid address to be rejected")
	}
}
