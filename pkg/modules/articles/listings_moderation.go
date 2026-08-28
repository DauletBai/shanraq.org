package articles

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"shanraq.org/pkg/modules/ai"
	"shanraq.org/pkg/modules/jobs"
	"shanraq.org/pkg/shanraq"
)

// JobModerateListing is the queue job that screens a newly posted listing.
const JobModerateListing = "ai_moderate_listing"

// The screening runs in the background rather than in the submission handler.
//
// The account this platform runs against allows three requests a minute, so a
// synchronous check would hold the form for up to twenty seconds on a busy
// evening — and the whole point of a classifieds site is that posting is quick.
// The listing appears at once and is read a moment later; anything the model
// objects to is flagged, which takes it out of search until a person looks.
//
// Nothing is deleted and nothing is refused at the door. A listing the model
// dislikes goes to 'flagged', which is where readers' reports already send
// listings — one queue, one appeal path, whoever raised the objection.

type moderateListingPayload struct {
	ListingID string `json:"listing_id"`
}

// RegisterJobs attaches the listing screening handler to the job queue.
func (m *Module) RegisterJobs(j *jobs.Module) {
	j.Handle(JobModerateListing, m.handleModerateListingJob)
}

// enqueueListingScreening files a listing for background screening. Failures are
// logged and swallowed: a listing that could not be queued is a listing nobody
// screened, which is the situation the site was in until now — not a reason to
// lose the author's work.
func (m *Module) enqueueListingScreening(ctx context.Context, listingID, authorID uuid.UUID) {
	if m.ai == nil || !m.ai.Enabled() {
		return
	}
	payload, err := json.Marshal(moderateListingPayload{ListingID: listingID.String()})
	if err != nil {
		m.rt.Logger.Warn("encode listing screening payload", zap.Error(err))
		return
	}
	if err := m.jobs.Enqueue(ctx, jobs.Job{
		ID:          uuid.New(),
		UserID:      authorID,
		Name:        JobModerateListing,
		Payload:     payload,
		RunAt:       time.Now(),
		MaxAttempts: 3,
	}); err != nil {
		m.rt.Logger.Warn("enqueue listing screening", zap.Error(err))
	}
}

func (m *Module) handleModerateListingJob(ctx context.Context, _ *shanraq.Runtime, job jobs.Job) error {
	if m.ai == nil || !m.ai.Enabled() {
		return nil
	}
	var payload moderateListingPayload
	if err := job.Decode(&payload); err != nil {
		return fmt.Errorf("decode payload: %w", err)
	}
	id, err := uuid.Parse(payload.ListingID)
	if err != nil {
		return fmt.Errorf("bad listing id: %w", err)
	}

	text, status, title, authorID, err := m.listingScreeningText(ctx, id)
	if err != nil {
		return err
	}
	// Already withdrawn, expired or flagged by readers: nothing to decide.
	if status != "published" {
		return nil
	}

	raw, err := m.ai.Check(ctx,
		"You are a strict but fair listings checker. Return JSON only.",
		listingReviewPrompt(listingNoteLang, text), 1500)
	if err != nil {
		if errors.Is(err, ai.ErrDisabled) {
			return nil
		}
		// A screening that could not run is not a verdict. Let the queue retry;
		// the listing stays visible meanwhile, which is where every listing was
		// before this existed.
		return fmt.Errorf("screen listing: %w", err)
	}
	findings, err := parseVerdict(raw)
	if err != nil {
		// An unreadable answer is not a pass, but it is not a reason to hide
		// somebody's flat either. Retry; if it keeps failing the queue says so.
		return fmt.Errorf("parse listing verdict: %w", err)
	}
	if len(findings) == 0 {
		return nil
	}

	blocking := false
	for _, f := range findings {
		if f.Blocking() {
			blocking = true
			break
		}
	}

	action := "warn"
	if blocking {
		action = "hidden"
		if _, err := m.rt.DB.Exec(ctx,
			`UPDATE listings SET status = 'flagged' WHERE id = $1 AND status = 'published'`, id); err != nil {
			return fmt.Errorf("flag listing: %w", err)
		}
	}
	// Whatever was decided, the author is told what and why — every objection,
	// with the words it objects to. A listing that vanishes from search without
	// a reason leaves its author guessing, and guessing is not moderation.
	m.logListingReview(ctx, id, authorID, title, action, findings)
	m.rt.Logger.Info("ai screened a listing",
		zap.String("listing_id", id.String()),
		zap.String("action", action),
		zap.Int("findings", len(findings)))
	return nil
}

// listingNoteLang is the language the checker writes its notes in.
//
// Listings carry no language of their own and accounts have no preferred one,
// so there is nothing to read it from. Russian is what almost every listing on
// this site is written in, and a note in the wrong language is still better
// than the silence this replaces.
const listingNoteLang = "ru"

