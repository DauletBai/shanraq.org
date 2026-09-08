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
import re

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

# One canvas for every map in the course. The maps used to be framed around
# their own contents, so each came out a different size and shape, the grid
# changed density from lesson to lesson, and the captions sat wherever the
# blocks left room. On a page that is a picture jumping about; for a reader it
# is a new frame to re-read every time.
#
# So the frame is fixed and the scene lives inside it: 16:9, one grid, and the
# two captions always at the same height. What changes between lessons is what
# stands on the grid -- which is the only thing that should.
CANVAS_W, CANVAS_H = 1216.0, 684.0
ORIGIN_Y = 152.0            # where iso(at(0, 0), 0) lands on screen
# The ground fills the frame edge to edge and stops just short of the two
# captions: in screen terms u only moves a point sideways and side only moves
# it up and down, so these two numbers are the rectangle itself.
GRID_U, GRID_S = 9.24, 5.7

# Caption baselines, as offsets from the centre of the canvas.
CAP_TOP, CAP_TOP_SUB = -278.0, -251.0
CAP_BOT, CAP_BOT_SUB = 256.0, 283.0

# Everything a scene draws has to stay inside this, or it collides with a
# caption or runs off the edge. render() checks and says which map broke it.
SAFE_X, SAFE_Y_UP, SAFE_Y_DOWN = 594.0, -232.0, 228.0

# How far a scene may be grown to fill the frame. Without a cap a map with
# three blocks would swell until its lettering dwarfed every other lesson's.
ZOOM_MAX = 1.25

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


def slab(u, z, half_u, half_side, h, top, left, right, side=0.0):
    """A block that can be wider along the road than across it.

    The side step is what makes a stack of these compact: a wide, shallow slab
    drops much less on screen than a cube of the same width, so several fit one
    above another without the gap swallowing the picture.
    """
    x, y = at(u, side)
    x -= half_u
    y -= half_side
    w, d = half_u * 2, half_side * 2
    p = lambda dx, dy, dz: iso(x + dx, y + dy, z + dz)
    return "".join([
        poly([p(0, 0, h), p(w, 0, h), p(w, d, h), p(0, d, h)], top),
        poly([p(0, d, 0), p(w, d, 0), p(w, d, h), p(0, d, h)], right),
        poly([p(w, 0, 0), p(w, d, 0), p(w, d, h), p(w, 0, h)], left),
    ])


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


def text_px(dy_from_centre, body, cls, dx=0.0):
    """A caption at a fixed place on the canvas, not measured off the blocks."""
    return (f'<text class="{cls}" x="{dx:.2f}" y="{ORIGIN_Y + dy_from_centre:.2f}">'
            f"{esc(body)}</text>")


def frame(s):
    """The grid and the four caption lines every map shares."""
    return "".join([
        ground(-GRID_U, GRID_U, -GRID_S, GRID_S, "g"),
        text_px(CAP_TOP, s["head"], "cap"),
        text_px(CAP_TOP_SUB, s["head_sub"], "mono"),
        text_px(CAP_BOT, s["foot"], "cap"),
        text_px(CAP_BOT_SUB, s["foot_sub"], "mono"),
    ])


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
/* Sizes are for the file, not the screen: a map 1216px wide is shown in an
   article at roughly half that, so everything here lands at about 0.6 of what
   it says. At the old sizes the caption arrived as 7px and the code line as 6,
   which is a squint, and a squint over forty lessons is a headache.

   The captions carry most of the increase because they stand in open space.
   A label on a block is held to what the block's face can hold: past that it
   overhangs the edges and reads as floating text rather than a name. */
