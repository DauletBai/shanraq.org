package articles

import (
	"strings"
	"testing"
)

// Reading and listening are counted apart. They do not mean the same thing:
// scrolling away halfway is giving up, stopping a recording halfway may be
// arriving at work, and averaged together each hides the other.
func TestListeningIsCountedApartFromReading(t *testing.T) {
	app := newTestApp(t)
	author := app.createUser("depth@example.com", "Password123!")
	id, slug := app.seedArticle(author, "published")

	post := func(query string) int {
		return app.do("POST", "/read/"+slug+"/progress?"+query, nil).Code
	}
	for _, q := range []string{"d=25", "d=50", "d=50&m=listen", "d=100&m=listen"} {
		if code := post(q); code != 204 {
			t.Fatalf("%s returned %d", q, code)
		}
	}

	got := map[string]int64{}
	rows, err := app.pool.Query(t.Context(),
		`SELECT mode, depth, count FROM reading_depth WHERE article_id = $1`, id)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var mode string
		var depth int
		var n int64
		if err := rows.Scan(&mode, &depth, &n); err != nil {
			t.Fatal(err)
		}
		got[mode+":"+string(rune('0'+depth/25))] = n
	}

	if got["read:1"] != 1 || got["read:2"] != 1 {
		t.Errorf("reading milestones not recorded: %v", got)
	}
	if got["listen:2"] != 1 || got["listen:4"] != 1 {
		t.Errorf("listening milestones not recorded: %v", got)
	}
	// The same milestone reached both ways stays two separate rows.
	if got["read:2"] != 1 || got["listen:2"] != 1 {
		t.Errorf("half-way was merged across modes: %v", got)
	}
}

// The mode parameter arrived after the beacon that sends it. A cached page
// still reporting without it must keep counting as reading rather than start
// dropping its reports.
func TestAnUnknownModeCountsAsReading(t *testing.T) {
	app := newTestApp(t)
	author := app.createUser("depth2@example.com", "Password123!")
	id, slug := app.seedArticle(author, "published")

	if code := app.do("POST", "/read/"+slug+"/progress?d=25&m=whatever", nil).Code; code != 204 {
		t.Fatalf("returned %d", code)
	}
	var mode string
	if err := app.pool.QueryRow(t.Context(),
		`SELECT mode FROM reading_depth WHERE article_id = $1`, id).Scan(&mode); err != nil {
		t.Fatal(err)
	}
	if mode != ModeRead {
		t.Errorf("recorded as %q, expected %q", mode, ModeRead)
	}
}

// The control is served hidden and revealed by the script only once a voice for
// this page's language is found. A button that produces silence is worse than
// no button, so the markup must not arrive visible.
// An article with no recording offers nothing to press.
//
// This used to assert the opposite: the control shipped on every article and
// the script revealed it once it had found a voice in the reader's browser.
// That path is gone -- browsers have no Kazakh voice, so the offer produced
// either silence or Kazakh spoken in Russian -- and a button that produces
// silence is worse than no button. The control is now rendered only where audio
// exists, which is why its absence here is the passing case.
func TestTheListenControlIsAbsentWithoutARecording(t *testing.T) {
	app := newTestApp(t)
	author := app.createUser("depth3@example.com", "Password123!")
	_, slug := app.seedArticle(author, "published")

	body := app.do("GET", "/read/"+slug+"?lang=ru", nil).Body.String()
	if strings.Contains(body, `id="listen"`) {
		t.Error("an article with no recording still offers a listen control")
	}
	if strings.Contains(body, `id="listen-audio"`) {
		t.Error("an article with no recording still carries an audio element")
	}
	if !strings.Contains(body, "listen.js") {
		t.Error("the article page does not load the reader script")
	}
}

// With a recording, the control appears and carries what the player needs.
//
// The pair with the test above is the contract: no audio, no button; audio, a
// button plus the cue map that lets the page follow along. Nothing else on the
// page decides this, so if the template ever renders the control unconditionally
// again, one of these two fails.
func TestTheListenControlAppearsWithARecording(t *testing.T) {
	app := newTestApp(t)
	author := app.createUser("depth4@example.com", "Password123!")
	id, slug := app.seedArticle(author, "published")

	cues := `[{"i":0,"a":0,"b":3.5},{"i":1,"a":3.5,"b":7.25}]`
	if _, err := app.pool.Exec(t.Context(),
		`INSERT INTO article_audio
		   (article_id, lang, storage_key, url, duration_sec, bytes, voice, text_sha256, cues)
		 VALUES ($1,'ru','audio/x.ogg','/media/audio/x.ogg',7,1234,'ru_RU-dmitri-medium.onnx','',$2::jsonb)`,
		id, cues); err != nil {
		t.Fatalf("seed narration: %v", err)
	}

	body := app.do("GET", "/read/"+slug+"?lang=ru", nil).Body.String()
	if !strings.Contains(body, `id="listen-audio"`) {
		t.Fatal("the recording is stored but the page carries no audio element")
	}
	if !strings.Contains(body, "/media/audio/x.ogg") {
		t.Error("the audio element does not point at the stored recording")
	}
	if !strings.Contains(body, `id="listen-play"`) {
		t.Error("the recording is present but there is nothing to press")
	}
	// Without cues the audio still plays; the highlight is what stops working,
	// and that is the part a reader notices as broken rather than absent.
	if !strings.Contains(body, "data-cues=") {
		t.Error("the cue map did not reach the page")
	}
}
