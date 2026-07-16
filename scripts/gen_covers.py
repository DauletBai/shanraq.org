#!/usr/bin/env python3
"""Generate Shanraq's house-style article covers: soft "watercolor" washes in
SVG (turbulence + displacement for bleeding edges, grain overlay), a minimal
white motif per theme, and the brand watermark. Self-contained, no raster.

Run: python3 scripts/gen_covers.py
"""
import pathlib

OUT = pathlib.Path("web/static/covers")
OUT.mkdir(parents=True, exist_ok=True)

W, H = 1200, 675  # 16:9, good for OG images

WM = ('<g transform="translate(1108,24) scale(0.16)" fill="none" stroke-linecap="round" opacity="0.9">'
      '<g stroke="#ffffff" stroke-opacity="0.5" stroke-width="45">'
      '<line x1="76" y1="76" x2="324" y2="324"/><line x1="137" y1="37" x2="363" y2="263"/>'
      '<line x1="37" y1="137" x2="263" y2="363"/><line x1="76" y1="324" x2="324" y2="76"/>'
      '<line x1="137" y1="363" x2="363" y2="137"/><line x1="37" y1="263" x2="263" y2="37"/></g>'
      '<g stroke="#ffffff" stroke-width="30">'
      '<line x1="76" y1="76" x2="324" y2="324"/><line x1="137" y1="37" x2="363" y2="263"/>'
      '<line x1="37" y1="137" x2="263" y2="363"/><line x1="76" y1="324" x2="324" y2="76"/>'
      '<line x1="137" y1="363" x2="363" y2="137"/><line x1="37" y1="263" x2="263" y2="37"/></g></g>')


def blobs(colors, seed):
    # A few large overlapping ellipses become a watercolor wash under the filter.
    pts = [(300, 260, 360, 300), (820, 360, 420, 340), (560, 520, 480, 300), (980, 140, 300, 260)]
    out = []
    for i, (cx, cy, rx, ry) in enumerate(pts):
        out.append(f'<ellipse cx="{cx}" cy="{cy}" rx="{rx}" ry="{ry}" fill="{colors[i % len(colors)]}" fill-opacity="0.7"/>')
    return "".join(out)


def motif_opinion():
    return ('<g fill="none" stroke="#ffffff" stroke-width="10" stroke-linecap="round" stroke-linejoin="round" opacity="0.92">'
            '<path d="M470 250 h260 a30 30 0 0 1 30 30 v120 a30 30 0 0 1 -30 30 h-160 l-70 60 v-60 h-30 a30 30 0 0 1 -30 -30 v-120 a30 30 0 0 1 30 -30 z"/>'
            '<line x1="510" y1="300" x2="700" y2="300"/><line x1="510" y1="340" x2="670" y2="340"/></g>')


def motif_network():
    n = [(470, 300), (650, 230), (760, 360), (600, 430), (500, 470)]
    lines = "".join(f'<line x1="{a[0]}" y1="{a[1]}" x2="{b[0]}" y2="{b[1]}"/>' for a, b in
                     [(n[0], n[1]), (n[1], n[2]), (n[2], n[3]), (n[3], n[4]), (n[4], n[0]), (n[0], n[3]), (n[1], n[3])])
    dots = "".join(f'<circle cx="{x}" cy="{y}" r="16" fill="#ffffff"/>' for x, y in n)
    return f'<g stroke="#ffffff" stroke-width="7" opacity="0.92">{lines}{dots}</g>'


def motif_gears():
    return ('<g fill="none" stroke="#ffffff" stroke-width="12" opacity="0.92">'
            '<circle cx="560" cy="340" r="90"/><circle cx="560" cy="340" r="34"/>'
            '<circle cx="720" cy="420" r="60"/><circle cx="720" cy="420" r="22"/></g>')


def motif_book():
    return ('<g fill="none" stroke="#ffffff" stroke-width="11" stroke-linejoin="round" opacity="0.92">'
            '<path d="M470 460 q120 -60 190 0 v-190 q-70 -60 -190 0 z"/>'
            '<path d="M660 460 q120 -60 190 0 v-190 q-120 -60 -190 0 z"/>'
            '<line x1="660" y1="270" x2="660" y2="460"/></g>')


def motif_mountains():
    return ('<g fill="none" stroke="#ffffff" stroke-width="11" stroke-linejoin="round" opacity="0.92">'
            '<path d="M440 470 L560 300 L650 400 L760 250 L880 470 Z"/>'
            '<circle cx="800" cy="220" r="34"/></g>')


