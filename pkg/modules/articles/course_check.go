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

// The languages a course's exercises are written in. A lesson belongs to one
// course, the course names one language, and every part of the check follows
// it: what tidies the code, what refuses it before the model is called, what
// colours it, and what the reviewer is told it is reading.
const (
	CodeGo     = "go"
	CodePython = "python"
	CodeSQL    = "sql"
	CodeRust   = "rust"
	CodeShell  = "shell"
	// CodeInformatics accepts observations, models, explanations, and project
	// evidence before the course reaches programming.
	CodeInformatics = "informatics"
	// CodeKazakh accepts language production and reflection as prose. Its
	// reviewer checks communicative meaning and forms without pretending that a
	// written answer proves pronunciation or listening.
	CodeKazakh = "kazakh"
	// CodeMath accepts a learner's reasoning in ordinary mathematical prose.
	// A mathematics course cannot pretend that every proof is a program merely
	// to reuse the exercise checker.
	CodeMath = "math"
)

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

	// Appealed records that this reader has already had an attempt back on this
	// exercise. Once per lesson: enough to undo a wrong verdict, not enough to
	// turn three checks into an unlimited supply.
	Appealed bool
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
		SELECT passed, attempts, note, solution, updated_at, appealed
		FROM course_progress WHERE user_id = $1 AND article_id = $2`, user, article).
		Scan(&p.Passed, &p.Attempts, &p.Note, &p.Solution, &p.Updated, &p.Appealed)
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

// Appeal gives one attempt back on a lesson whose verdict the reader says was
// wrong, and returns how many they have left. The false result means there was
// nothing to give: the lesson is already passed, no attempt has been spent, or
// this reader has appealed here before.
func (st *ProgressStore) Appeal(ctx context.Context, user, article uuid.UUID) (int, bool, error) {
	var attempts int
	err := st.db.QueryRow(ctx, `
		UPDATE course_progress
		   SET attempts = attempts - 1, appealed = true, updated_at = now()
		 WHERE user_id = $1 AND article_id = $2
		   AND NOT appealed AND NOT passed AND attempts > 0
		RETURNING attempts`, user, article).Scan(&attempts)
	if err != nil {
		return 0, false, nil
	}
	return attempts, true, nil
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
var exerciseHeads = []string{"## Задание", "## Тапсырма", "## Exercise",
	// An English lesson once headed it "The exercise", and the check box quietly
	// did not appear on it. The heading is uniform in the lessons now; this line
	// is what keeps a single stray heading from taking the feature away again.
	"## The exercise"}

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
func checkSystem(lang, codeLang string) string {
	if codeLang == CodeMath {
		return mathCheckSystem(lang)
	}
	if codeLang == CodeInformatics {
		return informaticsCheckSystem(lang)
	}
	if codeLang == CodeKazakh {
		return kazakhCheckSystem(lang)
	}
	name := "Go"
	switch codeLang {
	case CodePython:
		name = "Python"
	case CodeSQL:
		name = "SQL (SQLite dialect)"
	case CodeRust:
		name = "Rust"
	case CodeShell:
		name = "Shell (POSIX sh)"
	}
	common := `You review a beginner's solution to one exercise from a ` + name + ` course.

The solution is written in ` + name + `. Judge it as ` + name + `, and if it
would not run at all, name the line and what is wrong with it instead of
guessing at the intent.

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

