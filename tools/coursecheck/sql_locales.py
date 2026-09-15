#!/usr/bin/env python3
"""Ensure every SQL lesson has RU/KZ/EN copies with identical code fences."""

from __future__ import annotations

import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2] / "course/lessons/sql"
FENCES = re.compile(r"```(sql|python|text)\s*\n(.*?)```", re.DOTALL)


def blocks(path: Path) -> list[tuple[str, str]]:
    def structure(body: str) -> str:
        # Human-facing labels and sample notes must be translated. Ignore their
        # contents while pinning every operator, identifier and placeholder.
        return re.sub(r"'(?:''|[^'])*'|\"(?:\"\"|[^\"])*\"", "'…'", body.strip())
    return [(lang, structure(body)) for lang, body in FENCES.findall(path.read_text(encoding="utf-8"))]


def main() -> int:
    errors: list[str] = []
    originals = sorted(p for p in ROOT.glob("*.md") if not p.stem.endswith(("-kz", "-en")))
    for original in originals:
        expected = blocks(original)
        for lang in ("kz", "en"):
            translated = original.with_name(f"{original.stem}-{lang}.md")
            if not translated.exists():
                errors.append(f"missing: {translated.relative_to(ROOT.parent.parent.parent)}")
                continue
            got = blocks(translated)
            if got != expected:
                errors.append(f"code drift: {translated.name} has {len(got)} fences, expected {len(expected)}")
    if errors:
        print("\n".join(errors), file=sys.stderr)
        return 1
    print(f"SQL locales complete: {len(originals)} lessons × 3 languages")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