def motif_lamp():
    return ('<g fill="none" stroke="#ffffff" stroke-width="11" stroke-linecap="round" stroke-linejoin="round" opacity="0.92">'
            '<path d="M600 250 a90 90 0 0 1 90 90 c0 45 -35 70 -50 95 h-80 c-15 -25 -50 -50 -50 -95 a90 90 0 0 1 90 -90 z"/>'
            '<line x1="565" y1="470" x2="635" y2="470"/><line x1="575" y1="500" x2="625" y2="500"/></g>')


def motif_scales():
    return ('<g fill="none" stroke="#ffffff" stroke-width="10" stroke-linecap="round" stroke-linejoin="round" opacity="0.92">'
            '<line x1="600" y1="250" x2="600" y2="470"/><line x1="470" y1="290" x2="730" y2="290"/>'
            '<line x1="470" y1="290" x2="440" y2="360"/><line x1="470" y1="290" x2="500" y2="360"/>'
            '<path d="M440 360 a30 22 0 0 0 60 0"/>'
            '<line x1="730" y1="290" x2="700" y2="360"/><line x1="730" y1="290" x2="760" y2="360"/>'
            '<path d="M700 360 a30 22 0 0 0 60 0"/><line x1="540" y1="470" x2="660" y2="470"/></g>')


def motif_ball():
    import math
    cx, cy, r = 600, 360, 92
    pent = []
    for i in range(5):
        a = -math.pi / 2 + i * 2 * math.pi / 5
        pent.append((cx + 34 * math.cos(a), cy + 34 * math.sin(a)))
    poly = " ".join(f"{x:.0f},{y:.0f}" for x, y in pent)
    spokes = "".join(f'<line x1="{x:.0f}" y1="{y:.0f}" x2="{cx + (r) * (x - cx) / 34:.0f}" y2="{cy + (r) * (y - cy) / 34:.0f}"/>' for x, y in pent)
    return (f'<g fill="none" stroke="#ffffff" stroke-width="9" stroke-linejoin="round" opacity="0.92">'
            f'<circle cx="{cx}" cy="{cy}" r="{r}"/><polygon points="{poly}" fill="#ffffff"/>{spokes}</g>')


def motif_masks():
    return ('<g fill="none" stroke="#ffffff" stroke-width="9" stroke-linecap="round" opacity="0.92">'
            '<path d="M500 280 q60 -20 120 0 q10 90 -60 150 q-70 -60 -60 -150 z"/>'
            '<circle cx="540" cy="330" r="6" fill="#ffffff"/><circle cx="580" cy="330" r="6" fill="#ffffff"/>'
            '<path d="M545 380 q15 15 30 0"/>'
            '<path d="M640 320 q60 -20 120 0 q10 90 -60 150 q-70 -60 -60 -150 z"/>'
            '<circle cx="680" cy="370" r="6" fill="#ffffff"/><circle cx="720" cy="370" r="6" fill="#ffffff"/>'
            '<path d="M685 425 q15 -15 30 0"/></g>')


def motif_bars():
    return ('<g fill="none" stroke="#ffffff" stroke-width="10" stroke-linecap="round" stroke-linejoin="round" opacity="0.92">'
            '<line x1="480" y1="470" x2="760" y2="470"/>'
            '<rect x="510" y="400" width="46" height="70" fill="#ffffff" stroke="none"/>'
            '<rect x="580" y="360" width="46" height="110" fill="#ffffff" stroke="none"/>'
            '<rect x="650" y="310" width="46" height="160" fill="#ffffff" stroke="none"/>'
            '<path d="M500 420 L560 380 L620 350 L720 290"/><path d="M720 290 l-34 4 M720 290 l4 34"/></g>')


def motif_orbit():
    return ('<g fill="none" stroke="#ffffff" stroke-width="9" opacity="0.92">'
            '<circle cx="600" cy="360" r="26" fill="#ffffff"/>'
            '<ellipse cx="600" cy="360" rx="150" ry="60"/>'
            '<ellipse cx="600" cy="360" rx="150" ry="60" transform="rotate(60 600 360)"/>'
            '<ellipse cx="600" cy="360" rx="150" ry="60" transform="rotate(120 600 360)"/>'
            '<circle cx="750" cy="360" r="10" fill="#ffffff"/></g>')


