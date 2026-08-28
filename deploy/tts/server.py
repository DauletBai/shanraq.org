#!/usr/bin/env python3
"""A small speech synthesiser, kept beside the site rather than inside it.

The site is Go and the synthesiser is Python, so it lives in its own container
and answers over HTTP on a private network. Nothing here is reachable from
outside; the only caller is the article module, which sends the blocks of an
article and gets back one audio file plus the length of each block.

Those lengths are the point. A finished audio file says nothing about itself,
so the page could not follow along without them: each block is synthesised
separately and its duration recorded, which is what lets the reader see which
paragraph is sounding.

Voices are loaded on first use and kept. The Kazakh model is 122 MB on disk and
around 600 MB resident, so they are not all loaded at once on a 4 GB machine --
and a lock allows one synthesis at a time, because two Kazakh readings running
together is how this container gets killed for memory.
"""

import io
import re
import json
import os
import subprocess
import tempfile
import threading
import wave
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer

VOICES = {
    "kz": os.environ.get("VOICE_KZ", "/models/kk_KZ-issai-high.onnx"),
    "ru": os.environ.get("VOICE_RU", "/models/ru_RU-dmitri-medium.onnx"),
    "en": os.environ.get("VOICE_EN", "/models/en_US-lessac-medium.onnx"),
}
BITRATE = os.environ.get("TTS_BITRATE", "24k")

# How long each phoneme is held. Above 1.0 is slower.
#
# At the default pace the voice clips the ends of words -- Russian and Kazakh
# both carry their grammar in the endings, so a swallowed suffix is not a
# blemish but a lost case or a lost tense. Ten per cent slower is enough to land
# them and is still a natural reading speed; it costs the same ten per cent in
# file size and synthesis time.
#
# Kept as an environment variable because the right value is a judgement made by
# ear, and finding it should not require rebuilding an image.
LENGTH_SCALE = float(os.environ.get("TTS_LENGTH_SCALE", "1.1"))

# ---- mixed-language text ---------------------------------------------------
#
# Our articles are not monolingual. A Russian piece carries "Eurasian Resources
# Group", "Bloomberg", "The Insider"; a Kazakh one carries the same. Read by the
# wrong voice those come out as nonsense -- a Russian voice spelling its way
# through English letters -- and the listener hears a mistake rather than a name.
#
# So a block is cut into runs of one language and each run is spoken by its own
# voice. All three models share a 22 050 Hz sample rate, which is what makes the
# pieces joinable without resampling.
#
# What this can and cannot see, stated plainly:
#   * Latin letters are English. Reliable.
#   * The letters Kazakh has and Russian does not -- ә ғ қ ң ө ұ ү һ і -- mark
#     Kazakh. Reliable.
#   * Russian inside Kazakh is invisible here. Both are Cyrillic and a Russian
#     word carrying none of those letters is spelled exactly like a Kazakh one.
#     Telling them apart needs a dictionary, which is a different project; until
#     then such words are read in the article's own voice, as they are today.
KAZAKH_ONLY = set("әғқңөұүһіӘҒҚҢӨҰҮҺІ")
LATIN = re.compile(r"[A-Za-z]")
CYRILLIC = re.compile(r"[А-Яа-яЁё" + "".join(KAZAKH_ONLY) + "]")

# A run has to be worth switching for. One stray letter, an abbreviation inside
# a word, a lone Roman numeral: switching voice for those makes the sentence
# stutter between two people, which is worse than reading them slightly wrong.
MIN_RUN_CHARS = 4


def word_lang(word, base):
    """Which voice should say this word."""
    if LATIN.search(word):
        return "en"
    if any(ch in KAZAKH_ONLY for ch in word):
        return "kz"
    if CYRILLIC.search(word):
        return base
    return None  # digits, punctuation: they belong to whatever surrounds them


def language_runs(text, base):
    """Cuts one block into [(lang, text), ...] in reading order."""
    tokens = re.findall(r"\w+|\W+", text, re.UNICODE)
    runs = []
    for tok in tokens:
        lang = word_lang(tok, base) if tok.strip() else None
        if lang is None:
            if runs:
                runs[-1][1] += tok
            else:
                runs.append([base, tok])
            continue
        if runs and runs[-1][0] == lang:
            runs[-1][1] += tok
        else:
            runs.append([lang, tok])

    # Fold away runs too short to be worth a change of voice.
    merged = []
    for lang, txt in runs:
        letters = len(re.findall(r"\w", txt, re.UNICODE))
        if merged and (letters < MIN_RUN_CHARS or lang == merged[-1][0]):
            merged[-1][1] += txt
        else:
            merged.append([lang, txt])
    # A run whose voice is missing falls back to the article's own.
    return [(l if VOICES.get(l) and os.path.exists(VOICES[l]) else base, t)
            for l, t in merged if t.strip()]


_loaded = {}
_lock = threading.Lock()


