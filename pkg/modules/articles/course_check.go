package articles

import (
	"context"
	"encoding/json"
	"fmt"
	"go/format"
	"html/template"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Checking a lesson's exercise.
//
// The course page once promised progress marks with nothing behind them. This
// is what stands behind them now, and the mark is earned rather than scrolled
// to: a reader submits a solution, a model reads it against the exercise the
// lesson actually set, and only an accepted answer counts the lesson passed.
//
// The model is asked to review, not to grade. A learner who is told "wrong"
// learns nothing; one who is told which line does not do what they think can
// fix it themselves, which is the whole point of an exercise.

// CheckVerdict is one review of one submission.
type CheckVerdict struct {
	Passed bool   `json:"passed"`
	Note   string `json:"note"`
}

// Progress is what a reader has done with one lesson.
type Progress struct {
	Passed   bool
	Attempts int
	Note     string
	Solution string
	Updated  time.Time
}

// ProgressStore keeps per-reader lesson progress.
type ProgressStore struct{ db *pgxpool.Pool }

// NewProgressStore wires the store to the pool.
func NewProgressStore(db *pgxpool.Pool) *ProgressStore { return &ProgressStore{db: db} }

// Get returns what this reader has done with this lesson; a zero Progress when
// they have not tried it.
func (st *ProgressStore) Get(ctx context.Context, user, article uuid.UUID) (Progress, error) {
	var p Progress
	err := st.db.QueryRow(ctx, `
		SELECT passed, attempts, note, solution, updated_at
		FROM course_progress WHERE user_id = $1 AND article_id = $2`, user, article).
		Scan(&p.Passed, &p.Attempts, &p.Note, &p.Solution, &p.Updated)
	if err != nil {
		return Progress{}, nil
	}
	return p, nil
}

// Record stores an attempt. A lesson once passed stays passed: a later
// experiment that does not compile must not take the mark away.
func (st *ProgressStore) Record(ctx context.Context, user, article uuid.UUID, v CheckVerdict, solution string) error {
	_, err := st.db.Exec(ctx, `
		INSERT INTO course_progress (user_id, article_id, passed, attempts, note, solution, updated_at)
		VALUES ($1, $2, $3, 1, $4, $5, now())
		ON CONFLICT (user_id, article_id) DO UPDATE SET
			passed     = course_progress.passed OR EXCLUDED.passed,
			attempts   = course_progress.attempts + 1,
			note       = EXCLUDED.note,
			solution   = EXCLUDED.solution,
			updated_at = now()`,
		user, article, v.Passed, v.Note, solution)
	return err
}

// PassedIn returns the article ids of a course this reader has passed.
func (st *ProgressStore) PassedIn(ctx context.Context, user uuid.UUID, ids []uuid.UUID) (map[uuid.UUID]bool, error) {
	out := map[uuid.UUID]bool{}
	if len(ids) == 0 {
		return out, nil
	}
	rows, err := st.db.Query(ctx,
		`SELECT article_id FROM course_progress
		 WHERE user_id = $1 AND passed AND article_id = ANY($2)`, user, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		out[id] = true
	}
	return out, rows.Err()
}

// AttemptsSince counts this reader's submissions in the given window, which is
// what bounds the cost of the feature.
func (st *ProgressStore) AttemptsSince(ctx context.Context, user uuid.UUID, since time.Duration) (int, error) {
	var n int
	err := st.db.QueryRow(ctx,
		`SELECT COALESCE(SUM(attempts), 0) FROM course_progress
		 WHERE user_id = $1 AND updated_at > now() - $2::interval`,
		user, fmt.Sprintf("%d seconds", int(since.Seconds()))).Scan(&n)
	return n, err
}

// exerciseHeads are the lesson headings that introduce the exercise, in the
// three languages the course is written in.
var exerciseHeads = []string{"## Задание", "## Тапсырма", "## Exercise"}

// optionalHeads mark where the exercise stops being required. Everything after
// one of them is offered, not asked for -- and a reviewer that reads it can
// fail a reader for skipping work the lesson itself called voluntary, spending
// one of three attempts on it. All three languages are searched whatever the
// lesson's own language, because the cost of a stray match is nil and the cost
// of a missed one is a wrong verdict.
var optionalHeads = []string{"**По желанию", "**Қалауыңызша", "**Optional"}

// lessonExercise pulls the exercise out of a lesson's own text.
//
// The task is not stored separately on purpose: a second copy would drift from
// the lesson the moment either was edited, and the reader would be marked
// against something they never read.
func lessonExercise(body string) string {
	for _, head := range exerciseHeads {
		i := strings.Index(body, head)
		if i < 0 {
			continue
		}
		rest := body[i+len(head):]
		if j := strings.Index(rest, "\n## "); j >= 0 {
			rest = rest[:j]
		}
		for _, opt := range optionalHeads {
			if j := strings.Index(rest, opt); j >= 0 {
				rest = rest[:j]
			}
		}
		return strings.TrimSpace(rest)
	}
	return ""
}

// fenced strips a Markdown code fence a reader may have pasted around their
// answer, so the model reads code rather than a code block.
var fenced = regexp.MustCompile("(?s)^\\s*```[a-zA-Z]*\\n(.*?)\\n?```\\s*$")

func unfence(s string) string {
	if m := fenced.FindStringSubmatch(s); m != nil {
		return strings.TrimSpace(m[1])
	}
	return strings.TrimSpace(s)
}

// checkSystem is the reviewer's brief. It is deliberately not a grader's: the
// reader is a beginner who has met perhaps five lessons' worth of the language,
// and a verdict without a way forward wastes the one moment they are paying
// attention.
func checkSystem(lang string) string {
	common := `You review a beginner's solution to one exercise from a Go course.

Answer with JSON and nothing else: {"passed": true|false, "note": "..."}.

"passed" is true when the solution does what the exercise asked. Judge that and
nothing else. Do not fail a solution for style, for a name you would have chosen
differently, for missing error handling the exercise never asked for, or for not
using a feature the course has not taught yet.

"note" is two to four sentences addressed to the learner, in their language.
When the solution works, say briefly what it does right and, if there is one,
name a single thing worth knowing next. When it does not, point at the specific
line or expression that does not do what they think and say what it does
instead — never simply "wrong", and never the corrected code, which would take
the exercise away from them.

Plain text in "note", no Markdown, no code fences.`
	switch lang {
	case LangKZ:
		return common + "\n\nWrite \"note\" in Kazakh."
	case LangEN:
		return common + "\n\nWrite \"note\" in English."
	default:
		return common + "\n\nWrite \"note\" in Russian."
	}
}

// parseCheckVerdict reads the model's answer, tolerating the fence it sometimes puts
// around JSON despite being asked not to.
func parseCheckVerdict(raw string) (CheckVerdict, error) {
	s := unfence(raw)
	if i := strings.Index(s, "{"); i > 0 {
		s = s[i:]
	}
	if j := strings.LastIndex(s, "}"); j >= 0 {
		s = s[:j+1]
	}
	var v CheckVerdict
	if err := json.Unmarshal([]byte(s), &v); err != nil {
		return CheckVerdict{}, fmt.Errorf("verdict: %w", err)
	}
	v.Note = strings.TrimSpace(v.Note)
	if v.Note == "" {
		return CheckVerdict{}, fmt.Errorf("verdict: empty note")
	}
	return v, nil
}

// formatSolution runs the reader's code through gofmt.
//
// The standard library's own formatter, not an imitation: the reader sees
// exactly what their editor does to the same text, which is the point. A
// beginner has not learned the indentation rules yet and should not be spending
// attention on them while learning what a function is.
//
// It doubles as a free syntax check. Code that will not parse cannot be
// reviewed, and refusing it here costs nothing — the parser is local, the
// reviewer is a paid call — so a missing brace is answered instantly instead of
// spending one of three attempts on a verdict the reader could have got from
// their own editor.
func formatSolution(src string) (string, error) {
	out, err := format.Source([]byte(src))
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// syntaxHint turns a parser error into something a beginner can act on. Go's
// own message names the file it never had; the line and the complaint are what
// help.
func syntaxHint(err error) string {
	s := err.Error()
	if i := strings.Index(s, ".go:"); i >= 0 {
		s = s[i+4:]
	}
	return strings.TrimSpace(s)
}

// highlightGo renders code the way the lessons render theirs — through the same
// Markdown pipeline, so the colours in the box and the colours in the lesson
// are the same colours, and follow the same light/dark switch.
//
// The markup is produced by our own renderer from the reader's own text, and
// Chroma escapes every token it emits; it is also only ever shown back to the
// person who typed it.
func highlightGo(code string) template.HTML {
	return RenderMarkdown("```go\n" + code + "\n```")
}
