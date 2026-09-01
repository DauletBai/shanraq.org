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


def clear_above(z_top, half, gap_px=20.0):
    """The height a caption must sit at to clear the top corner of a block.

    A block's highest point on screen is its far top corner, which rises with
    the block's own footprint -- so the clearance depends on the block's size,
    not only on how tall it stands.
    """
    return z_top + half + gap_px / SCALE


def clear_below(z_bottom, half, gap_px=30.0):
    """The height a caption must sit at to clear the near bottom corner.

    The gap is larger than above because a caption's baseline is its underside:
    the lettering grows upward, back towards the block.
    """
    return z_bottom - half - gap_px / SCALE


def text(u, z, s, cls, side=0.0, dy=0.0):
    px, py = iso(*at(u, side), z)
    return f'<text class="{cls}" x="{px:.2f}" y="{py + dy:.2f}">{esc(s)}</text>'


def on_face(u, z_top, title, sub=""):
    """A name written across the middle of a block's top face.

    The centre of that face is the block's own top height taken at side 0, so
    the position follows the block instead of being guessed -- placed by eye the
    two names sat up by the far corner rather than in the rhombus.
    """
    out = [text(u, z_top, title, "lbl", dy=-2.0 if sub else 4.0)]
    if sub:
        out.append(text(u, z_top, sub, "sub", dy=11.0))
    return "".join(out)


