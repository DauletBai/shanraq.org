package articles

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"shanraq.org/pkg/modules/jobs"
	"shanraq.org/pkg/shanraq"
)

// JobNarrate is the queue name for recording one article in one language.
const JobNarrate = "article_narrate"

type narratePayload struct {
	ArticleID uuid.UUID `json:"article_id"`
	Lang      string    `json:"lang"`
}

// ttsClient talks to the synthesiser container.
//
// It is a container rather than a library because the synthesiser is Python
// with an ONNX runtime and three voice models, and the site is a small static
// Go binary. Keeping them apart means a deploy does not carry 240 MB of voices
// that change once a year, and a broken model cannot take the site down with
// it -- when this call fails, the article simply has no recording yet.
type ttsClient struct {
	base string
	http *http.Client
}

type ttsResult struct {
	Audio    []byte
	Cues     []byte
	Seconds  int
	Voice    string
	MIMEType string
}

func newTTSClient(base string) *ttsClient {
	return &ttsClient{
		base: strings.TrimRight(base, "/"),
		// Long, because it is meant to be. Kazakh synthesises at about twice
		// real time, so a twenty-five minute article is eleven minutes of work.
		// This is a queued background job; nobody is waiting on the socket.
		http: &http.Client{Timeout: 45 * time.Minute},
	}
}

func (c *ttsClient) synthesize(ctx context.Context, lang string, blocks []string) (*ttsResult, error) {
	body, err := json.Marshal(map[string]any{"lang": lang, "blocks": blocks})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.base+"/synthesize", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("tts unreachable: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		msg, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("tts %d: %s", resp.StatusCode, strings.TrimSpace(string(msg)))
	}
	audio, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	secs, _ := strconv.Atoi(resp.Header.Get("X-Audio-Seconds"))
	return &ttsResult{
		Audio:    audio,
		Cues:     []byte(resp.Header.Get("X-Audio-Cues")),
		Seconds:  secs,
		Voice:    resp.Header.Get("X-Audio-Voice"),
		MIMEType: resp.Header.Get("Content-Type"),
	}, nil
}

// enqueueNarration files an article for recording, in every language it has.
//
// Failures here are logged and dropped. An article that could not be queued is
// an article without audio, which is where every article started; losing the
// publish over it would be the worse trade.
func (m *Module) enqueueNarration(ctx context.Context, a *Article) {
	if m.jobs == nil || m.tts == nil {
		return
	}
	for _, lang := range a.AvailableLangs() {
		payload, err := json.Marshal(narratePayload{ArticleID: a.ID, Lang: lang})
		if err != nil {
			continue
		}
		if err := m.jobs.Enqueue(ctx, jobs.Job{
			ID:      uuid.New(),
			UserID:  a.AuthorID,
			Name:    JobNarrate,
			Payload: payload,
			RunAt:   time.Now(),
			// Three attempts, because the failures worth retrying are transient:
			// the synthesiser still starting, a container restarted mid-job.
			MaxAttempts: 3,
		}); err != nil {
			m.rt.Logger.Warn("enqueue narration", zap.String("lang", lang), zap.Error(err))
		}
	}
}

// handleNarrateJob records one article in one language and stores the result.
func (m *Module) handleNarrateJob(ctx context.Context, _ *shanraq.Runtime, job jobs.Job) error {
	var p narratePayload
	if err := json.Unmarshal(job.Payload, &p); err != nil {
		return fmt.Errorf("narrate payload: %w", err)
	}
	if m.tts == nil {
		return nil
	}

	a, err := m.store.GetByID(ctx, p.ArticleID, uuid.Nil)
	if err != nil {
		return fmt.Errorf("narrate load article: %w", err)
	}
	tr, ok := a.Translations[p.Lang]
	if !ok || strings.TrimSpace(tr.BodyMD) == "" {
		return nil // nothing in this language to read
	}

	// The digest is taken before the work, and checked against what is already
	// stored: an article re-published without a text change does not need
	// eleven minutes of Kazakh synthesis to arrive at the same recording.
	digest := TextDigest(tr.Title, tr.BodyMD)
	if have, err := m.audio.Get(ctx, a.ID, p.Lang, digest); err == nil && have != nil && !have.Stale {
		return nil
	}

	rendered, _ := RenderMarkdownTOC(tr.BodyMD)
	blocks := NarrationBlocks(string(rendered))
	if len(blocks) == 0 {
		return nil
	}

	res, err := m.tts.synthesize(ctx, p.Lang, blocks)
	if err != nil {
		return err
	}

	key := path.Join("audio", a.ID.String(), fmt.Sprintf("%s-%d%s", p.Lang, time.Now().Unix(), extFor(res.MIMEType)))
	url, err := m.media.SaveBlob(ctx, key, res.Audio, res.MIMEType)
	if err != nil {
		return fmt.Errorf("narrate store: %w", err)
	}

	replaced, err := m.audio.Upsert(ctx, a.ID, Narration{
		Lang: p.Lang, URL: url, StorageKey: key,
		DurationSec: res.Seconds, Bytes: int64(len(res.Audio)),
		Voice: res.Voice, TextSHA256: digest, Cues: res.Cues,
	})
	if err != nil {
		// The row is what makes the file reachable; without it the file is
		// litter from the moment it was written.
		_ = m.media.DeleteBlob(ctx, key)
		return err
	}
	if replaced != "" {
		_ = m.media.DeleteBlob(ctx, replaced)
	}

	m.rt.Logger.Info("article narrated",
		zap.String("slug", a.Slug), zap.String("lang", p.Lang),
		zap.Int("seconds", res.Seconds), zap.Int("bytes", len(res.Audio)))
	return nil
}

func extFor(mime string) string {
	switch {
	case strings.Contains(mime, "ogg"):
		return ".ogg"
	case strings.Contains(mime, "mp4"), strings.Contains(mime, "aac"):
		return ".m4a"
	case strings.Contains(mime, "mpeg"):
		return ".mp3"
	default:
		return ".wav"
	}
}

// narrateAfterPublish queues a reading of every language the article carries.
//
// This is the whole point of moving synthesis onto the server: an author
// publishes and the recording appears, in all three languages, without anyone
// touching a voice model or uploading a file. Nothing about publishing waits
// for it -- the job runs behind, and until it finishes the article simply has
// no listen button.
func (m *Module) narrateAfterPublish(ctx context.Context, id uuid.UUID) {
	if m.tts == nil || m.jobs == nil {
		return
	}
	a, err := m.store.GetByID(ctx, id, uuid.Nil)
	if err != nil {
		m.rt.Logger.Warn("narrate after publish", zap.Error(err))
		return
	}
	m.enqueueNarration(ctx, a)
}