def motif_heart():
    return ('<g fill="none" stroke="#ffffff" stroke-width="11" stroke-linecap="round" stroke-linejoin="round" opacity="0.92">'
            '<path d="M600 460 C 500 390 470 340 500 300 C 525 268 575 275 600 315 C 625 275 675 268 700 300 C 730 340 700 390 600 460 Z"/>'
            '<path d="M505 372 h55 l22 -40 l34 78 l22 -38 h60"/></g>')


def motif_rook():
    teeth = "".join(f'<rect x="{x}" y="286" width="15" height="26"/>' for x in (556, 583, 610, 627))
    return ('<g fill="#ffffff" opacity="0.92">'
            f'{teeth}'
            '<rect x="562" y="300" width="76" height="120" rx="4"/>'
            '<rect x="546" y="420" width="108" height="26" rx="5"/></g>')


def motif_leaf():
    return ('<g fill="none" stroke="#ffffff" stroke-width="10" stroke-linecap="round" stroke-linejoin="round" opacity="0.92">'
            '<path d="M600 296 C 516 338 508 432 600 476 C 692 432 684 338 600 296 Z"/>'
            '<line x1="600" y1="312" x2="600" y2="462"/>'
            '<path d="M600 356 l-42 -22 M600 400 l-42 -22 M600 356 l42 -22 M600 400 l42 -22"/></g>')


def motif_glove():
    return ('<g opacity="0.92">'
            '<circle cx="600" cy="350" r="70" fill="#ffffff"/>'
            '<circle cx="528" cy="362" r="26" fill="#ffffff"/>'
            '<rect x="556" y="406" width="88" height="44" rx="10" fill="#ffffff"/>'
            '<line x1="580" y1="418" x2="620" y2="418" stroke="#8a4a3a" stroke-width="5"/>'
            '<line x1="580" y1="436" x2="620" y2="436" stroke="#8a4a3a" stroke-width="5"/></g>')


def motif_note():
    return ('<g fill="none" stroke="#ffffff" stroke-width="11" opacity="0.92">'
            '<line x1="560" y1="300" x2="560" y2="432"/><line x1="684" y1="280" x2="684" y2="412"/>'
            '<line x1="560" y1="300" x2="684" y2="280"/></g>'
            '<g fill="#ffffff" opacity="0.92">'
            '<ellipse cx="544" cy="438" rx="27" ry="18" transform="rotate(-20 544 438)"/>'
            '<ellipse cx="668" cy="418" rx="27" ry="18" transform="rotate(-20 668 418)"/></g>')


def motif_robot():
    return ('<g fill="none" stroke="#ffffff" stroke-width="10" stroke-linejoin="round" opacity="0.92">'
            '<rect x="520" y="304" width="160" height="140" rx="20"/>'
            '<line x1="600" y1="304" x2="600" y2="272"/></g>'
            '<g fill="#ffffff" opacity="0.92">'
            '<circle cx="600" cy="262" r="11"/><circle cx="562" cy="360" r="17"/>'
            '<circle cx="638" cy="360" r="17"/><rect x="560" y="402" width="80" height="12" rx="6"/></g>')


def motif_wheat():
    rows = ""
    for y in (322, 352, 382, 412):
        rows += f'<line x1="600" y1="{y}" x2="566" y2="{y-26}"/><line x1="600" y1="{y}" x2="634" y2="{y-26}"/>'
    return ('<g fill="none" stroke="#ffffff" stroke-width="9" stroke-linecap="round" opacity="0.92">'
            '<line x1="600" y1="472" x2="600" y2="300"/>'
            '<line x1="600" y1="300" x2="582" y2="270"/><line x1="600" y1="300" x2="618" y2="270"/>'
            f'{rows}</g>')


def motif_tennis():
    strings = ""
    for x in (560, 580, 600):
        strings += f'<line x1="{x}" y1="278" x2="{x}" y2="402"/>'
    for y in (312, 340, 368):
        strings += f'<line x1="522" y1="{y}" x2="638" y2="{y}"/>'
    return ('<g fill="none" stroke="#ffffff" stroke-width="8" opacity="0.9">'
            '<ellipse cx="580" cy="340" rx="60" ry="72"/>'
            f'{strings}<line x1="580" y1="412" x2="580" y2="472" stroke-width="12" stroke-linecap="round"/></g>'
            '<circle cx="690" cy="404" r="20" fill="#ffffff" opacity="0.92"/>')


def motif_hanger():
    return ('<g fill="none" stroke="#ffffff" stroke-width="11" stroke-linecap="round" '
            'stroke-linejoin="round" opacity="0.92">'
            '<path d="M600 300 C 600 285 620 285 620 301 C 620 314 600 316 600 332"/>'
            '<path d="M600 332 L512 400 H688 Z"/></g>')