// logListingReview files the decision and every finding behind it into the same
// ledger article moderation uses, so the author reads it on the same page, with
// the same quotes and the same appeal form.
func (m *Module) logListingReview(ctx context.Context, id, author uuid.UUID, title, action string, findings []Finding) {
	// The ledger's reason code is the coarse one the page has a translation
	// for; which rules exactly were broken is in the findings underneath, each
	// with the words it objects to.
	actionID, err := m.mods.Log(ctx, ModAction{
		TargetType: "listing", TargetID: id.String(), Title: clip(title, 120),
		Action: action, ReasonCode: "rules_failed", ActorKind: "agent", ActorName: "AI Bake",
	}, &author, nil)
	if err != nil {
		m.rt.Logger.Error("log listing review", zap.Error(err))
		return
	}
	aid, err := uuid.Parse(actionID)
	if err != nil {
		return
	}
	if err := m.mods.SaveFindings(ctx, aid, findings); err != nil {
		m.rt.Logger.Error("save listing findings", zap.Error(err))
	}
}

// listingReviewPrompt asks for the same contract article review uses: a list of
// findings, each quoting the words it objects to and saying what to change.
//
// A one-line verdict was the first design and it was wrong. "Hidden: hidden
// advertising" tells an author nothing they can act on — they are left to guess
// which sentence offended, and a guess is not a correction. Every objection
// must name itself, quote the text, and say what to do.
func listingReviewPrompt(lang, listing string) string {
	var b strings.Builder
	b.WriteString("You are the listings checker for a Kazakhstani classifieds platform. ")
	b.WriteString("Text may be in Kazakh, Russian or English. Check the listing below and return JSON only.\n\n")
	b.WriteString("Rules, by code:\n")
	b.WriteString("- misleading: the description contradicts the declared fields — floor area, rooms, price or city — ")
	b.WriteString("or an agency presents itself as the owner\n")
	b.WriteString("- disguised_ad: promotion of some other business or service, an affiliate or referral link, ")
	b.WriteString("or a \"review\" whose purpose is to sell something other than this property\n")
	b.WriteString("- scam: fraud, fake earnings schemes, phishing, forged documents or certificates offered\n")
	b.WriteString("- spam: bulk promotion, link dumping, repeated advertising junk\n")
	b.WriteString("- prohibited_goods: goods or services that may not lawfully be advertised in Kazakhstan\n")
	b.WriteString("- personal_data: another person's private data published without their consent\n")
	b.WriteString("- illegal: other content plainly unlawful under the law of the Republic of Kazakhstan\n\n")
	b.WriteString("Severity: \"block\" takes the listing out of search until a person reviews it; ")
	b.WriteString("\"warn\" is advice that changes nothing. Use \"block\" for misleading, scam, ")
	b.WriteString("prohibited_goods, personal_data and illegal; use \"warn\" for disguised_ad and spam ")
	b.WriteString("unless the whole listing is an advertisement for something else.\n\n")
	b.WriteString("For every finding you MUST quote the exact words you object to in \"quote\". ")
	b.WriteString("A finding the author cannot locate is useless. Write \"note\" in ")
	b.WriteString(lang)
	b.WriteString(", addressed to the author, saying plainly what to change.\n\n")
	b.WriteString("Do not invent problems. A listing that breaks no rule must return an empty findings array. ")
	b.WriteString("A plain, dull, badly written or overpriced listing breaks no rule. ")
	b.WriteString("The seller's own contact details are not personal_data.\n\n")
	b.WriteString("Return exactly: {\"findings\":[{\"rule\":\"...\",\"severity\":\"block|warn\",\"quote\":\"...\",\"note\":\"...\"}]}\n\n")
	b.WriteString(listing)
	return b.String()
}

// listingScreeningText assembles what the model is asked to judge.
//
// The declared figures travel with the prose deliberately. Half of what makes a
// listing misleading is not in the description on its own but in the gap
// between it and the fields: a description promising a hundred and twenty
// square metres over an area of fifty-four is the complaint this site expects
// most, and the model cannot see it unless both halves are in front of it.
func (m *Module) listingScreeningText(ctx context.Context, id uuid.UUID) (text, status, title string, author uuid.UUID, err error) {
	var description, city string
	var area float64
	var price int64
	var rooms int
	err = m.rt.DB.QueryRow(ctx, `
		SELECT status, author_id, title, description, COALESCE(city, ''), COALESCE(area, 0),
		       COALESCE(price, 0), COALESCE(rooms, 0)
		FROM listings WHERE id = $1`, id).
		Scan(&status, &author, &title, &description, &city, &area, &price, &rooms)
	if err != nil {
		return "", "", "", uuid.Nil, fmt.Errorf("load listing: %w", err)
	}

	var b strings.Builder
	b.WriteString("Declared fields:\n")
	fmt.Fprintf(&b, "- area: %s m²\n", strconv.FormatFloat(area, 'f', -1, 64))
	fmt.Fprintf(&b, "- rooms: %d\n", rooms)
	fmt.Fprintf(&b, "- price: %d\n", price)
	if city != "" {
		fmt.Fprintf(&b, "- city: %s\n", city)
	}
	b.WriteString("\nTitle: ")
	b.WriteString(title)
	b.WriteString("\n\nDescription:\n")
	b.WriteString(clip(description, 4000))
	return b.String(), status, title, author, nil
}
