#!/usr/bin/env python3
"""Generate the isometric lesson maps for the Go course.

The maps are the course's memory aid: one screen per lesson, redrawable by hand.
They are generated rather than drawn so that all forty-two share one projection,
one palette and one set of shapes, and so that a translation is a change of
strings rather than a re-render -- three languages out of one scene.

SVG rather than a raster: the article renderer strips inline HTML, so a map is
loaded through <img>, and the media pipeline re-encodes uploads as JPEG, which
smears the thin lines a diagram is made of. A vector file avoids both, stays
sharp at any size, and carries its own light/dark palette.
"""

import math
import os

# Isometric projection: x runs right-and-down, y right-and-up, z straight up.
# A true isometric (30 degrees) keeps every edge measurable with a ruler, which
# is the point -- a reader has to be able to redraw this from memory.
COS30 = math.cos(math.radians(30))

# Scene units are chosen so a block is a few units across; SCALE turns them into
# pixels. Keeping the two separate means the geometry reads the same whatever
# size the map is finally drawn at.
SCALE = 38.0

# Everything in a scene is placed along a "road" rather than along the x axis.
# Laid out on x alone, an isometric scene runs from the top-left corner to the
# bottom-right one and leaves the other two corners empty -- which is exactly
# how the first draft of this looked. Stepping along +x and -y together moves
# straight across the picture instead, so the subject sits in the middle and the
# frame closes tightly around it.
ROAD = 4.0

_seen = []


def iso(x, y, z=0.0):
    """Project a point in scene space onto the drawing, in pixels."""
    p = ((x - y) * COS30 * SCALE, ((x + y) * 0.5 - z) * SCALE)
    _seen.append(p)
    return p


def at(u, side=0.0):
    """Scene coordinates u steps along the road, side steps across it."""
    return (ROAD + u + side, ROAD - u + side)


def esc(s):
    return s.replace("&", "&amp;").replace("<", "&lt;").replace(">", "&gt;")


def poly(points, cls):
    d = " ".join(f"{px:.2f},{py:.2f}" for px, py in points)
    return f'<polygon class="{cls}" points="{d}"/>'


def block(u, z, half, h, top, left, right, side=0.0):
    """A cube centred on the road at u, standing on height z."""
    x, y = at(u, side)
    x -= half
    y -= half
    w = d = half * 2
    p = lambda dx, dy, dz: iso(x + dx, y + dy, z + dz)
    return "".join([
        poly([p(0, 0, h), p(w, 0, h), p(w, d, h), p(0, d, h)], top),
        poly([p(0, d, 0), p(w, d, 0), p(w, d, h), p(0, d, h)], right),
        poly([p(w, 0, 0), p(w, d, 0), p(w, d, h), p(w, 0, h)], left),
    ])


def road(u0, u1, z, half_width, cls):
    """A flat band along the road -- the path a message travels."""
    corners = [at(u0, -half_width), at(u1, -half_width), at(u1, half_width), at(u0, half_width)]
    return poly([iso(x, y, z) for x, y in corners], cls)


def chevron(u, z, direction, cls, px=9.5):
    """A direction marker on a road. The road is horizontal on screen, so the
    marker is built in screen space: computed in scene units it came out
    pointing across the band instead of along it."""
    cx, cy = iso(*at(u), z)
    return poly([(cx + px * direction, cy), (cx - px * 0.45 * direction, cy - px * 0.72),
                 (cx - px * 0.45 * direction, cy + px * 0.72)], cls)


def text(u, z, s, cls, side=0.0, dy=0.0):
    px, py = iso(*at(u, side), z)
    return f'<text class="{cls}" x="{px:.2f}" y="{py + dy:.2f}">{esc(s)}</text>'


