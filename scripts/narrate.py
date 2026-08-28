#!/usr/bin/env python3
"""narrate.py — record an article and upload the reading.

Browsers have no Kazakh voice. The read-aloud button used to notice that and
offer a Russian one instead, so a Kazakh article was read in Russian. This
script replaces the guess with a recording.

It runs wherever Piper and the voice models are — a laptop, not the server. The
server keeps no synthesiser and no 122 MB model; it receives finished audio and
serves it as a file, which is why playback needs nothing from the reader's
system and survives a locked screen.

The cue map is the reason the page can follow along. A finished audio file says
nothing about itself, so each block is synthesised separately and its length
recorded. Cue N and block N must be the same paragraph, which is why the block
rule below is a copy of the one in listen.js. Change one and the highlight
drifts away from the sound.

Usage:
    scripts/narrate.py --slug <slug> --lang kz --model /path/kk_KZ-issai-high.onnx
    scripts/narrate.py --slug <slug> --lang kz --model ... --dry-run

Environment:
    SHANRAQ_BASE      site root, default https://shanraq.org
    SHANRAQ_API_KEY   key with the operator or admin role
    ESPEAK_DATA       espeak-ng data directory, if the bundled one is broken
"""

import argparse
import json
import os
import re
import subprocess
import sys
import tempfile
import wave
from html.parser import HTMLParser
from urllib import request as urlrequest

BASE = os.environ.get("SHANRAQ_BASE", "https://shanraq.org").rstrip("/")
BLOCK_TAGS = {"p", "h2", "h3", "h4", "li", "blockquote", "figcaption", "th", "td"}
SKIP_TAGS = {"pre", "code"}


class Blocks(HTMLParser):
    """Collects the article's readable blocks, outermost first.

    The page's own rule, transcribed: a block inside another block is read once,
    at the outer one; nothing inside code; nothing hidden; nothing shorter than
    a word.
    """

    def __init__(self):
        super().__init__(convert_charrefs=True)
        self.depth_prose = 0
        self.stack = []
        self.skip = 0
        self.buf = []
        self.out = []

    def handle_starttag(self, tag, attrs):
        a = dict(attrs)
        if self.depth_prose == 0:
            if "prose" in (a.get("class") or "").split():
                self.depth_prose = 1
            return
        if tag in SKIP_TAGS or a.get("aria-hidden") == "true":
            self.skip += 1
            return
        if tag in BLOCK_TAGS:
            # Nested block: the outer one is already collecting.
            self.stack.append(tag)

    def handle_endtag(self, tag):
        if self.depth_prose == 0:
            return
        if tag in SKIP_TAGS:
            self.skip = max(0, self.skip - 1)
            return
        if tag in BLOCK_TAGS and self.stack:
            self.stack.pop()
            if not self.stack:
                text = " ".join("".join(self.buf).split())
                if len(text) >= 2:
                    self.out.append(text)
                self.buf = []

    def handle_data(self, data):
        if self.depth_prose and self.stack and not self.skip:
            self.buf.append(data)


def fetch_blocks(slug: str, lang: str):
    url = f"{BASE}/read/{slug}?lang={lang}"
    with urlrequest.urlopen(url, timeout=60) as r:
        html = r.read().decode("utf-8", "replace")
    article_id = ""
    m = re.search(r'data-article-id="([0-9a-fA-F-]{36})"', html)
    if m:
        article_id = m.group(1)
    # The digest comes off the page rather than being recomputed here. The
    # server decides what counts as "the same text"; a second implementation of
    # that rule would drift, and every narration would report itself stale.
    digest = ""
    m = re.search(r'data-text-digest="([0-9a-f]{64})"', html)
    if m:
        digest = m.group(1)
    p = Blocks()
    p.feed(html)
    return article_id, digest, p.out


def load_voice(model: str):
    """Load Piper, working around a broken espeak data path in the wheel.

    The published wheel bakes in the path espeak-ng had on the build machine and
    looks for its data there, which exists on nobody's computer. Initialising the
    bridge once against a real directory, before anything else touches it, is
    enough; the phonemizer's own constructor is then stopped from undoing it.
    """
    from piper import espeakbridge

    data = os.environ.get("ESPEAK_DATA", "/opt/homebrew/share/espeak-ng-data")
    espeakbridge.initialize(data)
    import piper.phonemize_espeak as pe

    pe.EspeakPhonemizer.__init__ = lambda self, espeak_data_dir=None: None
    from piper import PiperVoice

    return PiperVoice.load(model)


