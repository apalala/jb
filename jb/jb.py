# Copyright (c) 2017-2026 Juancarlo Añez (apalala@gmail.com)
# SPDX-License-Identifier: MIT
from __future__ import annotations

import random
import re
import sys
import time
import urllib.request
from dataclasses import dataclass
from importlib.resources import files
from typing import Iterator

JB_HEADER = "# Johannes Blues - A view into great literary works"

STREAM_TIME = 5.0


@dataclass
class Work:
    id: str          # "pg1524", "pg2701"
    type: str        # "T" or "N"
    window_size: int # 15, 25


WORKS_DATABASE: list[Work] = [
    Work("pg1524", "T", 15),  # The Tragedy of Hamlet
    Work("pg2701", "N", 25),  # Moby-Dick; or, The Whale
    Work("pg1508", "T", 15),  # The Taming of the Shrew
    Work("pg84",   "N", 25),  # Frankenstein; or, The Modern Prometheus
]

THEATRE_CLEANING_PATTERNS = [
    r"(?m)^[A-Z0-9_\s]{2,15}[.:]\s*",
    r"(?m)[\[(].*?[\])]",
]

NOVEL_CLEANING_PATTERNS = [
    r"^(CHAPTER|C_H_A_P_T_E_R)\s+[IVX0-9]+.*",
]

CLEANING_PATTERNS = {
    "T": THEATRE_CLEANING_PATTERNS,
    "N": NOVEL_CLEANING_PATTERNS,
}

GUTENBERG_URL = "https://www.gutenberg.org/cache/epub/{gid}/pg{gid}.txt"


def _find_work_file(gid: str, suffix: str):
    for f in (files("jb") / "works").iterdir():
        if f.name.endswith(suffix) and gid in f.name:
            return f
    return None


def _extract_id(work_id: str) -> str:
    return work_id.removeprefix("pg")


def load_work(work_id: str) -> str:
    gid = _extract_id(work_id)
    padded = f"pg{gid.zfill(6)}"

    bmx = _find_work_file(padded, ".txt.bmx")
    if bmx is not None:
        from .bmx import unseal_text

        return unseal_text(bmx.read_text())

    txt = _find_work_file(padded, ".txt")
    if txt is not None:
        return txt.read_text(encoding="utf-8")

    url = GUTENBERG_URL.format(gid=gid)
    with urllib.request.urlopen(url) as response:
        return response.read().decode("utf-8")


def clean_work(work_type: str, text: str) -> str:
    raw_lines = [line.strip() for line in text.splitlines()]
    lines = [line for line in raw_lines if line and re.search(r"\w", line)]

    start_idx = 0
    for idx, line in enumerate(lines[:500]):
        upper = line.upper()
        if "START OF THE PROJECT" in upper or "START OF THIS PROJECT" in upper:
            start_idx = idx + 1
            break

    end_idx = len(lines)
    start = max(0, len(lines) - 1000)
    for idx, line in enumerate(lines[start:], start=start):
        upper = line.upper()
        if "END OF THE PROJECT" in upper or "END OF THIS PROJECT" in upper:
            end_idx = idx
            break

    body = lines[start_idx:end_idx]
    patterns = CLEANING_PATTERNS.get(work_type, [])
    cleaned: list[str] = []

    for line in body:
        result = line
        for pat in patterns:
            result = re.sub(pat, "", result)
        result = result.strip()
        if result and re.search(r"[a-zA-Z]", result):
            cleaned.append(result)

    return "\n".join(cleaned)


def print_work(text: str, window_size: int = 15, stream_time: float | None = None) -> None:
    verses = [line for line in text.split("\n") if line]
    if not verses:
        return

    stream = stream_blue_verses(verses, window_size=window_size)
    out = sys.stdout if sys.stdout.isatty() else sys.stderr
    print(JB_HEADER)
    start = _time()
    deadline = stream_time if stream_time is not None else STREAM_TIME
    while _time() - start < deadline:
        try:
            print(next(stream), file=out)
        except KeyboardInterrupt:
            continue
        except StopIteration:
            stream = stream_blue_verses(verses, window_size=window_size)
        except BrokenPipeError:
            break


def stream_blue_signal(beta: float = 0.65) -> Iterator[float]:
    last_white = random.normalvariate(0.0, 1.0)
    last_blue = 0.0

    while True:
        white = random.normalvariate(0.0, 1.0)
        blue = (white - last_white) + (beta * last_blue)
        last_white = white
        last_blue = blue
        yield blue


def stream_blue_verses(verses: list[str], window_size: int = 10) -> Iterator[str]:
    if not verses:
        raise ValueError("The text source array cannot be empty.")

    total_lines = len(verses)
    if total_lines == 1:
        while True:
            yield verses[0]

    window = min(window_size, total_lines - 1)
    if window < 1:
        window = 1
    noise_stream = stream_blue_signal()
    current_idx = random.randint(0, total_lines - 1)
    recent_verses: set[str] = set()
    recent_order: list[str] = []

    for _ in range(window):
        recent_verses.add("")
        recent_order.append("")

    while True:
        try:
            try:
                step = next(noise_stream)
            except StopIteration:
                continue

            jump = int(step * window_size)
            current_idx = (current_idx + jump) % total_lines

            verse = verses[current_idx]
            if verse in recent_verses:
                continue

            old = recent_order.pop(0)
            recent_verses.discard(old)
            recent_verses.add(verse)
            recent_order.append(verse)

            yield verse
        except KeyboardInterrupt:
            continue
        except BrokenPipeError:
            break


def _time() -> float:
    while True:
        try:
            return time.time()
        except KeyboardInterrupt:
            continue


def main() -> int:
    work = random.choice(WORKS_DATABASE)
    raw = load_work(work.id)
    cleaned = clean_work(work.type, raw)
    print_work(cleaned, window_size=work.window_size)
    return 0


if __name__ == "__main__":
    import sys

    try:
        sys.exit(main())
    except BrokenPipeError:
        sys.exit(0)
