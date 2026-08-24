package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// ProofVerdict is the checker's answer on one claimed typo.
//
// It deliberately carries a word and not a sentence. The caller applies the fix
// itself, replacing exactly one token inside exactly one sentence, so a model
// cannot quietly rewrite an author's argument while correcting their spelling.
// Anything wider than a word is refused by the caller as malformed.
type ProofVerdict struct {
	// Fix is true when the word really is misspelled and Fixed holds its
	// correct form.
	Fix bool `json:"fix"`
	// Fixed is the corrected word alone — no punctuation, no surrounding words.
	Fixed string `json:"fixed"`
	// Reason explains a refusal in the reader's own language. A refusal without
	// one is indistinguishable from the site being broken.
	Reason string `json:"reason"`
}

// Proofread judges one reader-reported typo.
//
// It is given the paragraph the sentence sits in, because a great many words are
// only wrong in context: agreement, case, a homophone, a name that is spelled one
// way throughout the piece. Judging the word alone would reject half the real
// typos and invent the other half.
//
// ErrDisabled is returned when the assistant is off; the caller then records the
// correction and tells the reader it is waiting rather than losing it.
func (m *Module) Proofread(ctx context.Context, lang, para, sentence, word string) (ProofVerdict, error) {
	c, model := m.moderateClient()
	if c == nil {
		return ProofVerdict{}, ErrDisabled
	}
	if strings.TrimSpace(sentence) == "" || strings.TrimSpace(word) == "" {
		return ProofVerdict{}, fmt.Errorf("proofread: empty claim")
	}
	user := "PARAGRAPH:\n" + para + "\n\nSENTENCE:\n" + sentence + "\n\nWORD:\n" + word
	raw, err := c.Complete(ctx, Request{
		Model:     model,
		System:    proofreadSystem(lang),
		User:      user,
		MaxTokens: 300,
	})
	if err != nil {
		return ProofVerdict{}, err
	}
	return parseProofVerdict(raw)
}

// parseProofVerdict extracts the verdict JSON, tolerating code fences or stray
// prose around it.
func parseProofVerdict(raw string) (ProofVerdict, error) {
	raw = strings.TrimSpace(raw)
	if i := strings.Index(raw, "{"); i >= 0 {
		if j := strings.LastIndex(raw, "}"); j > i {
			raw = raw[i : j+1]
		}
	}
	var v ProofVerdict
	if err := json.Unmarshal([]byte(raw), &v); err != nil {
		return ProofVerdict{}, fmt.Errorf("parse proofread verdict: %w", err)
	}
	v.Fixed = strings.TrimSpace(v.Fixed)
	v.Reason = strings.TrimSpace(v.Reason)
	if !v.Fix {
		v.Fixed = ""
	}
	return v, nil
}

// langName is how the checker is told which language it is reading.
func langName(lang string) string {
	switch lang {
	case "kz", "kk":
		return "Kazakh"
	case "en":
		return "English"
	}
	return "Russian"
}