def synth(voice, blocks, wav_path):
    """Writes one WAV of the whole article and returns the cue map."""
    from piper import SynthesisConfig

    cues, cursor = [], 0.0
    rate = channels = width = None
    frames = []
    for i, text in enumerate(blocks):
        with tempfile.NamedTemporaryFile(suffix=".wav", delete=True) as tmp:
            with wave.open(tmp.name, "wb") as w:
                voice.synthesize_wav(text, w, syn_config=SynthesisConfig())
            with wave.open(tmp.name, "rb") as w:
                rate = rate or w.getframerate()
                channels = channels or w.getnchannels()
                width = width or w.getsampwidth()
                n = w.getnframes()
                frames.append(w.readframes(n))
                dur = n / float(w.getframerate())
        cues.append({"i": i, "a": round(cursor, 3), "b": round(cursor + dur, 3)})
        cursor += dur
        print(f"  [{i + 1}/{len(blocks)}] {dur:5.1f} s  {text[:56]}", flush=True)
    with wave.open(wav_path, "wb") as out:
        out.setnchannels(channels or 1)
        out.setsampwidth(width or 2)
        out.setframerate(rate or 22050)
        for f in frames:
            out.writeframes(f)
    return cues, cursor


def encode(wav_path, out_path):
    """WAV to AAC. Speech at 48 kbit/s is about six megabytes for 18 minutes."""
    subprocess.run(
        ["afconvert", "-f", "m4af", "-d", "aac", "-b", "48000", "-q", "127", wav_path, out_path],
        check=True,
    )


def upload(article_id, lang, path, cues, seconds, voice_name, digest, key):
    with open(path, "rb") as f:
        body = f.read()
    url = (f"{BASE}/api/articles/{article_id}/audio/{lang}"
           f"?duration={int(seconds)}&voice={urlrequest.quote(voice_name)}&digest={digest}")
    req = urlrequest.Request(url, data=body, method="PUT")
    req.add_header("Content-Type", "audio/mp4")
    req.add_header("X-API-Key", key)
    req.add_header("X-Audio-Cues", json.dumps(cues, separators=(",", ":")))
    with urlrequest.urlopen(req, timeout=300) as r:
        return json.loads(r.read().decode())


def main() -> int:
    ap = argparse.ArgumentParser()
    ap.add_argument("--slug", required=True)
    ap.add_argument("--lang", required=True, choices=["kz", "ru", "en"])
    ap.add_argument("--model", required=True, help="Piper .onnx voice")
    ap.add_argument("--out", default="", help="keep the audio at this path too")
    ap.add_argument("--dry-run", action="store_true", help="synthesise, do not upload")
    ap.add_argument("--cues-out", default="", help="write the cue map here as JSON")
    args = ap.parse_args()

    article_id, digest, blocks = fetch_blocks(args.slug, args.lang)
    if not blocks:
        print("no readable blocks — wrong slug, or the page did not render", file=sys.stderr)
        return 1
    if not article_id and not args.dry_run:
        print("no article id on the page; cannot address the upload", file=sys.stderr)
        return 1
    print(f"{len(blocks)} blocks, {sum(len(b) for b in blocks)} characters")

    voice = load_voice(args.model)
    with tempfile.TemporaryDirectory() as tmp:
        wav = os.path.join(tmp, "full.wav")
        m4a = args.out or os.path.join(tmp, "full.m4a")
        cues, seconds = synth(voice, blocks, wav)
        encode(wav, m4a)
        size = os.path.getsize(m4a)
        print(f"total {seconds / 60:.1f} min, {size / 1048576:.1f} MB")
        if args.cues_out:
            # Writing the map out is what makes an upload-free install possible:
            # copy the audio into place, put the row in by hand, and the page
            # can still follow along. Useful before the first API key exists.
            with open(args.cues_out, "w", encoding="utf-8") as f:
                json.dump({"cues": cues, "digest": digest, "seconds": round(seconds, 1),
                           "article_id": article_id, "blocks": len(blocks)}, f,
                          ensure_ascii=False, separators=(",", ":"))
            print("cues at", args.cues_out)
        if args.dry_run:
            print("dry run: not uploaded")
            if args.out:
                print("audio at", args.out)
            return 0

        key = os.environ.get("SHANRAQ_API_KEY", "")
        if not key:
            print("SHANRAQ_API_KEY is not set", file=sys.stderr)
            return 1
        res = upload(article_id, args.lang, m4a, cues, seconds,
                     os.path.basename(args.model), digest, key)
        print("uploaded:", res.get("url"))
    return 0


if __name__ == "__main__":
    sys.exit(main())
