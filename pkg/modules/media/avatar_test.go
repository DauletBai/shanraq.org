package media

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"testing"

	"github.com/google/uuid"
)

func TestProcessAvatarSquare(t *testing.T) {
	m := &Module{}
	// A non-square source must come out as a centred avatarDim square.
	raw := solidPNG(t, 240, 120, color.RGBA{10, 20, 30, 255})
	out, err := m.processAvatar(raw)
	if err != nil {
		t.Fatalf("processAvatar: %v", err)
	}
	cfg, _, err := image.DecodeConfig(bytes.NewReader(out))
	if err != nil {
		t.Fatalf("decode result: %v", err)
	}
	if cfg.Width != avatarDim || cfg.Height != avatarDim {
		t.Errorf("avatar is %dx%d, want %dx%d square", cfg.Width, cfg.Height, avatarDim, avatarDim)
	}
}

func TestProcessAvatarRejectsNonImage(t *testing.T) {
	m := &Module{}
	if _, err := m.processAvatar([]byte("definitely not an image")); err == nil {
		t.Error("expected an error for non-image input")
	}
}

func TestProcessAndSaveAvatarStoresAndReturnsURL(t *testing.T) {
	store, err := NewFSStore(t.TempDir(), "media")
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	m := &Module{store: store}
	raw := solidPNG(t, 100, 100, color.RGBA{200, 100, 50, 255})
	url, err := m.ProcessAndSaveAvatar(context.Background(), uuid.New(), raw)
	if err != nil {
		t.Fatalf("ProcessAndSaveAvatar: %v", err)
	}
	if url == "" {
		t.Fatal("expected a non-empty avatar URL")
	}
	// The URL must sit under the avatar/ prefix (kept separate from content images).
	if want := "/media/avatar/"; len(url) < len(want) || url[:len(want)] != want {
		t.Errorf("avatar URL %q should start with %q", url, want)
	}
}
