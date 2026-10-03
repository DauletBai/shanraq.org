#!/usr/bin/env python3
"""Build the native-speaker recording manifest for the Kazakh course.

The lesson sources intentionally keep the Kazakh phrase in ``#speak-kz=``
fragments until a reviewed human recording exists.  This tool deduplicates the
phrases across all three lesson languages and assigns stable file names.  On
later runs existing ids never change; phrases from new blocks are appended.
"""

from __future__ import annotations

import json
import re
from pathlib import Path
from urllib.parse import unquote


ROOT = Path(__file__).resolve().parents[2]
LESSONS = ROOT / "course" / "lessons" / "kazakh-language"
AUDIO = ROOT / "web" / "static" / "course" / "kazakh-language" / "audio"
MANIFEST = AUDIO / "manifest.json"
SCRIPT = ROOT / "course" / "kazakh-language" / "audio-recording-script.md"
PHRASE_RE = re.compile(r"#speak-kz=([^\)]+)")


def existing_entries() -> list[dict]:
    if not MANIFEST.exists():
        return []
    data = json.loads(MANIFEST.read_text(encoding="utf-8"))
    return list(data.get("recordings", []))


def lesson_phrases() -> list[tuple[str, str]]:
    found: list[tuple[str, str]] = []
    for path in sorted(LESSONS.glob("*.md")):
        text = path.read_text(encoding="utf-8")
        for encoded in PHRASE_RE.findall(text):
            found.append((unquote(encoded), path.name))
    return found


def main() -> None:
    AUDIO.mkdir(parents=True, exist_ok=True)
    entries = existing_entries()
    by_phrase = {row["phrase"]: row for row in entries}

    for phrase, _ in lesson_phrases():
        if phrase not in by_phrase:
            number = len(entries) + 1
            row = {
                "id": f"kz-{number:03d}",
                "phrase": phrase,
                "file": f"kz-{number:03d}.wav",
                "status": "awaiting_native_recording",
            }
            entries.append(row)
            by_phrase[phrase] = row

    used_by: dict[str, set[str]] = {row["phrase"]: set() for row in entries}
    for phrase, lesson in lesson_phrases():
        used_by.setdefault(phrase, set()).add(lesson)
    for row in entries:
        row["lessons"] = sorted(used_by.get(row["phrase"], set()))

    payload = {
        "language": "kk-KZ",
        "format": "PCM WAV, 48 kHz, mono, 24-bit preferred",
        "review": "native_speaker_required",
        "recordings": entries,
    }
    MANIFEST.write_text(
        json.dumps(payload, ensure_ascii=False, indent=2) + "\n",
        encoding="utf-8",
    )

    lines = [
        "# Текст для записи казахского курса",
        "",
        "Читайте только фразу после имени файла. Темп спокойный, естественный; "
        "не проговаривайте знаки препинания и номер.",
        "",
    ]
    lines.extend(f"{row['file']} — {row['phrase']}" for row in entries)
    SCRIPT.write_text("\n".join(lines) + "\n", encoding="utf-8")

    print(f"{len(entries)} unique phrases -> {MANIFEST.relative_to(ROOT)}")


if __name__ == "__main__":
    main()
