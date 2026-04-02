#!/usr/bin/env python3
"""
gen-textures.py — generates cat-themed replacement textures for primary Minecraft mobs.

Run with:
    python3 scripts/gen-textures.py

Outputs PNG files into resource_pack/textures/entity/ at the same paths that
Minecraft Bedrock uses for vanilla mob textures.  Placing files at these paths
inside a resource pack automatically overrides the built-in textures.

Mobs and their cat themes
──────────────────────────
  zombie        → Orange tabby        (64 × 64)
  skeleton      → Cream/white cat     (64 × 64)
  creeper       → Black cat           (64 × 32)
  spider        → Siamese             (64 × 32)
  enderman      → Persian gray        (64 × 32)
  cow           → Tuxedo              (64 × 32)
  pig           → Calico              (64 × 32)
  chicken       → White cat           (64 × 32)
  sheep body    → Ragdoll tan         (64 × 32)
  sheep wool    → Ragdoll cream tabby (64 × 32)

To customize a skin, open the relevant PNG in any image editor (e.g. GIMP,
Aseprite, or Paint.NET) and repaint it — no code changes required.
"""

import os
import struct
import zlib

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
TEXTURES_DIR = os.path.join(ROOT, "resource_pack", "textures", "entity")


# ── Minimal PNG encoder (no external dependencies) ───────────────────────────

def _crc32(data: bytes) -> int:
    return zlib.crc32(data) & 0xFFFFFFFF


def _chunk(tag: bytes, data: bytes) -> bytes:
    payload = tag + data
    return struct.pack(">I", len(data)) + payload + struct.pack(">I", _crc32(payload))


def make_png(width: int, height: int, pixel_fn) -> bytes:
    """
    Build a valid RGBA PNG from a pixel function.

    pixel_fn(x, y) must return (r, g, b, a) with each component in 0–255.
    """
    raw = bytearray()
    for y in range(height):
        raw.append(0)          # PNG filter type: None
        for x in range(width):
            r, g, b, a = pixel_fn(x, y)
            raw += bytes([r & 0xFF, g & 0xFF, b & 0xFF, a & 0xFF])

    compressed = zlib.compress(bytes(raw), 9)
    ihdr = struct.pack(">II", width, height) + bytes([8, 6, 0, 0, 0])

    return (
        b"\x89PNG\r\n\x1a\n"
        + _chunk(b"IHDR", ihdr)
        + _chunk(b"IDAT", compressed)
        + _chunk(b"IEND", b"")
    )


# ── Pixel-pattern helpers ─────────────────────────────────────────────────────