// kazakhCheckSystem reviews what the current lesson actually taught. A text
// box can verify meaning, structure, and a learner's account of what they
// heard; it cannot honestly certify the sound of a voice it never received.
func kazakhCheckSystem(lang string) string {
	common := `You are the checking tutor for Shanraq's beginner Kazakh-language course.
The course teaches communication through a familiar situation, a short Kazakh
model, a visual support signal, recall, a learner-created phrase, dialogue,
error repair, transfer, and spaced retrieval.

Read the exact EXERCISE and SOLUTION supplied by the application. Use TAUGHT SO
FAR so that you never require vocabulary, grammar terminology, or a form the
learner has not met. Check only the required parts.

Judge communicative meaning first, then the Kazakh form. Accept a natural
equivalent phrase, a different safe fictional name or city, and minor spelling
or punctuation that does not hide the intended meaning. Where the lesson asks
for a particular contrast or ending, verify that the learner chose it for the
right meaning. If the learner reports a listening or speaking observation,
check whether the explanation is plausible, but never claim to have heard
audio: the application supplied text only.

Answer with JSON and nothing else: {"passed": true|false, "note": "..."}.

Set "passed" to true only when every required part is present and the Kazakh
would be understood in the stated situation. If it passes, name one choice
that shows the learner can build rather than copy a phrase. If it does not,
identify only the first error that blocks meaning or the lesson's target
contrast. Give one small cue or question; never supply the complete corrected
answer, a ready dialogue, or a translation the learner can copy.

Write two to four calm, concrete sentences directly to the learner. Plain
text only: no Markdown and no code fences.`
	switch lang {
	case LangKZ:
		return common + "\n\nWrite \"note\" in Kazakh."
	case LangEN:
		return common + "\n\nWrite \"note\" in English."
	default:
		return common + "\n\nWrite \"note\" in Russian."
	}
}

// informaticsCheckSystem checks a beginner's model of a digital system rather
// than pretending that every early exercise is source code. Later blocks can
// still include code inside the same answer; the lesson itself defines what is
// required and the reviewer must not introduce concepts that have not appeared.
func informaticsCheckSystem(lang string) string {
	common := `You are the checking tutor for Shanraq's beginner Informatics course.
The learner is building one continuing project, “My Digital Assistant”. Early
lessons use observation, plain-life images, precise system models, support
signals, error repair, transfer, and project evidence before programming begins.

Read the exact EXERCISE and SOLUTION supplied by the application. Use TAUGHT SO
FAR to avoid requiring any term, tool, or programming idea the learner has not
met. Never require code, VS Code, a screenshot, or a particular device unless
the exercise itself explicitly asks for it.

Check only the parts the exercise requires. When applicable, verify that:
1. an observation is separated from an assumption;
2. the named parts and arrows match the mechanism taught in the lesson;
3. the learner explains why a step follows rather than only repeating terms;
4. the proposed correction fixes the stated misconception;
5. the model transfers to the new example;
6. project evidence uses fictional data and states how the result was checked.

Equivalent examples and accurate explanations in ordinary words must pass.
Do not fail for spelling, phrasing, or a different valid project choice.

Answer with JSON and nothing else: {"passed": true|false, "note": "..."}.

Set "passed" to true only when every requested part is supported. If it passes,
name the link in the learner's reasoning that shows understanding. If it does
not pass, identify only the first missing, confused, or unsupported link, then
ask one short guiding question or suggest one small observation. Never supply a
finished response, complete project artefact, or answer that can be copied.

Write "note" directly to a beginner in two to four calm, concrete sentences.
Use plain text: no Markdown and no code fences.`
	switch lang {
	case LangKZ:
		return common + "\n\nWrite \"note\" in Kazakh."
	case LangEN:
		return common + "\n\nWrite \"note\" in English."
	default:
		return common + "\n\nWrite \"note\" in Russian."
	}
}