def motif_plane():
    return ('<g fill="#ffffff" opacity="0.92" stroke="#ffffff" stroke-width="4" '
            'stroke-linejoin="round">'
            '<path d="M690 298 L516 380 L600 372 L612 436 Z"/></g>'
            '<path d="M690 298 L600 372" fill="none" stroke="#cfe6f5" stroke-width="5" '
            'opacity="0.7"/>')


def motif_dove():
    return ('<g fill="none" stroke="#ffffff" stroke-width="12" stroke-linecap="round" '
            'opacity="0.92"><path d="M512 394 Q568 334 600 390 Q632 334 688 394"/></g>'
            '<g fill="#ffffff" opacity="0.82">'
            '<ellipse cx="586" cy="432" rx="15" ry="6" transform="rotate(32 586 432)"/>'
            '<ellipse cx="614" cy="432" rx="15" ry="6" transform="rotate(-32 614 432)"/></g>'
            '<line x1="600" y1="416" x2="600" y2="446" stroke="#ffffff" stroke-width="6" '
            'stroke-linecap="round" opacity="0.82"/>')


def motif_arch():
    return ('<g fill="none" stroke="#ffffff" stroke-width="11" stroke-linecap="round" '
            'stroke-linejoin="round" opacity="0.92">'
            '<path d="M528 456 V372 A72 72 0 0 1 672 372 V456"/>'
            '<path d="M502 456 H698"/>'
            '<path d="M552 456 V394 M600 456 V360 M648 456 V394"/></g>'
            '<path d="M586 296 h28 l9 28 h-46 z" fill="#ffffff" opacity="0.92"/>')


def motif_key():
    return ('<g opacity="0.92">'
            '<circle cx="600" cy="330" r="38" fill="none" stroke="#ffffff" stroke-width="12"/>'
            '<circle cx="600" cy="330" r="13" fill="none" stroke="#ffffff" stroke-width="9"/>'
            '<line x1="600" y1="368" x2="600" y2="456" stroke="#ffffff" stroke-width="12" '
            'stroke-linecap="round"/>'
            '<line x1="600" y1="456" x2="628" y2="456" stroke="#ffffff" stroke-width="12" '
            'stroke-linecap="round"/>'
            '<line x1="600" y1="430" x2="622" y2="430" stroke="#ffffff" stroke-width="12" '
            'stroke-linecap="round"/></g>')


def motif_candle():
    return ('<path d="M600 298 C 585 316 587 337 600 346 C 613 337 615 316 600 298 Z" '
            'fill="#ffffff" opacity="0.95"/>'
            '<line x1="600" y1="346" x2="600" y2="358" stroke="#ffffff" stroke-width="5"/>'
            '<rect x="576" y="358" width="48" height="100" rx="8" fill="none" '
            'stroke="#ffffff" stroke-width="11" opacity="0.92"/>'
            '<line x1="554" y1="458" x2="646" y2="458" stroke="#ffffff" stroke-width="11" '
            'stroke-linecap="round" opacity="0.92"/>')


def motif_laurel():
    return ('<g fill="none" stroke="#ffffff" stroke-width="9" stroke-linecap="round" opacity="0.92">'
            '<path d="M600 482 C 548 452 516 404 522 330"/>'
            '<path d="M600 482 C 652 452 684 404 678 330"/></g>'
            '<g fill="#ffffff" opacity="0.9">'
            '<ellipse cx="518" cy="352" rx="22" ry="10" transform="rotate(52 518 352)"/>'
            '<ellipse cx="526" cy="396" rx="22" ry="10" transform="rotate(36 526 396)"/>'
            '<ellipse cx="548" cy="436" rx="22" ry="10" transform="rotate(20 548 436)"/>'
            '<ellipse cx="682" cy="352" rx="22" ry="10" transform="rotate(-52 682 352)"/>'
            '<ellipse cx="674" cy="396" rx="22" ry="10" transform="rotate(-36 674 396)"/>'
            '<ellipse cx="652" cy="436" rx="22" ry="10" transform="rotate(-20 652 436)"/>'
            '<circle cx="600" cy="312" r="10"/></g>')