func proofreadSystem(lang string) string {
	name := langName(lang)
	return `You are the proof-reader for an independent Kazakhstani publishing platform. The text is in ` + name + `.

A reader has pointed at one WORD inside one SENTENCE of a published article and claims it is a spelling error or a typo. The PARAGRAPH around it is given so you can judge the word in context.

Decide whether that single word is genuinely misspelled, and if it is, give its correct form.

Correct it when the word is:
- a misspelling or a typo (transposed, doubled, missing or wrong letters);
- wrong in agreement, case or number for this sentence;
- the wrong member of a homophone pair, where the intended one is unambiguous from the paragraph;
- missing a hyphen or wrongly hyphenated as one word;
- in ` + name + ` specifically: a Kazakh word written without its required letters (ә, ғ, қ, ң, ө, ұ, ү, һ, і), or a Russian word confusing е/ё where the meaning changes.

Refuse — fix:false — in every one of these cases:
- the word is spelled correctly as it stands;
- it is a proper name, a brand, a place, a transliteration or a term you cannot verify;
- it is a deliberate stylisation, dialect, slang, or part of a quotation;
- correcting it would change the meaning of the sentence, not just its spelling;
- the reader is proposing a different word choice, a better style, or a factual correction rather than a spelling fix;
- more than one word would have to change;
- you are not sure. An article edited wrongly costs its author more than an uncorrected typo costs its reader, so when in doubt, refuse.

Rules for the answer:
- "fixed" is the corrected form of THAT ONE WORD and nothing else: no punctuation, no neighbouring words, no explanation, and it must differ from the word given.
- Keep the original capitalisation pattern: a word that began a sentence stays capitalised.
- "reason" is one short sentence in ` + name + `, addressed to the reader, saying why nothing was changed. Leave it empty when you do fix.

Respond with ONLY a JSON object and nothing else:
{"fix":true|false,"fixed":"<corrected word, empty when fix is false>","reason":"<short reason in ` + name + `, empty when fix is true>"}`
}

// ProofreadRefute puts a proposed fix back to the model as something to knock
// down, and reports whether it survived.
//
// This is the self-consistency check, and it runs only on the path that changes
// an author's text. A model asked "is this wrong?" and a model asked "show me
// this is not wrong" fail in different directions, and a fix that both readings
// agree on is much closer to being a real typo than one a single pass liked. A
// refusal gets no second opinion: nothing changes on a refusal, so there is
// nothing to be wrong about.
func (m *Module) ProofreadRefute(ctx context.Context, lang, para, sentence, word, fixed string) (bool, string, error) {
	c, model := m.moderateClient()
	if c == nil {
		return false, "", ErrDisabled
	}
	user := "PARAGRAPH:\n" + para + "\n\nSENTENCE:\n" + sentence +
		"\n\nWORD AS PUBLISHED:\n" + word + "\n\nPROPOSED REPLACEMENT:\n" + fixed
	raw, err := c.Complete(ctx, Request{
		Model:     model,
		System:    refuteSystem(lang),
		User:      user,
		MaxTokens: 250,
	})
	if err != nil {
		return false, "", err
	}
	var v struct {
		Stands bool   `json:"stands"`
		Reason string `json:"reason"`
	}
	raw = strings.TrimSpace(raw)
	if i := strings.Index(raw, "{"); i >= 0 {
		if j := strings.LastIndex(raw, "}"); j > i {
			raw = raw[i : j+1]
		}
	}
	if err := json.Unmarshal([]byte(raw), &v); err != nil {
		return false, "", fmt.Errorf("parse refutation: %w", err)
	}
	return v.Stands, strings.TrimSpace(v.Reason), nil
}

func refuteSystem(lang string) string {
	name := langName(lang)
	return `You are checking a proposed edit to a published article before it is applied. The text is in ` + name + `.

Someone has proposed replacing one word in a published sentence. Your job is to try to show that the replacement is WRONG and the published word should stay. Look for every reason to leave the text alone:

- the published word is a correct spelling in ` + name + `;
- it is a proper name, a brand, a place, a title or a transliteration;
- it is correct in this sentence's grammar even if it looks unusual out of context;
- it is inside a quotation, a citation, a term of art, dialect or deliberate stylisation;
- the replacement changes the meaning, the tense, the number or the register rather than the spelling;
- the replacement is a matter of taste, not of correctness;
- the replacement introduces its own error, including in capitalisation.

Answer "stands":true ONLY if you cannot find any such reason and the replacement is plainly the correct spelling of the published word. If you are unsure, answer false — a published article edited wrongly costs its author more than an uncorrected typo costs its reader.

"reason" is one short sentence in ` + name + `, addressed to the reader who reported it, saying why the text was left as it is. Leave it empty when stands is true.

Respond with ONLY a JSON object and nothing else:
{"stands":true|false,"reason":"<short reason in ` + name + `, empty when stands is true>"}`
}
