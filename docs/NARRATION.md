# Narration

How an article gets a recorded reading, and why the parts sit where they do.

## Why this exists

The Listen button used to be browser speech synthesis and nothing else. That is
free, needs no infrastructure, and fails in ways a reader cannot act on:

- **No Kazakh voice anywhere.** No browser ships one. The player detected the
  gap and offered a Russian voice instead, so a Kazakh article was read in
  Russian. That is worse than silence: it sounds like the site cannot tell the
  two languages apart.
- **A locked screen stops it.** On a phone, synthesis dies when the screen goes
  off — which is exactly when someone would want to listen.
- **The failures are invisible.** A missing voice, a blocked autoplay, a tab
  the browser throttled: the reader hears nothing and has nowhere to look.

A recording removes all of it. It is a file; browsers have played files for
twenty years.

## Where the pieces are

| Piece | Location |
|---|---|
| Table | `article_audio`, migration `20251108004400_article_audio.sql` |
| Store and digest | `pkg/modules/articles/audio.go` |
| Blocks and text cleaning | `pkg/modules/articles/narrate.go` |
| Client and background job | `pkg/modules/articles/narrate_job.go` |
| Synthesiser service | `deploy/tts/` |
| Blob storage | `media.Module.SaveBlob` / `DeleteBlob` |
| Player | `web/static/js/listen.js` |
| Manual upload (rarely needed) | `pkg/modules/articles/audio_http.go`, `scripts/narrate.py` |

## The synthesiser runs on the server

Narration used to be produced on a laptop and uploaded. That was wrong for a
site that gains an article a day: every publication would have waited on someone
remembering to run a script, and authors would have inherited a chore that has
nothing to do with writing.

So the synthesiser is a container beside the app, `deploy/tts`. It is separate
from the app rather than inside it because it is Python with an ONNX runtime and
240 MB of voice models, while the app is a small static Go binary; fusing them
would put a component that changes once a year into every deploy. It is not
published to the internet — only the app talks to it, over the compose network.

**Measured on this VPS** (4 cores, Broadwell, 3.9 GB):

| Language | Speed | Peak memory |
|---|---|---|
| Russian | 10.8× real time | 440 MB |
| Kazakh | 2.2× real time | 602 MB |

Kazakh is slower because its voice is the "high" model rather than "medium" —
122 MB against 60 MB. It is worth it: Kazakh is the language no browser could
speak at all, and the one we cannot afford to get wrong.

Memory, not speed, is the constraint. Two Kazakh readings at once is 1.2 GB of
model plus working memory on a machine with less than four, so the service holds
a lock and synthesises one article at a time. The container has a 1600 MB limit
so that a runaway job is killed rather than taking Postgres with it.

## Voices are a mounted volume

They are not in the image and not in git: 240 MB together, they change about
never, and the Kazakh one carries its own licence. Fill the volume once on the
host:

    deploy/tts/fetch-voices.sh /var/lib/shanraq/voices

## Mixed-language text

An article is not monolingual. A Russian piece carries "Eurasian Resources
Group"; a Kazakh one carries the same. Read by the wrong voice those come out as
a Russian speaker spelling through English letters, and the listener hears a
mistake rather than a name.

Each block is therefore cut into runs of one language and each run is spoken by
its own voice. All three models share a 22 050 Hz sample rate, which is what
lets the pieces be joined without resampling.

What the detection can and cannot see:

- **Latin letters are English.** Reliable.
- **The letters Kazakh has and Russian does not** — ә ғ қ ң ө ұ ү һ і — mark
  Kazakh. Reliable.
- **Russian inside Kazakh is invisible.** Both are Cyrillic, and a Russian word
  carrying none of those letters is spelled exactly like a Kazakh one. Telling
  them apart needs a dictionary. Until then such words are read in the article's
  own voice.

A run must be at least four letters to earn a change of voice. Below that the
sentence stutters between two speakers, which is worse than reading a short
abbreviation slightly wrong: "ERG" in a Russian voice is how it is said aloud
here anyway.