def motif_hands():
    return ('<g opacity="0.92">'
            '<path d="M600 300 C 560 262 505 275 505 322 C 505 366 560 396 600 424 '
            'C 640 396 695 366 695 322 C 695 275 640 262 600 300 Z" fill="#ffffff"/>'
            '<path d="M498 428 C 498 502 702 502 702 428" fill="none" stroke="#ffffff" '
            'stroke-width="13" stroke-linecap="round"/></g>'
            '<g stroke="#ffffff" stroke-width="7" stroke-linecap="round" opacity="0.75">'
            '<line x1="600" y1="250" x2="600" y2="228"/>'
            '<line x1="548" y1="264" x2="536" y2="244"/>'
            '<line x1="652" y1="264" x2="664" y2="244"/></g>')


def motif_family():
    # A Kazakh yurt (kiiz üi): domed roof, shanyraq crown, roof poles (uyq),
    # lattice wall band and a doorway — the hearth as a symbol of family.
    return ('<g fill="none" stroke="#ffffff" stroke-width="10" stroke-linecap="round" '
            'stroke-linejoin="round" opacity="0.92">'
            '<path d="M508 404 C 512 340 560 306 600 306 C 640 306 688 340 692 404"/>'  # dome
            '<path d="M508 404 L497 470 M692 404 L703 470"/>'                            # walls
            '<path d="M508 404 H692 M497 470 H703"/>'                                    # wall band + ground
            '<path d="M600 306 L520 404 M600 306 L560 404 M600 306 L640 404 '
            'M600 306 L680 404"/></g>'                                                   # roof poles (uyq)
            '<g opacity="0.92" stroke="#ffffff" stroke-linecap="round">'
            '<circle cx="600" cy="300" r="24" fill="none" stroke-width="9"/>'            # shanyraq
            '<path d="M600 278 V322 M578 300 H622 M584 284 L616 316 M616 284 L584 316" '
            'stroke-width="5"/></g>'                                                     # crown cross
            '<path d="M584 470 V436 Q600 424 616 436 V470 Z" fill="#ffffff" '
            'opacity="0.9"/>')                                                           # doorway


def motif_helix():
    import math
    cx, A, y0, y1, n = 600, 50, 252, 476, 44
    freq = 1.6

    def strand(phase):
        pts = []
        for i in range(n + 1):
            t = i / n
            y = y0 + (y1 - y0) * t
            x = cx + A * math.sin(2 * math.pi * freq * t + phase)
            pts.append(f"{x:.0f},{y:.0f}")
        return f"<polyline points='{' '.join(pts)}' fill='none' stroke='#ffffff' stroke-width='8' opacity='0.92'/>"

    rungs = ""
    for i in range(1, 8):
        t = i / 8.0
        y = y0 + (y1 - y0) * t
        x1 = cx + A * math.sin(2 * math.pi * freq * t)
        x2 = cx + A * math.sin(2 * math.pi * freq * t + math.pi)
        rungs += f"<line x1='{x1:.0f}' y1='{y:.0f}' x2='{x2:.0f}' y2='{y:.0f}' stroke='#ffffff' stroke-width='6' opacity='0.8'/>"
    return strand(0) + strand(math.pi) + rungs


