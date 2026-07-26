package articles

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/png"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func (a *testApp) userExists(id uuid.UUID) bool {
	var n int
	_ = a.pool.QueryRow(context.Background(), `SELECT count(*) FROM auth_users WHERE id = $1`, id).Scan(&n)
	return n > 0
}

func pngMultipart(t *testing.T, field string, w, h int) (*bytes.Buffer, string) {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, color.RGBA{100, 150, 200, 255})
		}
	}
	var raw bytes.Buffer
	if err := png.Encode(&raw, img); err != nil {
		t.Fatalf("encode png: %v", err)
	}
	var body bytes.Buffer
	mw := multipart.NewWriter(&body)
	fw, err := mw.CreateFormFile(field, "avatar.png")
	if err != nil {
		t.Fatalf("form file: %v", err)
	}
	_, _ = fw.Write(raw.Bytes())
	_ = mw.Close()
	return &body, mw.FormDataContentType()
}

func TestProfilePageAndAccountDeletion(t *testing.T) {
	app := newTestApp(t)
	email := "prof-" + uuid.NewString()[:8] + "@example.com"
	pass := "Str0ng-Pass-123"
	uid := app.createUser(email, pass)
	cookie := app.login(email, pass)

	// Profile page loads and exposes the avatar + delete forms.
	w := app.do(http.MethodGet, "/studio/profile", nil, withCookie(cookie))
	if w.Code != http.StatusOK {
		t.Fatalf("profile GET = %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, `action="/studio/avatar"`) || !strings.Contains(body, `action="/studio/delete"`) {
		t.Error("profile page missing avatar/delete forms")
	}

	// Wrong password must NOT delete the account.
	w = app.do(http.MethodPost, "/studio/delete", url.Values{"password": {"wrong-pass"}}, withCookie(cookie))
	if w.Code != http.StatusSeeOther || !strings.Contains(w.Header().Get("Location"), "del_badpass") {
		t.Fatalf("wrong-pass delete: code=%d loc=%s", w.Code, w.Header().Get("Location"))
	}
	if !app.userExists(uid) {
		t.Fatal("account must survive a wrong-password attempt")
	}

	// Correct password deletes the account and clears the session.
	w = app.do(http.MethodPost, "/studio/delete", url.Values{"password": {pass}}, withCookie(cookie))
	if w.Code != http.StatusSeeOther || !strings.Contains(w.Header().Get("Location"), "account_deleted") {
		t.Fatalf("delete: code=%d loc=%s", w.Code, w.Header().Get("Location"))
	}
	if app.userExists(uid) {
		t.Fatal("account should be gone after confirmed deletion")
	}
}

func TestAvatarUploadAndClear(t *testing.T) {
	app := newTestApp(t)
	email := "ava-" + uuid.NewString()[:8] + "@example.com"
	pass := "Str0ng-Pass-123"
	uid := app.createUser(email, pass)
	cookie := app.login(email, pass)

	// Upload a small PNG as the avatar.
	reqBody, ctype := pngMultipart(t, "avatar", 80, 60)
	r := httptest.NewRequest(http.MethodPost, "/studio/avatar", reqBody)
	r.Host = "localhost:8080"
	r.Header.Set("Origin", app.origin)
	r.Header.Set("Content-Type", ctype)
	r.AddCookie(cookie)
	w := httptest.NewRecorder()
	app.router.ServeHTTP(w, r)
	if w.Code != http.StatusSeeOther || !strings.Contains(w.Header().Get("Location"), "avatar_set") {
		t.Fatalf("avatar upload: code=%d loc=%s", w.Code, w.Header().Get("Location"))
	}
	if got := app.auth.Avatar(context.Background(), uid); got == "" {
		t.Fatal("avatar URL should be set after upload")
	}

	// Clearing the avatar reverts to none.
	w = app.do(http.MethodPost, "/studio/avatar/delete", nil, withCookie(cookie))
	if w.Code != http.StatusSeeOther {
		t.Fatalf("avatar delete: code=%d", w.Code)
	}
	if got := app.auth.Avatar(context.Background(), uid); got != "" {
		t.Fatalf("avatar should be cleared, got %q", got)
	}
}