## What the text loses on the way to the ear

A page is written for eyes. `speechText` in `pkg/modules/articles/narrate.go`
removes what means nothing aloud, and keeps what does:

| Removed | Why |
|---|---|
| `\ | * # ~ ^ < > { } [ ] _` | noise dropped into the middle of a sentence |
| `( ) « » „ " " ' '` | punctuation for the eye; **the words inside stay** |
| bare URLs | "h t t p s colon slash slash" for twenty characters |
| `·` `•` between sources | reads as "middle dot"; it is a pause |

| Rewritten | Why |
|---|---|
| `16 700` → `16700` | a thousands separator is read as a full stop: "sixteen, seven hundred" |
| ` — ` → `, ` | a spaced dash is a pause; a hyphen inside a word must survive |

A hyphen against a word is left alone, or compound words come apart.

## The cue map is the whole trick

The page highlights the paragraph being read and scrolls to follow it. Speech
synthesis gave that away free — it fires an event as it crosses each word. A
finished audio file says nothing at all about itself.

So the service synthesises **one block at a time**, records how long each took,
and returns the offsets alongside the audio:

    [{"i": 0, "a": 0.0, "b": 30.5}, {"i": 1, "a": 30.5, "b": 55.9}, ...]

`i` indexes the block sequence. Cue N and block N must be the same paragraph,
which means the block rule exists in more than one place and they have to agree:

- `NarrationBlocks` in `narrate.go` — what gets synthesised
- `blockEls()` in `listen.js` — what gets highlighted
- `Blocks` in `narrate.py` — only for the manual path

All say: outermost blocks of `p, h2, h3, h4, li, blockquote, figcaption, th,
td`; skip anything inside `pre` or `code`; skip `aria-hidden`; skip anything
shorter than two characters. **Change one and the highlight drifts away from the
sound** — silently, and only on articles with nested markup. `narrate_test.go`
pins both halves of the Go side.

## Staleness

An article edited after narration leaves the recording saying something the page
no longer says. It still plays; nothing errors.

`TextDigest` is computed in exactly one place and used by both sides: the job
stores it with the audio, and the page compares it against the text it is about
to serve. A second implementation of "the same text" would drift, and every
narration would report itself stale.

When they differ the page carries `data-stale="1"`. Re-publishing the article
queues a fresh reading; an unchanged article is skipped without spending the
eleven minutes of Kazakh synthesis to arrive at the same file.

## Two traps that were handled

**The media ledger would delete the audio.** The ledger meters uploads and
sweeps files nothing refers to, and it decides "refers to" by searching article
bodies for the key. Narration is referenced by a database column and never by
the prose, so a ledger entry would be swept as an orphan on the next pass —
quietly, and only for articles nobody had opened recently. `SaveBlob` therefore
writes outside the ledger, and the file's lifetime is tied to the article row
instead (`ON DELETE CASCADE`).

**Replacing a recording could delete the new one.** The old file is dropped only
after the row has been updated to point at the new one. If the row write fails,
the just-written file is removed instead, because from that moment nothing
references it.

## The espeak-ng workaround

The published `piper-tts` wheel bakes in the path espeak-ng's data had on the
maintainer's build machine and looks for it there. That path exists on nobody's
computer, and the error names it:

    Error processing file '/Users/runner/work/piper1-gpl/.../phontab'

`load_voice()` in the generator initialises the bridge once against a real
directory before anything else touches it, then stops the phonemizer's own
constructor from undoing that. Set `ESPEAK_DATA` if yours is not at
`/opt/homebrew/share/espeak-ng-data`.

This is one more reason the server has no synthesiser: a broken data path is a
two-line fix on a laptop and an unpleasant one inside a container.

## Licensing

Piper is MIT. The Kazakh voice is trained from scratch on
[IS2AI/Kazakh_TTS](https://github.com/IS2AI/Kazakh_TTS), released under
**CC-BY-4.0** — commercial use is permitted, attribution is required. Credit
ISSAI wherever the voice is presented as a feature.
