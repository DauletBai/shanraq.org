#!/usr/bin/env python3
"""Render review candidates with a dedicated Kazakh neural voice.

These files are candidates, not approved lesson audio.  The generated review
page lets a native speaker approve or reject every phrase before lesson links
are changed. Microsoft lists the selected voice specifically for ``kk-KZ``;
the owner's listening review remains the final pronunciation check.
"""

from __future__ import annotations

import asyncio
import html
import json
import subprocess
import tempfile
import wave
from pathlib import Path

import edge_tts


ROOT = Path(__file__).resolve().parents[2]
AUDIO = ROOT / "web" / "static" / "course" / "kazakh-language" / "audio"
MANIFEST = AUDIO / "manifest.json"
REVIEW = AUDIO / "review.html"
VOICE = "kk-KZ-AigulNeural"
RATE = "-8%"


def run(command: list[str]) -> None:
    subprocess.run(command, check=True, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)


def validate(path: Path) -> None:
    with wave.open(str(path), "rb") as wav:
        if wav.getnchannels() != 1:
            raise ValueError(f"{path.name}: expected mono")
        if wav.getframerate() != 48_000:
            raise ValueError(f"{path.name}: expected 48 kHz")
        if wav.getsampwidth() != 3:
            raise ValueError(f"{path.name}: expected 24-bit PCM")
        # Even the one-word samples are longer than this with the review padding.
        # The earlier local voice silently emitted no samples; checking duration
        # and the PCM peak makes that failure impossible to publish again.
        if wav.getnframes() < 36_000:
            raise ValueError(f"{path.name}: recording is unexpectedly short")
        raw = wav.readframes(wav.getnframes())
    peak = max(
        abs(int.from_bytes(raw[i : i + 3], "little", signed=True))
        for i in range(0, len(raw) - 2, 3)
    )
    if peak < 10_000:
        raise ValueError(f"{path.name}: recording has no audible signal")


async def render() -> None:
    data = json.loads(MANIFEST.read_text(encoding="utf-8"))
    rows = data["recordings"]
    candidates = [row for row in rows if row["status"] == "candidate_review"]
    if not candidates:
        raise ValueError("manifest has no candidate_review recordings to render")
    with tempfile.TemporaryDirectory(prefix="shanraq-kz-audio-") as tmp:
        tmpdir = Path(tmp)
        for index, row in enumerate(candidates, 1):
            source = tmpdir / f"{row['id']}.mp3"
            target = AUDIO / row["file"]
            if target.is_file():
                validate(target)
                print(f"[{index:02d}/{len(candidates)}] {target.name} (kept)", flush=True)
                continue
            await edge_tts.Communicate(row["phrase"], VOICE, rate=RATE).save(str(source))
            run([
                "ffmpeg", "-y", "-i", str(source), "-af",
                "adelay=180,apad=pad_dur=0.18,loudnorm=I=-18:TP=-2:LRA=7",
                "-ar", "48000", "-ac", "1", "-c:a", "pcm_s24le", str(target),
            ])
            validate(target)
            print(f"[{index:02d}/{len(candidates)}] {target.name}", flush=True)

    data["candidate"] = {
        "provider": "Microsoft neural text to speech",
        "voice": VOICE,
        "locale": "kk-KZ",
        "rate": RATE,
        "processing": "48 kHz mono PCM; loudness normalised",
        "approval": "pending_owner_and_native_speaker_review",
    }
    MANIFEST.write_text(json.dumps(data, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")

    cards = []
    for row in candidates:
        cards.append(
            '<article><b>{id}</b><p lang="kk">{phrase}</p>'
            '<audio controls preload="none" src="{file}"></audio></article>'.format(
                id=html.escape(row["id"]),
                phrase=html.escape(row["phrase"]),
                file=html.escape(row["file"]),
            )
        )
    REVIEW.write_text(
        """<!doctype html><html lang="ru"><meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>Проверка казахского произношения</title>
<style>body{font:16px system-ui;max-width:920px;margin:30px auto;padding:0 18px;color:#242424}
h1{font-size:30px}section{display:grid;grid-template-columns:repeat(auto-fit,minmax(270px,1fr));gap:14px}
article{border:1px solid #d8d8d8;border-radius:12px;padding:14px;box-shadow:0 3px 12px #0001}
p{font-size:20px;min-height:48px}audio{width:100%}</style>
<h1>__COUNT__ новых реплик: блоки 3–6</h1>
<p>Проверка произношения перед публикацией уроков 17–48: «Школа и город», «День и услуги», «Пространство и падежи» и «Действие и время». Голос Aigul, язык kk-KZ. Одобренные записи уроков 1–16 не изменялись.</p><section>"""
        .replace("__COUNT__", str(len(candidates)))
        + "\n".join(cards)
        + "</section></html>\n",
        encoding="utf-8",
    )
    print(f"review page -> {REVIEW.relative_to(ROOT)}")


if __name__ == "__main__":
    asyncio.run(render())