def tabby(base, stripe, stripe_width=4, offset=0):
    """Horizontal tabby stripes alternating between *base* and *stripe* colours."""
    br, bg, bb = base
    sr, sg, sb = stripe

    def _pixel(x, y):
        if ((y + offset) // stripe_width) % 2 == 1:
            return (sr, sg, sb, 255)
        return (br, bg, bb, 255)

    return _pixel


def solid(r, g, b, a=255):
    return lambda x, y: (r, g, b, a)


def siamese(width, height, body, points, blend_depth=10):
    """
    Light body with dark 'points' that fade in from the texture edges —
    the classic Siamese colour-point look.
    """
    br, bg, bb = body
    pr, pg, pb = points

    def _pixel(x, y):
        edge_dist = min(x, y, width - 1 - x, height - 1 - y)
        if edge_dist < blend_depth:
            t = 1.0 - edge_dist / blend_depth
            rr = int(br + t * (pr - br))
            gg = int(bg + t * (pg - bg))
            bb2 = int(bb + t * (pb - bb))
            return (rr, gg, bb2, 255)
        return (br, bg, bb, 255)

    return _pixel


def tuxedo(width, height):
    """Black-and-white tuxedo: black on the left/top, white on the right/bottom."""
    def _pixel(x, y):
        fx = x / width
        fy = y / height
        if (fx < 0.35 and fy < 0.55) or (fx > 0.65 and fy > 0.45):
            return (24, 24, 24, 255)
        return (240, 240, 240, 255)

    return _pixel


def calico(width, height):
    """
    Orange + white + black calico patches.
    The UV is divided into simple rectangular regions for each colour.
    """
    # Each entry: (x_frac_min, y_frac_min, x_frac_max, y_frac_max, rgba)
    patches = [
        (0.00, 0.00, 0.40, 0.50, (228, 96, 18, 255)),   # orange
        (0.40, 0.00, 1.00, 0.40, (240, 240, 240, 255)),  # white
        (0.00, 0.50, 0.50, 1.00, (240, 240, 240, 255)),  # white
        (0.50, 0.45, 1.00, 0.70, (28, 28, 28, 255)),     # black
        (0.50, 0.70, 1.00, 1.00, (228, 96, 18, 255)),    # orange
    ]

    def _pixel(x, y):
        fx = x / width
        fy = y / height
        for x0, y0, x1, y1, col in patches:
            if x0 <= fx < x1 and y0 <= fy < y1:
                return col
        return (240, 240, 240, 255)

    return _pixel


# ── Save helper ───────────────────────────────────────────────────────────────

def save(relative_path: str, width: int, height: int, pixel_fn):
    full_path = os.path.join(TEXTURES_DIR, relative_path)
    os.makedirs(os.path.dirname(full_path), exist_ok=True)
    with open(full_path, "wb") as fh:
        fh.write(make_png(width, height, pixel_fn))
    print(f"  ✔  textures/entity/{relative_path}  ({width}×{height})")


# ── Mob texture generators ────────────────────────────────────────────────────

def gen_zombie():
    """Orange tabby — Minecraft Bedrock zombie.png is 64 × 64."""
    save(
        "zombie/zombie.png", 64, 64,
        tabby(base=(224, 120, 48), stripe=(160, 72, 16), stripe_width=4),
    )


def gen_skeleton():
    """Cream/white cat — skeleton.png is 64 × 64."""
    save(
        "skeleton/skeleton.png", 64, 64,
        tabby(base=(240, 232, 208), stripe=(200, 185, 160), stripe_width=6),
    )


def gen_creeper():
    """Black cat with golden eyes — creeper.png is 64 × 32."""
    def _pixel(x, y):
        # Two golden eye patches on the face region of the UV map
        if (10 <= x <= 17 and 4 <= y <= 8) or (42 <= x <= 49 and 4 <= y <= 8):
            return (255, 210, 0, 255)
        return (22, 22, 22, 255)

    save("creeper/creeper.png", 64, 32, _pixel)


def gen_spider():
    """Siamese — spider.png is 64 × 32."""
    save(
        "spider/spider.png", 64, 32,
        siamese(64, 32, body=(245, 222, 179), points=(61, 31, 0), blend_depth=10),
    )


def gen_enderman():
    """Persian gray — enderman.png is 64 × 32."""
    save(
        "enderman/enderman.png", 64, 32,
        tabby(base=(88, 78, 108), stripe=(56, 48, 72), stripe_width=4),
    )


def gen_cow():
    """Tuxedo — cow.png is 64 × 32."""
    save("cow/cow.png", 64, 32, tuxedo(64, 32))


def gen_pig():
    """Calico — pig.png is 64 × 32."""
    save("pig/pig.png", 64, 32, calico(64, 32))


def gen_chicken():
    """White cat — chicken.png is 64 × 32."""
    save(
        "chicken.png", 64, 32,
        tabby(base=(250, 248, 240), stripe=(210, 200, 185), stripe_width=6),
    )


def gen_sheep():
    """Ragdoll — two files: sheep body (64 × 32) and sheep wool overlay (64 × 32)."""
    save("sheep/sheep.png", 64, 32, solid(160, 148, 132))
    save(
        "sheep/sheep_fur.png", 64, 32,
        tabby(base=(236, 216, 180), stripe=(200, 175, 135), stripe_width=5),
    )


# ── Entry point ───────────────────────────────────────────────────────────────

if __name__ == "__main__":
    print("Generating CatMob Madness textures…\n")
    gen_zombie()
    gen_skeleton()
    gen_creeper()
    gen_spider()
    gen_enderman()
    gen_cow()
    gen_pig()
    gen_chicken()
    gen_sheep()
    print("\nDone ✅  — edit the PNGs in resource_pack/textures/entity/ to customise each skin.")
