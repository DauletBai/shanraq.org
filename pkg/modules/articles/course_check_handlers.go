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

// maxSolution caps a submission. A beginner's answer to a lesson exercise is
// tens of lines; anything past this is a paste of something else.
const maxSolution = 8000

// checkResponse is what the lesson page gets back.
type checkResponse struct {
	Passed bool   `json:"passed"`
	Note   string `json:"note"`
	Error  string `json:"error,omitempty"`
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

	// The quota is checked before the model is called, not after: the point of
	// it is the bill, and a refusal that has already paid is not a refusal.
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
	reply(http.StatusOK, checkResponse{Passed: v.Passed, Note: v.Note})
}