.lbl{fill:var(--ink);font-size:18px;font-weight:600;text-anchor:middle}
.sub{fill:var(--soft);font-size:14px;text-anchor:middle}
.tag{fill:var(--tagink);font-size:14px;font-weight:700;text-anchor:middle;letter-spacing:.3px}
.cap{fill:var(--ink);font-size:21px;font-weight:700;text-anchor:middle}\n.lbl-acc{fill:var(--tagink);font-size:18px;font-weight:600;text-anchor:middle}\n.sub-acc{fill:var(--tagink);opacity:.85;font-size:14px;text-anchor:middle}
.mono{font-family:ui-monospace,SFMono-Regular,Menlo,monospace;font-size:18px;fill:var(--soft);text-anchor:middle}
"""


def map00(s):
    """Lesson 0: a client asks, a server answers.

    The two roads run at different heights on purpose: a reader has to see that
    the question and the answer are separate journeys, not one line with arrows
    at both ends.
    """
    parts = []

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

    parts.append(on_face(-5.3, 1.7, s["client"], s["client_sub"]))
    parts.append(on_face(5.3, 2.95, s["server"], s["server_sub"]))
    return "".join(parts)


def map02(s):
    """Lesson: the three tools and what each is answerable for.

    A beginner meets the editor, the terminal and the language on the same day
    and blames the wrong one when something breaks. The map exists to fix which
    is which before any of them misbehaves.
    """
    parts = []

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
    return "".join(parts)


def map01(s):
    """Lesson 1: the two halves of a first program, and where to look at it.

    Three blocks in a row rather than a stack: in isometric a block standing on
    another hides the lower one's top face, and both halves have to carry a
    name. Read left to right they are also the order the program runs in.
    """
    parts = []

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
    return "".join(parts)


def map03(s):
    """Lesson: the four types that carry a blog, and that a type is chosen once.

    Four boxes rather than a list, because the point is that each value sits in
    a container of a fixed kind. The caption underneath is the half beginners
    trip over: the box does not change shape later.
    """
    parts = []

    half, top = 1.6, 1.15
    seats = (-5.7, -1.9, 1.9, 5.7)
    keys = ("t1", "t2", "t3", "t4")
    for u, k in zip(seats, keys):
        parts.append(block(u, 0, half, top, "t", "l", "r"))
        parts.append(on_face(u, top, s[k], s[k + "_val"]))
    return "".join(parts)


def map04(s):
    """Lesson: a call goes in on one line and two answers come back.

    The two results sit side by side across the road rather than in a queue
    along it: in a row they would read as one following the other, and the
    whole point is that both arrive from the same call.
    """
    parts = []

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
    return "".join(parts)


def map05(s):
    """The announcement cover: the arc of the course, and the road out of it.

    The blocks used to be blank, on the argument that numbering them would claim
    a lesson count that goes stale. It does — but five anonymous cubes say
    nothing at all, which is worse. Named after the course's modules they carry
    the shape of the whole thing. The road past the last block used to mean
    lessons still unwritten; the course is written now, and it means what comes
    after it: the reader's own blog.
    """
    parts = []

    half, top = 1.0, 0.95
    for u, key in ((-5.4, "m1"), (-3.2, "m2"), (-1.0, "m3"), (1.2, "m4"), (3.4, "m5")):
        parts.append(block(u, 0, half, top, "gt", "gl", "gr"))
        parts.append(text(u, top, s[key], "tag", dy=4.0))

    # The road past the last block: what the course leads to rather than what
    # is left of it. Drawn as a path with nothing standing on it, because what
    # stands there is built by the reader.
    parts.append(road(4.6, 7.4, top * 0.5, 0.44, "band-flow"))
    for u in (5.3, 6.2, 7.1):
        parts.append(chevron(u, top * 0.5 + 0.02, +1, "arw-flow"))
    return "".join(parts)


def map06(s):
    """Lesson: one loop word, three shapes.

    Three blocks rather than a list, because the point is that they are the same
    word wearing different clothes — and the caption underneath names the three
    keywords Go does not have, which is what the reader is really being told.
    """
    parts = []

    half, top = 1.6, 1.15
    seats = (-5.0, 0.0, 5.0)
    for u, k in zip(seats, ("f1", "f2", "f3")):
        parts.append(block(u, 0, half, top, "t", "l", "r"))
        parts.append(on_face(u, top, s[k], s[k + "_sub"]))
    for a, b in ((-3.3, -1.7), (1.7, 3.3)):
        parts.append(road(a, b, top * 0.55, 0.42, "band-flow"))
    return "".join(parts)


def map07(s):
    """Lesson: a byte is not a letter.

    Seven letters, and the wide ones cost two bytes where the narrow ones cost
    one. The widths are the argument: a row of equal boxes would illustrate the
    opposite of what the lesson says.
    """
    parts = []
    # G o (space) т і л і — the Latin half one byte each, the Kazakh half two.
    letters = [("G", 1), ("o", 1), ("␣", 1), ("т", 2), ("і", 2), ("л", 2), ("і", 2)]
    unit, gap, top = 0.62, 0.16, 1.0
    total = sum(n * unit * 2 for _, n in letters) + gap * (len(letters) - 1)
    u = -total / 2
    offset = 0
    marks = []
    for ch, n in letters:
        halfu = n * unit
        centre = u + halfu
        wide = n == 2
        parts.append(slab(centre, 0, halfu, 0.62, top,
                          "rt" if wide else "t", "rl" if wide else "l", "rr" if wide else "r"))
        parts.append(text(centre, top, ch, "lbl-acc" if wide else "lbl", dy=4.0))
        marks.append((centre, offset))
        offset += n
        u += halfu * 2 + gap

    # The byte offset under each letter: the number range hands back, and the
    # reason it jumps by two. Below the blocks, not behind them — set level with
    # the row it belongs to, half of them disappeared behind the fronts.
    offsets_at = clear_below(0, 0.62, gap_px=30)
    for centre, off in marks:
        parts.append(text(centre, offsets_at, str(off), "mono"))

    return "".join(parts)


def map08(s):
    """Lesson: length is what is in it, capacity is what fits before it moves.

    Three full cells and one empty one, because the whole idea is the gap
    between the two numbers. Drawn at len == cap the picture would say nothing.
    """
    parts = []

    half, top = 1.15, 1.0
    seats = (-3.6, -1.2, 1.2, 3.6)
    for i, u in enumerate(seats):
        full = i < 3
        parts.append(block(u, 0, half, top,
                           "gt" if full else "t", "gl" if full else "l", "gr" if full else "r"))
        parts.append(text(u, top, s["c%d" % (i + 1)],
                          "lbl-acc" if full else "sub", dy=4.0))
    return "".join(parts)


def map09(s):
    """Lesson: a key hands you its value at once, and a missing key hands zero.

    The palette is the one lesson 0 set: the question is red, the answer green.
    Three keys stand over the value they open; the fourth stands over a flat
    neutral tile with a nought on it, because a key that is not there is not an
    error -- it is an answer nobody put anything into.
    """
    parts = []

    half, top = 1.15, 0.9
    tile, tile_h, tile_z = 0.9, 0.45, 1.05
    seats = (-3.6, -1.2, 1.2, 3.6)
    vals = ("v1", "v2", "v3", "v4")
    keys = ("k1", "k2", "k3", "k4")

    # The value goes UNDER its box, not on its top face: the key tile sits over
    # that face and hides whatever is written there. Read downwards, a column is
    # the whole idea in three lines -- key, box, value.
    line_val = clear_below(0, half, gap_px=16.0)
    for i, u in enumerate(seats):
        if i < 3:
            parts.append(block(u, 0, half, top, "gt", "gl", "gr"))
        else:
            # Flat, not a cube: nothing was ever put here, and the nought under
            # it is the whole point of the lesson's second half.
            parts.append(block(u, 0, half, 0.16, "t", "l", "r"))
        parts.append(block(u, tile_z, tile, tile_h, "rt", "rl", "rr"))
        parts.append(text(u, tile_z + tile_h, s[keys[i]], "tag", dy=4.0))
        parts.append(text(u, line_val, s[vals[i]], "lbl"))
    return "".join(parts)


def map10(s):
    """Lesson: the test stands beside the code, calls it, and answers twice.

    Red asks and green answers, the same way round as lesson 0: here the test
    is the question and the program is what answers it. The two outcomes are
    drawn as tiles rather than described, because a beginner meets FAIL long
    before PASS and needs to recognise it as a report, not a punishment.
    """
    parts = []

    half, top = 1.3, 1.0

    # The test calls the code: a band pointing back down the road, laid before
    # the blocks so it runs behind them rather than across their faces.
    parts.append(road(-2.4, -0.3, 0.55, 0.5, "band-req"))
    for u in (-2.0, -1.0):
        parts.append(chevron(u, 0.57, -1, "arw-req"))

    parts.append(block(-3.7, 0, half, top, "gt", "gl", "gr"))
    parts.append(on_face_side(-3.7, 0.0, top, s["code"], s["code_sub"], accent=True))

    parts.append(block(0.0, 0, half, top, "rt", "rl", "rr"))
    parts.append(on_face_side(0.0, 0.0, top, s["test"], s["test_sub"], accent=True))

    # The two outcomes, one above the other on screen. They are set apart across
    # the road rather than stacked in z: a step sideways moves a block straight
    # down the picture without moving it along, so the pair reads as a choice
    # between two answers and neither tile covers the other's face.
    tile, tile_h = 0.85, 0.4
    for side, key, palette in ((-1.9, "good", ("gt", "gl", "gr")), (1.9, "bad", ("rt", "rl", "rr"))):
        parts.append(block(3.9, 0, tile, tile_h, *palette, side=side))
        parts.append(text(3.9, tile_h, s[key], "tag", side=side, dy=4.0))
    return "".join(parts)


def map11(s):
    """Preface: three things you need, and one you do not.

    The same row the lesson maps use, so a reader who has seen one recognises
    the grammar before reading a word: green is what you must have, the flat
    neutral tile is the thing people wrongly think they need and turn back over.
    """
    parts = []

    half, top = 1.15, 0.9
    seats = (-3.6, -1.2, 1.2, 3.6)
    line_lbl = clear_below(0, half, gap_px=16.0)
    for i, u in enumerate(seats):
        need = i < 3
        if need:
            parts.append(block(u, 0, half, top, "gt", "gl", "gr"))
            parts.append(text(u, top, s["n%d" % (i + 1)], "lbl-acc", dy=4.0))
        else:
            parts.append(block(u, 0, half, 0.16, "t", "l", "r"))
            parts.append(text(u, 0.16, s["n4"], "sub", dy=4.0))
        parts.append(text(u, line_lbl, s["s%d" % (i + 1)], "sub"))
    return "".join(parts)


def map12(s):
    """Lesson: three parallel maps collapse into one type.

    The left side is what the blog looked like before -- a title here, a word
    count there, a language somewhere else, all kept in step by hand. The right
    side is the same data as one value. The footer carries the trap that costs
    beginners an afternoon: range hands you a copy, so writing to it changes
    nothing.
    """
    parts = []

    # Far enough apart that each card is its own card: at the old spacing the
    # three overlapped and read as one leaning stack, which is the opposite of
    # the point — they are three separate things kept in step by hand.
    tile, tile_h = 0.9, 0.4
    for side, key in ((-2.5, "p1"), (0.0, "p2"), (2.5, "p3")):
        parts.append(block(-3.7, 0, tile, tile_h, "t", "l", "r", side=side))
        parts.append(text(-3.7, tile_h, s[key], "sub", side=side, dy=4.0))

    parts.append(road(-2.3, -0.7, 0.2, 0.45, "band-flow"))
    for u in (-1.9, -1.1):
        parts.append(chevron(u, 0.22, +1, "arw-flow"))

    half, top = 1.5, 1.15
    parts.append(block(1.9, 0, half, top, "gt", "gl", "gr"))
    parts.append(on_face_side(1.9, 0.0, top, s["name"], s["fields"], accent=True))
    return "".join(parts)


def map13(s):
    """Lesson: a copy is a dead end, an address leads back to the value.

    Two roads, as in lesson 0, because the reader has to see two separate
    journeys rather than one line with arrows at both ends. The upper one
    carries a copy away from the value and stops at a grey tile: whatever is
    written there stays there. The lower one carries an address back INTO the
    value, and the arrows point that way on purpose.
    """
    parts = []

    # Roads first: drawn after the blocks they would lie across their faces.
    parts.append(road(-2.4, 4.4, 3.15, 0.58, "band-flow"))
    for u in (-1.4, 0.4, 2.2):
        parts.append(chevron(u, 3.17, +1, "arw-flow"))
    parts.append(road(-2.4, 4.4, 0.35, 0.58, "band-res"))
    for u in (-1.2, 0.6, 2.4):
        parts.append(chevron(u, 0.37, -1, "arw-res"))

    # Lifted to sit level with the middle of the two roads: standing on the
    # ground it hung below both of them and looked like a third thing rather
    # than the value they lead to and from.
    half, top = 1.45, 1.6
    base = 0.95
    parts.append(block(-4.3, base, half, top, "gt", "gl", "gr"))
    parts.append(on_face_side(-4.3, 0.0, base + top, s["value"], s["value_sub"], accent=True))

    # The two roads end at two tiles, one above the other, so the gap between
    # them has to clear a whole tile — height and both halves. At the first
    # spacing the white one sat on the red one's shoulder.
    tile, tile_h = 0.95, 0.42
    parts.append(block(5.6, 3.0, tile, tile_h, "t", "l", "r"))
    parts.append(text(5.6, 3.0 + tile_h, s["copy"], "lbl", dy=4.0))

    parts.append(block(5.6, 0.0, tile, tile_h, "rt", "rl", "rr"))
    parts.append(text(5.6, tile_h, s["addr"], "tag", dy=4.0))

    # What each road does is written on the road, not on the tile it ends at:
    # a caption laid over a tile's face fights the fill for contrast and loses,
    # and on the red one it lost badly.
    parts.append(text(1.0, 3.72, s["copy_sub"], "sub"))
    parts.append(text(1.0, 0.92, s["addr_sub"], "sub"))
    return "".join(parts)


def map14(s):
    """Lesson: one function, any store, because the contract is a list of skills.

    The code that does the work stands on the left and knows nothing about what
    is on the right. Between them is the contract -- two method names on a red
    tile, red because it is the thing being asked for. Three stores satisfy it,
    and one of them exists only in a test, which is the point that sells
    interfaces to a beginner.
    """
    parts = []

    parts.append(road(-2.5, 0.6, 0.55, 0.5, "band-req"))
    for u in (-2.0, -1.0, 0.0):
        parts.append(chevron(u, 0.57, +1, "arw-req"))

    half, top = 1.4, 1.4
    parts.append(block(-4.0, 0, half, top, "gt", "gl", "gr"))
    parts.append(on_face_side(-4.0, 0.0, top, s["fn"], s["fn_sub"], accent=True))

    gate, gate_h = 1.3, 1.1
    parts.append(block(1.7, 0, gate, gate_h, "rt", "rl", "rr"))
    parts.append(on_face_side(1.7, 0.0, gate_h, s["iface"], s["iface_sub"], accent=True))

    tile, tile_h = 1.0, 0.42
    for side, key in ((-2.7, "s1"), (0.0, "s2"), (2.7, "s3")):
        parts.append(block(5.0, 0, tile, tile_h, "t", "l", "r", side=side))
        parts.append(text(5.0, tile_h, s[key], "sub", side=side, dy=4.0))
    return "".join(parts)


def map15(s):
    """Lesson: a function has two ways out, and the second is checked first.

    Two roads again, because by now the reader reads them without help: the
    upper green one carries the value, the lower red one carries the error.
    They leave the same function, which is the point -- an error in Go is not a
    siren somewhere else in the building, it is the second thing handed back.
    """
    parts = []

    parts.append(road(-1.7, 3.4, 2.15, 0.55, "band-res"))
    for u in (-1.2, 0.2, 1.6):
        parts.append(chevron(u, 2.17, +1, "arw-res"))
    parts.append(road(-1.7, 3.4, 0.35, 0.55, "band-req"))
    for u in (-1.2, 0.2, 1.6):
        parts.append(chevron(u, 0.37, +1, "arw-req"))

    half, top = 1.45, 1.5
    parts.append(block(-3.3, 0, half, top, "t", "l", "r"))
    parts.append(on_face_side(-3.3, 0.0, top, s["fn"], s["fn_sub"]))

    tile, tile_h = 1.15, 0.48
    parts.append(block(4.6, 2.45, tile, tile_h, "gt", "gl", "gr"))
    parts.append(text(4.6, 2.45 + tile_h, s["good"], "tag", dy=4.0))
    parts.append(block(4.6, 0.0, tile, tile_h, "rt", "rl", "rr"))
    parts.append(text(4.6, tile_h, s["bad"], "tag", dy=4.0))

    parts.append(text(0.9, 2.72, s["good_sub"], "sub"))
    parts.append(text(0.9, 0.92, s["bad_sub"], "sub"))
    return "".join(parts)


def map16(s):
    """Lesson: a package is a folder, and a capital letter is its door.

    main on the left knows only what the package lets out. The two tiles beside
    the package say which is which: green for what a capital letter exports,
    grey for what stays inside. Drawn this way the export rule stops being a
    style convention and becomes a wall you can see.
    """
    parts = []

    parts.append(road(-2.1, 0.6, 0.5, 0.5, "band-flow"))
    for u in (-1.6, -0.7, 0.2):
        parts.append(chevron(u, 0.52, +1, "arw-flow"))

    half, top = 1.35, 1.3
    parts.append(block(-3.6, 0, half, top, "t", "l", "r"))
    parts.append(on_face_side(-3.6, 0.0, top, s["main"], s["main_sub"]))

    parts.append(block(1.9, 0, half, top, "gt", "gl", "gr"))
    parts.append(on_face_side(1.9, 0.0, top, s["pkg"], s["pkg_sub"], accent=True))

    # Wide enough for the names written on them: a label that overhangs its own
    # tile reads as two labels, which is what the first render looked like.
    tile, tile_h = 1.5, 0.42
    parts.append(block(5.2, 0, tile, tile_h, "gt", "gl", "gr", side=-2.2))
    parts.append(text(5.2, tile_h, s["out"], "tag", side=-2.2, dy=4.0))
    parts.append(block(5.2, 0, tile, tile_h, "t", "l", "r", side=2.2))
    parts.append(text(5.2, tile_h, s["inn"], "sub", side=2.2, dy=4.0))
    return "".join(parts)


def map17(s):
    """Lesson: two steps, not one -- add chooses, commit records.

    Beginners lose files between those two commands, so the map draws them as
    two separate journeys along one road rather than one arrow. Green is what is
    kept forever; the middle block is the waiting room and neither one thing nor
    the other, which is exactly how the index behaves.
    """
    parts = []

    for u0, u1, mid in ((-2.5, -1.0, -1.75), (1.0, 2.5, 1.75)):
        parts.append(road(u0, u1, 0.55, 0.5, "band-flow"))
        parts.append(chevron(mid, 0.57, +1, "arw-flow"))

    half, top = 1.3, 1.15
    parts.append(block(-4.0, 0, half, top, "t", "l", "r"))
    parts.append(on_face_side(-4.0, 0.0, top, s["work"], s["work_sub"]))

    parts.append(block(0.0, 0, half, top, "rt", "rl", "rr"))
    parts.append(on_face_side(0.0, 0.0, top, s["index"], s["index_sub"], accent=True))

    parts.append(block(4.0, 0, half, top, "gt", "gl", "gr"))
    parts.append(on_face_side(4.0, 0.0, top, s["hist"], s["hist_sub"], accent=True))

    # Above the blocks, not level with them: at block height the second label
    # started exactly where the red block ends and read as part of it.
    parts.append(text(-1.75, 1.75, s["cmd1"], "mono"))
    parts.append(text(1.75, 1.75, s["cmd2"], "mono"))
    return "".join(parts)


def map18(s):
    """Lesson: the same history, in two places, kept in step by two commands.

    Not a backup and not a different thing -- the same records, one copy on the
    desk and one on the network. push sends, pull brings back, and the two roads
    run in opposite directions because that is the whole of it.
    """
    parts = []

    parts.append(road(-2.2, 2.2, 2.1, 0.55, "band-res"))
    for u in (-1.5, 0.0, 1.5):
        parts.append(chevron(u, 2.12, +1, "arw-res"))
    parts.append(road(-2.2, 2.2, 0.35, 0.55, "band-req"))
    for u in (-1.5, 0.0, 1.5):
        parts.append(chevron(u, 0.37, -1, "arw-req"))

    half, top = 1.45, 1.45
    parts.append(block(-4.2, 0, half, top, "gt", "gl", "gr"))
    parts.append(on_face_side(-4.2, 0.0, top, s["local"], s["local_sub"], accent=True))

    parts.append(block(4.2, 0, half, top, "t", "l", "r"))
    parts.append(on_face_side(4.2, 0.0, top, s["remote"], s["remote_sub"]))

    parts.append(text(0, 2.72, s["push"], "mono"))
    parts.append(text(0, 0.92, s["pull"], "mono"))
    return "".join(parts)


def map19(s):
    """Capstone: nothing new on the table, everything already learned.

    Six tiles for six lessons and one green block for the program they make.
    The point of the picture is that the tiles are ordinary and grey -- the
    reader has held every one of them before, and the only new thing today is
    that they fit together.
    """
    parts = []

    parts.append(road(-1.0, 1.0, 0.5, 0.5, "band-flow"))
    for u in (-0.6, 0.2):
        parts.append(chevron(u, 0.52, +1, "arw-flow"))

    # Spread far enough that no tile lands on its neighbour: a step across the
    # road moves a block straight down the screen by that step, so the gap has
    # to clear the tile's own height. The first draft did not, and the six
    # familiar parts arrived looking like a pile of rubble.
    tile, tile_h = 1.0, 0.4
    left = (("p1", -4.9, -3.1), ("p2", -4.9, 0.0), ("p3", -4.9, 3.1),
            ("p4", -2.3, -3.1), ("p5", -2.3, 0.0), ("p6", -2.3, 3.1))
    for key, u, side in left:
        parts.append(block(u, 0, tile, tile_h, "t", "l", "r", side=side))
        parts.append(text(u, tile_h, s[key], "sub", side=side, dy=4.0))

    half, top = 1.6, 1.5
    parts.append(block(2.8, 0, half, top, "gt", "gl", "gr"))
    parts.append(on_face_side(2.8, 0.0, top, s["prog"], s["prog_sub"], accent=True))
    return "".join(parts)


def map20(s):
    """Lesson: a response goes out in one direction and does not come back.

    Three blocks in the order they must happen, and the middle one is red
    because it is the point of no return: once the status code is out, the
    headers are already on their way and nothing written after it will travel.
    """
    parts = []

    for u0, u1, mid in ((-2.6, -1.1, -1.85), (1.1, 2.6, 1.85)):
        parts.append(road(u0, u1, 0.55, 0.5, "band-flow"))
        parts.append(chevron(mid, 0.57, +1, "arw-flow"))

    half, top = 1.35, 1.15
    parts.append(block(-4.2, 0, half, top, "gt", "gl", "gr"))
    parts.append(on_face_side(-4.2, 0.0, top, s["h1"], s["h1_sub"], accent=True))

    parts.append(block(0.0, 0, half, top, "rt", "rl", "rr"))
    parts.append(on_face_side(0.0, 0.0, top, s["h2"], s["h2_sub"], accent=True))

    parts.append(block(4.2, 0, half, top, "t", "l", "r"))
    parts.append(on_face_side(4.2, 0.0, top, s["h3"], s["h3_sub"]))
    return "".join(parts)


def map21(s):
    """Lesson: one entrance, several labelled doors, and the router decides.

    The routes stand apart on the right with the method written into each,
    because in Go 1.22 the method is part of the pattern rather than an if at
    the top of the handler. What the reader wrote by hand last lesson is now
    the box in the middle.
    """
    parts = []

    parts.append(road(-2.6, -0.6, 0.5, 0.5, "band-flow"))
    for u in (-2.1, -1.2):
        parts.append(chevron(u, 0.52, +1, "arw-flow"))

    half, top = 1.35, 1.2
    parts.append(block(-4.2, 0, half, top, "t", "l", "r"))
    parts.append(on_face_side(-4.2, 0.0, top, s["req"], s["req_sub"]))

    parts.append(block(0.6, 0, half, top, "rt", "rl", "rr"))
    parts.append(on_face_side(0.6, 0.0, top, s["mux"], s["mux_sub"], accent=True))

    # Two routes rather than three: a third slab either overlapped its
    # neighbours or had to be spread so far that the picture went hollow, and
    # two are enough to show that the method rides in the pattern. The one that
    # did not fit is named in the caption instead.
    tile_u, tile_s, tile_h = 1.95, 1.0, 0.4
    for side, key in ((-1.75, "r1"), (1.75, "r2")):
        parts.append(slab(4.4, 0, tile_u, tile_s, tile_h, "gt", "gl", "gr", side=side))
        parts.append(text(4.4, tile_h, s[key], "tag", side=side, dy=4.0))
    return "".join(parts)


def map22(s):
    """Lesson: a wrapper the request passes through twice.

    The row is what the request meets on the way in; the lower road is the same
    row on the way out. Drawing both makes the point that the wrapper runs code
    before and after the handler -- which is the whole reason it can time a
    response or record its status code.
    """
    parts = []

    # Clear of the row, not behind it: laid level with the blocks the bands ran
    # behind them and six of the eight arrows disappeared, which made the
    # journey look like it stopped at the first box.
    parts.append(road(-5.9, 5.9, 3.75, 0.55, "band-req"))
    for u in (-4.4, -1.6, 1.2, 4.0):
        parts.append(chevron(u, 3.77, +1, "arw-req"))
    parts.append(road(-5.9, 5.9, -0.7, 0.55, "band-res"))
    for u in (-4.0, -1.2, 1.6, 4.4):
        parts.append(chevron(u, -0.68, -1, "arw-res"))

    half, top, base = 1.2, 1.0, 1.15
    row = ((-4.4, "b1", ("t", "l", "r"), False),
           (-1.5, "b2", ("rt", "rl", "rr"), True),
           (1.5, "b3", ("t", "l", "r"), False),
           (4.4, "b4", ("gt", "gl", "gr"), True))
    for u, key, palette, accent in row:
        parts.append(block(u, base, half, top, *palette))
        parts.append(on_face_side(u, 0.0, base + top, s[key], s[key + "_sub"], accent=accent))
    return "".join(parts)


# Rough width of a rendered string, in px. Enough to catch a label that has
# outgrown the block it names -- exact metrics would need the font itself.
_CHAR_W = {"lbl": 0.56, "lbl-acc": 0.56, "sub": 0.52, "sub-acc": 0.52,
           "tag": 0.58, "cap": 0.58, "mono": 0.60}
_FONT_PX = {"lbl": 18, "lbl-acc": 18, "sub": 14, "sub-acc": 14,
            "tag": 14, "cap": 21, "mono": 18}


def map23(s):
    """Lesson: a form and some data make a page.

    Two roads meeting one block, because that is the whole shape of it: the
    template is written once and the data changes on every request. The footer
    carries what html/template does that a beginner would otherwise have to be
    bitten by first.
    """
    parts = []

    # Two roads meeting: the template comes down from above, the data up from
    # below, and both end at the same page.
    parts.append(road(-2.4, 1.0, 2.75, 0.5, "band-flow"))
    for u in (-1.9, -0.9, 0.1):
        parts.append(chevron(u, 2.77, +1, "arw-flow"))
    parts.append(road(-2.4, 1.0, -1.65, 0.5, "band-res"))
    for u in (-1.9, -0.9, 0.1):
        parts.append(chevron(u, -1.63, +1, "arw-res"))

    # Far enough apart to clear each other: two blocks stacked in z need a gap
    # wider than a whole block, both halves and its height together.
    half, top = 1.3, 1.05
    parts.append(block(-4.0, 2.2, half, top, "t", "l", "r"))
    parts.append(on_face_side(-4.0, 0.0, 2.2 + top, s["tpl"], s["tpl_sub"]))

    parts.append(block(-4.0, -2.2, half, top, "gt", "gl", "gr"))
    parts.append(on_face_side(-4.0, 0.0, -2.2 + top, s["data"], s["data_sub"], accent=True))

    parts.append(block(3.1, 0.0, 1.5, 1.35, "gt", "gl", "gr"))
    parts.append(on_face_side(3.1, 0.0, 1.35, s["page"], s["page_sub"], accent=True))
    return "".join(parts)


def map24(s):
    """Lesson: the reader sends something, and the answer is a redirect.

    Three stops rather than two, because the third is the point: the handler
    that changed something does not draw a page, it sends the browser to one.
    Without that stop a refresh submits the form again, which is the bug this
    picture exists to prevent.
    """
    parts = []

    parts.append(road(-2.6, -1.45, 0.55, 0.5, "band-req"))
    parts.append(chevron(-2.03, 0.57, +1, "arw-req"))
    parts.append(road(1.75, 2.9, 0.55, 0.5, "band-res"))
    parts.append(chevron(2.33, 0.57, +1, "arw-res"))

    half, top = 1.3, 1.15
    parts.append(block(-4.2, 0, half, top, "t", "l", "r"))
    parts.append(on_face_side(-4.2, 0.0, top, s["b1"], s["b1_sub"]))

    parts.append(block(0.15, 0, half, top, "rt", "rl", "rr"))
    parts.append(on_face_side(0.15, 0.0, top, s["b2"], s["b2_sub"], accent=True))

    parts.append(block(4.5, 0, half, top, "gt", "gl", "gr"))
    parts.append(on_face_side(4.5, 0.0, top, s["b3"], s["b3_sub"], accent=True))

    parts.append(text(-2.03, 2.15, s["c1"], "mono"))
    parts.append(text(2.33, 2.15, s["c2"], "mono"))
    return "".join(parts)


def map25(s):
    """Lesson: a page is assembled, a file is not.

    The road runs straight from the browser to the finished page, and the two
    blocks set off on either side of it are the two ways the server can answer:
    build something this time, or hand over a file exactly as it lies on disk.
    They sit at the same point on the road because they are the same request
    for the reader -- only the work behind them differs.
    """
    parts = []

    parts.append(road(-2.6, 0.2, 0.55, 0.5, "band-req"))
    parts.append(road(0.2, 3.0, 0.55, 0.5, "band-res"))
    parts.append(chevron(-1.7, 0.57, +1, "arw-req"))
    parts.append(chevron(2.1, 0.57, +1, "arw-res"))

    half, top = 1.25, 1.1
    parts.append(block(-4.2, 0, half, top, "t", "l", "r"))
    parts.append(on_face_side(-4.2, 0.0, top, s["b1"], s["b1_sub"]))

    parts.append(block(0.2, 0, half, top, "rt", "rl", "rr", side=-2.9))
    parts.append(on_face_side(0.2, -2.9, top, s["b2"], s["b2_sub"], accent=True))

    parts.append(block(0.2, 0, half, top, "gt", "gl", "gr", side=2.9))
    parts.append(on_face_side(0.2, 2.9, top, s["b3"], s["b3_sub"], accent=True))

    parts.append(block(4.6, 0, half, top, "t", "l", "r"))
    parts.append(on_face_side(4.6, 0.0, top, s["b4"], s["b4_sub"]))

    return "".join(parts)


def map26(s):
    """Lesson: what the reader gets when a handler blows up.

    Three beats along one road, because the reader's experience is a sequence,
    not a set: something breaks, something catches it, and something reaches
    the screen. The footer carries the limit -- once the answer has left, the
    catch is too late.
    """
    parts = []

    parts.append(road(-2.6, -1.45, 0.55, 0.5, "band-req"))
    parts.append(chevron(-2.03, 0.57, +1, "arw-req"))
    parts.append(road(1.75, 2.9, 0.55, 0.5, "band-res"))
    parts.append(chevron(2.33, 0.57, +1, "arw-res"))

    half, top = 1.3, 1.15
    parts.append(block(-4.2, 0, half, top, "rt", "rl", "rr"))
    parts.append(on_face_side(-4.2, 0.0, top, s["b1"], s["b1_sub"], accent=True))

    parts.append(block(0.15, 0, half, top, "t", "l", "r"))
    parts.append(on_face_side(0.15, 0.0, top, s["b2"], s["b2_sub"]))

    parts.append(block(4.5, 0, half, top, "gt", "gl", "gr"))
    parts.append(on_face_side(4.5, 0.0, top, s["b3"], s["b3_sub"], accent=True))

    parts.append(text(-2.03, 2.15, s["c1"], "mono"))
    parts.append(text(2.33, 2.15, s["c2"], "mono"))
    return "".join(parts)


def map27(s):
    """Lesson: one setting, three places it can come from.

    Left to right is the order of precedence, which is also the order the
    program reads them in: the default written in the code, then the variable
    in the environment, then the flag on the command line. The last block is
    accented because the last one wins.
    """
    parts = []

    parts.append(road(-2.6, -1.45, 0.55, 0.5, "band-req"))
    parts.append(chevron(-2.03, 0.57, +1, "arw-req"))
    parts.append(road(1.75, 2.9, 0.55, 0.5, "band-res"))
    parts.append(chevron(2.33, 0.57, +1, "arw-res"))

    half, top = 1.3, 1.15
    parts.append(block(-4.2, 0, half, top, "t", "l", "r"))
    parts.append(on_face_side(-4.2, 0.0, top, s["b1"], s["b1_sub"]))

    parts.append(block(0.15, 0, half, top, "rt", "rl", "rr"))
    parts.append(on_face_side(0.15, 0.0, top, s["b2"], s["b2_sub"], accent=True))

    parts.append(block(4.5, 0, half, top, "gt", "gl", "gr"))
    parts.append(on_face_side(4.5, 0.0, top, s["b3"], s["b3_sub"], accent=True))

    parts.append(text(-2.03, 2.15, s["c1"], "mono"))
    parts.append(text(2.33, 2.15, s["c2"], "mono"))
    return "".join(parts)


def map28(s):
    """Lesson: two tables and the key that ties them.

    The middle block is accented because the foreign key is the whole idea: a
    comment does not contain its article, it points at it by id, and the join
    is what puts them back together for one answer.
    """
    parts = []

    parts.append(road(-2.6, -1.45, 0.55, 0.5, "band-req"))
    parts.append(chevron(-2.03, 0.57, +1, "arw-req"))
    parts.append(road(1.75, 2.9, 0.55, 0.5, "band-res"))
    parts.append(chevron(2.33, 0.57, +1, "arw-res"))

    half, top = 1.3, 1.15
    parts.append(block(-4.2, 0, half, top, "t", "l", "r"))
    parts.append(on_face_side(-4.2, 0.0, top, s["b1"], s["b1_sub"]))

    parts.append(block(0.15, 0, half, top, "rt", "rl", "rr"))
    parts.append(on_face_side(0.15, 0.0, top, s["b2"], s["b2_sub"], accent=True))

    parts.append(block(4.5, 0, half, top, "gt", "gl", "gr"))
    parts.append(on_face_side(4.5, 0.0, top, s["b3"], s["b3_sub"], accent=True))

    parts.append(text(-2.03, 2.15, s["c1"], "mono"))
    parts.append(text(2.33, 2.15, s["c2"], "mono"))
    return "".join(parts)


def map29(s):
    """Lesson: database/sql hands out a pool, not a connection.

    The middle block is the pool, accented because it is the thing beginners do
    not know is there: a pragma sent as a query lands on one connection of it,
    and the rest of the pool never hears about it.
    """
    parts = []

    parts.append(road(-2.6, -1.45, 0.55, 0.5, "band-req"))
    parts.append(chevron(-2.03, 0.57, +1, "arw-req"))
    parts.append(road(1.75, 2.9, 0.55, 0.5, "band-res"))
    parts.append(chevron(2.33, 0.57, +1, "arw-res"))

    half, top = 1.3, 1.15
    parts.append(block(-4.2, 0, half, top, "t", "l", "r"))
    parts.append(on_face_side(-4.2, 0.0, top, s["b1"], s["b1_sub"]))

    parts.append(block(0.15, 0, half, top, "rt", "rl", "rr"))
    parts.append(on_face_side(0.15, 0.0, top, s["b2"], s["b2_sub"], accent=True))

    parts.append(block(4.5, 0, half, top, "gt", "gl", "gr"))
    parts.append(on_face_side(4.5, 0.0, top, s["b3"], s["b3_sub"], accent=True))

    parts.append(text(-2.03, 2.15, s["c1"], "mono"))
    parts.append(text(2.33, 2.15, s["c2"], "mono"))
    return "".join(parts)


def map30(s):
    """Lesson: one table, a row put in and a row taken back out.

    Deliberately the plainest scene in the course, because the lesson it belongs
    to is the plainest idea in the module: a table is a shape, insert puts a row
    into it, select asks for rows back.
    """
    parts = []

    parts.append(road(-2.6, -1.45, 0.55, 0.5, "band-req"))
    parts.append(chevron(-2.03, 0.57, +1, "arw-req"))
    parts.append(road(1.75, 2.9, 0.55, 0.5, "band-res"))
    parts.append(chevron(2.33, 0.57, +1, "arw-res"))

    half, top = 1.3, 1.15
    parts.append(block(-4.2, 0, half, top, "t", "l", "r"))
    parts.append(on_face_side(-4.2, 0.0, top, s["b1"], s["b1_sub"]))

    parts.append(block(0.15, 0, half, top, "rt", "rl", "rr"))
    parts.append(on_face_side(0.15, 0.0, top, s["b2"], s["b2_sub"], accent=True))

    parts.append(block(4.5, 0, half, top, "gt", "gl", "gr"))
    parts.append(on_face_side(4.5, 0.0, top, s["b3"], s["b3_sub"], accent=True))

    parts.append(text(-2.03, 2.15, s["c1"], "mono"))
    parts.append(text(2.33, 2.15, s["c2"], "mono"))
    return "".join(parts)


def map31(s):
    """Lesson: four actions, and the one that answers silently.

    The middle block is accented because it is where the surprise lives: the
    database was asked to change a row that does not exist, and it did exactly
    that -- nothing -- without calling it an error.
    """
    parts = []

    parts.append(road(-2.6, -1.45, 0.55, 0.5, "band-req"))
    parts.append(chevron(-2.03, 0.57, +1, "arw-req"))
    parts.append(road(1.75, 2.9, 0.55, 0.5, "band-res"))
    parts.append(chevron(2.33, 0.57, +1, "arw-res"))

    half, top = 1.3, 1.15
    parts.append(block(-4.2, 0, half, top, "t", "l", "r"))
    parts.append(on_face_side(-4.2, 0.0, top, s["b1"], s["b1_sub"]))

    parts.append(block(0.15, 0, half, top, "rt", "rl", "rr"))
    parts.append(on_face_side(0.15, 0.0, top, s["b2"], s["b2_sub"], accent=True))

    parts.append(block(4.5, 0, half, top, "gt", "gl", "gr"))
    parts.append(on_face_side(4.5, 0.0, top, s["b3"], s["b3_sub"], accent=True))

    parts.append(text(-2.03, 2.15, s["c1"], "mono"))
    parts.append(text(2.33, 2.15, s["c2"], "mono"))
    return "".join(parts)



def map32(s):
    """Lesson: numbered files, a journal, and a schema that changes in steps.

    The middle block is accented because the journal is what makes the whole
    thing work: the files are only files until something remembers which of
    them this database has already seen.
    """
    parts = []

    parts.append(road(-2.6, -1.45, 0.55, 0.5, "band-req"))
    parts.append(chevron(-2.03, 0.57, +1, "arw-req"))
    parts.append(road(1.75, 2.9, 0.55, 0.5, "band-res"))
    parts.append(chevron(2.33, 0.57, +1, "arw-res"))

    half, top = 1.3, 1.15
    parts.append(block(-4.2, 0, half, top, "t", "l", "r"))
    parts.append(on_face_side(-4.2, 0.0, top, s["b1"], s["b1_sub"]))

    parts.append(block(0.15, 0, half, top, "rt", "rl", "rr"))
    parts.append(on_face_side(0.15, 0.0, top, s["b2"], s["b2_sub"], accent=True))

    parts.append(block(4.5, 0, half, top, "gt", "gl", "gr"))
    parts.append(on_face_side(4.5, 0.0, top, s["b3"], s["b3_sub"], accent=True))

    parts.append(text(-2.03, 2.15, s["c1"], "mono"))
    parts.append(text(2.33, 2.15, s["c2"], "mono"))
    return "".join(parts)


def map33(s):
    """Lesson: the value travels beside the query, never inside it.

    The middle block is the accent because that is the whole lesson: the place
    where a value goes in is a place, not a hole in the text.
    """
    parts = []

    parts.append(road(-2.6, -1.45, 0.55, 0.5, "band-req"))
    parts.append(chevron(-2.03, 0.57, +1, "arw-req"))
    parts.append(road(1.75, 2.9, 0.55, 0.5, "band-res"))
    parts.append(chevron(2.33, 0.57, +1, "arw-res"))

    half, top = 1.3, 1.15
    parts.append(block(-4.2, 0, half, top, "t", "l", "r"))
    parts.append(on_face_side(-4.2, 0.0, top, s["b1"], s["b1_sub"]))

    parts.append(block(0.15, 0, half, top, "rt", "rl", "rr"))
    parts.append(on_face_side(0.15, 0.0, top, s["b2"], s["b2_sub"], accent=True))

    parts.append(block(4.5, 0, half, top, "gt", "gl", "gr"))
    parts.append(on_face_side(4.5, 0.0, top, s["b3"], s["b3_sub"], accent=True))

    parts.append(text(-2.03, 2.15, s["c1"], "mono"))
    parts.append(text(2.33, 2.15, s["c2"], "mono"))
    return "".join(parts)


def map34(s):
    """Lesson: two tables and the third one that holds the pairs.

    The middle block is the accent because it is the one the reader would not
    have thought to build: it stores nothing of its own, only meetings.
    """
    parts = []

    parts.append(road(-2.6, -1.45, 0.55, 0.5, "band-req"))
    parts.append(chevron(-2.03, 0.57, +1, "arw-req"))
    parts.append(road(1.75, 2.9, 0.55, 0.5, "band-res"))
    parts.append(chevron(2.2, 0.57, -1, "arw-res"))

    half, top = 1.3, 1.15
    parts.append(block(-4.2, 0, half, top, "t", "l", "r"))
    parts.append(on_face_side(-4.2, 0.0, top, s["b1"], s["b1_sub"]))

    parts.append(block(0.15, 0, half, top, "rt", "rl", "rr"))
    parts.append(on_face_side(0.15, 0.0, top, s["b2"], s["b2_sub"], accent=True))

    parts.append(block(4.5, 0, half, top, "gt", "gl", "gr"))
    parts.append(on_face_side(4.5, 0.0, top, s["b3"], s["b3_sub"], accent=True))

    parts.append(text(-2.03, 2.15, s["c1"], "mono"))
    parts.append(text(2.33, 2.15, s["c2"], "mono"))
    return "".join(parts)


def map35(s):
    """Lesson: what a person typed, the index, and the order of the answer.

    Both later blocks are accented: the index is the thing that replaces a
    scan, and the order is the thing a reader actually judges the search by.
    """
    parts = []

    parts.append(road(-2.6, -1.45, 0.55, 0.5, "band-req"))
    parts.append(chevron(-2.03, 0.57, +1, "arw-req"))
    parts.append(road(1.75, 2.9, 0.55, 0.5, "band-res"))
    parts.append(chevron(2.33, 0.57, +1, "arw-res"))

    half, top = 1.3, 1.15
    parts.append(block(-4.2, 0, half, top, "t", "l", "r"))
    parts.append(on_face_side(-4.2, 0.0, top, s["b1"], s["b1_sub"]))

    parts.append(block(0.15, 0, half, top, "rt", "rl", "rr"))
    parts.append(on_face_side(0.15, 0.0, top, s["b2"], s["b2_sub"], accent=True))

    parts.append(block(4.5, 0, half, top, "gt", "gl", "gr"))
    parts.append(on_face_side(4.5, 0.0, top, s["b3"], s["b3_sub"], accent=True))

    parts.append(text(-2.03, 2.15, s["c1"], "mono"))
    parts.append(text(2.33, 2.15, s["c2"], "mono"))
    return "".join(parts)


def map36(s):
    """Lesson: pages fill a frame, and the frame is what gets executed.

    The frame is the accent: it is the block that does the drawing, which is
    the thing a reader keeps getting backwards.
    """
    parts = []

    parts.append(road(-2.6, -1.45, 0.55, 0.5, "band-req"))
    parts.append(chevron(-2.03, 0.57, +1, "arw-req"))
    parts.append(road(1.75, 2.9, 0.55, 0.5, "band-res"))
    parts.append(chevron(2.33, 0.57, +1, "arw-res"))

    half, top = 1.3, 1.15
    parts.append(block(-4.2, 0, half, top, "t", "l", "r"))
    parts.append(on_face_side(-4.2, 0.0, top, s["b1"], s["b1_sub"]))

    parts.append(block(0.15, 0, half, top, "rt", "rl", "rr"))
    parts.append(on_face_side(0.15, 0.0, top, s["b2"], s["b2_sub"], accent=True))

    parts.append(block(4.5, 0, half, top, "gt", "gl", "gr"))
    parts.append(on_face_side(4.5, 0.0, top, s["b3"], s["b3_sub"], accent=True))

    parts.append(text(-2.03, 2.15, s["c1"], "mono"))
    parts.append(text(2.33, 2.15, s["c2"], "mono"))
    return "".join(parts)


def map37(s):
    """Lesson: what goes into the database in place of a password.

    bcrypt is the accent because the whole lesson is about paying for its
    slowness on purpose.
    """
    parts = []

    parts.append(road(-2.6, -1.45, 0.55, 0.5, "band-req"))
    parts.append(chevron(-2.03, 0.57, +1, "arw-req"))
    parts.append(road(1.75, 2.9, 0.55, 0.5, "band-res"))
    parts.append(chevron(2.33, 0.57, +1, "arw-res"))

    half, top = 1.3, 1.15
    parts.append(block(-4.2, 0, half, top, "t", "l", "r"))
    parts.append(on_face_side(-4.2, 0.0, top, s["b1"], s["b1_sub"]))

    parts.append(block(0.15, 0, half, top, "rt", "rl", "rr"))
    parts.append(on_face_side(0.15, 0.0, top, s["b2"], s["b2_sub"], accent=True))

    parts.append(block(4.5, 0, half, top, "gt", "gl", "gr"))
    parts.append(on_face_side(4.5, 0.0, top, s["b3"], s["b3_sub"], accent=True))

    parts.append(text(-2.03, 2.15, s["c1"], "mono"))
    parts.append(text(2.33, 2.15, s["c2"], "mono"))
    return "".join(parts)


def map38(s):
    """Lesson: the browser carries a ticket, the database keeps everything else.

    The token is the accent: it is the only thing that travels, and the whole
    lesson is about it meaning nothing on its own.
    """
    parts = []

    parts.append(road(-2.6, -1.45, 0.55, 0.5, "band-req"))
    parts.append(chevron(-2.03, 0.57, +1, "arw-req"))
    parts.append(road(1.75, 2.9, 0.55, 0.5, "band-res"))
    parts.append(chevron(2.33, 0.57, +1, "arw-res"))

    half, top = 1.3, 1.15
    parts.append(block(-4.2, 0, half, top, "t", "l", "r"))
    parts.append(on_face_side(-4.2, 0.0, top, s["b1"], s["b1_sub"]))

    parts.append(block(0.15, 0, half, top, "rt", "rl", "rr"))
    parts.append(on_face_side(0.15, 0.0, top, s["b2"], s["b2_sub"], accent=True))

    parts.append(block(4.5, 0, half, top, "gt", "gl", "gr"))
    parts.append(on_face_side(4.5, 0.0, top, s["b3"], s["b3_sub"], accent=True))

    parts.append(text(-2.03, 2.15, s["c1"], "mono"))
    parts.append(text(2.33, 2.15, s["c2"], "mono"))
    return "".join(parts)


def map39(s):
    """Lesson: three questions in order, and three different refusals.

    The last block is the accent: ownership is the question a beginner skips,
    and the one an attacker asks first.
    """
    parts = []

    parts.append(road(-2.6, -1.45, 0.55, 0.5, "band-req"))
    parts.append(chevron(-2.03, 0.57, +1, "arw-req"))
    parts.append(road(1.75, 2.9, 0.55, 0.5, "band-res"))
    parts.append(chevron(2.33, 0.57, +1, "arw-res"))

    half, top = 1.3, 1.15
    parts.append(block(-4.2, 0, half, top, "t", "l", "r"))
    parts.append(on_face_side(-4.2, 0.0, top, s["b1"], s["b1_sub"]))

    parts.append(block(0.15, 0, half, top, "rt", "rl", "rr"))
    parts.append(on_face_side(0.15, 0.0, top, s["b2"], s["b2_sub"], accent=True))

    parts.append(block(4.5, 0, half, top, "gt", "gl", "gr"))
    parts.append(on_face_side(4.5, 0.0, top, s["b3"], s["b3_sub"], accent=True))

    parts.append(text(-2.03, 2.15, s["c1"], "mono"))
    parts.append(text(2.33, 2.15, s["c2"], "mono"))
    return "".join(parts)


def map40(s):
    """Lesson: what comes from outside, what escapes it, and what proves the form is ours.

    Escaping is the accent: it is the one of the three that is already working
    and the one a reader can switch off by accident.
    """
    parts = []

    parts.append(road(-2.6, -1.45, 0.55, 0.5, "band-req"))
    parts.append(chevron(-2.03, 0.57, +1, "arw-req"))
    parts.append(road(1.75, 2.9, 0.55, 0.5, "band-res"))
    parts.append(chevron(2.33, 0.57, +1, "arw-res"))

    half, top = 1.3, 1.15
    parts.append(block(-4.2, 0, half, top, "t", "l", "r"))
    parts.append(on_face_side(-4.2, 0.0, top, s["b1"], s["b1_sub"]))

    parts.append(block(0.15, 0, half, top, "rt", "rl", "rr"))
    parts.append(on_face_side(0.15, 0.0, top, s["b2"], s["b2_sub"], accent=True))

    parts.append(block(4.5, 0, half, top, "gt", "gl", "gr"))
    parts.append(on_face_side(4.5, 0.0, top, s["b3"], s["b3_sub"], accent=True))

    parts.append(text(-2.03, 2.15, s["c1"], "mono"))
    parts.append(text(2.33, 2.15, s["c2"], "mono"))
    return "".join(parts)


def map41(s):
    """Lesson: what arrives, what decides it is an image, and what it is called.

    The check is the accent: it is the step a beginner replaces with a look at
    the file extension.
    """
    parts = []

    parts.append(road(-2.6, -1.45, 0.55, 0.5, "band-req"))
    parts.append(chevron(-2.03, 0.57, +1, "arw-req"))
    parts.append(road(1.75, 2.9, 0.55, 0.5, "band-res"))
    parts.append(chevron(2.33, 0.57, +1, "arw-res"))

    half, top = 1.3, 1.15
    parts.append(block(-4.2, 0, half, top, "t", "l", "r"))
    parts.append(on_face_side(-4.2, 0.0, top, s["b1"], s["b1_sub"]))

    parts.append(block(0.15, 0, half, top, "rt", "rl", "rr"))
    parts.append(on_face_side(0.15, 0.0, top, s["b2"], s["b2_sub"], accent=True))

    parts.append(block(4.5, 0, half, top, "gt", "gl", "gr"))
    parts.append(on_face_side(4.5, 0.0, top, s["b3"], s["b3_sub"], accent=True))

    parts.append(text(-2.03, 2.15, s["c1"], "mono"))
    parts.append(text(2.33, 2.15, s["c2"], "mono"))
    return "".join(parts)


def map42(s):
    """Lesson: cases in a list, a name on each, and a report that points.

    t.Run is the accent: the name is what turns a failure from a fact into an
    address.
    """
    parts = []

    parts.append(road(-2.6, -1.45, 0.55, 0.5, "band-req"))
    parts.append(chevron(-2.03, 0.57, +1, "arw-req"))
    parts.append(road(1.75, 2.9, 0.55, 0.5, "band-res"))
    parts.append(chevron(2.33, 0.57, +1, "arw-res"))

    half, top = 1.3, 1.15
    parts.append(block(-4.2, 0, half, top, "t", "l", "r"))
    parts.append(on_face_side(-4.2, 0.0, top, s["b1"], s["b1_sub"]))

    parts.append(block(0.15, 0, half, top, "rt", "rl", "rr"))
    parts.append(on_face_side(0.15, 0.0, top, s["b2"], s["b2_sub"], accent=True))

    parts.append(block(4.5, 0, half, top, "gt", "gl", "gr"))
    parts.append(on_face_side(4.5, 0.0, top, s["b3"], s["b3_sub"], accent=True))

    parts.append(text(-2.03, 2.15, s["c1"], "mono"))
    parts.append(text(2.33, 2.15, s["c2"], "mono"))
    return "".join(parts)


def map43(s):
    """Lesson: the server was already concurrent; the memory is what needs a lock.

    Shared memory is the accent: it is the one box a reader did not know they
    had built.
    """
    parts = []

    parts.append(road(-2.6, -1.45, 0.55, 0.5, "band-req"))
    parts.append(chevron(-2.03, 0.57, +1, "arw-req"))
    parts.append(road(1.75, 2.9, 0.55, 0.5, "band-res"))
    parts.append(chevron(2.33, 0.57, +1, "arw-res"))

    half, top = 1.3, 1.15
    parts.append(block(-4.2, 0, half, top, "t", "l", "r"))
    parts.append(on_face_side(-4.2, 0.0, top, s["b1"], s["b1_sub"]))

    parts.append(block(0.15, 0, half, top, "rt", "rl", "rr"))
    parts.append(on_face_side(0.15, 0.0, top, s["b2"], s["b2_sub"], accent=True))

    parts.append(block(4.5, 0, half, top, "gt", "gl", "gr"))
    parts.append(on_face_side(4.5, 0.0, top, s["b3"], s["b3_sub"], accent=True))

    parts.append(text(-2.03, 2.15, s["c1"], "mono"))
    parts.append(text(2.33, 2.15, s["c2"], "mono"))
    return "".join(parts)


def map44(s):
    """Lesson: a cancellation travelling down the chain.

    The context is the accent: it is the thing the reader did not create and
    cannot see, and it is what carries the news that nobody is waiting.
    """
    parts = []

    parts.append(road(-2.6, -1.45, 0.55, 0.5, "band-req"))
    parts.append(chevron(-2.03, 0.57, +1, "arw-req"))
    parts.append(road(1.75, 2.9, 0.55, 0.5, "band-res"))
    parts.append(chevron(2.33, 0.57, +1, "arw-res"))

    half, top = 1.3, 1.15
    parts.append(block(-4.2, 0, half, top, "t", "l", "r"))
    parts.append(on_face_side(-4.2, 0.0, top, s["b1"], s["b1_sub"]))

    parts.append(block(0.15, 0, half, top, "rt", "rl", "rr"))
    parts.append(on_face_side(0.15, 0.0, top, s["b2"], s["b2_sub"], accent=True))

    parts.append(block(4.5, 0, half, top, "gt", "gl", "gr"))
    parts.append(on_face_side(4.5, 0.0, top, s["b3"], s["b3_sub"], accent=True))

    parts.append(text(-2.03, 2.15, s["c1"], "mono"))
    parts.append(text(2.33, 2.15, s["c2"], "mono"))
    return "".join(parts)


def map45(s):
    """Lesson: what leaves the program, and how it is dressed on the way out.

    The tags are the accent: they are the only place where the shape of the
    answer is decided, and the one a reader forgets to think about.
    """
    parts = []

    parts.append(road(-2.6, -1.45, 0.55, 0.5, "band-req"))
    parts.append(chevron(-2.03, 0.57, +1, "arw-req"))
    parts.append(road(1.75, 2.9, 0.55, 0.5, "band-res"))
    parts.append(chevron(2.33, 0.57, +1, "arw-res"))

    half, top = 1.3, 1.15
    parts.append(block(-4.2, 0, half, top, "t", "l", "r"))
    parts.append(on_face_side(-4.2, 0.0, top, s["b1"], s["b1_sub"]))

    parts.append(block(0.15, 0, half, top, "rt", "rl", "rr"))
    parts.append(on_face_side(0.15, 0.0, top, s["b2"], s["b2_sub"], accent=True))

    parts.append(block(4.5, 0, half, top, "gt", "gl", "gr"))
    parts.append(on_face_side(4.5, 0.0, top, s["b3"], s["b3_sub"], accent=True))

    parts.append(text(-2.03, 2.15, s["c1"], "mono"))
    parts.append(text(2.33, 2.15, s["c2"], "mono"))
    return "".join(parts)

def map46(s):
    """Lesson: one setting, one file, one line -- and the blog is yours.

    The accent is on the setting: the name used to stand in eleven templates,
    and that is the block a reader has to see feeding all of them.
    """
    parts = []

    parts.append(road(-2.6, -1.45, 0.55, 0.5, "band-req"))
    parts.append(chevron(-2.03, 0.57, +1, "arw-req"))
    parts.append(road(1.75, 2.9, 0.55, 0.5, "band-res"))
    parts.append(chevron(2.33, 0.57, +1, "arw-res"))

    half, top = 1.3, 1.15
    parts.append(block(-4.2, 0, half, top, "rt", "rl", "rr"))
    parts.append(on_face_side(-4.2, 0.0, top, s["b1"], s["b1_sub"], accent=True))

    parts.append(block(0.15, 0, half, top, "t", "l", "r"))
    parts.append(on_face_side(0.15, 0.0, top, s["b2"], s["b2_sub"]))

    parts.append(block(4.5, 0, half, top, "gt", "gl", "gr"))
    parts.append(on_face_side(4.5, 0.0, top, s["b3"], s["b3_sub"], accent=True))

    parts.append(text(-2.03, 2.15, s["c1"], "mono"))
    parts.append(text(2.33, 2.15, s["c2"], "mono"))
    return "".join(parts)

def map47(s):
    """Lesson: what a host decides for us -- the port, the interface, the disk.

    The accent is on the interface: a service that listens to itself alone is
    the one failure every first deploy runs into.
    """
    parts = []

    parts.append(road(-2.6, -1.45, 0.55, 0.5, "band-req"))
    parts.append(chevron(-2.03, 0.57, +1, "arw-req"))
    parts.append(road(1.75, 2.9, 0.55, 0.5, "band-res"))
    parts.append(chevron(2.33, 0.57, +1, "arw-res"))

    half, top = 1.3, 1.15
    parts.append(block(-4.2, 0, half, top, "t", "l", "r"))
    parts.append(on_face_side(-4.2, 0.0, top, s["b1"], s["b1_sub"]))

    parts.append(block(0.15, 0, half, top, "rt", "rl", "rr"))
    parts.append(on_face_side(0.15, 0.0, top, s["b2"], s["b2_sub"], accent=True))

    parts.append(block(4.5, 0, half, top, "gt", "gl", "gr"))
    parts.append(on_face_side(4.5, 0.0, top, s["b3"], s["b3_sub"], accent=True))

    parts.append(text(-2.03, 2.15, s["c1"], "mono"))
    parts.append(text(2.33, 2.15, s["c2"], "mono"))
    return "".join(parts)

def map48(s):
    """Lesson: a name instead of digits, and a channel nobody can read.

    The accent is on TLS: the lesson exists because the same password is a
    plain word in one channel and nothing at all in the other.
    """
    parts = []

    parts.append(road(-2.6, -1.45, 0.55, 0.5, "band-req"))
    parts.append(chevron(-2.03, 0.57, +1, "arw-req"))
    parts.append(road(1.75, 2.9, 0.55, 0.5, "band-res"))
    parts.append(chevron(2.33, 0.57, +1, "arw-res"))

    half, top = 1.3, 1.15
    parts.append(block(-4.2, 0, half, top, "t", "l", "r"))
    parts.append(on_face_side(-4.2, 0.0, top, s["b1"], s["b1_sub"]))

    parts.append(block(0.15, 0, half, top, "gt", "gl", "gr"))
    parts.append(on_face_side(0.15, 0.0, top, s["b2"], s["b2_sub"], accent=True))

    parts.append(block(4.5, 0, half, top, "rt", "rl", "rr"))
    parts.append(on_face_side(4.5, 0.0, top, s["b3"], s["b3_sub"], accent=True))

    parts.append(text(-2.03, 2.15, s["c1"], "mono"))
    parts.append(text(2.33, 2.15, s["c2"], "mono"))
    return "".join(parts)

def map49(s):
    """Lesson: a signal, the answers still in flight, and who restarts us.

    The accent is on the middle block: the whole lesson is the difference
    between cutting a reader off and letting the answer finish.
    """
    parts = []

    parts.append(road(-2.6, -1.45, 0.55, 0.5, "band-req"))
    parts.append(chevron(-2.03, 0.57, +1, "arw-req"))
    parts.append(road(1.75, 2.9, 0.55, 0.5, "band-res"))
    parts.append(chevron(2.33, 0.57, +1, "arw-res"))

    half, top = 1.3, 1.15
    parts.append(block(-4.2, 0, half, top, "t", "l", "r"))
    parts.append(on_face_side(-4.2, 0.0, top, s["b1"], s["b1_sub"]))

    parts.append(block(0.15, 0, half, top, "rt", "rl", "rr"))
    parts.append(on_face_side(0.15, 0.0, top, s["b2"], s["b2_sub"], accent=True))

    parts.append(block(4.5, 0, half, top, "gt", "gl", "gr"))
    parts.append(on_face_side(4.5, 0.0, top, s["b3"], s["b3_sub"], accent=True))

    parts.append(text(-2.03, 2.15, s["c1"], "mono"))
    parts.append(text(2.33, 2.15, s["c2"], "mono"))
    return "".join(parts)

def map50(s):
    """Lesson: the three things a blog is not opened to strangers without.

    The accent is on the timeouts: a server without them is held open by one
    slow line, and that is the item people skip.
    """
    parts = []

    parts.append(road(-2.6, -1.45, 0.55, 0.5, "band-req"))
    parts.append(chevron(-2.03, 0.57, +1, "arw-req"))
    parts.append(road(1.75, 2.9, 0.55, 0.5, "band-res"))
    parts.append(chevron(2.33, 0.57, +1, "arw-res"))

    half, top = 1.3, 1.15
    parts.append(block(-4.2, 0, half, top, "rt", "rl", "rr"))
    parts.append(on_face_side(-4.2, 0.0, top, s["b1"], s["b1_sub"], accent=True))

    parts.append(block(0.15, 0, half, top, "t", "l", "r"))
    parts.append(on_face_side(0.15, 0.0, top, s["b2"], s["b2_sub"]))

    parts.append(block(4.5, 0, half, top, "gt", "gl", "gr"))
    parts.append(on_face_side(4.5, 0.0, top, s["b3"], s["b3_sub"], accent=True))

    parts.append(text(-2.03, 2.15, s["c1"], "mono"))
    parts.append(text(2.33, 2.15, s["c2"], "mono"))
    return "".join(parts)

def map51(s):
    """Last lesson: a type parameter, the condition on it, and the road on.

    The accent is on the condition: a type parameter without one is a promise
    nobody checks, and the compiler is the whole point of the feature.
    """
    parts = []

    parts.append(road(-2.6, -1.45, 0.55, 0.5, "band-req"))
    parts.append(chevron(-2.03, 0.57, +1, "arw-req"))
    parts.append(road(1.75, 2.9, 0.55, 0.5, "band-res"))
    parts.append(chevron(2.33, 0.57, +1, "arw-res"))

    half, top = 1.3, 1.15
    parts.append(block(-4.2, 0, half, top, "t", "l", "r"))
    parts.append(on_face_side(-4.2, 0.0, top, s["b1"], s["b1_sub"]))

    parts.append(block(0.15, 0, half, top, "rt", "rl", "rr"))
    parts.append(on_face_side(0.15, 0.0, top, s["b2"], s["b2_sub"], accent=True))

    parts.append(block(4.5, 0, half, top, "gt", "gl", "gr"))
    parts.append(on_face_side(4.5, 0.0, top, s["b3"], s["b3_sub"], accent=True))

    parts.append(text(-2.03, 2.15, s["c1"], "mono"))
    parts.append(text(2.33, 2.15, s["c2"], "mono"))
    return "".join(parts)

def pymap01(s):
    """Python lesson 1: what the numbers say before any of them is explained.

    Three blocks: what prices did, what prices elsewhere did, and the
    question the course exists to answer. The accent is the middle one --
    the comparison is what makes the first block impossible to wave away.
    """
    parts = []

    parts.append(road(-2.6, -1.45, 0.55, 0.5, "band-req"))
    parts.append(chevron(-2.03, 0.57, +1, "arw-req"))
    parts.append(road(1.75, 2.9, 0.55, 0.5, "band-res"))
    parts.append(chevron(2.33, 0.57, +1, "arw-res"))

    half, top = 1.3, 1.15
    parts.append(block(-4.2, 0, half, top, "t", "l", "r"))
    parts.append(on_face_side(-4.2, 0.0, top, s["b1"], s["b1_sub"]))

    parts.append(block(0.15, 0, half, top, "rt", "rl", "rr"))
    parts.append(on_face_side(0.15, 0.0, top, s["b2"], s["b2_sub"], accent=True))

    parts.append(block(4.5, 0, half, top, "gt", "gl", "gr"))
    parts.append(on_face_side(4.5, 0.0, top, s["b3"], s["b3_sub"], accent=True))

    parts.append(text(-2.03, 2.15, s["c1"], "mono"))
    parts.append(text(2.33, 2.15, s["c2"], "mono"))
    return "".join(parts)

def pymap02(s):
    """Python lesson 2: the system, the project's environment and the project.

    The accent is the middle block: the whole lesson is that a package installed
    there is invisible outside it, which is the difference between a workplace
    and a machine full of somebody else's versions.
    """
    parts = []

    parts.append(road(-2.6, -1.45, 0.55, 0.5, "band-req"))
    parts.append(chevron(-2.03, 0.57, +1, "arw-req"))
    parts.append(road(1.75, 2.9, 0.55, 0.5, "band-res"))
    parts.append(chevron(2.33, 0.57, +1, "arw-res"))

    half, top = 1.3, 1.15
    parts.append(block(-4.2, 0, half, top, "t", "l", "r"))
    parts.append(on_face_side(-4.2, 0.0, top, s["b1"], s["b1_sub"]))

    parts.append(block(0.15, 0, half, top, "rt", "rl", "rr"))
    parts.append(on_face_side(0.15, 0.0, top, s["b2"], s["b2_sub"], accent=True))

    parts.append(block(4.5, 0, half, top, "gt", "gl", "gr"))
    parts.append(on_face_side(4.5, 0.0, top, s["b3"], s["b3_sub"], accent=True))

    parts.append(text(-2.03, 2.15, s["c1"], "mono"))
    parts.append(text(2.33, 2.15, s["c2"], "mono"))
    return "".join(parts)

def pymap03(s):
    """Python lesson 3: numbers, strings and the way out of the program.

    The accent is the last block: a number printed without a format is a number
    nobody reads, and the lesson's own example is thirteen digits long.
    """
    parts = []

    parts.append(road(-2.6, -1.45, 0.55, 0.5, "band-req"))
    parts.append(chevron(-2.03, 0.57, +1, "arw-req"))
    parts.append(road(1.75, 2.9, 0.55, 0.5, "band-res"))
    parts.append(chevron(2.33, 0.57, +1, "arw-res"))

    half, top = 1.3, 1.15
    parts.append(block(-4.2, 0, half, top, "t", "l", "r"))
    parts.append(on_face_side(-4.2, 0.0, top, s["b1"], s["b1_sub"]))

    parts.append(block(0.15, 0, half, top, "t", "l", "r"))
    parts.append(on_face_side(0.15, 0.0, top, s["b2"], s["b2_sub"]))

    parts.append(block(4.5, 0, half, top, "rt", "rl", "rr"))
    parts.append(on_face_side(4.5, 0.0, top, s["b3"], s["b3_sub"], accent=True))

    parts.append(text(-2.03, 2.15, s["c1"], "mono"))
    parts.append(text(2.33, 2.15, s["c2"], "mono"))
    return "".join(parts)

def pymap04(s):
    """Python lesson 4: text, number and truth.

    The accent is the middle block: money in float is the mistake this lesson
    exists to prevent, and 0.1 + 0.2 is the shortest proof of it.
    """
    parts = []

    parts.append(road(-2.6, -1.45, 0.55, 0.5, "band-req"))
    parts.append(chevron(-2.03, 0.57, +1, "arw-req"))
    parts.append(road(1.75, 2.9, 0.55, 0.5, "band-res"))
    parts.append(chevron(2.33, 0.57, +1, "arw-res"))

    half, top = 1.3, 1.15
    parts.append(block(-4.2, 0, half, top, "t", "l", "r"))
    parts.append(on_face_side(-4.2, 0.0, top, s["b1"], s["b1_sub"]))

    parts.append(block(0.15, 0, half, top, "rt", "rl", "rr"))
    parts.append(on_face_side(0.15, 0.0, top, s["b2"], s["b2_sub"], accent=True))

    parts.append(block(4.5, 0, half, top, "gt", "gl", "gr"))
    parts.append(on_face_side(4.5, 0.0, top, s["b3"], s["b3_sub"], accent=True))

    parts.append(text(-2.03, 2.15, s["c1"], "mono"))
    parts.append(text(2.33, 2.15, s["c2"], "mono"))
    return "".join(parts)

def pymap05(s):
    """Python lesson 5: a list, a slice and a tuple.

    The accent is the middle block: the slice is where the off-by-one lives,
    and the lesson is built so the reader meets it on purpose.
    """
    parts = []

    parts.append(road(-2.6, -1.45, 0.55, 0.5, "band-req"))
    parts.append(chevron(-2.03, 0.57, +1, "arw-req"))
    parts.append(road(1.75, 2.9, 0.55, 0.5, "band-res"))
    parts.append(chevron(2.33, 0.57, +1, "arw-res"))

    half, top = 1.3, 1.15
    parts.append(block(-4.2, 0, half, top, "t", "l", "r"))
    parts.append(on_face_side(-4.2, 0.0, top, s["b1"], s["b1_sub"]))

    parts.append(block(0.15, 0, half, top, "rt", "rl", "rr"))
    parts.append(on_face_side(0.15, 0.0, top, s["b2"], s["b2_sub"], accent=True))

    parts.append(block(4.5, 0, half, top, "gt", "gl", "gr"))
    parts.append(on_face_side(4.5, 0.0, top, s["b3"], s["b3_sub"], accent=True))

    parts.append(text(-2.03, 2.15, s["c1"], "mono"))
    parts.append(text(2.33, 2.15, s["c2"], "mono"))
    return "".join(parts)

def pymap06(s):
    """Python lesson 6: a key, a value and a counter.

    The accent is the middle block: a missing key is an error rather than a
    zero, and that distinction is what keeps a calculation honest.
    """
    parts = []

    parts.append(road(-2.6, -1.45, 0.55, 0.5, "band-req"))
    parts.append(chevron(-2.03, 0.57, +1, "arw-req"))
    parts.append(road(1.75, 2.9, 0.55, 0.5, "band-res"))
    parts.append(chevron(2.33, 0.57, +1, "arw-res"))

    half, top = 1.3, 1.15
    parts.append(block(-4.2, 0, half, top, "t", "l", "r"))
    parts.append(on_face_side(-4.2, 0.0, top, s["b1"], s["b1_sub"]))

    parts.append(block(0.15, 0, half, top, "rt", "rl", "rr"))
    parts.append(on_face_side(0.15, 0.0, top, s["b2"], s["b2_sub"], accent=True))

    parts.append(block(4.5, 0, half, top, "gt", "gl", "gr"))
    parts.append(on_face_side(4.5, 0.0, top, s["b3"], s["b3_sub"], accent=True))

    parts.append(text(-2.03, 2.15, s["c1"], "mono"))
    parts.append(text(2.33, 2.15, s["c2"], "mono"))
    return "".join(parts)

def pymap07(s):
    """Python lesson 7: the fork, the chain and the gap in the data.

    The accent is the last block: a missing value is not a zero, and a course
    about data has to say so at the first condition the reader writes.
    """
    parts = []

    parts.append(road(-2.6, -1.45, 0.55, 0.5, "band-req"))
    parts.append(chevron(-2.03, 0.57, +1, "arw-req"))
    parts.append(road(1.75, 2.9, 0.55, 0.5, "band-res"))
    parts.append(chevron(2.33, 0.57, +1, "arw-res"))

    half, top = 1.3, 1.15
    parts.append(block(-4.2, 0, half, top, "t", "l", "r"))
    parts.append(on_face_side(-4.2, 0.0, top, s["b1"], s["b1_sub"]))

    parts.append(block(0.15, 0, half, top, "rt", "rl", "rr"))
    parts.append(on_face_side(0.15, 0.0, top, s["b2"], s["b2_sub"], accent=True))

    parts.append(block(4.5, 0, half, top, "gt", "gl", "gr"))
    parts.append(on_face_side(4.5, 0.0, top, s["b3"], s["b3_sub"], accent=True))

    parts.append(text(-2.03, 2.15, s["c1"], "mono"))
    parts.append(text(2.33, 2.15, s["c2"], "mono"))
    return "".join(parts)

def pymap08(s):
    """Python lesson 8: the walk, the skip and the edge.

    The accent is the last block: a loop that says nothing when it runs out of
    data is the trap of this lesson -- zip cuts to the shorter series without a
    word, and a while runs past the end of one.
    """
    parts = []

    parts.append(road(-2.6, -1.45, 0.55, 0.5, "band-req"))
    parts.append(chevron(-2.03, 0.57, +1, "arw-req"))
    parts.append(road(1.75, 2.9, 0.55, 0.5, "band-res"))
    parts.append(chevron(2.33, 0.57, +1, "arw-res"))

    half, top = 1.3, 1.15
    parts.append(block(-4.2, 0, half, top, "t", "l", "r"))
    parts.append(on_face_side(-4.2, 0.0, top, s["b1"], s["b1_sub"]))

    parts.append(block(0.15, 0, half, top, "rt", "rl", "rr"))
    parts.append(on_face_side(0.15, 0.0, top, s["b2"], s["b2_sub"], accent=True))

    parts.append(block(4.5, 0, half, top, "gt", "gl", "gr"))
    parts.append(on_face_side(4.5, 0.0, top, s["b3"], s["b3_sub"], accent=True))

    parts.append(text(-2.03, 2.15, s["c1"], "mono"))
    parts.append(text(2.33, 2.15, s["c2"], "mono"))
    return "".join(parts)

def check_labels(svg_body, name, lang):
    """Warn when a label is wider than the face it is written on.

    A label that overhangs its block stops naming it and turns into text lying
    over the picture. Eyeballing sixty-six files does not scale, so the width is
    estimated here and the map that broke it is named.

    The face is a rhombus, so its widest chord is the one through the middle and
    every line above or below it has less room -- which is exactly where a
    second line of label goes. That is why the check narrows the allowance by
    how far the text sits from the centre.
    """
    faces = []
    for m in re.finditer(r'<polygon class="(gt|rt|t)" points="([^"]+)"/>', svg_body):
        pts = [tuple(map(float, pt.split(","))) for pt in m.group(2).split()]
        xs, ys = [p[0] for p in pts], [p[1] for p in pts]
        cx, cy = sum(xs) / len(xs), sum(ys) / len(ys)
        faces.append((cx, cy, max(xs) - min(xs), (max(ys) - min(ys)) / 2))
    if not faces:
        return
    for m in re.finditer(r'<text class="([a-z-]+)" x="([-\d.]+)" y="([-\d.]+)">([^<]*)</text>', svg_body):
        cls, x, y = m.group(1), float(m.group(2)), float(m.group(3))
        body = (m.group(4).replace("&lt;", "<").replace("&gt;", ">")
                .replace("&amp;", "&"))
        if cls not in ("lbl", "lbl-acc", "sub", "sub-acc", "tag") or not body:
            continue
        cx, cy, fw, fh = min(faces, key=lambda f: (f[0] - x) ** 2 + (f[1] - y) ** 2)
        if abs(cx - x) > 40 or abs(cy - y) > fh + 20:
            continue  # not written on a block at all
        room = fw * max(0.25, 1.0 - abs(y - cy) / max(fh, 1.0))
        want = len(body) * _FONT_PX[cls] * _CHAR_W[cls]
        if want > room * 0.98:
            print(f"  ШИРЕ ГРАНИ: {name}-{lang} «{body}» {want:.0f}px, место {room:.0f}px")


def render(scene, strings, name="", lang=""):
    """Draw the shared frame, then the scene centred inside it, 16:9."""
    _seen.clear()
    grid_and_caps = frame(strings)

    # Only the scene's own points decide where it sits, so the grid and the
    # captions are drawn first and their points forgotten.
    _seen.clear()
    body = scene(strings)
    xs = [x for x, _ in _seen]
    ys = [y for _, y in _seen]

    # Scenes were each composed inside their own tight frame, so their natural
    # centres are all over the place. Rather than nudging twenty-two of them by
    # hand, the whole scene is shifted so its middle lands on the middle of the
    # canvas: the captions stay put and the picture stops sliding about.
    cx, cy = (min(xs) + max(xs)) / 2, (min(ys) + max(ys)) / 2
    half_w = max((max(xs) - min(xs)) / 2, 1.0)
    half_h = max((max(ys) - min(ys)) / 2, 1.0)

    # And grown to use the room it has. A scene composed for a tight frame
    # looks marooned in the middle of a fixed one, so it is scaled up until it
    # nearly touches the safe area. The cap keeps a sparse map from blowing its
    # three blocks up to twice the size of everyone else's.
    k = min(SAFE_X * 0.94 / half_w, SAFE_Y_DOWN * 0.94 / half_h, ZOOM_MAX)
    body = (f'<g transform="translate(0,{ORIGIN_Y:.2f}) scale({k:.3f}) '
            f'translate({-cx:.2f},{-cy:.2f})">{body}</g>')

    if half_w * k > SAFE_X or half_h * k > SAFE_Y_DOWN:
        print(f"  НЕ ВЛЕЗАЕТ: {name}-{lang} {half_w * 2 * k:.0f}x{half_h * 2 * k:.0f}px")
    check_labels(body, name, lang)

    x0, y0 = -CANVAS_W / 2, ORIGIN_Y - CANVAS_H / 2
    return (
        f'<svg xmlns="http://www.w3.org/2000/svg" '
        f'viewBox="{x0:.0f} {y0:.0f} {CANVAS_W:.0f} {CANVAS_H:.0f}" '
        f'width="{CANVAS_W:.0f}" height="{CANVAS_H:.0f}" '
        f'role="img" aria-label="{esc(strings["alt"])}">'
        f"<style>{STYLE}</style>"
        f"{grid_and_caps}{body}</svg>"
    )


L = {
    "kz": dict(
        alt="Клиент сұраныс жібереді, сервер жауап қайтарады",
        client="Клиент", client_sub="браузер, curl",
        server="Сервер", server_sub="сұрақ күтіп тұр",
        head="СҰРАНЫС", head_sub="GET /read/salem",
        foot="ЖАУАП", foot_sub="200 + HTML",
        req1="ӘДІС", req2="ЖОЛ", req3="ТАҚЫРЫПТАР",
        res1="КОД", res2="ТАҚЫРЫПТАР", res3="ДЕНЕ",
    ),
    "ru": dict(
        alt="Клиент отправляет запрос, сервер возвращает ответ",
        client="Клиент", client_sub="браузер, curl",
        server="Сервер", server_sub="ждёт вопроса",
        head="ЗАПРОС", head_sub="GET /read/salem",
        foot="ОТВЕТ", foot_sub="200 + HTML",
        req1="МЕТОД", req2="ПУТЬ", req3="ЗАГОЛОВКИ",
        res1="КОД", res2="ЗАГОЛОВКИ", res3="ТЕЛО",
    ),
    "en": dict(
        alt="A client sends a request, a server returns a response",
        client="Client", client_sub="browser, curl",
        server="Server", server_sub="waiting to be asked",
        head="REQUEST", head_sub="GET /read/salem",
        foot="RESPONSE", foot_sub="200 + HTML",
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
        head="БІР ПӘРМЕН", head_sub="go run .",
        foot="ЖАУАП", foot_sub="200 · Сәлем",
    ),
    "ru": dict(
        alt="Программа знает, что отвечать и где слушать, а браузер видит ответ",
        b1="ЧТО ОТВЕЧАТЬ", b1_sub="http.HandleFunc",
        b2="ГДЕ СЛУШАТЬ", b2_sub="http.ListenAndServe",
        b3="БРАУЗЕР", b3_sub="localhost:8080",
        head="ОДНА КОМАНДА", head_sub="go run .",
        foot="ОТВЕТ", foot_sub="200 · Сәлем",
    ),
    "en": dict(
        alt="The program knows what to answer and where to listen; the browser sees the answer",
        b1="WHAT TO ANSWER", b1_sub="http.HandleFunc",
        b2="WHERE TO LISTEN", b2_sub="http.ListenAndServe",
        b3="BROWSER", b3_sub="localhost:8080",
        head="ONE COMMAND", head_sub="go run .",
        foot="RESPONSE", foot_sub="200 · Salem",
    ),
}

L02 = {
    "kz": dict(
        alt="Редактор, терминал және Go: қайсысы не үшін жауап береді",
        b1="РЕДАКТОР", b1_sub="VS Code",
        b2="ТЕРМИНАЛ", b2_sub="go run .",
        b3="GO", b3_sub="жинап, қосады",
        head="ҚАЙСЫСЫ НЕ ҮШІН ЖАУАП БЕРЕДІ", head_sub="үш құрал, үш жауапкершілік",
        foot="ЖОБА ҚАЛТАСЫ", foot_sub="go-oqu/sabaq-01 · go.mod + main.go",
    ),
    "ru": dict(
        alt="Редактор, терминал и Go: кто за что отвечает",
        b1="РЕДАКТОР", b1_sub="VS Code",
        b2="ТЕРМИНАЛ", b2_sub="go run .",
        b3="GO", b3_sub="собирает и запускает",
        head="КТО ЗА ЧТО ОТВЕЧАЕТ", head_sub="три инструмента, три зоны вины",
        foot="ПАПКА ПРОЕКТА", foot_sub="go-oqu/sabaq-01 · go.mod + main.go",
    ),
    "en": dict(
        alt="Editor, terminal and Go: which one is answerable for what",
        b1="EDITOR", b1_sub="VS Code",
        b2="TERMINAL", b2_sub="go run .",
        b3="GO", b3_sub="builds it and runs it",
        head="WHICH ONE IS ANSWERABLE", head_sub="three tools, three kinds of fault",
        foot="PROJECT FOLDER", foot_sub="go-oqu/sabaq-01 · go.mod + main.go",
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
    "kz": dict(alt="Курс модульдері тұр, жол одан әрі оқырманның блогына апарады",
               m1="ТІЛ", m2="ВЕБ", m3="ДЕРЕКТЕР", m4="АДАМДАР", m5="ІСКЕ ҚОСУ",
               head="ТЕГІН КУРС: GO НӨЛДЕН", head_sub="50 сабақ · тіркелусіз · үш тілде",
               foot="КУРС ТОЛЫҚ ЖАЗЫЛҒАН", foot_sub="соңында — интернеттегі өз блогың"),
    "ru": dict(alt="Модули курса стоят, дорога ведёт дальше — к вашему блогу",
               m1="ЯЗЫК", m2="ВЕБ", m3="ДАННЫЕ", m4="ЛЮДИ", m5="ЗАПУСК",
               head="БЕСПЛАТНЫЙ КУРС: GO С НУЛЯ", head_sub="50 уроков · без регистрации · три языка",
               foot="КУРС НАПИСАН ЦЕЛИКОМ", foot_sub="в конце — ваш блог в интернете"),
    "en": dict(alt="The course's modules stand, and the road leads on to your blog",
               m1="LANGUAGE", m2="WEB", m3="DATA", m4="PEOPLE", m5="LAUNCH",
               head="A FREE COURSE: GO FROM SCRATCH", head_sub="50 lessons · no account · three languages",
               foot="THE COURSE IS WRITTEN IN FULL", foot_sub="at the end, your own blog online"),
}

L06 = {
    "kz": dict(
        alt="Бір for сөзі, үш түрлі пішін: санауыш, шарт және шексіз",
        f1="САНАУЫШ", f1_sub="for i := 1; i <= 3; i++",
        f2="ШАРТ", f2_sub="for words < 600",
        f3="ШЕКСІЗ", f3_sub="for { … break }",
        head="БІР СӨЗ, ҮШ ПІШІН", head_sub="Go-дағы жалғыз цикл — for",
        foot="GO-ДА БҰЛАР ЖОҚ", foot_sub="while · do-while · foreach",
    ),
    "ru": dict(
        alt="Одно слово for, три формы: счётчик, условие и бесконечный",
        f1="СЧЁТЧИК", f1_sub="for i := 1; i <= 3; i++",
        f2="УСЛОВИЕ", f2_sub="for words < 600",
        f3="БЕЗ УСЛОВИЯ", f3_sub="for { … break }",
        head="ОДНО СЛОВО, ТРИ ФОРМЫ", head_sub="единственный цикл в Go — for",
        foot="В GO ИХ НЕТ", foot_sub="while · do-while · foreach",
    ),
    "en": dict(
        alt="One word, for, in three shapes: counter, condition and endless",
        f1="COUNTER", f1_sub="for i := 1; i <= 3; i++",
        f2="CONDITION", f2_sub="for words < 600",
        f3="ENDLESS", f3_sub="for { … break }",
        head="ONE WORD, THREE SHAPES", head_sub="the only loop in Go is for",
        foot="GO HAS NONE OF THESE", foot_sub="while · do-while · foreach",
    ),
}

L07 = {
    "kz": dict(alt="Жеті әріп, он бір байт: латын әрпі бір, қазақ әрпі екі байт",
               head="«Go тілі» — ЖЕТІ ӘРІП", head_sub="len(\"Go тілі\") = 11",
               foot="БАЙТ ≠ ӘРІП", foot_sub="astyndagy sandar — байт ығысуы"),
    "ru": dict(alt="Семь букв, одиннадцать байт: латинская буква один, казахская два",
               head="«Go тілі» — СЕМЬ БУКВ", head_sub="len(\"Go тілі\") = 11",
               foot="БАЙТ ≠ БУКВА", foot_sub="числа под буквами — смещение в байтах"),
    "en": dict(alt="Seven letters, eleven bytes: a Latin letter costs one, a Kazakh letter two",
               head="\"Go тілі\" — SEVEN LETTERS", head_sub="len(\"Go тілі\") = 11",
               foot="A BYTE IS NOT A LETTER", foot_sub="the numbers below are byte offsets"),
}
L07["kz"]["foot_sub"] = "әріп астындағы сандар — байт ығысуы"

L08 = {
    "kz": dict(alt="Ұзындық үш, сыйымдылық төрт: үш ұяшық толы, біреуі бос",
               c1="Шаңырақ", c2="Go тілі", c3="Дала", c4="бос",
               head="len = 3 — НЕШЕУІ БАР", head_sub="titles[0] … titles[2]",
               foot="cap = 4 — КӨШПЕЙ НЕШЕУІ СЫЯДЫ", foot_sub="орын жетпесе, append көшіріп, жаңа мекенжай береді"),
    "ru": dict(alt="Длина три, ёмкость четыре: три ячейки заняты, одна свободна",
               c1="Шаңырақ", c2="Go тілі", c3="Дала", c4="свободно",
               head="len = 3 — СКОЛЬКО ЛЕЖИТ", head_sub="titles[0] … titles[2]",
               foot="cap = 4 — СКОЛЬКО ВЛЕЗЕТ БЕЗ ПЕРЕЕЗДА", foot_sub="места не хватило — append переселит и вернёт новый адрес"),
    "en": dict(alt="Length three, capacity four: three cells taken, one free",
               c1="Шаңырақ", c2="Go тілі", c3="Дала", c4="free",
               head="len = 3 — WHAT IS IN IT", head_sub="titles[0] … titles[2]",
               foot="cap = 4 — WHAT FITS BEFORE IT MOVES", foot_sub="out of room, append moves it and hands back a new address"),
}

L09 = {
    "kz": dict(alt="Үш кілт өз мәнін ашады, төртінші кілтке ештеңе салынбаған: нөл",
               k1='"go"', k2='"веб"', k3='"дала"', k4='"музыка"',
               v1="121", v2="30", v3="7", v4="0",
               head="КІЛТ МӘНДІ БІРДЕН БЕРЕДІ", head_sub='views := map[string]int{…}',
               foot="ЖОҚ КІЛТ — ҚАТЕ ЕМЕС, НӨЛ",
               foot_sub='n, ok := views["музыка"]   // 0, false'),
    "ru": dict(alt="Три ключа открывают своё значение, в четвёртый никто ничего не клал: ноль",
               k1='"go"', k2='"веб"', k3='"дала"', k4='"музыка"',
               v1="121", v2="30", v3="7", v4="0",
               head="КЛЮЧ ОТДАЁТ ЗНАЧЕНИЕ СРАЗУ", head_sub='views := map[string]int{…}',
               foot="КЛЮЧА НЕТ — ЭТО НЕ ОШИБКА, А НОЛЬ",
               foot_sub='n, ok := views["музыка"]   // 0, false'),
    "en": dict(alt="Three keys open their value, the fourth had nothing put in it: zero",
               k1='"go"', k2='"веб"', k3='"дала"', k4='"музыка"',
               v1="121", v2="30", v3="7", v4="0",
               head="A KEY HANDS BACK ITS VALUE AT ONCE", head_sub='views := map[string]int{…}',
               foot="A MISSING KEY IS NOT AN ERROR, IT IS ZERO",
               foot_sub='n, ok := views["музыка"]   // 0, false'),
}

L10 = {
    "kz": dict(alt="Тест кодтың қасында тұрып, оны шақырады және екі жауаптың бірін береді",
               code="main.go", code_sub="бағдарлама", test="main_test.go", test_sub="тексеру",
               good="PASS", bad="FAIL",
               head="go test — ЖАЗЫЛЫП ҚОЙҒАН ТЕКСЕРУ", head_sub="func TestCount(t *testing.T)",
               foot="_test.go ДАЙЫН БАҒДАРЛАМАҒА КІРМЕЙДІ", foot_sub="ok  sabaq09   ·   main_test.go:9: алдық 2, күттік 3"),
    "ru": dict(alt="Тест стоит рядом с кодом, вызывает его и даёт один из двух ответов",
               code="main.go", code_sub="программа", test="main_test.go", test_sub="проверка",
               good="PASS", bad="FAIL",
               head="go test — ЗАПИСАННАЯ ПРОВЕРКА", head_sub="func TestCount(t *testing.T)",
               foot="_test.go В ГОТОВУЮ ПРОГРАММУ НЕ ПОПАДАЕТ", foot_sub="ok  sabaq09   ·   main_test.go:9: получили 2, ждали 3"),
    "en": dict(alt="The test stands beside the code, calls it and gives one of two answers",
               code="main.go", code_sub="the program", test="main_test.go", test_sub="the check",
               good="PASS", bad="FAIL",
               head="go test — A CHECK THAT STAYS WRITTEN DOWN", head_sub="func TestCount(t *testing.T)",
               foot="_test.go NEVER SHIPS WITH THE PROGRAM", foot_sub="ok  sabaq09   ·   main_test.go:9: got 2, wanted 3"),
}

L11 = {
    "kz": dict(alt="Керек үш нәрсе — компьютер, орнату құқығы, күніне жарты сағат — және керек емес біреуі",
               n1="КОМПЬЮТЕР", n2="ОРНАТУ ҚҰҚЫҒЫ", n3="ЖАРТЫ САҒАТ", n4="АҒЫЛШЫН",
               s1="телефон жетпейді", s2="бағдарлама орнату", s3="күніне", s4="керек емес",
               head="БАСТАУ ҮШІН НЕ КЕРЕК", head_sub="go.dev/dl · code.visualstudio.com",
               foot="ОРНАТУҒА ТЫЙЫМ САЛЫНСА — go.dev/play",
               foot_sub="сервер сабағынан бастап өз компьютері керек"),
    "ru": dict(alt="Три нужные вещи — компьютер, право ставить программы, полчаса в день — и одна ненужная",
               n1="КОМПЬЮТЕР", n2="ПРАВО СТАВИТЬ", n3="ПОЛЧАСА", n4="АНГЛИЙСКИЙ",
               s1="телефона не хватит", s2="программы", s3="в день", s4="не нужен",
               head="ЧТО НУЖНО, ЧТОБЫ НАЧАТЬ", head_sub="go.dev/dl · code.visualstudio.com",
               foot="ЕСЛИ СТАВИТЬ ЗАПРЕЩЕНО — go.dev/play",
               foot_sub="с урока про сервер понадобится свой компьютер"),
    "en": dict(alt="Three things you need — a computer, the right to install, half an hour a day — and one you do not",
               n1="A COMPUTER", n2="CAN INSTALL", n3="HALF AN HOUR", n4="ENGLISH",
               s1="a phone will not do", s2="software", s3="a day", s4="not required",
               head="WHAT YOU NEED TO START", head_sub="go.dev/dl · code.visualstudio.com",
               foot="IF INSTALLING IS BLOCKED — go.dev/play",
               foot_sub="from the server lesson on you need your own machine"),
}

L12 = {
    "kz": dict(alt="Үш бөлек сөздік бір типке жиналады: Article",
               p1="titles", p2="words", p3="lang",
               name="Article", fields="Title · Words · Lang",
               head="ҮШЕУДІҢ ОРНЫНА БІР МӘН", head_sub="type Article struct { … }",
               foot="ҚҰРЫЛЫМ КӨШІРМЕМЕН БЕРІЛЕДІ",
               foot_sub="for _, a := range blog { a.Title = … }  // ештеңе өзгермейді"),
    "ru": dict(alt="Три отдельных словаря собираются в один тип: Article",
               p1="titles", p2="words", p3="lang",
               name="Article", fields="Title · Words · Lang",
               head="ОДНО ЗНАЧЕНИЕ ВМЕСТО ТРЁХ", head_sub="type Article struct { … }",
               foot="СТРУКТУРА ПЕРЕДАЁТСЯ КОПИЕЙ",
               foot_sub="for _, a := range blog { a.Title = … }  // ничего не изменит"),
    "en": dict(alt="Three separate maps collapse into one type: Article",
               p1="titles", p2="words", p3="lang",
               name="Article", fields="Title · Words · Lang",
               head="ONE VALUE INSTEAD OF THREE", head_sub="type Article struct { … }",
               foot="A STRUCT IS PASSED AS A COPY",
               foot_sub="for _, a := range blog { a.Title = … }  // changes nothing"),
}

L13 = {
    "kz": dict(alt="Көшірме тұйыққа кетеді, мекенжай мәнге қайта әкеледі",
               value="Article", value_sub="жадыдағы мән",
               copy="көшірме", copy_sub="жазғаныңыз сонда қалады",
               addr="&a", addr_sub="жазғаныңыз жетеді",
               head="ӘДІС ЕКЕУДІҢ БІРІН АЛАДЫ", head_sub="func (a Article)   ·   func (a *Article)",
               foot="ӨЗГЕРТСЕ — СІЛТЕГІШ, ТЕК ОҚИТЫН БОЛСА — МӘН",
               foot_sub="for i := range blog { blog[i].Publish() }"),
    "ru": dict(alt="Копия уходит в тупик, адрес возвращает к самому значению",
               value="Article", value_sub="значение в памяти",
               copy="копия", copy_sub="правка остаётся там",
               addr="&a", addr_sub="правка доходит",
               head="МЕТОД БЕРЁТ ОДНО ИЗ ДВУХ", head_sub="func (a Article)   ·   func (a *Article)",
               foot="МЕНЯЕТ — УКАЗАТЕЛЬ, ТОЛЬКО ЧИТАЕТ — ЗНАЧЕНИЕ",
               foot_sub="for i := range blog { blog[i].Publish() }"),
    "en": dict(alt="A copy is a dead end, an address leads back to the value itself",
               value="Article", value_sub="the value in memory",
               copy="a copy", copy_sub="the edit stays there",
               addr="&a", addr_sub="the edit arrives",
               head="A METHOD TAKES ONE OF THE TWO", head_sub="func (a Article)   ·   func (a *Article)",
               foot="CHANGES IT — POINTER; ONLY READS IT — VALUE",
               foot_sub="for i := range blog { blog[i].Publish() }"),
}

L14 = {
    "kz": dict(alt="Бір функция, кез келген қойма: келісім — біліктер тізімі",
               fn="report", fn_sub="қойманы білмейді",
               iface="Store", iface_sub="Add · All",
               s1="MemoryStore", s2="LastTwoStore", s3="fakeStore",
               head="ҚҰРЫЛЫСЫ ЕМЕС, ӘДІСТЕРІ", head_sub="type Store interface { Add · All }",
               foot="ІСКЕ АСЫРУ ЖАРИЯЛАНБАЙДЫ — ӨЗІНЕН-ӨЗІ ШЫҒАДЫ",
               foot_sub="func (s *MemoryStore) Add(a Article)"),
    "ru": dict(alt="Одна функция, любое хранилище: договор — это список умений",
               fn="report", fn_sub="не знает какое",
               iface="Store", iface_sub="Add · All",
               s1="MemoryStore", s2="LastTwoStore", s3="fakeStore",
               head="УМЕНИЯ, А НЕ УСТРОЙСТВО", head_sub="type Store interface { Add · All }",
               foot="РЕАЛИЗАЦИЯ НЕ ОБЪЯВЛЯЕТСЯ — ОНА ПОЛУЧАЕТСЯ САМА",
               foot_sub="func (s *MemoryStore) Add(a Article)"),
    "en": dict(alt="One function, any store: the contract is a list of skills",
               fn="report", fn_sub="knows no store",
               iface="Store", iface_sub="Add · All",
               s1="MemoryStore", s2="LastTwoStore", s3="fakeStore",
               head="WHAT IT CAN DO, NOT WHAT IT IS", head_sub="type Store interface { Add · All }",
               foot="IMPLEMENTING IS NOT DECLARED — IT SIMPLY HAPPENS",
               foot_sub="func (s *MemoryStore) Add(a Article)"),
}

L15 = {
    "kz": dict(alt="Функцияның екі шығуы бар: мән және қате",
               fn="Get(slug)", fn_sub="екі мән",
               good="Article", bad="error",
               good_sub="қате болмаса — мән", bad_sub="қате болса — nil емес",
               head="ҚАТЕ — ЕРЕКШЕ ЖАҒДАЙ ЕМЕС, ЕКІНШІ МӘН",
               head_sub="func (s *Store) Get(slug string) (Article, error)",
               foot="ЕКІНШІСІ БІРІНШІ ТЕКСЕРІЛЕДІ",
               foot_sub="if err != nil { return … }   ·   errors.Is(err, ErrNotFound)"),
    "ru": dict(alt="У функции два выхода: значение и ошибка",
               fn="Get(slug)", fn_sub="два значения",
               good="Article", bad="error",
               good_sub="нет ошибки — есть значение", bad_sub="есть ошибка — не nil",
               head="ОШИБКА — НЕ ИСКЛЮЧЕНИЕ, А ВТОРОЕ ЗНАЧЕНИЕ",
               head_sub="func (s *Store) Get(slug string) (Article, error)",
               foot="ВТОРОЕ ПРОВЕРЯЮТ ПЕРВЫМ",
               foot_sub="if err != nil { return … }   ·   errors.Is(err, ErrNotFound)"),
    "en": dict(alt="A function has two ways out: a value and an error",
               fn="Get(slug)", fn_sub="two values",
               good="Article", bad="error",
               good_sub="no error — a value", bad_sub="an error — not nil",
               head="AN ERROR IS NOT AN EXCEPTION, IT IS THE SECOND VALUE",
               head_sub="func (s *Store) Get(slug string) (Article, error)",
               foot="THE SECOND ONE IS CHECKED FIRST",
               foot_sub="if err != nil { return … }   ·   errors.Is(err, ErrNotFound)"),
}

L16 = {
    "kz": dict(alt="Пакет — қалта: бас әріппен жазылғаны сыртқа көрінеді, қалғаны ішінде қалады",
               main="main", main_sub="рұқсат етілгенді",
               pkg="blog", pkg_sub="blog/ қалтасы",
               out="Article · Get", inn="items · count",
               head="ПАКЕТ — БҰЛ ҚАЛТА", head_sub='import "sabaq14/blog"',
               foot="БАС ӘРІП — СЫРТҚА АШЫЛАТЫН ЕСІК",
               foot_sub="s.Get(…) көрінеді   ·   s.count() көрінбейді"),
    "ru": dict(alt="Пакет — это папка: с заглавной буквы видно снаружи, остальное остаётся внутри",
               main="main", main_sub="видит разрешённое",
               pkg="blog", pkg_sub="папка blog/",
               out="Article · Get", inn="items · count",
               head="ПАКЕТ — ЭТО ПАПКА", head_sub='import "sabaq14/blog"',
               foot="ЗАГЛАВНАЯ БУКВА — ДВЕРЬ НАРУЖУ",
               foot_sub="s.Get(…) видно   ·   s.count() не видно"),
    "en": dict(alt="A package is a folder: a capital letter is visible outside, the rest stays in",
               main="main", main_sub="sees the exports",
               pkg="blog", pkg_sub="the blog/ folder",
               out="Article · Get", inn="items · count",
               head="A PACKAGE IS A FOLDER", head_sub='import "sabaq14/blog"',
               foot="A CAPITAL LETTER IS THE DOOR OUT",
               foot_sub="s.Get(…) is visible   ·   s.count() is not"),
}

L17 = {
    "kz": dict(alt="Екі қадам: git add таңдайды, git commit жазып қояды",
               work="жұмыс қалтасы", work_sub="өзгертесіз",
               index="индекс", index_sub="не сақталады",
               hist="тарих", hist_sub="оралу нүктесі",
               cmd1="git add", cmd2="git commit",
               head="ЕКІ ҚАДАМ, БІРЕУ ЕМЕС", head_sub="git status — қазір қай кезеңде тұрғаныңыз",
               foot="ТАРИХ — ҚАЙТА ОРАЛУҒА БОЛАТЫН НҮКТЕЛЕР",
               foot_sub="git log --oneline   ·   git restore --source=HEAD~1 main.go"),
    "ru": dict(alt="Два шага: git add выбирает, git commit записывает",
               work="рабочая папка", work_sub="правите файлы",
               index="индекс", index_sub="что запишем",
               hist="история", hist_sub="точка возврата",
               cmd1="git add", cmd2="git commit",
               head="ДВА ШАГА, А НЕ ОДИН", head_sub="git status — на каком шаге вы сейчас",
               foot="ИСТОРИЯ — ЭТО ТОЧКИ, В КОТОРЫЕ МОЖНО ВЕРНУТЬСЯ",
               foot_sub="git log --oneline   ·   git restore --source=HEAD~1 main.go"),
    "en": dict(alt="Two steps: git add chooses, git commit records",
               work="working folder", work_sub="you edit them",
               index="the index", index_sub="what gets in",
               hist="history", hist_sub="a return point",
               cmd1="git add", cmd2="git commit",
               head="TWO STEPS, NOT ONE", head_sub="git status — which step you are on",
               foot="HISTORY IS A SET OF POINTS YOU CAN RETURN TO",
               foot_sub="git log --oneline   ·   git restore --source=HEAD~1 main.go"),
}

L18 = {
    "kz": dict(alt="Бір тарих, екі жерде: push жібереді, pull кері әкеледі",
               local="компьютеріңіз", local_sub=".git қалтасы",
               remote="GitHub", remote_sub="желідегі көшірме",
               push="git push", pull="git pull",
               head="СОЛ ТАРИХ, ЕКІ ЖЕРДЕ", head_sub="git remote add origin … · git push -u origin main",
               foot="БҰЛ САҚТЫҚ КӨШІРМЕ ЕМЕС — КӨРІНЕТІН ЖҰМЫС",
               foot_sub="README.md — адам ең алдымен оқитын нәрсе"),
    "ru": dict(alt="Одна история в двух местах: push отправляет, pull приносит обратно",
               local="ваш компьютер", local_sub="папка .git",
               remote="GitHub", remote_sub="копия в сети",
               push="git push", pull="git pull",
               head="ТА ЖЕ ИСТОРИЯ, В ДВУХ МЕСТАХ", head_sub="git remote add origin … · git push -u origin main",
               foot="ЭТО НЕ РЕЗЕРВНАЯ КОПИЯ — ЭТО ВИДИМАЯ РАБОТА",
               foot_sub="README.md — первое, что читает человек"),
    "en": dict(alt="One history in two places: push sends, pull brings it back",
               local="your computer", local_sub="the .git folder",
               remote="GitHub", remote_sub="the network copy",
               push="git push", pull="git pull",
               head="THE SAME HISTORY, IN TWO PLACES", head_sub="git remote add origin … · git push -u origin main",
               foot="THIS IS NOT A BACKUP — IT IS WORK PEOPLE CAN SEE",
               foot_sub="README.md — the first thing a person reads"),
}

L19 = {
    "kz": dict(alt="Алты таныс бөлшек бір бағдарламаға жиналады",
               p1="struct", p2="slice", p3="map", p4="әдістер", p5="қателер", p6="пакеттер",
               prog="жазба кітапшасы", prog_sub="аяқталған бағдарлама",
               head="ЖАҢА ЕШТЕҢЕ ЖОҚ — БӘРІ ТАНЫС", head_sub="add · list · find · top · del",
               foot="ЖИНАҚТАУ — БҰЛ ДА ШЕБЕРЛІК",
               foot_sub="go run .   ·   go test ./..."),
    "ru": dict(alt="Шесть знакомых деталей собираются в одну программу",
               p1="struct", p2="slice", p3="map", p4="методы", p5="ошибки", p6="пакеты",
               prog="блокнот в консоли", prog_sub="законченная программа",
               head="НИЧЕГО НОВОГО — ВСЁ ЗНАКОМОЕ", head_sub="add · list · find · top · del",
               foot="СОБРАТЬ — ТОЖЕ УМЕНИЕ",
               foot_sub="go run .   ·   go test ./..."),
    "en": dict(alt="Six familiar parts come together into one program",
               p1="struct", p2="slice", p3="map", p4="methods", p5="errors", p6="packages",
               prog="a console notebook", prog_sub="a finished program",
               head="NOTHING NEW — ALL OF IT FAMILIAR", head_sub="add · list · find · top · del",
               foot="PUTTING IT TOGETHER IS A SKILL TOO",
               foot_sub="go run .   ·   go test ./..."),
}

L20 = {
    "kz": dict(alt="Жауап бір бағытта кетеді: тақырыптар, күй коды, дене",
               h1="тақырыптар", h1_sub="w.Header().Set(…)",
               h2="күй коды", h2_sub="w.WriteHeader(404)",
               h3="дене", h3_sub="fmt.Fprint(w, …)",
               head="ЖАУАП БІР БАҒЫТТА КЕТЕДІ", head_sub="кодтан кейін тақырып жіберілмейді",
               foot="ДЕНЕНІ БІРІНШІ ЖАЗСАҢЫЗ — КОД 200 БОЛЫП ҚАЛАДЫ",
               foot_sub="http: superfluous response.WriteHeader call"),
    "ru": dict(alt="Ответ уходит в одну сторону: заголовки, код состояния, тело",
               h1="заголовки", h1_sub="w.Header().Set(…)",
               h2="код состояния", h2_sub="w.WriteHeader(404)",
               h3="тело", h3_sub="fmt.Fprint(w, …)",
               head="ОТВЕТ УХОДИТ В ОДНУ СТОРОНУ", head_sub="после кода заголовок уже не уедет",
               foot="НАПИСАЛИ ТЕЛО ПЕРВЫМ — КОД НАВСЕГДА 200",
               foot_sub="http: superfluous response.WriteHeader call"),
    "en": dict(alt="A response goes out one way: headers, status code, body",
               h1="headers", h1_sub="w.Header().Set(…)",
               h2="status code", h2_sub="w.WriteHeader(404)",
               h3="body", h3_sub="fmt.Fprint(w, …)",
               head="A RESPONSE GOES OUT ONE WAY", head_sub="after the code, a header no longer travels",
               foot="WRITE THE BODY FIRST AND THE CODE IS 200 FOR GOOD",
               foot_sub="http: superfluous response.WriteHeader call"),
}

L21 = {
    "kz": dict(alt="Бір кіреберіс, бірнеше белгіленген есік: бағыттауыш таңдайды",
               req="сұраныс", req_sub="әдіс + жол",
               mux="ServeMux", mux_sub="кімге береді",
               r1="GET /{$}", r2="GET /read/{slug}", r3="POST /read/{slug}/like",
               head="ӘДІС ПЕН ЖОЛ — БІР ЖОЛДА", head_sub='mux.HandleFunc("POST /read/{slug}/like", …)',
               foot="ТАППАСА — 404, ӘДІСІ БАСҚА БОЛСА — 405, ӨЗІ",
               foot_sub="Allow: GET, HEAD"),
    "ru": dict(alt="Один вход, несколько подписанных дверей: выбирает маршрутизатор",
               req="запрос", req_sub="метод + путь",
               mux="ServeMux", mux_sub="выбирает, кому",
               r1="GET /{$}", r2="GET /read/{slug}", r3="POST /read/{slug}/like",
               head="МЕТОД И ПУТЬ — В ОДНОЙ СТРОКЕ", head_sub='mux.HandleFunc("POST /read/{slug}/like", …)',
               foot="НЕ НАШЁЛ — 404, НЕ ТОТ МЕТОД — 405, САМ",
               foot_sub="Allow: GET, HEAD"),
    "en": dict(alt="One entrance, several labelled doors: the router chooses",
               req="a request", req_sub="method + path",
               mux="ServeMux", mux_sub="picks who gets it",
               r1="GET /{$}", r2="GET /read/{slug}", r3="POST /read/{slug}/like",
               head="METHOD AND PATH ON ONE LINE", head_sub='mux.HandleFunc("POST /read/{slug}/like", …)',
               foot="NOT FOUND — 404; WRONG METHOD — 405; ON ITS OWN",
               foot_sub="Allow: GET, HEAD"),
}

L22 = {
    "kz": dict(alt="Орама: сұраныс өңдеушіге дейін де, кейін де сол қабаттан өтеді",
               b1="сұраныс", b1_sub="кіріп келеді",
               b2="logging", b2_sub="орама",
               b3="ServeMux", b3_sub="бағыттауыш",
               b4="өңдеуші", b4_sub="жауап жазады",
               head="ОРАМА — БАРЛЫҒЫНЫҢ АЙНАЛАСЫНДА", head_sub="http.ListenAndServe(\":8080\", logging(mux))",
               foot="СҰРАНЫС ОДАН ЕКІ РЕТ ӨТЕДІ: БАРҒАНДА ЖӘНЕ ҚАЙТҚАНДА",
               foot_sub="next.ServeHTTP(w, r) — оның алдында да, соңында да код бар"),
    "ru": dict(alt="Обёртка: запрос проходит через неё до обработчика и после него",
               b1="запрос", b1_sub="входит",
               b2="logging", b2_sub="обёртка",
               b3="ServeMux", b3_sub="маршрутизатор",
               b4="обработчик", b4_sub="пишет ответ",
               head="ОБЁРТКА — ВОКРУГ ВСЕГО СРАЗУ", head_sub="http.ListenAndServe(\":8080\", logging(mux))",
               foot="ЗАПРОС ПРОХОДИТ ЧЕРЕЗ НЕЁ ДВАЖДЫ: ТУДА И ОБРАТНО",
               foot_sub="next.ServeHTTP(w, r) — код есть и до него, и после"),
    "en": dict(alt="A wrapper: the request passes through it before the handler and after it",
               b1="a request", b1_sub="comes in",
               b2="logging", b2_sub="the wrapper",
               b3="ServeMux", b3_sub="the router",
               b4="the handler", b4_sub="writes it",
               head="ONE WRAPPER AROUND EVERYTHING", head_sub="http.ListenAndServe(\":8080\", logging(mux))",
               foot="THE REQUEST GOES THROUGH IT TWICE: IN AND OUT",
               foot_sub="next.ServeHTTP(w, r) — there is code before it and after it"),
}

L23 = {
    "kz": dict(alt="Үлгі мен деректер бір бетке қосылады",
               tpl="үлгі", tpl_sub="{{ .Title }}",
               data="деректер", data_sub="[]Article",
               page="бет", page_sub="дайын HTML",
               head="БЛАНК ПЕН ДЕРЕКТЕР — БЕТ БОЛАДЫ",
               head_sub="template.Must(template.ParseFS(files, …))",
               foot="БӨТЕН МӘТІН ӨЗІНЕН-ӨЗІ ЗАЛАЛСЫЗДАНАДЫ",
               foot_sub="«<script>» бетте мәтін болып қалады, код емес"),
    "ru": dict(alt="Шаблон и данные складываются в одну страницу",
               tpl="шаблон", tpl_sub="{{ .Title }}",
               data="данные", data_sub="[]Article",
               page="страница", page_sub="готовый HTML",
               head="БЛАНК И ДАННЫЕ ДАЮТ СТРАНИЦУ",
               head_sub="template.Must(template.ParseFS(files, …))",
               foot="ЧУЖОЙ ТЕКСТ ОБЕЗВРЕЖИВАЕТСЯ САМ",
               foot_sub="«<script>» на странице останется текстом, а не кодом"),
    "en": dict(alt="A template and some data add up to one page",
               tpl="a template", tpl_sub="{{ .Title }}",
               data="the data", data_sub="[]Article",
               page="a page", page_sub="finished HTML",
               head="A FORM PLUS DATA MAKES A PAGE",
               head_sub="template.Must(template.ParseFS(files, …))",
               foot="SOMEBODY ELSE'S TEXT IS MADE SAFE FOR YOU",
               foot_sub="a «<script>» on the page stays text and never runs"),
}

L24 = {
    "kz": dict(alt="Оқырман дерек жібереді, өңдеуші оны браузерді бетке жібереді",
               b1="браузер", b1_sub="форманы толтырады",
               b2="POST /add", b2_sub="өзгертеді",
               b3="GET /", b3_sub="көрсетеді",
               c1="POST", c2="303",
               head="ОҚЫРМАН ДЕРЕК ЖІБЕРЕДІ", head_sub='title := r.FormValue("title")',
               foot="POST → 303 → GET: ЖАҢАРТУ ФОРМАНЫ ҚАЙТА ЖІБЕРМЕЙДІ",
               foot_sub="http.Redirect(w, r, \"/\", http.StatusSeeOther)"),
    "ru": dict(alt="Читатель отправляет данные, обработчик отправляет браузер на страницу",
               b1="браузер", b1_sub="заполняет форму",
               b2="POST /add", b2_sub="меняет",
               b3="GET /", b3_sub="показывает",
               c1="POST", c2="303",
               head="ЧИТАТЕЛЬ ОТПРАВЛЯЕТ ДАННЫЕ", head_sub='title := r.FormValue("title")',
               foot="POST → 303 → GET: ОБНОВЛЕНИЕ НЕ ОТПРАВИТ ФОРМУ ПОВТОРНО",
               foot_sub="http.Redirect(w, r, \"/\", http.StatusSeeOther)"),
    "en": dict(alt="The reader sends data, and the handler sends the browser to a page",
               b1="the browser", b1_sub="fills the form in",
               b2="POST /add", b2_sub="changes things",
               b3="GET /", b3_sub="shows them",
               c1="POST", c2="303",
               head="THE READER SENDS SOMETHING", head_sub='title := r.FormValue("title")',
               foot="POST → 303 → GET: A REFRESH WILL NOT RESUBMIT",
               foot_sub="http.Redirect(w, r, \"/\", http.StatusSeeOther)"),
}

L25 = {
    "kz": dict(alt="Бет жиналады, ал файл сол күйінде беріледі",
               b1="браузер", b1_sub="бетті сұрайды",
               b2="GET /", b2_sub="үлгі + дерек",
               b3="GET /static/", b3_sub="файл сол күйінде",
               b4="бет", b4_sub="мәнерімен",
               head="САЙТ ЖҰМЫС ІСТЕП ҚАНА ҚОЙМАЙ, КӨРІНЕДІ",
               head_sub='mux.Handle("GET /static/", http.FileServerFS(files))',
               foot="ЖИНАЛАТЫН БЕТ ПЕН ӨЗГЕРМЕЙТІН ФАЙЛ — ЕКІ БӨЛЕК ЖОЛ",
               foot_sub="Content-Type: text/css; charset=utf-8"),
    "ru": dict(alt="Страница собирается, а файл отдаётся как есть",
               b1="браузер", b1_sub="просит страницу",
               b2="GET /", b2_sub="шаблон + данные",
               b3="GET /static/", b3_sub="файл как есть",
               b4="страница", b4_sub="со стилем",
               head="САЙТ НЕ ТОЛЬКО РАБОТАЕТ, НО И ВЫГЛЯДИТ",
               head_sub='mux.Handle("GET /static/", http.FileServerFS(files))',
               foot="СОБРАННАЯ СТРАНИЦА И НЕИЗМЕННЫЙ ФАЙЛ — ДВА РАЗНЫХ ПУТИ",
               foot_sub="Content-Type: text/css; charset=utf-8"),
    "en": dict(alt="A page is assembled, a file is handed over unchanged",
               b1="the browser", b1_sub="asks for a page",
               b2="GET /", b2_sub="template + data",
               b3="GET /static/", b3_sub="the file as it is",
               b4="the page", b4_sub="with its style",
               head="THE SITE NOW LOOKS LIKE SOMETHING, NOT ONLY WORKS",
               head_sub='mux.Handle("GET /static/", http.FileServerFS(files))',
               foot="AN ASSEMBLED PAGE AND AN UNCHANGING FILE TAKE DIFFERENT PATHS",
               foot_sub="Content-Type: text/css; charset=utf-8"),
}

L26 = {
    "kz": dict(alt="Паника ұсталады да, оқырман код пен нөмір алады",
               b1="паника", b1_sub="нөл-сөздікке жазу",
               b2="recover", b2_sub="аралық қабат",
               b3="500 + нөмір", b3_sub="журналда стек",
               c1="ұсталды", c2="жауап",
               head="АҚ ЭКРАННЫҢ ОРНЫНА — КОД ПЕН НӨМІР",
               head_sub="if v := recover(); v != nil {",
               foot="ЖАУАП КЕТІП ҮЛГЕРСЕ, ҚҰТҚАРУ КЕШ: 200 ЖӘНЕ ЖАРТЫ БЕТ",
               foot_sub="superfluous response.WriteHeader call"),
    "ru": dict(alt="Панику ловят, и читатель получает код и номер",
               b1="паника", b1_sub="nil-словарь",
               b2="recover", b2_sub="прослойка",
               b3="500 + номер", b3_sub="в журнале стек",
               c1="поймали", c2="ответ",
               head="ВМЕСТО БЕЛОГО ЭКРАНА — КОД И НОМЕР",
               head_sub="if v := recover(); v != nil {",
               foot="ЕСЛИ ОТВЕТ УЖЕ УШЁЛ, СПАСАТЬ ПОЗДНО: 200 И ПОЛСТРАНИЦЫ",
               foot_sub="superfluous response.WriteHeader call"),
    "en": dict(alt="The panic is caught, and the reader gets a code and a number",
               b1="a panic", b1_sub="a nil map",
               b2="recover", b2_sub="middleware",
               b3="500 + a number", b3_sub="stack in the log",
               c1="caught", c2="answer",
               head="A CODE AND A NUMBER INSTEAD OF A BLANK SCREEN",
               head_sub="if v := recover(); v != nil {",
               foot="ONCE THE ANSWER HAS LEFT IT IS TOO LATE: 200 AND HALF A PAGE",
               foot_sub="superfluous response.WriteHeader call"),
}

L27 = {
    "kz": dict(alt="Бір баптау үш жерден келеді, соңғысы жеңеді",
               b1="код", b1_sub="әдепкі мән",
               b2="орта айнымалысы", b2_sub="кодты басады",
               b3="флаг", b3_sub="екеуін де басады",
               c1="басады", c2="басады",
               head="ҚҰПИЯ СӨЗ КОДТА ТҰРМАЙДЫ",
               head_sub='flag.StringVar(&c.dsn, "dsn", env("BLOG_DSN", ""), …)',
               foot="БАПТАУ ЖОҚ БОЛСА — БІРІНШІ СҰРАНЫСТА ЕМЕС, ІСКЕ ҚОСЫЛҒАНДА ҚҰЛАЙМЫЗ",
               foot_sub="баптау: дерекқор жолы берілмеген"),
    "ru": dict(alt="Одна настройка приходит из трёх мест, побеждает последнее",
               b1="код", b1_sub="по умолчанию",
               b2="окружение", b2_sub="подменяет код",
               b3="флаг", b3_sub="подменяет обе",
               c1="подменяет", c2="подменяет",
               head="ПАРОЛЬ НЕ ЖИВЁТ В КОДЕ",
               head_sub='flag.StringVar(&c.dsn, "dsn", env("BLOG_DSN", ""), …)',
               foot="НЕТ НАСТРОЙКИ — ПАДАЕМ ПРИ СТАРТЕ, А НЕ НА ПЕРВОМ ЗАПРОСЕ",
               foot_sub="настройки: не задана строка подключения"),
    "en": dict(alt="One setting arrives from three places, and the last one wins",
               b1="the code", b1_sub="the default",
               b2="the environment", b2_sub="overrides code",
               b3="a flag", b3_sub="overrides both",
               c1="overrides", c2="overrides",
               head="THE PASSWORD DOES NOT LIVE IN THE CODE",
               head_sub='flag.StringVar(&c.dsn, "dsn", env("BLOG_DSN", ""), …)',
               foot="A MISSING SETTING STOPS THE START, NOT THE FIRST REQUEST",
               foot_sub="config: no database string given"),
}

L28 = {
    "kz": dict(alt="Екі кесте және оларды байланыстыратын кілт",
               b1="comments", b1_sub="article_id",
               b2="articles", b2_sub="id",
               b3="join", b3_sub="бір жауап",
               c1="сілтейді", c2="бірге",
               head="ДЕРЕК ҚАЙТА ҚОСУДАН АМАН ҚАЛАДЫ",
               head_sub="join comments c on c.article_id = a.id",
               foot="WHERE ҰМЫТЫЛСА — БАРЛЫҚ ЖОЛ ӨЗГЕРЕДІ",
               foot_sub="select changes(); → 4"),
    "ru": dict(alt="Две таблицы и ключ, который их связывает",
               b1="comments", b1_sub="article_id",
               b2="articles", b2_sub="id",
               b3="join", b3_sub="один ответ",
               c1="ссылается", c2="вместе",
               head="ДАННЫЕ ПЕРЕЖИВАЮТ ПЕРЕЗАПУСК",
               head_sub="join comments c on c.article_id = a.id",
               foot="ЗАБЫЛИ WHERE — ИЗМЕНИЛИСЬ ВСЕ СТРОКИ",
               foot_sub="select changes(); → 4"),
    "en": dict(alt="Two tables and the key that ties them together",
               b1="comments", b1_sub="article_id",
               b2="articles", b2_sub="id",
               b3="join", b3_sub="one answer",
               c1="points at", c2="together",
               head="THE DATA SURVIVES A RESTART",
               head_sub="join comments c on c.article_id = a.id",
               foot="FORGET THE WHERE AND EVERY ROW CHANGES",
               foot_sub="select changes(); → 4"),
}

L29 = {
    "kz": dict(alt="database/sql бір қосылым емес, пул береді",
               b1="бағдарлама", b1_sub="database/sql",
               b2="пул", b2_sub="бірнеше қосылым",
               b3="blog.db", b3_sub="бір файл",
               c1="сұраныс", c2="ашады",
               head="АШУ ӘЛІ ҚОСЫЛУ ДЕГЕН СӨЗ ЕМЕС",
               head_sub="db.Ping()",
               foot="ПРАГМА СҰРАНЫСПЕН ТӨРТТІҢ БІРІН ҒАНА БАПТАЙДЫ",
               foot_sub="foreign_keys: 1 0 0 0"),
    "ru": dict(alt="database/sql даёт пул соединений, а не одно соединение",
               b1="программа", b1_sub="database/sql",
               b2="пул", b2_sub="много соединений",
               b3="blog.db", b3_sub="один файл",
               c1="запрос", c2="открывает",
               head="ОТКРЫТЬ — ЕЩЁ НЕ ЗНАЧИТ ПОДКЛЮЧИТЬСЯ",
               head_sub="db.Ping()",
               foot="PRAGMA ЗАПРОСОМ НАСТРОИТ ОДНО СОЕДИНЕНИЕ ИЗ ЧЕТЫРЁХ",
               foot_sub="foreign_keys: 1 0 0 0"),
    "en": dict(alt="database/sql hands you a pool, not a single connection",
               b1="the program", b1_sub="database/sql",
               b2="the pool", b2_sub="many connections",
               b3="blog.db", b3_sub="one file",
               c1="a query", c2="opens",
               head="OPENING IS NOT YET CONNECTING",
               head_sub="db.Ping()",
               foot="A PRAGMA SENT AS A QUERY REACHES ONE CONNECTION OF FOUR",
               foot_sub="foreign_keys: 1 0 0 0"),
}

L30 = {
    "kz": dict(alt="Бір кесте: жол салынады да, сұраныспен қайтарылады",
               b1="insert", b1_sub="жолды салады",
               b2="articles", b2_sub="бір кесте",
               b3="select", b3_sub="жолды алады",
               c1="жаңа жол", c2="сұраныс",
               head="ДЕРЕК ҚАЙТА ҚОСУДАН АМАН ҚАЛАДЫ",
               head_sub="insert … returning id;",
               foot="WHERE ҰМЫТЫЛСА — БАРЛЫҚ ЖОЛ ӨЗГЕРЕДІ",
               foot_sub="select changes(); → 4"),
    "ru": dict(alt="Одна таблица: строку кладут и достают запросом",
               b1="insert", b1_sub="кладёт строку",
               b2="articles", b2_sub="одна таблица",
               b3="select", b3_sub="достаёт строку",
               c1="новая строка", c2="запрос",
               head="ДАННЫЕ ПЕРЕЖИВАЮТ ПЕРЕЗАПУСК",
               head_sub="insert … returning id;",
               foot="ЗАБЫЛИ WHERE — ИЗМЕНИЛИСЬ ВСЕ СТРОКИ",
               foot_sub="select changes(); → 4"),
    "en": dict(alt="One table: a row goes in and comes back out with a query",
               b1="insert", b1_sub="puts a row in",
               b2="articles", b2_sub="one table",
               b3="select", b3_sub="takes a row out",
               c1="a new row", c2="a query",
               head="THE DATA SURVIVES A RESTART",
               head_sub="insert … returning id;",
               foot="FORGET THE WHERE AND EVERY ROW CHANGES",
               foot_sub="select changes(); → 4"),
}

L31 = {
    "kz": dict(alt="Төрт әрекет, ал біреуі үнсіз жауап береді",
               b1="Exec", b1_sub="update … where",
               b2="дерекқор", b2_sub="жол табылмады",
               b3="RowsAffected", b3_sub="жалғыз белгі",
               c1="сұраныс", c2="сан",
               head="МАҚАЛА ҮСТІНДЕГІ ТӨРТ ӘРЕКЕТ",
               head_sub="insert · select · update · delete",
               foot="ЖОҚ НӘРСЕНІ ЖАҢАРТУ — ҚАТЕ ЕМЕС",
               foot_sub="қате=<nil>   RowsAffected=0"),
    "ru": dict(alt="Четыре действия, и одно отвечает молча",
               b1="Exec", b1_sub="update … where",
               b2="база", b2_sub="строк не нашлось",
               b3="RowsAffected", b3_sub="весь сигнал",
               c1="запрос", c2="число",
               head="ЧЕТЫРЕ ДЕЙСТВИЯ НАД СТАТЬЁЙ",
               head_sub="insert · select · update · delete",
               foot="ОБНОВИТЬ ТО, ЧЕГО НЕТ, — НЕ ОШИБКА",
               foot_sub="ошибка=<nil>   RowsAffected=0"),
    "en": dict(alt="Four actions, and one of them answers in silence",
               b1="Exec", b1_sub="update … where",
               b2="the database", b2_sub="no rows matched",
               b3="RowsAffected", b3_sub="the only signal",
               c1="a query", c2="a number",
               head="FOUR ACTIONS ON AN ARTICLE",
               head_sub="insert · select · update · delete",
               foot="UPDATING WHAT IS NOT THERE IS NOT AN ERROR",
               foot_sub="err=<nil>   RowsAffected=0"),
}

L32 = {
    "kz": dict(alt="Нөмірленген файлдар, журнал және бір транзакция",
               b1="migrations/", b1_sub="нөмірлі файлдар",
               b2="журнал", b2_sub="не қолданылған",
               b3="схема", b3_sub="қадаммен өзгереді",
               c1="кезекпен", c2="бір рет",
               head="КӨШІ-ҚОН: ДЕРЕКҚОР ТАРИХЫ",
               head_sub="0001 · 0002 · 0003",
               foot="ҚОЛДАНЫЛҒАН ФАЙЛ ТҮЗЕТІЛМЕЙДІ",
               foot_sub="келесісі жазылады"),
    "ru": dict(alt="Нумерованные файлы, журнал и одна транзакция",
               b1="migrations/", b1_sub="файлы по номерам",
               b2="журнал", b2_sub="что уже применено",
               b3="схема", b3_sub="меняется шагами",
               c1="по порядку", c2="один раз",
               head="МИГРАЦИИ: ИСТОРИЯ БАЗЫ",
               head_sub="0001 · 0002 · 0003",
               foot="ПРИМЕНЁННЫЙ ФАЙЛ НЕ ПРАВЯТ",
               foot_sub="пишут следующий"),
    "en": dict(alt="Numbered files, a journal and one transaction",
               b1="migrations/", b1_sub="files by number",
               b2="the journal", b2_sub="what is applied",
               b3="the schema", b3_sub="changes in steps",
               c1="in order", c2="once",
               head="MIGRATIONS: THE HISTORY OF A DATABASE",
               head_sub="0001 · 0002 · 0003",
               foot="AN APPLIED FILE IS NEVER EDITED",
               foot_sub="you write the next one"),
}

L33 = {
    "kz": dict(alt="Мән сұраныстан бөлек жүреді, ал екі әрекет бірге",
               b1="сұраныс", b1_sub="мәтіні — біздікі",
               b2="?", b2_sub="мән — бөтен",
               b3="дерекқор", b3_sub="бөлек талдайды",
               c1="мәтін", c2="мән",
               head="СҰРАНЫС — ЖАБЫСТЫРЫЛАТЫН ЖОЛ ЕМЕС",
               head_sub="select … where slug = ?",
               foot="НЕ ЕКІ СҰРАНЫС, НЕ БІРДЕ-БІРІ",
               foot_sub="Begin · Commit · Rollback"),
    "ru": dict(alt="Значение идёт отдельно от запроса, а два действия — вместе",
               b1="запрос", b1_sub="текст пишем мы",
               b2="?", b2_sub="значение чужое",
               b3="база", b3_sub="разбирает порознь",
               c1="текст", c2="значение",
               head="ЗАПРОС — НЕ СТРОКА, КОТОРУЮ СКЛЕИВАЮТ",
               head_sub="select … where slug = ?",
               foot="ЛИБО ОБА ЗАПРОСА, ЛИБО НИ ОДНОГО",
               foot_sub="Begin · Commit · Rollback"),
    "en": dict(alt="The value travels apart from the query, the two actions together",
               b1="the query", b1_sub="its text is ours",
               b2="?", b2_sub="the value is not",
               b3="the database", b3_sub="reads them apart",
               c1="text", c2="a value",
               head="A QUERY IS NOT A STRING YOU GLUE",
               head_sub="select … where slug = ?",
               foot="EITHER BOTH STATEMENTS OR NEITHER",
               foot_sub="Begin · Commit · Rollback"),
}

L34 = {
    "kz": dict(alt="Екі кесте, ал байланыс үшіншісінде",
               b1="articles", b1_sub="мақалалар",
               b2="article_tags", b2_sub="жұптар",
               b3="tags", b3_sub="тегтер",
               c1="сілтеме", c2="сілтеме",
               head="ТЕГ: КӨПКЕ-КӨП БАЙЛАНЫС",
               head_sub="primary key (article_id, tag_id)",
               foot="БАЙЛАНЫС ҮШІНШІ КЕСТЕДЕ ТҰРАДЫ",
               foot_sub="тізімі бар бағанда емес"),
    "ru": dict(alt="Две таблицы, а связь между ними в третьей",
               b1="articles", b1_sub="статьи",
               b2="article_tags", b2_sub="пары",
               b3="tags", b3_sub="теги",
               c1="ссылка", c2="ссылка",
               head="ТЕГИ: СВЯЗЬ МНОГИЕ-КО-МНОГИМ",
               head_sub="primary key (article_id, tag_id)",
               foot="СВЯЗЬ ЖИВЁТ В ТРЕТЬЕЙ ТАБЛИЦЕ",
               foot_sub="а не в колонке со списком"),
    "en": dict(alt="Two tables, and the link between them in a third",
               b1="articles", b1_sub="the articles",
               b2="article_tags", b2_sub="the pairs",
               b3="tags", b3_sub="the tags",
               c1="a reference", c2="a reference",
               head="TAGS: A MANY-TO-MANY LINK",
               head_sub="primary key (article_id, tag_id)",
               foot="THE LINK LIVES IN A THIRD TABLE",
               foot_sub="not in a column holding a list"),
}

L35 = {
    "kz": dict(alt="Сөз, көрсеткіш және реттелген жауап",
               b1="сұрау", b1_sub="адам терген",
               b2="FTS5", b2_sub="сөз көрсеткіші",
               b3="жауап", b3_sub="ең жақсысы алда",
               c1="дайындау", c2="rank",
               head="ІЗДЕУ: LIKE ЕМЕС, КӨРСЕТКІШ",
               head_sub="match … order by rank",
               foot="ПАРАМЕТР SQL-ДІ ҚОРҒАЙДЫ, ІЗДЕУ СҰРАУЫН ЕМЕС",
               foot_sub="сөз тырнақшаға, жұлдызша сыртқа"),
    "ru": dict(alt="Слово, указатель и ответ по порядку",
               b1="запрос", b1_sub="что набрали",
               b2="FTS5", b2_sub="указатель слов",
               b3="ответ", b3_sub="лучшее первым",
               c1="подготовка", c2="rank",
               head="ПОИСК: НЕ LIKE, А УКАЗАТЕЛЬ",
               head_sub="match … order by rank",
               foot="ПАРАМЕТР СПАСАЕТ SQL, НО НЕ ПОИСКОВЫЙ ЗАПРОС",
               foot_sub="слово в кавычки, звёздочку наружу"),
    "en": dict(alt="A word, an index and an answer in order",
               b1="the query", b1_sub="what was typed",
               b2="FTS5", b2_sub="an index of words",
               b3="the answer", b3_sub="best one first",
               c1="preparing", c2="rank",
               head="SEARCH: AN INDEX, NOT A LIKE",
               head_sub="match … order by rank",
               foot="A PARAMETER SAVES THE SQL, NOT THE SEARCH QUERY",
               foot_sub="quote the word, keep the star outside"),
}

L36 = {
    "kz": dict(alt="Бір жақтау, ортақ бөліктер және беттер",
               b1="pages", b1_sub="title мен main",
               b2="base", b2_sub="жақтау",
               b3="бет", b3_sub="дайын HTML",
               c1="блоктар", c2="base орындалады",
               head="МАКЕТ: БАРЛЫҚ БЕТКЕ БІР ЖАҚТАУ",
               head_sub="define · block · template",
               foot="ӘР БЕТКЕ ӨЗ ҮЛГІ ЖИЫНЫ",
               foot_sub="аттар жиын ішінде ортақ"),
    "ru": dict(alt="Одна рамка, общие куски и страницы",
               b1="pages", b1_sub="title и main",
               b2="base", b2_sub="рамка",
               b3="страница", b3_sub="готовый HTML",
               c1="блоки", c2="выполняется base",
               head="МАКЕТ: ОДНА РАМКА НА ВСЕ СТРАНИЦЫ",
               head_sub="define · block · template",
               foot="СВОЙ НАБОР ШАБЛОНОВ НА КАЖДУЮ СТРАНИЦУ",
               foot_sub="имена внутри набора общие"),
    "en": dict(alt="One frame, the shared pieces and the pages",
               b1="pages", b1_sub="title and main",
               b2="base", b2_sub="the frame",
               b3="the page", b3_sub="finished HTML",
               c1="blocks", c2="base is executed",
               head="LAYOUT: ONE FRAME FOR EVERY PAGE",
               head_sub="define · block · template",
               foot="A TEMPLATE SET OF ITS OWN PER PAGE",
               foot_sub="names are shared inside a set"),
}

L37 = {
    "kz": dict(alt="Құпиясөз, тұз және құны",
               b1="құпиясөз", b1_sub="сақталмайды",
               b2="bcrypt", b2_sub="тұз ішінде",
               b3="хеш", b3_sub="дерекқорда",
               c1="10 хеш/с", c2="кері жол жоқ",
               head="ҚҰПИЯСӨЗДІҢ ОРНЫНА НЕ САҚТАЙДЫ",
               head_sub="$2a$10$…",
               foot="БІР ҚҰПИЯСӨЗ — ӘР ЖОЛЫ БАСҚА ХЕШ",
               foot_sub="72 байт = 36 қазақ әрпі"),
    "ru": dict(alt="Пароль, соль и стоимость",
               b1="пароль", b1_sub="не хранится",
               b2="bcrypt", b2_sub="соль внутри",
               b3="хеш", b3_sub="лежит в базе",
               c1="10 хешей/с", c2="обратно нельзя",
               head="ЧТО ХРАНЯТ ВМЕСТО ПАРОЛЯ",
               head_sub="$2a$10$…",
               foot="ОДИН ПАРОЛЬ — КАЖДЫЙ РАЗ ДРУГОЙ ХЕШ",
               foot_sub="72 байта = 36 казахских букв"),
    "en": dict(alt="A password, a salt and a cost",
               b1="the password", b1_sub="never stored",
               b2="bcrypt", b2_sub="salt inside",
               b3="the hash", b3_sub="in the database",
               c1="10 hashes/s", c2="no way back",
               head="WHAT IS KEPT INSTEAD OF A PASSWORD",
               head_sub="$2a$10$…",
               foot="ONE PASSWORD, A DIFFERENT HASH EVERY TIME",
               foot_sub="72 bytes = 36 Kazakh letters"),
}

L38 = {
    "kz": dict(alt="Браузерде нөмірше, қалғаны дерекқорда",
               b1="браузер", b1_sub="нөмірше ғана",
               b2="токен", b2_sub="32 кездейсоқ байт",
               b3="дерекқор", b3_sub="хеш пен иесі",
               c1="Cookie", c2="sha256",
               head="СЕССИЯ: КІМ КІРГЕНІН ЕСТЕ САҚТАУ",
               head_sub="HttpOnly · Secure · SameSite",
               foot="ШЫҒУ — ДЕРЕКҚОРДАҒЫ ЖОЛДЫ ӨШІРУ",
               foot_sub="печеньеге сенбейді"),
    "ru": dict(alt="Номерок в браузере, всё остальное в базе",
               b1="браузер", b1_sub="только номерок",
               b2="токен", b2_sub="32 байта наугад",
               b3="база", b3_sub="хеш и хозяин",
               c1="Cookie", c2="sha256",
               head="СЕССИЯ: КАК СЕРВЕР ПОМНИТ ВОШЕДШЕГО",
               head_sub="HttpOnly · Secure · SameSite",
               foot="ВЫХОД — ЭТО УДАЛИТЬ СТРОКУ В БАЗЕ",
               foot_sub="печенью не верят"),
    "en": dict(alt="A ticket in the browser, the rest in the database",
               b1="the browser", b1_sub="a ticket only",
               b2="the token", b2_sub="32 random bytes",
               b3="the database", b3_sub="hash and owner",
               c1="Cookie", c2="sha256",
               head="SESSIONS: HOW A SERVER REMEMBERS YOU",
               head_sub="HttpOnly · Secure · SameSite",
               foot="SIGNING OUT DELETES THE ROW",
               foot_sub="the cookie is never believed"),
}

L39 = {
    "kz": dict(alt="Кірді ме, табылды ма, оныкі ме",
               b1="кірді ме", b1_sub="жоқ — кіруге",
               b2="табылды ма", b2_sub="жоқ болса — 404",
               b3="оныкі ме", b3_sub="жоқ болса — 403",
               c1="сессия", c2="author_id",
               head="ҚҰҚЫҚ: КІМ НЕНІ ӨЗГЕРТЕ АЛАДЫ",
               head_sub="MayEdit — жалғыз орын",
               foot="БАТЫРМАНЫ ЖАСЫРУ — ҚОРҒАНЫС ЕМЕС",
               foot_sub="сұраныс онсыз да жіберіледі"),
    "ru": dict(alt="Вошёл ли, нашли ли, его ли",
               b1="вошёл ли", b1_sub="нет — на вход",
               b2="нашли ли", b2_sub="нет — 404",
               b3="его ли", b3_sub="нет — 403",
               c1="сессия", c2="author_id",
               head="ПРАВА: КТО ЧТО МОЖЕТ МЕНЯТЬ",
               head_sub="MayEdit — одно место",
               foot="СПРЯТАТЬ КНОПКУ — НЕ ЗАЩИТА",
               foot_sub="запрос отправят и без неё"),
    "en": dict(alt="Signed in, found, and theirs",
               b1="signed in?", b1_sub="no — to sign-in",
               b2="found?", b2_sub="no — 404",
               b3="theirs?", b3_sub="no — 403",
               c1="the session", c2="author_id",
               head="PERMISSIONS: WHO MAY CHANGE WHAT",
               head_sub="MayEdit — one place",
               foot="HIDING A BUTTON IS NOT A GUARD",
               foot_sub="the request goes without it"),
}

L40 = {
    "kz": dict(alt="Экрандау, жетон және үш тақырып",
               b1="бөтен мәтін", b1_sub="оқырманнан",
               b2="экрандау", b2_sub="орнына қарай",
               b3="жетон", b3_sub="формада",
               c1="html/template", c2="POST",
               head="CSRF ПЕН XSS: ҮШ ҚОРҒАНЫС",
               head_sub="экрандау · жетон · тақырыптар",
               foot="БӨТЕН САЙТ БЕТІҢІЗДІ ОҚИ АЛМАЙДЫ",
               foot_sub="сондықтан жетонды білмейді"),
    "ru": dict(alt="Экранирование, жетон и три заголовка",
               b1="чужой текст", b1_sub="от читателя",
               b2="экранирование", b2_sub="по месту",
               b3="жетон", b3_sub="в форме",
               c1="html/template", c2="POST",
               head="CSRF И XSS: ТРИ ЗАЩИТЫ",
               head_sub="экранирование · жетон · заголовки",
               foot="ЧУЖОЙ САЙТ НЕ ПРОЧТЁТ ВАШУ СТРАНИЦУ",
               foot_sub="потому и не знает жетона"),
    "en": dict(alt="Escaping, a token and three headers",
               b1="from outside", b1_sub="a reader typed it",
               b2="escaping", b2_sub="by context",
               b3="the token", b3_sub="in the form",
               c1="html/template", c2="POST",
               head="CSRF AND XSS: THREE GUARDS",
               head_sub="escaping · token · headers",
               foot="ANOTHER SITE CANNOT READ YOUR PAGE",
               foot_sub="which is why it cannot know the token"),
}

L41 = {
    "kz": dict(alt="Өлшем, түр, ат және орын",
               b1="форма", b1_sub="бөтен файл",
               b2="тексеру", b2_sub="алғашқы 512 байт",
               b3="дискідегі ат", b3_sub="біз ойлап табамыз",
               c1="2 МиБ шек", c2="uploads/",
               head="СУРЕТ ЖҮКТЕУ: БӨТЕН ФАЙЛ",
               head_sub="өлшем · түр · ат · орын",
               foot="КЕҢЕЙТІМГЕ ДЕ, БРАУЗЕРГЕ ДЕ СЕНБЕЙДІ",
               foot_sub="мазмұны шешеді"),
    "ru": dict(alt="Размер, тип, имя и место",
               b1="форма", b1_sub="чужой файл",
               b2="проверка", b2_sub="первые 512 байт",
               b3="имя на диске", b3_sub="придумали мы",
               c1="предел 2 МиБ", c2="uploads/",
               head="ЗАГРУЗКА КАРТИНОК: ЧУЖОЙ ФАЙЛ",
               head_sub="размер · тип · имя · место",
               foot="НИ РАСШИРЕНИЮ, НИ БРАУЗЕРУ НЕ ВЕРЯТ",
               foot_sub="решает содержимое"),
    "en": dict(alt="Size, type, name and place",
               b1="the form", b1_sub="a stranger's file",
               b2="the check", b2_sub="first 512 bytes",
               b3="name on disk", b3_sub="we make it up",
               c1="2 MiB cap", c2="uploads/",
               head="UPLOADS: SOMEBODY ELSE'S FILE",
               head_sub="size · type · name · place",
               foot="NEITHER EXTENSION NOR BROWSER IS BELIEVED",
               foot_sub="the content decides"),
}

L42 = {
    "kz": dict(alt="Жағдайлар кестесі, аттар және есеп",
               b1="кесте", b1_sub="жағдай — бір жол",
               b2="t.Run", b2_sub="әрқайсысына ат",
               b3="есеп", b3_sub="қайсысы сынды",
               c1="цикл", c2="go test -v",
               head="ТЕСТ: ТҮЗЕТКЕНДЕ НЕ СЫНҒАНЫН БІЛУ",
               head_sub="кестелік тест",
               foot="«ТЕСТ ҚҰЛАДЫ» ЕМЕС, «ҚАЙСЫСЫ ҚҰЛАДЫ»",
               foot_sub="аты бар жағдай өзін айтады"),
    "ru": dict(alt="Таблица случаев, имена и отчёт",
               b1="таблица", b1_sub="случай — строка",
               b2="t.Run", b2_sub="каждому имя",
               b3="отчёт", b3_sub="что сломалось",
               c1="цикл", c2="go test -v",
               head="ТЕСТЫ: ЧТО СЛОМАЛОСЬ ПРИ ПРАВКЕ",
               head_sub="табличный тест",
               foot="НЕ «ТЕСТ УПАЛ», А «УПАЛ ВОТ ЭТОТ СЛУЧАЙ»",
               foot_sub="имя случая говорит само"),
    "en": dict(alt="A table of cases, names and a report",
               b1="the table", b1_sub="a case per line",
               b2="t.Run", b2_sub="a name each",
               b3="the report", b3_sub="what broke",
               c1="a loop", c2="go test -v",
               head="TESTS: WHAT AN EDIT BROKE",
               head_sub="table-driven tests",
               foot="NOT THE TEST FAILED BUT THIS CASE FAILED",
               foot_sub="the name of a case speaks"),
}

L43 = {
    "kz": dict(alt="Әр сұраныс — өз горутинасы, ортақ жады құлып астында",
               b1="сұраныстар", b1_sub="әрқайсысы бөлек",
               b2="ортақ жад", b2_sub="құлыпсыз — жарыс",
               b3="жауап", b3_sub="дұрыс сан",
               c1="go func", c2="-race",
               head="ГОРУТИНА МЕН АРНА: ВЕБКЕ НЕ ҮШІН",
               head_sub="Mutex · chan · WaitGroup",
               foot="ЖАРЫС ЖҮКТЕМЕНІ КҮТЕДІ, СІЗДІ ЕМЕС",
               foot_sub="сондықтан тестте -race"),
    "ru": dict(alt="Каждый запрос — своя горутина, общая память под замком",
               b1="запросы", b1_sub="каждый в горутине",
               b2="общая память", b2_sub="без замка — гонка",
               b3="ответ", b3_sub="верное число",
               c1="go func", c2="-race",
               head="ГОРУТИНЫ И КАНАЛЫ: ЗАЧЕМ ОНИ ВЕБУ",
               head_sub="Mutex · chan · WaitGroup",
               foot="ГОНКА ЖДЁТ НАГРУЗКИ, А НЕ ВАС",
               foot_sub="поэтому -race в тестах"),
    "en": dict(alt="Each request its own goroutine, shared memory under a lock",
               b1="requests", b1_sub="one each",
               b2="shared memory", b2_sub="no lock, a race",
               b3="the answer", b3_sub="the right number",
               c1="go func", c2="-race",
               head="GOROUTINES AND CHANNELS ON THE WEB",
               head_sub="Mutex · chan · WaitGroup",
               foot="A RACE WAITS FOR LOAD, NOT FOR YOU",
               foot_sub="which is why -race runs in tests"),
}

L44 = {
    "kz": dict(alt="Оқырман кетті — тізбек бойымен тоқтау",
               b1="оқырман", b1_sub="қойындыны жапты",
               b2="r.Context()", b2_sub="өзі жабылады",
               b3="дерекқор", b3_sub="сұраныс үзіледі",
               c1="Done()", c2="...Context",
               head="CONTEXT: СҰРАНЫСТЫ УАҚЫТЫНДА ТОҚТАТУ",
               head_sub="cancel · timeout · value",
               foot="КЕТКЕН ОҚЫРМАНҒА ЖҰМЫС ІСТЕМЕЙДІ",
               foot_sub="кілт — өз типімен"),
    "ru": dict(alt="Читатель ушёл — остановка по всей цепочке",
               b1="читатель", b1_sub="закрыл вкладку",
               b2="r.Context()", b2_sub="закрылся сам",
               b3="база", b3_sub="запрос оборван",
               c1="Done()", c2="...Context",
               head="CONTEXT: ВОВРЕМЯ ОСТАНОВИТЬ ЗАПРОС",
               head_sub="cancel · timeout · value",
               foot="НА УШЕДШЕГО ЧИТАТЕЛЯ НЕ РАБОТАЮТ",
               foot_sub="ключ — своего типа"),
    "en": dict(alt="The reader left, and the whole chain stops",
               b1="the reader", b1_sub="closed the tab",
               b2="r.Context()", b2_sub="closes itself",
               b3="the database", b3_sub="query cut off",
               c1="Done()", c2="...Context",
               head="CONTEXT: STOPPING WORK IN TIME",
               head_sub="cancel · timeout · value",
               foot="NOBODY WORKS FOR A READER WHO LEFT",
               foot_sub="a key of your own type"),
}

L45 = {
    "kz": dict(alt="Құрылым, тегтер және ағын",
               b1="құрылым", b1_sub="бас әріппен ғана",
               b2="тегтер", b2_sub="ат пен omitempty",
               b3="JSON", b3_sub="ағынға жазылады",
               c1="Marshal", c2="Encoder",
               head="JSON ЖӘНЕ ӨЗ API-ІҢ",
               head_sub="blog-ты браузерсіз оқу",
               foot="ҚАТЕ ДЕ JSON БОЛУҒА ТИІС",
               foot_sub="оны бағдарлама оқиды"),
    "ru": dict(alt="Структура, теги и поток",
               b1="структура", b1_sub="только с большой",
               b2="теги", b2_sub="имя и omitempty",
               b3="JSON", b3_sub="пишется в поток",
               c1="Marshal", c2="Encoder",
               head="JSON И СВОЁ API",
               head_sub="блог читают не браузером",
               foot="ОШИБКА ТОЖЕ ДОЛЖНА БЫТЬ JSON",
               foot_sub="её читает программа"),
    "en": dict(alt="A struct, its tags and a stream",
               b1="the struct", b1_sub="exported only",
               b2="the tags", b2_sub="name, omitempty",
               b3="JSON", b3_sub="into a stream",
               c1="Marshal", c2="Encoder",
               head="JSON AND AN API OF YOUR OWN",
               head_sub="a blog read without a browser",
               foot="AN ERROR HAS TO BE JSON TOO",
               foot_sub="a program is reading it"),
}

L46 = {
    "kz": dict(alt="Баптау, белгі және бетбелгі",
               b1="блог аты", b1_sub="бір баптау",
               b2="үлгілер", b2_sub="он бір файл",
               b3="белгіше", b3_sub="қойындыда тұрады",
               c1="BLOG_NAME", c2="/favicon.ico",
               head="ӨЗ БЕЛГІҢ, ӨЗ АТЫҢ",
               head_sub="blog енді сенікі",
               foot="SVG — БҰЛ МӘТІН",
               foot_sub="193 байт, кез келген өлшемде анық"),
    "ru": dict(alt="Настройка, знак и вкладка",
               b1="имя блога", b1_sub="одна настройка",
               b2="шаблоны", b2_sub="11 файлов",
               b3="иконка", b3_sub="видна на вкладке",
               c1="BLOG_NAME", c2="/favicon.ico",
               head="СВОЙ ЗНАК, СВОЁ ИМЯ",
               head_sub="блог становится вашим",
               foot="SVG — ЭТО ТЕКСТ",
               foot_sub="193 байта, чёткий в любом размере"),
    "en": dict(alt="A setting, a mark and a tab",
               b1="the name", b1_sub="one setting",
               b2="templates", b2_sub="eleven files",
               b3="the icon", b3_sub="seen on the tab",
               c1="BLOG_NAME", c2="/favicon.ico",
               head="YOUR OWN MARK AND NAME",
               head_sub="the blog becomes yours",
               foot="SVG IS TEXT",
               foot_sub="193 bytes, sharp at any size"),
}

L47 = {
    "kz": dict(alt="Порт, интерфейс және диск",
               b1="порт", b1_sub="PORT-тан келеді",
               b2="интерфейс", b2_sub="барлық мекенжай",
               b3="диск", b3_sub="деплойда жоғалады",
               c1="/healthz", c2="бір файл",
               head="БЛОГТЫ АДАМДАРҒА ЖЕТКІЗУ",
               head_sub="тегін тарифке қою",
               foot="ТЕГІН ДИСК САҚТАМАЙДЫ",
               foot_sub="дерек сыртқы базада тұрсын"),
    "ru": dict(alt="Порт, интерфейс и диск",
               b1="порт", b1_sub="приходит в PORT",
               b2="интерфейс", b2_sub="все адреса",
               b3="диск", b3_sub="деплой стирает",
               c1="/healthz", c2="один файл",
               head="БЛОГ ТУДА, ГДЕ ЕГО ОТКРОЮТ",
               head_sub="бесплатный тариф",
               foot="БЕСПЛАТНЫЙ ДИСК НЕ ХРАНИТ",
               foot_sub="данные — во внешнюю базу"),
    "en": dict(alt="A port, an interface and a disk",
               b1="the port", b1_sub="arrives in PORT",
               b2="the interface", b2_sub="all, not just one",
               b3="the disk", b3_sub="wiped on deploy",
               c1="/healthz", c2="one file",
               head="PUTTING THE BLOG WITHIN REACH",
               head_sub="on a free plan",
               foot="A FREE DISK KEEPS NOTHING",
               foot_sub="data belongs in a database"),
}

L48 = {
    "kz": dict(alt="Домен, прокси және TLS",
               b1="домен", b1_sub="A-жазба → IP",
               b2="прокси", b2_sub="куәлікті ұстайды",
               b3="TLS", b3_sub="құпиясөз жасырын",
               c1="308 → https", c2="Let's Encrypt",
               head="ӨЗ ДОМЕНІҢ ЖӘНЕ HTTPS",
               head_sub="блог санмен емес, атпен",
               foot="TLS-СІЗ ҚҰПИЯСӨЗ КӨРІНЕДІ",
               foot_sub="өлшенді: 222 байт ашық мәтін"),
    "ru": dict(alt="Домен, прокси и TLS",
               b1="домен", b1_sub="A-запись → IP",
               b2="прокси", b2_sub="держит сертификат",
               b3="TLS", b3_sub="пароль скрыт",
               c1="308 → https", c2="Let's Encrypt",
               head="СВОЙ ДОМЕН И HTTPS",
               head_sub="блог по имени, а не по цифрам",
               foot="БЕЗ TLS ПАРОЛЬ ВИДЕН",
               foot_sub="измерено: 222 байта открытым текстом"),
    "en": dict(alt="A domain, a proxy and TLS",
               b1="the domain", b1_sub="an A record → IP",
               b2="the proxy", b2_sub="holds the cert",
               b3="TLS", b3_sub="nothing readable",
               c1="308 → https", c2="Let's Encrypt",
               head="A DOMAIN OF YOUR OWN AND HTTPS",
               head_sub="a blog by name, not by digits",
               foot="WITHOUT TLS THE PASSWORD SHOWS",
               foot_sub="measured: 222 bytes in the clear"),
}

L49 = {
    "kz": dict(alt="Сигнал, Shutdown және қызмет",
               b1="сигнал", b1_sub="контекст жабылды",
               b2="Shutdown", b2_sub="жауаптар бітсін",
               b3="systemd", b3_sub="қайта қосады",
               c1="SIGTERM", c2="Restart=always",
               head="ҚАЙТА ЖҮКТЕУДЕН АМАН ҚЫЗМЕТ",
               head_sub="blog енді ұқыпты тоқтайды",
               foot="ҮЗІЛУ МЕН ЖАУАП — БІР ЖОЛДЫҢ АЙЫРМАСЫ",
               foot_sub="өлшенді: EOF пен 200"),
    "ru": dict(alt="Сигнал, Shutdown и служба",
               b1="сигнал", b1_sub="контекст закрыт",
               b2="Shutdown", b2_sub="дописать ответы",
               b3="systemd", b3_sub="поднимет снова",
               c1="SIGTERM", c2="Restart=always",
               head="СЛУЖБА, ПЕРЕЖИВАЮЩАЯ ПЕРЕЗАГРУЗКУ",
               head_sub="блог останавливается аккуратно",
               foot="ОБРЫВ ИЛИ ОТВЕТ — РАЗНИЦА В СТРОКЕ",
               foot_sub="измерено: EOF против 200"),
    "en": dict(alt="A signal, Shutdown and a service",
               b1="the signal", b1_sub="context closed",
               b2="Shutdown", b2_sub="finish answers",
               b3="systemd", b3_sub="starts it again",
               c1="SIGTERM", c2="Restart=always",
               head="A SERVICE THAT SURVIVES A REBOOT",
               head_sub="the blog stops tidily now",
               foot="A CUT OR AN ANSWER: ONE LINE APART",
               foot_sub="measured: EOF against 200"),
}

L50 = {
    "kz": dict(alt="Таймаут, көшірме және сипаттама",
               b1="таймаут", b1_sub="баяу ұстамайды",
               b2="көшірме", b2_sub="тірі базадан",
               b3="сипаттама", b3_sub="іздеу көреді",
               c1="ReadHeaderTimeout", c2="vacuum into",
               head="ІСКЕ ҚОСАР АЛДЫНДА",
               head_sub="тексеру парағы",
               foot="ТАЙМАУТСЫЗ СЕРВЕР МӘҢГІ КҮТЕДІ",
               foot_sub="өлшенді: 2 сек күттік те, бас тарттық"),
    "ru": dict(alt="Таймауты, копия и описание",
               b1="таймауты", b1_sub="медленный уйдёт",
               b2="копия", b2_sub="с живой базы",
               b3="описание", b3_sub="его видит поиск",
               c1="ReadHeaderTimeout", c2="vacuum into",
               head="ПЕРЕД ЗАПУСКОМ",
               head_sub="чек-лист, а не надежда",
               foot="СЕРВЕР БЕЗ ТАЙМАУТОВ ЖДЁТ ВЕЧНО",
               foot_sub="измерено: ждали 2 с и сдались первыми"),
    "en": dict(alt="Timeouts, a copy and a description",
               b1="timeouts", b1_sub="the slow let go",
               b2="a backup", b2_sub="from a live db",
               b3="description", b3_sub="search shows it",
               c1="ReadHeaderTimeout", c2="vacuum into",
               head="BEFORE THE LAUNCH",
               head_sub="a checklist, not a hope",
               foot="A SERVER WITHOUT TIMEOUTS WAITS FOREVER",
               foot_sub="measured: we gave up first, after 2 s"),
}

L51 = {
    "kz": dict(alt="Тип параметрі, шарт және әрі қарайғы жол",
               b1="тип параметрі", b1_sub="бір код, көп тип",
               b2="шарт", b2_sub="компилятор көреді",
               b3="әрі қарай", b3_sub="жол бітпейді",
               c1="Sum[T Number]", c2="slices, maps",
               head="ӘРІ ҚАРАЙ ҚАЙДА",
               head_sub="курстың соңғы сабағы",
               foot="ДЖЕНЕРИК — ҚАЙТАЛАУДЫҢ ЕМІ",
               foot_sub="өлшенді: 534 нс пен 1801 нс"),
    "ru": dict(alt="Параметр типа, условие и дорога дальше",
               b1="параметр типа", b1_sub="код один на все",
               b2="условие", b2_sub="проверит сборка",
               b3="дальше", b3_sub="дорога длиннее",
               c1="Sum[T Number]", c2="slices, maps",
               head="КУДА ДАЛЬШЕ",
               head_sub="последний урок курса",
               foot="ДЖЕНЕРИКИ — ЛЕКАРСТВО ОТ ПОВТОРА",
               foot_sub="измерено: 534 нс против 1801 нс"),
    "en": dict(alt="A type parameter, its condition and the road on",
               b1="type parameter", b1_sub="any type at all",
               b2="the condition", b2_sub="the build checks",
               b3="onwards", b3_sub="the road goes on",
               c1="Sum[T Number]", c2="slices, maps",
               head="WHERE TO GO NEXT",
               head_sub="the last lesson of the course",
               foot="GENERICS CURE REPETITION",
               foot_sub="measured: 534 ns against 1801 ns"),
}

PY01 = {
    "kz": dict(alt="Бағалар, салыстыру және сұрақ",
               b1="бағалар", b1_sub="3,48 есе өсті",
               b2="басқа елдер", b2_sub="сан басқа",
               b3="сұрақ", b3_sub="өзіңіз санайсыз",
               c1="FP.CPI.TOTL", c2="urllib + json",
               head="PYTHON: ДЕРЕКТЕН ӨЗ ЕСЕБІҢІЗГЕ",
               head_sub="бірінші сабақ: не үшін санаймыз",
               foot="1000 ТЕҢГЕ ТҰРҒАН НӘРСЕ — 3481",
               foot_sub="өлшенді: бағалар индексі 100 → 348,1"),
    "ru": dict(alt="Цены, сравнение и вопрос",
               b1="цены", b1_sub="в 3,48 раза",
               b2="другие страны", b2_sub="цифры другие",
               b3="вопрос", b3_sub="считаете вы сами",
               c1="FP.CPI.TOTL", c2="urllib + json",
               head="PYTHON: ОТ ДАННЫХ ДО СВОЕЙ СВОДКИ",
               head_sub="первый урок: зачем мы считаем",
               foot="ЧТО СТОИЛО 1000 ТЕНГЕ, СТОИТ 3481",
               foot_sub="измерено: индекс цен 100 → 348,1"),
    "en": dict(alt="Prices, the comparison and the question",
               b1="prices", b1_sub="up 3.48 times",
               b2="other countries", b2_sub="other numbers",
               b3="the question", b3_sub="you count it",
               c1="FP.CPI.TOTL", c2="urllib + json",
               head="PYTHON: FROM DATA TO YOUR OWN DIGEST",
               head_sub="lesson one: why we count",
               foot="WHAT COST 1000 TENGE NOW COSTS 3481",
               foot_sub="measured: the price index 100 to 348.1"),
}

PY02 = {
    "kz": dict(alt="Жүйе, орта және жоба",
               b1="жүйе", b1_sub="ортақ Python",
               b2="орта", b2_sub=".venv, тек жобаға",
               b3="жоба", b3_sub="код пен тізім",
               c1="python3 -m venv", c2="requirements.txt",
               head="ЖҰМЫС ОРНЫ: ЖОБАНЫҢ ОРТАСЫ",
               head_sub="екінші сабақ: не үшін .venv",
               foot="ОРТАДАН ТЫС ПАКЕТ ЖОҚ",
               foot_sub="өлшенді: ModuleNotFoundError"),
    "ru": dict(alt="Система, окружение и проект",
               b1="система", b1_sub="общий Python",
               b2="окружение", b2_sub=".venv проекта",
               b3="проект", b3_sub="код и список",
               c1="python3 -m venv", c2="requirements.txt",
               head="РАБОЧЕЕ МЕСТО: ОКРУЖЕНИЕ ПРОЕКТА",
               head_sub="второй урок: зачем нужен .venv",
               foot="ВНЕ ОКРУЖЕНИЯ ПАКЕТА НЕТ",
               foot_sub="измерено: ModuleNotFoundError"),
    "en": dict(alt="The system, the environment and the project",
               b1="the system", b1_sub="shared Python",
               b2="the venv", b2_sub="this project only",
               b3="the project", b3_sub="code and a list",
               c1="python3 -m venv", c2="requirements.txt",
               head="A WORKPLACE: THE PROJECT'S ENVIRONMENT",
               head_sub="lesson two: what a .venv is for",
               foot="OUTSIDE IT THE PACKAGE IS GONE",
               foot_sub="measured: ModuleNotFoundError"),
}

PY03 = {
    "kz": dict(alt="Сандар, жолдар және шығыс",
               b1="сандар", b1_sub="екі бөлу",
               b2="жолдар", b2_sub="тазалау, кесу",
               b3="f-жол", b3_sub="сан мәтін ішінде",
               c1="// және %", c2="{x:,.0f}",
               head="САНДАР, ЖОЛДАР ЖӘНЕ ШЫҒЫС",
               head_sub="үшінші сабақ: чектен бастаймыз",
               foot="ЕКІ БӨЛУ — ЕКІ ТҮРЛІ ЖАУАП",
               foot_sub="өлшенді: 7/2 = 3.5, 7//2 = 3"),
    "ru": dict(alt="Числа, строки и вывод",
               b1="числа", b1_sub="два деления",
               b2="строки", b2_sub="чистим и режем",
               b3="f-строка", b3_sub="число в тексте",
               c1="// и %", c2="{x:,.0f}",
               head="ЧИСЛА, СТРОКИ И ВЫВОД",
               head_sub="третий урок: начинаем с чека",
               foot="ДВА ДЕЛЕНИЯ — РАЗНЫЕ ОТВЕТЫ",
               foot_sub="измерено: 7/2 = 3.5, 7//2 = 3"),
    "en": dict(alt="Numbers, strings and the output",
               b1="numbers", b1_sub="two divisions",
               b2="strings", b2_sub="clean and cut",
               b3="an f-string", b3_sub="a number inside",
               c1="// and %", c2="{x:,.0f}",
               head="NUMBERS, STRINGS AND OUTPUT",
               head_sub="lesson three: start with a receipt",
               foot="TWO DIVISIONS, TWO ANSWERS",
               foot_sub="measured: 7/2 = 3.5, 7//2 = 3"),
}

PY04 = {
    "kz": dict(alt="Мәтін, сан және ақиқат",
               b1="мәтін", b1_sub="файлдан келеді",
               b2="сан", b2_sub="ақшаға Decimal",
               b3="ақиқат", b3_sub="None нөл емес",
               c1="int(\"260\")", c2="0.1 + 0.2",
               head="АЙНЫМАЛЫЛАР ЖӘНЕ ТИПТЕР",
               head_sub="төртінші сабақ: мәтін мен сан",
               foot="0.1 + 0.2 ҮШ ОНДЫҚҚА ТЕҢ ЕМЕС",
               foot_sub="өлшенді: 0.30000000000000004"),
    "ru": dict(alt="Текст, число и истина",
               b1="текст", b1_sub="так придёт файл",
               b2="число", b2_sub="деньги — Decimal",
               b3="истина", b3_sub="None — не ноль",
               c1="int(\"260\")", c2="0.1 + 0.2",
               head="ПЕРЕМЕННЫЕ И ТИПЫ",
               head_sub="четвёртый урок: текст и число",
               foot="0.1 + 0.2 НЕ РАВНО 0.3",
               foot_sub="измерено: 0.30000000000000004"),
    "en": dict(alt="Text, number and truth",
               b1="text", b1_sub="files give this",
               b2="a number", b2_sub="money: Decimal",
               b3="truth", b3_sub="None is not zero",
               c1="int(\"260\")", c2="0.1 + 0.2",
               head="VARIABLES AND TYPES",
               head_sub="lesson four: text and number",
               foot="0.1 + 0.2 IS NOT 0.3",
               foot_sub="measured: 0.30000000000000004"),
}

PY05 = {
    "kz": dict(alt="Тізім, тілім және кортеж",
               b1="тізім", b1_sub="реті сақталады",
               b2="тілім", b2_sub="оң шек кірмейді",
               b3="кортеж", b3_sub="өзгермейді",
               c1="sum, max, min", c2="sorted(key=)",
               head="ТІЗІМ ЖӘНЕ КОРТЕЖ",
               head_sub="бесінші сабақ: чек өзін санайды",
               foot="sorted() ЖАҢАСЫН, .sort() ӨЗІН",
               foot_sub="өлшенді: .sort() None қайтарады"),
    "ru": dict(alt="Список, срез и кортеж",
               b1="список", b1_sub="порядок хранится",
               b2="срез", b2_sub="правый край вне",
               b3="кортеж", b3_sub="не меняется",
               c1="sum, max, min", c2="sorted(key=)",
               head="СПИСОК И КОРТЕЖ",
               head_sub="пятый урок: чек считает себя",
               foot="sorted() — НОВЫЙ, .sort() — СЕБЯ",
               foot_sub="измерено: .sort() возвращает None"),
    "en": dict(alt="A list, a slice and a tuple",
               b1="a list", b1_sub="order is kept",
               b2="a slice", b2_sub="right end is out",
               b3="a tuple", b3_sub="does not change",
               c1="sum, max, min", c2="sorted(key=)",
               head="A LIST AND A TUPLE",
               head_sub="lesson five: the receipt counts itself",
               foot="sorted() MAKES ONE, .sort() CHANGES ONE",
               foot_sub="measured: .sort() returns None"),
}

PY06 = {
    "kz": dict(alt="Кілт, мән және санағыш",
               b1="кілт", b1_sub="іздеу осымен",
               b2="мән", b2_sub="жоқ кілт — қате",
               b3="санағыш", b3_sub="Counter бір жолда",
               c1=".get(кілт, 0)", c2=".items()",
               head="СӨЗДІК: РЕТКЕ ЕМЕС, БАЙЛАНЫСҚА",
               head_sub="алтыншы сабақ: кілт пен мән",
               foot="ЖОҚ КІЛТ — БОС ЕМЕС, ҚАТЕ",
               foot_sub="өлшенді: KeyError"),
    "ru": dict(alt="Ключ, значение и счётчик",
               b1="ключ", b1_sub="по нему ищут",
               b2="значение", b2_sub="нет ключа — сбой",
               b3="счётчик", b3_sub="Counter в строку",
               c1=".get(ключ, 0)", c2=".items()",
               head="СЛОВАРЬ: СВЯЗЬ ВМЕСТО ПОРЯДКА",
               head_sub="шестой урок: ключ и значение",
               foot="НЕТ КЛЮЧА — НЕ ПУСТОТА, А ОШИБКА",
               foot_sub="измерено: KeyError"),
    "en": dict(alt="A key, a value and a counter",
               b1="a key", b1_sub="what you look by",
               b2="a value", b2_sub="no key: an error",
               b3="a counter", b3_sub="Counter, one line",
               c1=".get(key, 0)", c2=".items()",
               head="A DICTIONARY: A LINK, NOT AN ORDER",
               head_sub="lesson six: key and value",
               foot="A MISSING KEY IS AN ERROR, NOT A ZERO",
               foot_sub="measured: KeyError"),
}

PY07 = {
    "kz": dict(alt="Айыр, тізбек және деректегі олқылық",
               b1="салыстыру", b1_sub="if — бір шарт",
               b2="тізбек", b2_sub="бірінде тоқтайды",
               b3="олқылық", b3_sub="None — нөл емес",
               c1="if / elif / else", c2="is None",
               head="ШАРТТАР: САН ЖОЛДЫ ТАҢДАЙДЫ",
               head_sub="жетінші сабақ: if, elif, else",
               foot="ЖОҚ САН — НӨЛ ЕМЕС",
               foot_sub="өлшенді: 2025 — 11,4%, әлем — 3,0%"),
    "ru": dict(alt="Развилка, цепочка и пропуск в данных",
               b1="сравнение", b1_sub="if — одно условие",
               b2="цепочка", b2_sub="стоп на первом",
               b3="пропуск", b3_sub="None — не ноль",
               c1="if / elif / else", c2="is None",
               head="УСЛОВИЯ: ЧИСЛО ВЫБИРАЕТ ПУТЬ",
               head_sub="седьмой урок: if, elif, else",
               foot="ЧИСЛА НЕТ — ЭТО НЕ НОЛЬ",
               foot_sub="измерено: 2025 — 11,4%, мир — 3,0%"),
    "en": dict(alt="A fork, a chain and a gap in the data",
               b1="a comparison", b1_sub="if: one condition",
               b2="a chain", b2_sub="stops at first",
               b3="a gap", b3_sub="None is not zero",
               c1="if / elif / else", c2="is None",
               head="CONDITIONS: THE NUMBER PICKS THE ROAD",
               head_sub="lesson seven: if, elif, else",
               foot="NO NUMBER IS NOT A ZERO",
               foot_sub="measured: 2025 at 11.4%, the world at 3.0%"),
}

PY08 = {
    "kz": dict(alt="Аралау, аттап өту және шет",
               b1="аралау", b1_sub="ереже бір рет",
               b2="олқылық", b2_sub="есепке кірмейді",
               b3="шет", b3_sub="zip үнсіз қияды",
               c1="for / range", c2="zip / while",
               head="ЦИКЛ: ҚАТАРДЫ БАҒДАРЛАМА АРАЛАЙДЫ",
               head_sub="сегізінші сабақ: for, range, while",
               foot="БЕС ЖЫЛДЫҢ ОРТАШАСЫ — 11,52%",
               foot_sub="өлшенді: алтыдан бесеуінде сан бар"),
    "ru": dict(alt="Проход, пропуск и край",
               b1="проход", b1_sub="правило один раз",
               b2="пропуск", b2_sub="в счёт не идёт",
               b3="край", b3_sub="zip молча режет",
               c1="for / range", c2="zip / while",
               head="ЦИКЛЫ: РЯД ПРОХОДИТ ПРОГРАММА",
               head_sub="восьмой урок: for, range, while",
               foot="СРЕДНЕЕ ЗА ПЯТЬ ЛЕТ — 11,52%",
               foot_sub="измерено: пять лет с числами из шести"),
    "en": dict(alt="The walk, the skip and the edge",
               b1="the walk", b1_sub="the rule once",
               b2="the skip", b2_sub="left out",
               b3="the edge", b3_sub="zip cuts, quietly",
               c1="for / range", c2="zip / while",
               head="LOOPS: THE PROGRAM WALKS THE SERIES",
               head_sub="lesson eight: for, range, while",
               foot="FIVE-YEAR AVERAGE: 11.52%",
               foot_sub="measured: five of six years have a figure"),
}



if __name__ == "__main__":
    out = os.path.join(os.path.dirname(__file__), "..", "..", "web", "static", "course", "go")
    out = os.path.normpath(out)
    os.makedirs(out, exist_ok=True)
    # The second course keeps its maps apart: same shapes and palette, its own
    # folder, so a lesson number in one course never collides with the other.
    py_out = os.path.normpath(os.path.join(os.path.dirname(__file__), "..", "..",
                                           "web", "static", "course", "py"))
    os.makedirs(py_out, exist_ok=True)
    for name, scene, table in (("prices", pymap01, PY01),
                               ("workspace", pymap02, PY02),
                               ("numbers", pymap03, PY03),
                               ("types", pymap04, PY04),
                               ("list", pymap05, PY05),
                               ("dict", pymap06, PY06),
                               ("conditions", pymap07, PY07),
                               ("loops", pymap08, PY08)):
        for lang, strings in table.items():
            path = os.path.join(py_out, f"map-{name}-{lang}.svg")
            with open(path, "w", encoding="utf-8") as f:
                f.write(render(scene, strings, name, lang))
            print("wrote", path)
    for name, scene, table in (("internet", map00, L), ("first-server", map01, L01),
                               ("workspace", map02, L02),
                               ("types", map03, L03),
                               ("functions", map04, L04),
                               ("launch", map05, L05),
                               ("loops", map06, L06),
                               ("bytes", map07, L07),
                               ("slices", map08, L08),
                               ("maps", map09, L09),
                               ("test", map10, L10),
                               ("start", map11, L11),
                               ("struct", map12, L12),
                               ("pointers", map13, L13),
                               ("interface", map14, L14),
                               ("errors", map15, L15),
                               ("packages", map16, L16),
                               ("git", map17, L17),
                               ("github", map18, L18),
                               ("capstone", map19, L19),
                               ("http", map20, L20),
                               ("routes", map21, L21),
                               ("middleware", map22, L22),
                               ("templates", map23, L23),
                               ("forms", map24, L24),
                               ("static", map25, L25),
                               ("errors2", map26, L26),
                               ("config", map27, L27),
                               ("sql", map28, L28),
                               ("dbsql", map29, L29),
                               ("sqlone", map30, L30),
                               ("crud", map31, L31),
                               ("migrate", map32, L32),
                               ("tx", map33, L33),
                               ("tags", map34, L34),
                               ("search", map35, L35),
                               ("layout", map36, L36),
                               ("passwords", map37, L37),
                               ("sessions", map38, L38),
                               ("perms", map39, L39),
                               ("csrf", map40, L40),
                               ("upload", map41, L41),
                               ("tests", map42, L42),
                               ("race", map43, L43),
                               ("context", map44, L44),
                               ("json", map45, L45),
                               ("brand", map46, L46),
                               ("deploy", map47, L47),
                               ("https", map48, L48),
                               ("systemd", map49, L49),
                               ("checklist", map50, L50),
                               ("next", map51, L51)):
        for lang, strings in table.items():
            path = os.path.join(out, f"map-{name}-{lang}.svg")
            with open(path, "w", encoding="utf-8") as f:
                f.write(render(scene, strings, name, lang))
            print("wrote", path)