COVERS = {
    "opinion":  (["#ef8a5a", "#e0653f", "#f2b06a", "#c8492e"], motif_opinion, 11),
    "ai":       (["#4f8fd6", "#3f5fb0", "#6fb3d6", "#2f3f8f"], motif_network, 21),
    "labor":    (["#5aa86a", "#3f8f57", "#8fc47a", "#2e7d47"], motif_gears, 31),
    "language": (["#d97aa0", "#c05680", "#e6a0bd", "#a83f68"], motif_book, 41),
    "world":    (["#4f9fd6", "#3f79b0", "#6fc0d6", "#2f5f9f"], motif_mountains, 51),
    "education":(["#e6a24f", "#d67f3f", "#f2c06a", "#c8622e"], motif_lamp, 61),
    "politics": (["#5b6b9a", "#3f4f7a", "#7a8ab0", "#2f3a5f"], motif_scales, 71),
    "football": (["#3f9a5a", "#2e7d47", "#6fb87a", "#256b3a"], motif_ball, 81),
    "culture":  (["#8a5aa8", "#6f3f8f", "#a87ac0", "#5a2f7a"], motif_masks, 91),
    "economy":  (["#d99a4f", "#c07f3f", "#e6b86a", "#a8622e"], motif_bars, 101),
    "space":    (["#3f4f8a", "#2f3a66", "#5f6fa8", "#232a4a"], motif_orbit, 111),
    "health":   (["#4fa89a", "#3f8f82", "#7ac0b4", "#2e7d70"], motif_heart, 121),
    "chess":    (["#8a6f4f", "#6f573f", "#a88f6f", "#5a442e"], motif_rook, 131),
    "ecology":  (["#5aa070", "#3f8f57", "#86c48f", "#2e7d47"], motif_leaf, 141),
    "boxing":   (["#c0563f", "#9a3f2e", "#d9856f", "#7a2e20"], motif_glove, 151),
    "music":    (["#9a5aa8", "#7f3f8f", "#bd82c4", "#6a2f7a"], motif_note, 161),
    "robot":    (["#5a6b8a", "#3f4f6f", "#7a8aa8", "#2f3a52"], motif_robot, 171),
    "agriculture": (["#d9a84f", "#c08f3f", "#e6c46a", "#a8752e"], motif_wheat, 181),
    "tennis":   (["#7aa83f", "#5f8f2e", "#a8c46f", "#4a7a20"], motif_tennis, 191),
    "biotech":  (["#3f9a8a", "#2e7d70", "#6fc0b0", "#25706a"], motif_helix, 201),
    "athletics":(["#4f9a5a", "#3f8f4f", "#7ac06f", "#2e7d3f"], motif_laurel, 211),
    "charity":  (["#d9607a", "#c04060", "#e68aa0", "#a83f58"], motif_hands, 221),
    "family":   (["#e0975a", "#c77a3f", "#f2b87a", "#a8622e"], motif_family, 231),
    "architecture": (["#b89a6f", "#9a7f57", "#d0b487", "#7a6240"], motif_arch, 241),
    "crime":    (["#4f7a9a", "#3f6180", "#6f9ac0", "#2f4f6a"], motif_key, 251),
    "holidays": (["#9a5a7a", "#7f3f60", "#bd82a0", "#6a2f4f"], motif_candle, 261),
    "fashion":  (["#b0548a", "#8f3f6f", "#c87aa8", "#7a2f5a"], motif_hanger, 271),
    "aviation": (["#4f9fd6", "#3f7fb8", "#7ac0e0", "#2f6fa8"], motif_plane, 281),
    "defense":  (["#6a8aa8", "#4f6f8f", "#8fa8c0", "#3f5a7a"], motif_dove, 291),
}


# Rubric (category) each cover belongs to, so covers are filed under
# web/static/covers/<category>/<name>.svg instead of one flat directory.
CAT = {
    "opinion": "opinion", "ai": "it",
    "labor": "economy", "language": "culture", "world": "world",
    "education": "society", "politics": "politics", "football": "sport",
    "culture": "culture", "economy": "economy", "space": "technology",
    "health": "society", "chess": "sport", "ecology": "society",
    "boxing": "sport", "music": "culture", "robot": "technology",
    "agriculture": "economy", "tennis": "sport", "biotech": "technology",
    "athletics": "sport", "charity": "society", "family": "society",
    "architecture": "culture", "crime": "society", "holidays": "society",
    "fashion": "culture", "aviation": "technology", "defense": "politics",
}


def build(name, colors, motif, seed):
    svg = f'''<svg width="{W}" height="{H}" viewBox="0 0 {W} {H}" xmlns="http://www.w3.org/2000/svg">
  <defs>
    <filter id="wc" x="-15%" y="-15%" width="130%" height="130%">
      <feTurbulence type="fractalNoise" baseFrequency="0.011 0.014" numOctaves="3" seed="{seed}" result="n"/>
      <feDisplacementMap in="SourceGraphic" in2="n" scale="26" xChannelSelector="R" yChannelSelector="G"/>
      <feGaussianBlur stdDeviation="2"/>
    </filter>
    <filter id="grain"><feTurbulence type="fractalNoise" baseFrequency="0.9" numOctaves="2" seed="{seed}"/>
      <feColorMatrix type="matrix" values="0 0 0 0 0  0 0 0 0 0  0 0 0 0 0  0 0 0 0.05 0"/></filter>
  </defs>
  <rect width="{W}" height="{H}" fill="{colors[0]}"/>
  <g filter="url(#wc)">{blobs(colors, seed)}</g>
  <rect width="{W}" height="{H}" fill="#ffffff" filter="url(#grain)"/>
  {motif()}
  {WM}
</svg>
'''
    d = OUT / CAT[name]
    d.mkdir(parents=True, exist_ok=True)
    (d / f"{name}.svg").write_text(svg)
    print("wrote", d / f"{name}.svg")


for name, (colors, motif, seed) in COVERS.items():
    build(name, colors, motif, seed)
