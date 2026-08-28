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
| Upload endpoint | `pkg/modules/articles/audio_http.go` |
| Blob storage | `media.Module.SaveBlob` / `DeleteBlob` |
| Player | `web/static/js/listen.js`, the recorded-reading branch |
| Generator | `scripts/narrate.py` |

## The server has no synthesiser

Deliberately. Piper needs a voice model — 122 MB for Kazakh, 60 MB each for
Russian and English — plus espeak-ng and an ONNX runtime. None of that belongs
on a small VPS to render audio that changes once per article.

So narration is produced wherever the models are, usually a laptop, and posted
back. The server stores the file and serves it. Nothing about playback depends
on the generator being reachable, or existing at all.

## Running the generator

Piper and the voices are not vendored. Set up once:

    python3 -m venv ~/.piper-venv
    ~/.piper-venv/bin/pip install piper-tts
    brew install espeak-ng          # see the espeak note below

Download voices from `rhasspy/piper-voices` on Hugging Face:

| Language | Voice | Size |
|---|---|---|
| Kazakh | `kk/kk_KZ/issai/high` | 122 MB |
| Russian | `ru/ru_RU/dmitri/medium` | 60 MB |
| English | `en/en_US/lessac/medium` | 60 MB |

Then, for one article and one language:

    export SHANRAQ_BASE=https://shanraq.org
    export SHANRAQ_API_KEY=<key with the operator or admin role>
    ~/.piper-venv/bin/python scripts/narrate.py \
        --slug sel-nad-almaty-... --lang kz \
        --model ~/voices/kk_KZ-issai-high.onnx

Add `--dry-run --out /tmp/preview.m4a` to hear it before anything is uploaded.

Roughly what to expect: a Kazakh article synthesises at about twice real time
because the Kazakh voice is a high-quality model; Russian and English run at
fourteen to sixteen times real time on the same machine. An eighteen-minute
article lands at about six megabytes at 48 kbit/s.

## The cue map is the whole trick

The page highlights the paragraph being read and scrolls to follow it. Speech
synthesis gave that away free — it fires an event as it crosses each word. A
finished audio file says nothing at all about itself.

So the generator synthesises **one block at a time**, records how long each
took, and uploads the offsets alongside the audio:

    [{"i": 0, "a": 0.0, "b": 30.5}, {"i": 1, "a": 30.5, "b": 55.9}, ...]

`i` indexes the block sequence. Cue N and block N must be the same paragraph,
which means the block rule exists in two places and they have to agree:

- `blockEls()` in `listen.js`
- `Blocks` in `narrate.py`

Both say: outermost blocks of `p, h2, h3, h4, li, blockquote, figcaption, th,
td`; skip anything inside `pre` or `code`; skip `aria-hidden`; skip anything
shorter than two characters. **Change one and the highlight drifts away from
the sound** — silently, and only on articles with nested markup.

The map travels in the `X-Audio-Cues` header rather than the query string: a
long article has a cue per block, which is kilobytes of JSON, and query strings
get truncated by proxies and written into access logs.

## Staleness

An article edited after narration leaves the recording saying something the page
no longer says. It still plays; nothing errors.

The page computes a digest of the served title and body, hands it to the reader
in `data-text-digest`, and the generator echoes it back with the audio. The
digest is computed in exactly one place on purpose: a second implementation of
"the same text" would drift, and every narration would report itself stale.

When they differ the page carries `data-stale="1"`. Re-run the generator.

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
