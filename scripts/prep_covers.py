#!/usr/bin/env python3
"""Turn freely-licensed source images into Shanraq article covers.

House style (see waves 11-13): a real photograph or public-domain artwork,
cropped to 16:9 and served as WebP — realistic, not stylised. Processing is
deliberately light: a smart crop, a modest contrast/saturation lift, and a
gentle vignette so white cover text stays readable.

Sources MUST be public domain, CC0, or CC BY (never share-alike — the covers
ship inside a proprietary layout). Record the credit in the article body.

Usage:
    python3 scripts/prep_covers.py <src_dir> <manifest.json>

manifest.json: {"<rubric>": {"file": "world.jpg", "out": "world/hormuz.webp",
                             "focus": "center|top|bottom",
                             "trim": [left, top, right, bottom]}}

"trim" is optional and takes fractions of the source (0-1) to cut away before
cropping — use it to remove broadcast chyrons, station bugs, or mount borders.
"""
import json
import pathlib
import sys

from PIL import Image, ImageEnhance, ImageDraw, ImageFilter

W, H = 1200, 675  # 16:9, also a good Open Graph size
QUALITY = 84


def trim_edges(im: Image.Image, trim) -> Image.Image:
    """Cut fractional borders (e.g. a broadcast chyron) before the 16:9 crop."""
    if not trim:
        return im
    l, t, r, b = trim
    w, h = im.size
    return im.crop((int(w * l), int(h * t), int(w * (1 - r)), int(h * (1 - b))))


def smart_crop(im: Image.Image, focus: str) -> Image.Image:
    """Crop to 16:9 around the requested focus band, then resize."""
    im = im.convert("RGB")
    target = W / H
    w, h = im.size
    if w / h > target:  # too wide -> trim sides, keep the middle
        new_w = int(h * target)
        left = (w - new_w) // 2
        im = im.crop((left, 0, left + new_w, h))
    else:  # too tall -> trim vertically around the focus band
        new_h = int(w / target)
        if focus == "top":
            top = 0
        elif focus == "bottom":
            top = h - new_h
        else:
            top = (h - new_h) // 2
        im = im.crop((0, top, w, top + new_h))
    return im.resize((W, H), Image.LANCZOS)


def grade(im: Image.Image) -> Image.Image:
    """Light editorial grade — keep the image reading as a real photograph."""
    im = ImageEnhance.Contrast(im).enhance(1.06)
    im = ImageEnhance.Color(im).enhance(1.04)
    im = ImageEnhance.Sharpness(im).enhance(1.08)
    return im


def vignette(im: Image.Image) -> Image.Image:
    """Darken the bottom so the overlaid title/badge keeps contrast."""
    grad = Image.new("L", (1, H), 0)
    for y in range(H):
        t = y / (H - 1)
        # flat for the top two-thirds, easing to ~40% darkness at the very bottom
        v = 0 if t < 0.55 else int(((t - 0.55) / 0.45) ** 1.5 * 100)
        grad.putpixel((0, y), v)
    mask = grad.resize((W, H)).filter(ImageFilter.GaussianBlur(8))
    shade = Image.new("RGB", (W, H), (18, 12, 10))
    return Image.composite(shade, im, mask.point(lambda p: p))


def main() -> int:
    if len(sys.argv) != 3:
        print(__doc__)
        return 2
    src_dir = pathlib.Path(sys.argv[1])
    manifest = json.loads(pathlib.Path(sys.argv[2]).read_text())
    out_root = pathlib.Path("web/static/covers")

    for rubric, spec in manifest.items():
        src = src_dir / spec["file"]
        if not src.exists():
            print(f"!! {rubric}: missing {src}")
            continue
        im = Image.open(src)
        im = trim_edges(im, spec.get("trim"))
        im = smart_crop(im, spec.get("focus", "center"))
        im = grade(im)
        im = vignette(im)
        dest = out_root / spec["out"]
        dest.parent.mkdir(parents=True, exist_ok=True)
        im.save(dest, "WEBP", quality=QUALITY, method=6)
        kb = dest.stat().st_size // 1024
        print(f"{rubric:11} -> {dest}  ({kb} KB)")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