def ground(u0, u1, s0, s1, cls, step=1.1):
    """A ground plane cut to the objects standing on it. The first draft drew a
    square grid from the origin; it reached far past the scene and was what put
    all the empty space in the corners."""
    out = []
    n = int(round((u1 - u0) / step))
    m = int(round((s1 - s0) / step))
    for i in range(n + 1):
        u = u0 + i * step
        a, b = iso(*at(u, s0)), iso(*at(u, s1))
        out.append(f'<line class="{cls}" x1="{a[0]:.2f}" y1="{a[1]:.2f}" x2="{b[0]:.2f}" y2="{b[1]:.2f}"/>')
    for j in range(m + 1):
        sd = s0 + j * step
        a, b = iso(*at(u0, sd)), iso(*at(u1, sd))
        out.append(f'<line class="{cls}" x1="{a[0]:.2f}" y1="{a[1]:.2f}" x2="{b[0]:.2f}" y2="{b[1]:.2f}"/>')
    return "".join(out)


STYLE = """
:root{
  --ink:#3a3a3a; --soft:#6b6b6b; --line:#d6d3ce; --ground:#efece7;
  --top:#ffffff; --lft:#e6e1d9; --rgt:#d9d3c9;
  --red:#d32f2f; --red-d:#a92525; --red-l:#e86a63;
  --grn:#2e7d32; --grn-d:#24632a; --grn-l:#5fa363;
  --card:#faf8f4;
  --tagink:#ffffff;
}
@media (prefers-color-scheme: dark){
  :root{
    --ink:#ddd9d3; --soft:#a8a49d; --line:#4a4844; --ground:#343330;
    --top:#4a4844; --lft:#3c3b38; --rgt:#333230;
    --red:#ef5350; --red-d:#b93b38; --red-l:#ff8a84;
    --grn:#66bb6a; --grn-d:#478a4a; --grn-l:#8fd193;
    --card:#2b2a28;
    --tagink:#241f1f;
  }
}
.bg{fill:var(--card)}
.g{stroke:var(--line);stroke-width:.7;fill:none;opacity:.55}
.t{fill:var(--top)}   .l{fill:var(--lft)}   .r{fill:var(--rgt)}
.rt{fill:var(--red-l)} .rl{fill:var(--red)}  .rr{fill:var(--red-d)}
.gt{fill:var(--grn-l)} .gl{fill:var(--grn)}  .gr{fill:var(--grn-d)}
.band-req{fill:var(--red);opacity:.16}
.band-res{fill:var(--grn);opacity:.16}
.arw-req{fill:var(--red)} .arw-res{fill:var(--grn)}
text{font-family:ui-sans-serif,system-ui,-apple-system,"Segoe UI",Roboto,sans-serif}
.lbl{fill:var(--ink);font-size:13px;font-weight:600;text-anchor:middle}
.sub{fill:var(--soft);font-size:11px;text-anchor:middle}
.tag{fill:var(--tagink);font-size:10px;font-weight:700;text-anchor:middle;letter-spacing:.3px}
.cap{fill:var(--ink);font-size:12.5px;font-weight:700;text-anchor:middle}
.mono{font-family:ui-monospace,SFMono-Regular,Menlo,monospace;font-size:11px;fill:var(--soft);text-anchor:middle}
"""


def map00(s):
    """Lesson 0: a client asks, a server answers.

    The two roads run at different heights on purpose: a reader has to see that
    the question and the answer are separate journeys, not one line with arrows
    at both ends.
    """
    parts = [ground(-6.6, 6.6, -1.8, 1.8, "g")]

    # Painter's order matters here. The roads go down before the blocks, so a
    # road runs behind the server rack instead of across its face -- drawn last,
    # the bands lay a translucent stripe over the very thing they arrive at.
    parts.append(road(-4.6, 4.6, 2.15, 0.62, "band-req"))
    for u in (-2.9, -1.0, 0.9, 2.8):
        parts.append(chevron(u, 2.17, +1, "arw-req"))
    parts.append(road(-4.6, 4.6, 0.3, 0.62, "band-res"))
    for u in (-2.8, -0.9, 1.0, 2.9):
        parts.append(chevron(u, 0.32, -1, "arw-res"))

    # Client near the left end, server near the right one -- same height on
    # screen, so the eye travels straight across rather than down a diagonal.
    parts.append(block(-5.3, 0, 1.45, 1.7, "t", "l", "r"))
    for i in range(3):
        parts.append(block(5.3, i * 1.05, 1.45, 0.85, "t", "l", "r"))

    # The three parts of a request ride above the upper road; the three parts of
    # a response sit below the lower one. Same shape both times: a reader who
    # has learned to read the red row can read the green one unprompted.
    for i, key in enumerate(("req1", "req2", "req3")):
        u = -2.45 + i * 2.45
        parts.append(block(u, 3.05, 1.02, 0.5, "rt", "rl", "rr"))
        parts.append(text(u, 3.55, s[key], "tag", dy=4.0))
    for i, key in enumerate(("res1", "res2", "res3")):
        u = -2.45 + i * 2.45
        parts.append(block(u, -1.15, 1.02, 0.5, "gt", "gl", "gr"))
        parts.append(text(u, -0.65, s[key], "tag", dy=4.0))

    # Captions sit clear of the geometry: the two headings on the centre line
    # above and below, the two blocks named beside their own tops.
    parts.append(text(0, 4.55, s["req"], "cap"))
    parts.append(text(0, 4.15, s["req_ex"], "mono"))
    parts.append(text(0, -2.45, s["res"], "cap"))
    parts.append(text(0, -2.85, s["res_ex"], "mono"))
    parts.append(text(-5.3, 2.55, s["client"], "lbl"))
    parts.append(text(-5.3, 2.15, s["client_sub"], "sub"))
    parts.append(text(5.3, 3.75, s["server"], "lbl"))
    parts.append(text(5.3, 3.35, s["server_sub"], "sub"))
    return "".join(parts)