def voice_for(lang):
    """Loads a voice once and keeps it.

    Keeping every voice resident would cost more than this machine has, so they
    arrive as they are asked for. In practice a site serving three languages
    ends up holding all three, which is why the container is given headroom
    rather than exactly one model's worth.
    """
    if lang in _loaded:
        return _loaded[lang]
    path = VOICES.get(lang)
    if not path or not os.path.exists(path):
        raise FileNotFoundError(f"no voice for {lang} at {path}")
    from piper import PiperVoice

    v = PiperVoice.load(path)
    _loaded[lang] = v
    return v


def synthesize(lang, blocks):
    """One WAV of the whole article, and where each block sits inside it."""
    from piper import SynthesisConfig

    cfg = SynthesisConfig(length_scale=LENGTH_SCALE)
    frames, cues, cursor = [], [], 0.0
    rate = channels = width = None

    for i, text in enumerate(blocks):
        if not text.strip():
            continue
        # A block is not necessarily one language. Each run gets its own voice
        # and the pieces are joined; a block with no foreign words is a single
        # run and behaves exactly as before.
        started = cursor
        for run_lang, run_text in language_runs(text, lang):
            voice = voice_for(run_lang)
            with tempfile.NamedTemporaryFile(suffix=".wav") as tmp:
                with wave.open(tmp.name, "wb") as w:
                    voice.synthesize_wav(run_text, w, syn_config=cfg)
                with wave.open(tmp.name, "rb") as w:
                    rate = rate or w.getframerate()
                    channels = channels or w.getnchannels()
                    width = width or w.getsampwidth()
                    n = w.getnframes()
                    frames.append(w.readframes(n))
                    cursor += n / float(w.getframerate())
        # The cue spans the whole block however many voices it took, because the
        # page highlights paragraphs, not runs.
        cues.append({"i": i, "a": round(started, 3), "b": round(cursor, 3)})

    if not frames:
        raise ValueError("nothing to say")

    buf = io.BytesIO()
    with wave.open(buf, "wb") as out:
        out.setnchannels(channels or 1)
        out.setsampwidth(width or 2)
        out.setframerate(rate or 22050)
        for f in frames:
            out.writeframes(f)
    return buf.getvalue(), cues, cursor


def encode(wav_bytes):
    """WAV to Opus in an Ogg container.

    Opus at 24 kbit/s is transparent for speech and roughly half the size of the
    AAC this used to produce. That matters here in a way it would not elsewhere:
    the archive grows by an article a day in three languages, and halving the
    per-article cost halves how fast the disk fills.
    """
    p = subprocess.run(
        ["ffmpeg", "-hide_banner", "-loglevel", "error", "-i", "pipe:0",
         "-c:a", "libopus", "-b:a", BITRATE, "-application", "voip",
         "-f", "ogg", "pipe:1"],
        input=wav_bytes, stdout=subprocess.PIPE, stderr=subprocess.PIPE, check=False,
    )
    if p.returncode != 0:
        raise RuntimeError(p.stderr.decode()[:400])
    return p.stdout


class Handler(BaseHTTPRequestHandler):
    def log_message(self, *_):
        pass  # the caller logs; this would only duplicate it

    def do_GET(self):
        if self.path == "/health":
            langs = [l for l, p in VOICES.items() if p and os.path.exists(p)]
            self._json(200, {"ok": True, "voices": sorted(langs),
                             "loaded": sorted(_loaded.keys())})
            return
        self._json(404, {"error": "not found"})

    def do_POST(self):
        if self.path != "/synthesize":
            self._json(404, {"error": "not found"})
            return
        try:
            n = int(self.headers.get("Content-Length") or 0)
            req = json.loads(self.rfile.read(n) or b"{}")
            lang = str(req.get("lang") or "")
            blocks = req.get("blocks") or []
            if lang not in VOICES or not isinstance(blocks, list) or not blocks:
                self._json(400, {"error": "lang and blocks are required"})
                return
        except Exception as e:
            self._json(400, {"error": str(e)[:200]})
            return

        # One at a time. Two Kazakh readings at once is 1.2 GB of model plus
        # working memory on a machine with less than four, and the kernel
        # resolves that argument by killing this process.
        with _lock:
            try:
                wav, cues, seconds = synthesize(lang, blocks)
                audio = encode(wav)
            except Exception as e:
                self._json(500, {"error": str(e)[:400]})
                return

        self.send_response(200)
        self.send_header("Content-Type", "audio/ogg")
        self.send_header("Content-Length", str(len(audio)))
        self.send_header("X-Audio-Cues", json.dumps(cues, separators=(",", ":")))
        self.send_header("X-Audio-Seconds", str(int(seconds)))
        self.send_header("X-Audio-Voice", os.path.basename(VOICES[lang]))
        self.end_headers()
        self.wfile.write(audio)

    def _json(self, code, obj):
        body = json.dumps(obj).encode()
        self.send_response(code)
        self.send_header("Content-Type", "application/json")
        self.send_header("Content-Length", str(len(body)))
        self.end_headers()
        self.wfile.write(body)


if __name__ == "__main__":
    port = int(os.environ.get("TTS_PORT", "8081"))
    ThreadingHTTPServer(("0.0.0.0", port), Handler).serve_forever()
