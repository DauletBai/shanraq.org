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


COVERS = {
    "opinion":  (["#ef8a5a", "#e0653f", "#f2b06a", "#c8492e"], motif_opinion, 11),
    "ai":       (["#4f8fd6", "#3f5fb0", "#6fb3d6", "#2f3f8f"], motif_network, 21),
    "labor":    (["#5aa86a", "#3f8f57", "#8fc47a", "#2e7d47"], motif_gears, 31),
    "language": (["#d97aa0", "#c05680", "#e6a0bd", "#a83f68"], motif_book, 41),
    "world":    (["#4f9fd6", "#3f79b0", "#6fc0d6", "#2f5f9f"], motif_mountains, 51),
    "education":(["#e6a24f", "#d67f3f", "#f2c06a", "#c8622e"], motif_lamp, 61),
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
    (OUT / f"cover-{name}.svg").write_text(svg)
    print("wrote", OUT / f"cover-{name}.svg")


for name, (colors, motif, seed) in COVERS.items():
    build(name, colors, motif, seed)