def render(scene, strings, pad=46):
    """Draw the scene, then frame it around its own bounding box."""
    _seen.clear()
    body = scene(strings)
    xs = [p[0] for p in _seen]
    ys = [p[1] for p in _seen]
    x0, x1 = min(xs) - pad, max(xs) + pad
    y0, y1 = min(ys) - pad, max(ys) + pad
    w, h = x1 - x0, y1 - y0
    return (
        f'<svg xmlns="http://www.w3.org/2000/svg" viewBox="{x0:.1f} {y0:.1f} {w:.1f} {h:.1f}" '
        f'width="{w:.0f}" height="{h:.0f}" role="img" aria-label="{esc(strings["alt"])}">'
        f"<style>{STYLE}</style>"
        f'<rect class="bg" x="{x0:.1f}" y="{y0:.1f}" width="{w:.1f}" height="{h:.1f}" rx="16"/>'
        f"{body}</svg>"
    )


L = {
    "kz": dict(
        alt="Клиент сұраныс жібереді, сервер жауап қайтарады",
        client="Клиент", client_sub="браузер, curl",
        server="Сервер", server_sub="сұрақ күтіп тұр",
        req="СҰРАНЫС", req_ex="GET /read/salem",
        res="ЖАУАП", res_ex="200 + HTML",
        req1="ӘДІС", req2="ЖОЛ", req3="ТАҚЫРЫПТАР",
        res1="КОД", res2="ТАҚЫРЫПТАР", res3="ДЕНЕ",
    ),
    "ru": dict(
        alt="Клиент отправляет запрос, сервер возвращает ответ",
        client="Клиент", client_sub="браузер, curl",
        server="Сервер", server_sub="ждёт вопроса",
        req="ЗАПРОС", req_ex="GET /read/salem",
        res="ОТВЕТ", res_ex="200 + HTML",
        req1="МЕТОД", req2="ПУТЬ", req3="ЗАГОЛОВКИ",
        res1="КОД", res2="ЗАГОЛОВКИ", res3="ТЕЛО",
    ),
    "en": dict(
        alt="A client sends a request, a server returns a response",
        client="Client", client_sub="browser, curl",
        server="Server", server_sub="waiting to be asked",
        req="REQUEST", req_ex="GET /read/salem",
        res="RESPONSE", res_ex="200 + HTML",
        req1="METHOD", req2="PATH", req3="HEADERS",
        res1="CODE", res2="HEADERS", res3="BODY",
    ),
}

if __name__ == "__main__":
    out = os.path.join(os.path.dirname(__file__), "..", "..", "web", "static", "course", "go")
    out = os.path.normpath(out)
    os.makedirs(out, exist_ok=True)
    for lang, strings in L.items():
        path = os.path.join(out, f"map-00-{lang}.svg")
        with open(path, "w", encoding="utf-8") as f:
            f.write(render(map00, strings))
        print("wrote", path)
