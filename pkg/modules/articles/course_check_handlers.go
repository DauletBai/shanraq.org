package articles

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// checkWindow and checkQuota bound what one reader can spend of somebody else's
// money in an hour. Generous enough that nobody working through a lesson meets
// it, low enough that a script cannot run up a bill overnight.
const (
	checkWindow = time.Hour
	checkQuota  = 30
)

// maxAttempts is how many times one reader may have one exercise checked.
//
// Three, and then the machine stops answering. Not to save money — three checks
// cost less than half a cent — but because a fourth is where the exercise turns
// into a guessing game against a reviewer, and guessing is the opposite of what
// the lesson is for. What comes after the third is deliberately not another
// model call: re-read the lesson, or ask a person in the comments.
const maxAttempts = 3

// maxSolution caps a submission. A beginner's answer to a lesson exercise is
// tens of lines; anything past this is a paste of something else.
const maxSolution = 8000

// checkResponse is what the lesson page gets back.
type checkResponse struct {
	Passed bool   `json:"passed"`
	Note   string `json:"note"`
	// Code is the submission after gofmt, and HTML the same thing coloured. The
	// reader gets their own text back tidied: the layout rules are learned by
	// watching them applied, not by being told about them.
	Code string `json:"code,omitempty"`
	HTML string `json:"html,omitempty"`
	// Syntax is a parse error. It is answered separately from a review because
	// it is not one, and because it must not cost an attempt.
	Syntax string `json:"syntax,omitempty"`
	// Left is how many checks of this exercise remain. The reader is told before
	// spending one, not after: a limit nobody can see is a trap.
	Left  int    `json:"left"`
	Error string `json:"error,omitempty"`
}

