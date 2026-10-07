#!/usr/bin/env python3
"""Render Kazakh lesson audio with the established Kazakh neural voice.

The review page remains available for spot checks and uncertain phrases.
After technical and content checks, a complete block can be published without
waiting for a separate owner review. Microsoft lists the voice for ``kk-KZ``.
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
        approved_new = [row for row in rows if row["status"] == "approved" and int(row["id"][3:]) >= 317]
        cards = "\n".join(
            '<article><b>{id}</b><p lang="kk">{phrase}</p>'
            '<audio controls preload="none" src="{file}"></audio></article>'.format(
                id=html.escape(row["id"]), phrase=html.escape(row["phrase"]),
                file=html.escape(row["file"]),
            ) for row in approved_new
        )
        REVIEW.write_text(
            """<!doctype html><html lang="ru"><meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>Казахский курс — одобренные записи</title>
<style>body{font:16px system-ui;max-width:920px;margin:30px auto;padding:0 18px;color:#242424}h1{font-size:30px}section{display:grid;grid-template-columns:repeat(auto-fit,minmax(270px,1fr));gap:14px}article{border:1px solid #d8d8d8;border-radius:12px;padding:14px;box-shadow:0 3px 12px #0001}p{font-size:20px;min-height:48px}audio{width:100%}</style>
<h1>Казахские аудиозаписи курса</h1><p>Записи доступны для выборочной проверки произношения. Для новых блоков используется ранее одобренный голос Aigul; аудиофайлы проходят техническую проверку перед выпуском.</p><section>"""
            + cards + "</section></html>\n",
            encoding="utf-8",
        )
        print(f"review page -> {REVIEW.relative_to(ROOT)} (no pending candidates)")
        return
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
        "approval": "technical_validation_before_release",
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
<h1>__COUNT__ новых реплик: следующий блок курса</h1>
<p>Проверка произношения перед публикацией следующего блока. Голос Aigul, язык kk-KZ. Ранее одобренные записи не изменялись.</p><section>"""
        .replace("__COUNT__", str(len(candidates)))
        + "\n".join(cards)
        + "</section></html>\n",
        encoding="utf-8",
    )
    print(f"review page -> {REVIEW.relative_to(ROOT)}")


if __name__ == "__main__":
    asyncio.run(render())
