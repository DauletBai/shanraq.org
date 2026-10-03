#!/usr/bin/env python3
"""Refresh lesson usage in the Kazakh course audio manifest.

Published lessons link only to approved WAV assets. New recordings are added
to the manifest by the recording workflow, reviewed, and only then linked from
the lesson generators. Existing ids never change.
"""

from __future__ import annotations

import json
import re
from pathlib import Path


ROOT = Path(__file__).resolve().parents[2]
LESSONS = ROOT / "course" / "lessons" / "kazakh-language"
AUDIO = ROOT / "web" / "static" / "course" / "kazakh-language" / "audio"
MANIFEST = AUDIO / "manifest.json"
SCRIPT = ROOT / "course" / "kazakh-language" / "audio-recording-script.md"
AUDIO_RE = re.compile(
    r"/static/course/kazakh-language/audio/(kz-\d{3}\.wav)"
)


def existing_entries() -> list[dict]:
    if not MANIFEST.exists():
        return []
    data = json.loads(MANIFEST.read_text(encoding="utf-8"))
    return list(data.get("recordings", []))


def lesson_recordings() -> list[tuple[str, str]]:
    found: list[tuple[str, str]] = []
    for path in sorted(LESSONS.glob("*.md")):
        text = path.read_text(encoding="utf-8")
        for filename in AUDIO_RE.findall(text):
            found.append((filename, path.name))
    return found


def main() -> None:
    AUDIO.mkdir(parents=True, exist_ok=True)
    entries = existing_entries()
    by_file = {row["file"]: row for row in entries}
    linked = lesson_recordings()
    unknown = sorted({filename for filename, _ in linked} - set(by_file))
    if unknown:
        raise ValueError(f"lesson links recordings absent from manifest: {unknown}")

    used_by: dict[str, set[str]] = {row["file"]: set() for row in entries}
    for filename, lesson in linked:
        used_by[filename].add(lesson)
    for row in entries:
        row["lessons"] = sorted(used_by[row["file"]])

    payload = json.loads(MANIFEST.read_text(encoding="utf-8"))
    payload["recordings"] = entries
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

    print(f"{len(entries)} recordings, {len(linked)} lesson links -> {MANIFEST.relative_to(ROOT)}")


if __name__ == "__main__":
    main()
