#!/usr/bin/env python3
"""Build compact, reproducible support maps for all Kazakh AI lessons."""

from __future__ import annotations

import html
import re
import textwrap
from pathlib import Path


ROOT = Path(__file__).resolve().parents[2]
LESSONS = ROOT / "course/lessons/kazakh-ai"
OUTPUT = ROOT / "web/static/course/kazakh-ai"
MAP_HEADING = re.compile(
    r"^## (?:Опорная карта(?: и проверка)?|Тірек сызба(?: және тексеру)?|"
    r"Memory map(?: and check)?|Support map|Recall map)$"
)


def locale(path: Path) -> str:
    if path.stem.endswith("-kz"):
        return "kz"
    if path.stem.endswith("-en"):
        return "en"
    return "ru"


def base_stem(path: Path) -> str:
    return re.sub(r"-(?:kz|en)$", "", path.stem)


def plain(value: str) -> str:
    value = re.sub(r"[`*_]", "", value)
    value = re.sub(r"\[([^]]+)]\([^)]+\)", r"\1", value)
    return re.sub(r"\s+", " ", value).strip()


def nodes_from(signal: str) -> list[str]:
    sentences = re.split(r"\.\s+", signal)
    if len(sentences) > 1:
        nodes = [plain(part).replace("→", "⇒").strip(" .;") for part in sentences]
    else:
        nodes = [plain(part).strip(" .;") for part in signal.split("→")]
    nodes = [node for node in nodes if node]
    if len(nodes) > 6:
        nodes = nodes[:5] + ["; ".join(nodes[5:])]
    return nodes or [plain(signal)]


def lines(value: str, width: int = 42) -> list[str]:
    wrapped = textwrap.wrap(value, width=width, break_long_words=True,
                            break_on_hyphens=True)
    return wrapped or [value]


def svg(title: str, signal: str, lang: str) -> str:
    nodes = nodes_from(signal)
    canvas_width = 760
    center_x = canvas_width / 2
    card_x = 40
    card_width = 680
    card_gap = 58
    card_top = 180
    wrapped_nodes = [lines(node) for node in nodes]
    card_heights = [max(92, 42 + len(node_lines) * 32) for node_lines in wrapped_nodes]
    positions = []
    cursor = card_top
    for card_height in card_heights:
        positions.append(cursor)
        cursor += card_height + card_gap
    footer_y = cursor + 8
    canvas_height = footer_y + 100
    labels = {"ru": "ОПОРНАЯ СХЕМА", "kz": "ТІРЕК СЫЗБАСЫ", "en": "SUPPORT MAP"}
    limits = {
        "ru": "Красная рамка — граница шага · зелёная стрелка — порядок",
        "kz": "Қызыл жиек — қадам шегі · жасыл жебе — реті",
        "en": "Red border — step boundary · green arrow — order",
    }
    out = [f'''<svg xmlns="http://www.w3.org/2000/svg" width="{canvas_width}" height="{canvas_height}" viewBox="0 0 {canvas_width} {canvas_height}" role="img" aria-label="{html.escape(title)}">
<defs><marker id="arrow" markerWidth="8" markerHeight="8" refX="7" refY="3.5" orient="auto"><path d="M0 0L0 7L8 3.5Z" fill="#2f7d45"/></marker></defs>
<style>
:root{{--bg:#f7f5f1;--card:#fff;--ink:#303030;--soft:#5b5b5b;--red:#d83232;--green:#2f7d45;--line:#e1ddd6;--badge:#fbe9e7}}
@media (prefers-color-scheme:dark){{:root{{--bg:#242321;--card:#302f2c;--ink:#f1efeb;--soft:#c8c3bb;--red:#ff6b66;--green:#75c98b;--line:#514e49;--badge:#482f2d}}}}
text{{font-family:Arial,"Noto Sans",sans-serif;fill:var(--ink)}}
</style>
<rect width="{canvas_width}" height="{canvas_height}" rx="28" fill="var(--bg)"/>
<text x="{center_x}" y="66" text-anchor="middle" font-size="18" font-weight="700" fill="var(--red)">{labels[lang]}</text>
<text x="{center_x}" y="108" text-anchor="middle" font-size="29" font-weight="800">{html.escape(title)}</text>''']
    for index, (node_lines, y, card_height) in enumerate(zip(wrapped_nodes, positions, card_heights)):
        if index:
            previous_bottom = positions[index - 1] + card_heights[index - 1]
            out.append(f'<path d="M{center_x:.0f} {previous_bottom + 10:.1f} V{y - 12:.1f}" stroke="var(--green)" stroke-width="4" stroke-linecap="round" marker-end="url(#arrow)"/>')
        out.append(f'<rect x="{card_x}" y="{y}" width="{card_width}" height="{card_height}" rx="22" fill="var(--card)" stroke="var(--red)" stroke-width="4"/>')
        badge_y = y + card_height / 2
        out.append(f'<circle cx="80" cy="{badge_y:.1f}" r="25" fill="var(--badge)" stroke="var(--red)" stroke-width="3"/>')
        out.append(f'<text x="80" y="{badge_y + 8:.1f}" text-anchor="middle" font-size="22" font-weight="800" fill="var(--red)">{index + 1}</text>')
        first_y = y + card_height / 2 - (len(node_lines) - 1) * 16 + 8
        for line_no, line in enumerate(node_lines):
            out.append(f'<text x="410" y="{first_y + line_no * 32:.1f}" text-anchor="middle" font-size="23" font-weight="700">{html.escape(line)}</text>')
    out.append(f'''<path d="M40 {footer_y} H720" stroke="var(--line)" stroke-width="3"/>
<text x="{center_x}" y="{footer_y + 52}" text-anchor="middle" font-size="18" fill="var(--soft)">{limits[lang]}</text>
</svg>''')
    return "\n".join(out) + "\n"


def process(path: Path) -> None:
    body = path.read_text(encoding="utf-8")
    rows = body.splitlines()
    heading_index = next((i for i, row in enumerate(rows) if MAP_HEADING.fullmatch(row)), None)
    if heading_index is None:
        raise ValueError(f"support-map heading not found: {path}")
    signal_index = next(i for i in range(heading_index + 1, len(rows)) if rows[i].strip())
    signal = rows[signal_index].strip()
    lang = locale(path)
    stem = base_stem(path)
    number = stem[:2]
    title = {"ru": f"Урок {int(number)}", "kz": f"{int(number)}-сабақ", "en": f"Lesson {int(number)}"}[lang]
    filename = f"map-{stem}-{lang}.svg"
    (OUTPUT / filename).write_text(svg(title, signal, lang), encoding="utf-8")
    image = {
        "ru": f"![Опорная схема урока {int(number)}](/static/course/kazakh-ai/{filename})",
        "kz": f"![{int(number)}-сабақтың тірек сызбасы](/static/course/kazakh-ai/{filename})",
        "en": f"![Lesson {int(number)} support map](/static/course/kazakh-ai/{filename})",
    }[lang]
    if image not in body:
        rows[signal_index + 1:signal_index + 1] = ["", image]
        path.write_text("\n".join(rows) + "\n", encoding="utf-8")


def main() -> None:
    OUTPUT.mkdir(parents=True, exist_ok=True)
    paths = sorted(LESSONS.glob("[0-9][0-9]-*.md"))
    for path in paths:
        process(path)
    print(f"generated {len(paths)} support maps in {OUTPUT}")


if __name__ == "__main__":
    main()