// handleCourseCheck reviews a reader's solution to a lesson's exercise.
//
// Answers JSON because the page asks without reloading — a reader who has just
// typed thirty lines into a textarea must not lose them to a redirect.
func (m *Module) handleCourseCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	lang := m.resolveLang(w, r)
	reply := func(code int, res checkResponse) {
		w.WriteHeader(code)
		_ = json.NewEncoder(w).Encode(res)
	}

	user, ok := m.authorID(r)
	if !ok {
		reply(http.StatusUnauthorized, checkResponse{Error: T(lang, "chk.login")})
		return
	}
	if m.ai == nil || !m.ai.Enabled() {
		reply(http.StatusServiceUnavailable, checkResponse{Error: T(lang, "chk.off")})
		return
	}

	a, err := m.store.GetPublishedBySlug(r.Context(), chi.URLParam(r, "slug"))
	if err != nil {
		reply(http.StatusNotFound, checkResponse{Error: T(lang, "chk.no_lesson")})
		return
	}

	if err := r.ParseForm(); err != nil {
		reply(http.StatusBadRequest, checkResponse{Error: T(lang, "chk.bad")})
		return
	}
	solution := unfence(r.FormValue("solution"))
	if solution == "" {
		reply(http.StatusBadRequest, checkResponse{Error: T(lang, "chk.empty")})
		return
	}
	if len(solution) > maxSolution {
		reply(http.StatusRequestEntityTooLarge, checkResponse{Error: T(lang, "chk.too_long")})
		return
	}

	// Tidy first, and refuse code that will not parse before anything is spent
	// on it. A missing brace is not something to ask a reviewer about, and the
	// reader's own editor would have said so — which is the lesson here.
	formatted, ferr := formatSolution(solution)
	if ferr != nil {
		reply(http.StatusUnprocessableEntity, checkResponse{
			Syntax: syntaxHint(ferr),
			Error:  T(lang, "chk.syntax"),
		})
		return
	}
	solution = formatted

	// Both limits are read before the model is called, not after: the point of a
	// limit is what it prevents, and a refusal that has already paid prevents
	// nothing.
	pr, _ := m.progress.Get(r.Context(), user, a.ID)
	if pr.Passed {
		reply(http.StatusOK, checkResponse{Passed: true, Note: pr.Note, Left: 0})
		return
	}
	if pr.Attempts >= maxAttempts {
		reply(http.StatusTooManyRequests, checkResponse{Error: T(lang, "chk.spent"), Left: 0})
		return
	}
	if n, err := m.progress.AttemptsSince(r.Context(), user, checkWindow); err == nil && n >= checkQuota {
		reply(http.StatusTooManyRequests, checkResponse{Error: T(lang, "chk.quota")})
		return
	}

	// The exercise is read out of the lesson the reader actually saw, in the
	// language they are reading it in.
	tr, served := a.Translation(lang)
	if tr == nil {
		reply(http.StatusNotFound, checkResponse{Error: T(lang, "chk.no_lesson")})
		return
	}
	task := lessonExercise(tr.BodyMD)
	if task == "" {
		reply(http.StatusNotFound, checkResponse{Error: T(lang, "chk.no_task")})
		return
	}

	var b strings.Builder
	b.WriteString("LESSON: ")
	b.WriteString(tr.Title)
	b.WriteString("\n\nEXERCISE:\n")
	b.WriteString(task)
	b.WriteString("\n\nSOLUTION:\n")
	b.WriteString(solution)

	raw, err := m.ai.Check(r.Context(), checkSystem(served), b.String(), 700)
	if err != nil {
		m.rt.Logger.Warn("course check", zap.Error(err))
		reply(http.StatusBadGateway, checkResponse{Error: T(lang, "chk.failed")})
		return
	}
	v, err := parseCheckVerdict(raw)
	if err != nil {
		m.rt.Logger.Warn("course check verdict", zap.Error(err), zap.String("raw", clip(raw, 200)))
		reply(http.StatusBadGateway, checkResponse{Error: T(lang, "chk.failed")})
		return
	}

	if err := m.progress.Record(r.Context(), user, a.ID, v, solution); err != nil {
		// The review is what the reader came for; losing the bookkeeping is not
		// a reason to withhold it.
		m.rt.Logger.Warn("course progress record", zap.Error(err))
	}
	left := maxAttempts - (pr.Attempts + 1)
	if v.Passed || left < 0 {
		left = 0
	}
	reply(http.StatusOK, checkResponse{
		Passed: v.Passed, Note: v.Note, Left: left,
		Code: solution, HTML: string(highlightGo(solution)),
	})
}

// handleCourseFormat tidies a submission and hands it back coloured, without
// asking the reviewer anything.
//
// Separate from the check on purpose: formatting is local, instant and free, so
// a reader may lean on it as often as they like. Tying it to the three attempts
// would have taught them to avoid the very tool the lesson wants them using.
func (m *Module) handleCourseFormat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	lang := m.resolveLang(w, r)
	reply := func(code int, res checkResponse) {
		w.WriteHeader(code)
		_ = json.NewEncoder(w).Encode(res)
	}
	if _, ok := m.authorID(r); !ok {
		reply(http.StatusUnauthorized, checkResponse{Error: T(lang, "chk.login")})
		return
	}
	if err := r.ParseForm(); err != nil {
		reply(http.StatusBadRequest, checkResponse{Error: T(lang, "chk.bad")})
		return
	}
	src := unfence(r.FormValue("solution"))
	if src == "" {
		reply(http.StatusBadRequest, checkResponse{Error: T(lang, "chk.empty")})
		return
	}
	if len(src) > maxSolution {
		reply(http.StatusRequestEntityTooLarge, checkResponse{Error: T(lang, "chk.too_long")})
		return
	}
	out, err := formatSolution(src)
	if err != nil {
		reply(http.StatusUnprocessableEntity, checkResponse{
			Syntax: syntaxHint(err),
			Error:  T(lang, "chk.syntax"),
		})
		return
	}
	reply(http.StatusOK, checkResponse{Code: out, HTML: string(highlightGo(out))})
}