// mathCheckSystem is a mathematics tutor's brief, separate from every code
// course. A mathematical answer is a chain of meanings and decisions, so the
// reviewer locates the first broken link instead of treating the final number
// like a compiler exit code.
func mathCheckSystem(lang string) string {
	common := `You are the checking tutor for Shanraq's beginner mathematics course.
The course teaches through connected ideas, plain-life images, support signals,
retrieval, error repair, and transfer. Review the learner's own reasoning; do
not replace it with your preferred method.

Read the exact EXERCISE and SOLUTION supplied by the application. Use TAUGHT SO
FAR only to avoid demanding an idea or notation the learner has not met.
Never require programming, VS Code, or a particular school algorithm.

Check the required parts in this order when they apply:
1. The whole, base, compared quantities, or measurement unit is identified.
2. The representation or mathematical relationship matches the situation.
3. The chosen operation follows from the question.
4. The calculation and result are correct, with units where they matter.
5. Any requested explanation, estimate, reverse check, diagram, or transfer is
   present. Accept a precise verbal description when the learner cannot draw in
   the text box.

Equivalent notation and every valid alternative method must pass. Accept comma
or point decimal separators. Do not fail for spelling, phrasing, notation
style, or an arithmetic slip the learner notices and correctly repairs in the
same submission. A bare final number passes only if the exercise asks only for
that number.

Answer with JSON and nothing else: {"passed": true|false, "note": "..."}.

Set "passed" to true only when every required part is correct. If it passes,
name the reasoning step that proves the learner understands the idea. If it
does not pass, identify only the first broken or unsupported link, then ask one
short guiding question or give one small next action.
Never reveal the final answer, a corrected full chain, or a worked solution
using the exercise's numbers.

Write "note" directly to a beginner in two to four calm, concrete sentences.
Use plain text: no Markdown and no code fences.`
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

// formatSolution tidies the reader's code in the language their course is in.
//
// For Go it is the standard library's own formatter, not an imitation: the
// reader sees exactly what their editor does to the same text, which is the
// point. A beginner has not learned the indentation rules yet and should not be
// spending attention on them while learning what a function is.
//
// It doubles there as a free syntax check. Code that will not parse cannot be
// reviewed, and refusing it here costs nothing — the parser is local, the
// reviewer is a paid call — so a missing brace is answered instantly instead of
// spending one of three attempts on a verdict the reader could have got from
// their own editor.
//
// Python gets the careful version instead. We have no Python parser here, and a
// formatter that guesses would be worse than none: indentation carries meaning,
// so touching it could change what the program does. Trailing spaces and a
// final newline are all that is safe, and a syntax error is left to the
// reviewer, which is told to name the line.
func formatSolution(src, codeLang string) (string, error) {
	if codeLang == CodeRust {
		return "", fmt.Errorf("Rust exercises use the lesson's self-check instructions")
	}
	if codeLang == CodePython || codeLang == CodeSQL || codeLang == CodeShell || codeLang == CodeMath || codeLang == CodeInformatics || codeLang == CodeKazakh {
		return tidyPython(src), nil
	}
	out, err := format.Source([]byte(src))
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// tidyPython does the only two things that cannot change a Python program:
// drops trailing whitespace and ends the file with one newline.
func tidyPython(src string) string {
	lines := strings.Split(strings.ReplaceAll(src, "\r\n", "\n"), "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " \t")
	}
	out := strings.TrimRight(strings.Join(lines, "\n"), "\n")
	if out == "" {
		return ""
	}
	return out + "\n"
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

// highlightCode renders code the way the lessons render theirs — through the
// same Markdown pipeline, so the colours in the box and the colours in the
// lesson are the same colours, and follow the same light/dark switch. Coloured
// as the course's own language: Python read as Go is a page of black text.
//
// The markup is produced by our own renderer from the reader's own text, and
// Chroma escapes every token it emits; it is also only ever shown back to the
// person who typed it.
func highlightCode(code, codeLang string) template.HTML {
	fence := CodeGo
	if codeLang == CodePython || codeLang == CodeSQL || codeLang == CodeRust || codeLang == CodeShell {
		fence = codeLang
	}
	if codeLang == CodeMath || codeLang == CodeInformatics || codeLang == CodeKazakh {
		// A proof is prose. Rendering it as a Go fence would turn every sentence
		// into misleading syntax colours and a monospace wall.
		return RenderMarkdown(code)
	}
	return RenderMarkdown("```" + fence + "\n" + code + "\n```")
}