def on_face_side(u, side, z_top, title, sub="", accent=False):
    """on_face for a block set off across the road, optionally an accent one."""
    lbl, sb = ("lbl-acc", "sub-acc") if accent else ("lbl", "sub")
    out = [text(u, z_top, title, lbl, side=side, dy=-2.0 if sub else 4.0)]
    if sub:
        out.append(text(u, z_top, sub, sb, side=side, dy=11.0))
    return "".join(out)


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
  --ink:#3a3a3a; --soft:#5c5c5c; --line:#d6d3ce;
  --top:#ded8cd; --lft:#bcb6aa; --rgt:#9d978d;
  --red:#d32f2f; --red-d:#a92525; --red-l:#d9564f;
  --grn:#2e7d32; --grn-d:#24632a; --grn-l:#4e8f52;
  --tagink:#ffffff;
}
@media (prefers-color-scheme: dark){
  :root{
    --ink:#d8d6d2; --soft:#c2bfb8; --line:#454545;
    --top:#4f4c47; --lft:#403d3a; --rgt:#34322f;
    --red:#ef5350; --red-d:#b93b38; --red-l:#ff8a84;
    --grn:#66bb6a; --grn-d:#478a4a; --grn-l:#8fd193;
    --tagink:#241f1f;
  }
}
.g{stroke:#8a857e;stroke-width:.7;fill:none;opacity:.35}
.t{fill:var(--top)}   .l{fill:var(--lft)}   .r{fill:var(--rgt)}
.rt{fill:var(--red-l)} .rl{fill:var(--red)}  .rr{fill:var(--red-d)}
.gt{fill:var(--grn-l)} .gl{fill:var(--grn)}  .gr{fill:var(--grn-d)}
.band-req{fill:var(--red);opacity:.16}
.band-res{fill:var(--grn);opacity:.16}\n.band-flow{fill:var(--ink);opacity:.13}\n.arw-flow{fill:var(--soft)}
.arw-req{fill:var(--red)} .arw-res{fill:var(--grn)}
text{font-family:ui-sans-serif,system-ui,-apple-system,"Segoe UI",Roboto,sans-serif}
.lbl{fill:var(--ink);font-size:13px;font-weight:600;text-anchor:middle}
.sub{fill:var(--soft);font-size:11px;text-anchor:middle}
.tag{fill:var(--tagink);font-size:10px;font-weight:700;text-anchor:middle;letter-spacing:.3px}
.cap{fill:var(--ink);font-size:12.5px;font-weight:700;text-anchor:middle}\n.lbl-acc{fill:var(--tagink);font-size:13px;font-weight:600;text-anchor:middle}\n.sub-acc{fill:var(--tagink);opacity:.8;font-size:11px;text-anchor:middle}
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
    tile, tile_h = 0.92, 0.5
    req_z, res_z = 3.05, -1.15
    for i, key in enumerate(("req1", "req2", "req3")):
        u = -2.45 + i * 2.45
        parts.append(block(u, req_z, tile, tile_h, "rt", "rl", "rr"))
        parts.append(text(u, req_z + tile_h, s[key], "tag", dy=4.0))
    for i, key in enumerate(("res1", "res2", "res3")):
        u = -2.45 + i * 2.45
        parts.append(block(u, res_z, tile, tile_h, "gt", "gl", "gr"))
        parts.append(text(u, res_z + tile_h, s[key], "tag", dy=4.0))

    # Captions are measured off the tiles rather than placed by eye, so a longer
    # scene or a taller block cannot push a heading back onto the drawing.
    req_line = clear_above(req_z + tile_h, tile)
    res_line = clear_below(res_z, tile)
    parts.append(text(0, req_line + 0.42, s["req"], "cap"))
    parts.append(text(0, req_line, s["req_ex"], "mono"))
    parts.append(text(0, res_line, s["res"], "cap"))
    parts.append(text(0, res_line - 0.42, s["res_ex"], "mono"))
    parts.append(on_face(-5.3, 1.7, s["client"], s["client_sub"]))
    parts.append(on_face(5.3, 2.95, s["server"], s["server_sub"]))
    return "".join(parts)


def map02(s):
    """Lesson: the three tools and what each is answerable for.

    A beginner meets the editor, the terminal and the language on the same day
    and blames the wrong one when something breaks. The map exists to fix which
    is which before any of them misbehaves.
    """
    parts = [ground(-6.7, 6.7, -1.5, 1.5, "g")]

    half, top = 1.5, 1.15
    seats = (-5.0, 0.0, 5.0)
    for u in seats:
        parts.append(block(u, 0, half, top, "t", "l", "r"))
    for a, b in ((-3.3, -1.7), (1.7, 3.3)):
        parts.append(road(a, b, top * 0.55, 0.42, "band-flow"))
        for u in (a + 0.45, b - 0.45):
            parts.append(chevron(u, top * 0.55 + 0.02, +1, "arw-flow"))

    for u, k in zip(seats, ("b1", "b2", "b3")):
        parts.append(on_face(u, top, s[k], s[k + "_sub"]))

    line_up = clear_above(top, half)
    line_dn = clear_below(0, half)
    parts.append(text(0, line_up + 0.42, s["head"], "cap"))
    parts.append(text(0, line_up, s["head_sub"], "mono"))
    parts.append(text(0, line_dn, s["dir"], "cap"))
    parts.append(text(0, line_dn - 0.42, s["dir_sub"], "mono"))
    return "".join(parts)


def map01(s):
    """Lesson 1: the two halves of a first program, and where to look at it.

    Three blocks in a row rather than a stack: in isometric a block standing on
    another hides the lower one's top face, and both halves have to carry a
    name. Read left to right they are also the order the program runs in.
    """
    parts = [ground(-6.7, 6.7, -1.5, 1.5, "g")]

    half, top = 1.5, 1.15
    seats = (-5.0, 0.0, 5.0)
    for u in seats:
        parts.append(block(u, 0, half, top, "t", "l", "r"))

    # Connectors sit in the gaps only. Run as one band across the whole row it
    # read as a road passing straight through the middle block -- the blocks
    # hide it, so what showed was a stripe entering one side and leaving the
    # other. Green is the colour the response already wore a lesson ago.
    for a, b in ((-3.3, -1.7), (1.7, 3.3)):
        parts.append(road(a, b, top * 0.55, 0.42, "band-res"))
        for u in (a + 0.45, b - 0.45):
            parts.append(chevron(u, top * 0.55 + 0.02, +1, "arw-res"))

    for u, k in zip(seats, ("b1", "b2", "b3")):
        parts.append(on_face(u, top, s[k], s[k + "_sub"]))

    line_up = clear_above(top, half)
    line_dn = clear_below(0, half)
    parts.append(text(0, line_up + 0.42, s["run"], "cap"))
    parts.append(text(0, line_up, s["run_sub"], "mono"))
    parts.append(text(0, line_dn, s["out"], "cap"))
    parts.append(text(0, line_dn - 0.42, s["out_sub"], "mono"))
    return "".join(parts)


def map03(s):
    """Lesson: the four types that carry a blog, and that a type is chosen once.

    Four boxes rather than a list, because the point is that each value sits in
    a container of a fixed kind. The caption underneath is the half beginners
    trip over: the box does not change shape later.
    """
    parts = [ground(-7.4, 7.4, -1.5, 1.5, "g")]

    half, top = 1.6, 1.15
    seats = (-5.7, -1.9, 1.9, 5.7)
    keys = ("t1", "t2", "t3", "t4")
    for u, k in zip(seats, keys):
        parts.append(block(u, 0, half, top, "t", "l", "r"))
        parts.append(on_face(u, top, s[k], s[k + "_val"]))

    line_up = clear_above(top, half)
    line_dn = clear_below(0, half)
    parts.append(text(0, line_up + 0.42, s["head"], "cap"))
    parts.append(text(0, line_up, s["head_sub"], "mono"))
    parts.append(text(0, line_dn, s["foot"], "cap"))
    parts.append(text(0, line_dn - 0.42, s["foot_sub"], "mono"))
    return "".join(parts)


def map04(s):
    """Lesson: a call goes in on one line and two answers come back.

    The two results sit side by side across the road rather than in a queue
    along it: in a row they would read as one following the other, and the
    whole point is that both arrive from the same call.
    """
    parts = [ground(-6.8, 6.2, -3.5, 3.5, "g")]

    half, top = 1.45, 1.15
    parts.append(road(-3.6, -1.9, top * 0.55, 0.42, "band-flow"))
    parts.append(chevron(-2.75, top * 0.55 + 0.02, +1, "arw-flow"))
    parts.append(road(1.1, 2.9, top * 0.55, 0.42, "band-res"))
    parts.append(chevron(2.0, top * 0.55 + 0.02, +1, "arw-res"))

    parts.append(block(-5.2, 0, half, top, "t", "l", "r"))
    parts.append(block(-0.4, 0, half, top, "t", "l", "r"))
    parts.append(block(4.5, 0, half, top, "gt", "gl", "gr", side=-1.95))
    parts.append(block(4.5, 0, half, top, "gt", "gl", "gr", side=1.95))

    parts.append(on_face(-5.2, top, s["in"], s["in_sub"]))
    parts.append(on_face(-0.4, top, s["fn"], s["fn_sub"]))
    parts.append(on_face_side(4.5, -1.95, top, s["o1"], s["o1_sub"], accent=True))
    parts.append(on_face_side(4.5, 1.95, top, s["o2"], s["o2_sub"], accent=True))

    line_up = clear_above(top, half + 1.95)
    line_dn = clear_below(0, half + 1.95)
    parts.append(text(0, line_up + 0.42, s["head"], "cap"))
    parts.append(text(0, line_up, s["head_sub"], "mono"))
    parts.append(text(0, line_dn, s["foot"], "cap"))
    parts.append(text(0, line_dn - 0.42, s["foot_sub"], "mono"))
    return "".join(parts)


def map05(s):
    """The announcement cover: five lessons stand, the road keeps going.

    Five solid blocks and an empty road after them says what a sentence would
    have to spell out — that the course is real and unfinished at once — and it
    cannot go stale the way a number in a headline does.
    """
    parts = [ground(-7.0, 7.0, -1.4, 1.4, "g")]

    half, top = 0.86, 0.95
    for i, u in enumerate((-5.6, -3.6, -1.6, 0.4, 2.4)):
        parts.append(block(u, 0, half, top, "gt", "gl", "gr"))
        parts.append(text(u, top, str(i + 1), "tag", dy=4.0))

    # The road past the last block: lessons that exist as a plan, not yet as
    # pages. Drawn as a path with nothing standing on it.
    parts.append(road(3.5, 6.4, top * 0.5, 0.44, "band-flow"))
    for u in (4.2, 5.1, 6.0):
        parts.append(chevron(u, top * 0.5 + 0.02, +1, "arw-flow"))

    line_up = clear_above(top, half)
    line_dn = clear_below(0, half)
    parts.append(text(0, line_up + 0.42, s["head"], "cap"))
    parts.append(text(0, line_up, s["head_sub"], "mono"))
    parts.append(text(0, line_dn, s["foot"], "cap"))
    parts.append(text(0, line_dn - 0.42, s["foot_sub"], "mono"))
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

L01 = {
    "kz": dict(
        alt="Бағдарлама нені жауап беретінін біледі, қай жерде тыңдайтынын біледі, ал браузер жауапты көреді",
        b1="НЕНІ ЖАУАП БЕРУ", b1_sub="http.HandleFunc",
        b2="ҚАЙДА ТЫҢДАУ", b2_sub="http.ListenAndServe",
        b3="БРАУЗЕР", b3_sub="localhost:8080",
        run="БІР ПӘРМЕН", run_sub="go run main.go",
        out="ЖАУАП", out_sub="200 · Сәлем",
    ),
    "ru": dict(
        alt="Программа знает, что отвечать и где слушать, а браузер видит ответ",
        b1="ЧТО ОТВЕЧАТЬ", b1_sub="http.HandleFunc",
        b2="ГДЕ СЛУШАТЬ", b2_sub="http.ListenAndServe",
        b3="БРАУЗЕР", b3_sub="localhost:8080",
        run="ОДНА КОМАНДА", run_sub="go run main.go",
        out="ОТВЕТ", out_sub="200 · Сәлем",
    ),
    "en": dict(
        alt="The program knows what to answer and where to listen; the browser sees the answer",
        b1="WHAT TO ANSWER", b1_sub="http.HandleFunc",
        b2="WHERE TO LISTEN", b2_sub="http.ListenAndServe",
        b3="BROWSER", b3_sub="localhost:8080",
        run="ONE COMMAND", run_sub="go run main.go",
        out="RESPONSE", out_sub="200 · Salem",
    ),
}

L02 = {
    "kz": dict(
        alt="Редактор, терминал және Go: қайсысы не үшін жауап береді",
        b1="РЕДАКТОР", b1_sub="VS Code — мәтін жазасыз",
        b2="ТЕРМИНАЛ", b2_sub="go run . — пәрмен бересіз",
        b3="GO", b3_sub="жинайды және іске қосады",
        head="ҚАЙСЫСЫ НЕ ҮШІН ЖАУАП БЕРЕДІ", head_sub="үш құрал, үш жауапкершілік",
        dir="ЖОБА ҚАЛТАСЫ", dir_sub="go-oqu/sabaq-01 · go.mod + main.go",
    ),
    "ru": dict(
        alt="Редактор, терминал и Go: кто за что отвечает",
        b1="РЕДАКТОР", b1_sub="VS Code — пишете текст",
        b2="ТЕРМИНАЛ", b2_sub="go run . — отдаёте команду",
        b3="GO", b3_sub="собирает и запускает",
        head="КТО ЗА ЧТО ОТВЕЧАЕТ", head_sub="три инструмента, три зоны вины",
        dir="ПАПКА ПРОЕКТА", dir_sub="go-oqu/sabaq-01 · go.mod + main.go",
    ),
    "en": dict(
        alt="Editor, terminal and Go: which one is answerable for what",
        b1="EDITOR", b1_sub="VS Code — you write text",
        b2="TERMINAL", b2_sub="go run . — you give the order",
        b3="GO", b3_sub="builds it and runs it",
        head="WHICH ONE IS ANSWERABLE", head_sub="three tools, three kinds of fault",
        dir="PROJECT FOLDER", dir_sub="go-oqu/sabaq-01 · go.mod + main.go",
    ),
}

L03 = {
    "kz": dict(
        alt="Төрт тип: мәтін, бүтін сан, бөлшек сан және иә-жоқ",
        t1="STRING", t1_val="мәтін · «Сәлем, әлем!»",
        t2="INT", t2_val="бүтін сан · 0",
        t3="FLOAT64", t3_val="бөлшек сан · 4.5",
        t4="BOOL", t4_val="иә-жоқ · true",
        head="ӘЗІРГЕ ЖЕТЕТІН ТӨРТ ТИП", head_sub="title · views · rating · published",
        foot="ТИП БІР РЕТ ТАҢДАЛАДЫ", foot_sub='views = "көп"  →  жинақтау қатесі',
    ),
    "ru": dict(
        alt="Четыре типа: текст, целое число, дробное число и да-нет",
        t1="STRING", t1_val="текст · «Сәлем, әлем!»",
        t2="INT", t2_val="целое · 0",
        t3="FLOAT64", t3_val="дробное · 4.5",
        t4="BOOL", t4_val="да-нет · true",
        head="ЧЕТЫРЕ ТИПА, КОТОРЫХ ПОКА ХВАТИТ", head_sub="title · views · rating · published",
        foot="ТИП ВЫБИРАЕТСЯ ОДИН РАЗ", foot_sub='views = "много"  →  ошибка при сборке',
    ),
    "en": dict(
        alt="Four types: text, whole number, fractional number and yes-no",
        t1="STRING", t1_val="text · \u0022Salem, alem!\u0022",
        t2="INT", t2_val="whole · 0",
        t3="FLOAT64", t3_val="fractional · 4.5",
        t4="BOOL", t4_val="yes-no · true",
        head="FOUR TYPES THAT WILL DO FOR NOW", head_sub="title · views · rating · published",
        foot="A TYPE IS CHOSEN ONCE", foot_sub='views = "a lot"  \u2192  build error',
    ),
}

L04 = {
    "kz": dict(
        alt="Бір шақыру — екі жауап: нәтиже және шықты ма деген белгі",
        in_="", fn="ФУНКЦИЯ", fn_sub="divide",
        o1="НӘТИЖЕ", o1_sub="2.5",
        o2="ШЫҚТЫ МА", o2_sub="true",
        head="БІР ШАҚЫРУ — ЕКІ ЖАУАП", head_sub="result, ok := divide(10, 4)",
        foot="ЕКІНШІ ЖАУАПТЫ ЕЛЕМЕУГЕ БОЛМАЙДЫ", foot_sub="result, _ := divide(10, 4)",
    ),
    "ru": dict(
        alt="Один вызов — два ответа: результат и признак, получилось ли",
        fn="ФУНКЦИЯ", fn_sub="divide",
        o1="РЕЗУЛЬТАТ", o1_sub="2.5",
        o2="ПОЛУЧИЛОСЬ ЛИ", o2_sub="true",
        head="ОДИН ВЫЗОВ — ДВА ОТВЕТА", head_sub="result, ok := divide(10, 4)",
        foot="ВТОРОЙ ОТВЕТ ИГНОРИРОВАТЬ НЕЛЬЗЯ", foot_sub="result, _ := divide(10, 4)",
    ),
    "en": dict(
        alt="One call, two answers: the result and whether it worked",
        fn="FUNCTION", fn_sub="divide",
        o1="RESULT", o1_sub="2.5",
        o2="DID IT WORK", o2_sub="true",
        head="ONE CALL, TWO ANSWERS", head_sub="result, ok := divide(10, 4)",
        foot="THE SECOND ANSWER CANNOT BE IGNORED", foot_sub="result, _ := divide(10, 4)",
    ),
}
L04["kz"]["in"] = "НЕ ӘКЕЛДІҢІЗ"; L04["kz"]["in_sub"] = "a, b float64"
L04["ru"]["in"] = "ЧТО ПРИНЕСЛИ"; L04["ru"]["in_sub"] = "a, b float64"
L04["en"]["in"] = "WHAT YOU BRING"; L04["en"]["in_sub"] = "a, b float64"

L05 = {
    "kz": dict(alt="Бес сабақ жарияланды, жол әрі қарай созылады",
               head="ТЕГІН КУРС: GO НӨЛДЕН", head_sub="43 сабақ · тіркелусіз · үш тілде",
               foot="БЕС САБАҚ ЖАРИЯЛАНДЫ", foot_sub="қалғаны — жазылу ретімен"),
    "ru": dict(alt="Пять уроков опубликованы, дорога идёт дальше",
               head="БЕСПЛАТНЫЙ КУРС: GO С НУЛЯ", head_sub="43 урока · без регистрации · три языка",
               foot="ОПУБЛИКОВАНО ПЯТЬ", foot_sub="остальные — по мере написания"),
    "en": dict(alt="Five lessons published, the road keeps going",
               head="A FREE COURSE: GO FROM SCRATCH", head_sub="43 lessons · no account · three languages",
               foot="FIVE ARE PUBLISHED", foot_sub="the rest as they are written"),
}

if __name__ == "__main__":
    out = os.path.join(os.path.dirname(__file__), "..", "..", "web", "static", "course", "go")
    out = os.path.normpath(out)
    os.makedirs(out, exist_ok=True)
    for name, scene, table in (("internet", map00, L), ("first-server", map01, L01),
                               ("workspace", map02, L02),
                               ("types", map03, L03),
                               ("functions", map04, L04),
                               ("launch", map05, L05)):
        for lang, strings in table.items():
            path = os.path.join(out, f"map-{name}-{lang}.svg")
            with open(path, "w", encoding="utf-8") as f:
                f.write(render(scene, strings))
            print("wrote", path)
