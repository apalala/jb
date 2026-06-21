# Copyright (c) 2017-2026 Juancarlo Añez (apalala@gmail.com)
# SPDX-License-Identifier: MIT
from __future__ import annotations

import random
import re
import time
import urllib.request
from pathlib import Path
from typing import Iterator

STREAM_TIME = 5.0

ROOT_PATH = Path(__file__).parent

WORKS_DATABASE = {
    "hamlet": Path("works/pg1524.txt.bmx"),
    "mobi_dick": Path("works/pg2700.txt.bmx"),
}


# Project Gutenberg Stable UTF-8 Raw Text URLs
HAMLET_URL = "https://www.gutenberg.org/cache/epub/1524/pg1524.txt"
MOBY_DICK_URL = "https://www.gutenberg.org/cache/epub/2701/pg2701.txt"

# Agnostic Cleaning Presets
THEATRE_CLEANING_PATTERNS = [
    # 1. Matches leading speaker tags: e.g., "HAMLET.", "HORATIO:", "HAM_1."
    r"(?m)^[A-Z0-9_\s]{2,15}[.:]\s*",
    # 2. Matches stage directions inside brackets or parentheses: e.g., "[Exit Ghost]"
    r"(?m)[\[(].*?[\])]",
]

NOVEL_CLEANING_PATTERNS = [
    # Matches typical Gutenberg chapter headers or illustration markers
    r"^(CHAPTER|C_H_A_P_T_E_R)\s+[IVX0-9]+.*",
]


def fetch_and_parse_verses(
    url: str,
    cleaning_patterns: list[str] | None = None,
) -> list[str]:
    """Downloads a text file from Gutenberg, strips metadata, and applies cleaning filters."""
    # print(f"Fetching source text from {url}...")
    with urllib.request.urlopen(url) as response:
        raw_text = response.read().decode("utf-8")

    return parse_verses(raw_text, cleaning_patterns)


def load_work(
    name: str,
    cleaning_patterns: list[str],
) -> list[str]:
    from .bmx import unseal_text

    file_path = WORKS_DATABASE[name]
    path = ROOT_PATH / file_path
    if not path.exists():
        path = ROOT_PATH.parent / file_path

    compressed = path.read_text()
    unsealed = unseal_text(compressed)
    return parse_verses(unsealed, cleaning_patterns=cleaning_patterns)


def parse_verses(
    raw_text: str,
    cleaning_patterns: list[str],
) -> list[str]:
    raw_lines = [line.strip() for line in raw_text.splitlines()]

    # Basic cleanup: remove empty lines and metadata lines
    lines = [line for line in raw_lines if line and re.search(r"\w", line)]

    # Slice out Gutenberg metadata blocks (headers/footers)
    start_idx = 0
    for idx, line in enumerate(lines[:500]):
        if (
            "START OF THE PROJECT" in line.upper()
            or "START OF THIS PROJECT" in line.upper()
        ):
            start_idx = idx + 1
            break

    end_idx = len(lines)
    for idx, line in enumerate(lines[-1000:]):
        if (
            "END OF THE PROJECT" in line.upper()
            or "END OF THIS PROJECT" in line.upper()
        ):
            end_idx = len(lines) - 1000 + idx
            break

    active_body = lines[start_idx:end_idx]
    patterns = cleaning_patterns
    cleaned_verses: list[str] = []

    for line in active_body:
        # Apply the provided sequence of agnostic filters
        for pattern in patterns:
            line = re.sub(pattern, "", line)

        line = line.strip()

        # Keep line if it still contains readable text post-cleansing
        if line and re.search(r"[a-zA-Z]", line):
            cleaned_verses.append(line)

    return cleaned_verses


def stream_blue_signal(beta: float = 0.65) -> Iterator[float]:
    """Generates an infinite stream of high-passed 1D Blue Noise scalar values."""
    last_white = random.normalvariate(0.0, 1.0)
    last_blue = 0.0

    while True:
        white = random.normalvariate(0.0, 1.0)
        blue = (white - last_white) + (beta * last_blue)
        last_white = white
        last_blue = blue
        yield blue


def stream_blue_verses(verses: list[str], window_size: int = 10) -> Iterator[str]:
    """Streams lines from a text corpus using a blue noise index shaper.

    Uses a rolling window mechanism to prevent immediate local clustering
    while traveling smoothly through the available line matrix.
    """
    if not verses:
        raise ValueError("The text source array cannot be empty.")

    total_lines = len(verses)
    noise_stream = stream_blue_signal()
    # Start our tracking index in the middle of the play
    current_idx = random.randint(0, total_lines - 1)
    recent_verses = set([""] * window_size)
    while True:
        # Pull next raw blue noise step (-3.0 to +3.0 scale roughly)
        step = next(noise_stream)

        # Translate the noise value into a bounded structural jump index
        jump = int(step * window_size)

        # Apply the shift and wrap around seamlessly using modulo arithmetic
        current_idx = (current_idx + jump) % total_lines

        verse = verses[current_idx]
        if verse in recent_verses:
            continue
        recent_verses.pop()
        recent_verses.add(verse)
        yield verse


def print_hamlet_verses() -> None:
    try:
        hamlet_lines = load_work("hamlet", cleaning_patterns=THEATRE_CLEANING_PATTERNS)
    except FileNotFoundError:
        hamlet_lines = fetch_and_parse_verses(
            HAMLET_URL, cleaning_patterns=THEATRE_CLEANING_PATTERNS
        )
    hamlet_stream = stream_blue_verses(hamlet_lines, window_size=15)

    out = sys.stdout if sys.stdout.isatty() else sys.stderr
    start = time.time()
    while time.time() - start < STREAM_TIME:
        print(next(hamlet_stream), file=out)


def print_moby_verses() -> None:
    try:
        moby_lines = load_work("moby_dick", cleaning_patterns=NOVEL_CLEANING_PATTERNS)
    except FileNotFoundError:
        moby_lines = fetch_and_parse_verses(
            MOBY_DICK_URL, cleaning_patterns=NOVEL_CLEANING_PATTERNS
        )
    moby_stream = stream_blue_verses(moby_lines, window_size=25)

    start = time.time()
    while time.time() - start < STREAM_TIME:
        print(next(moby_stream))


def main() -> int:
    print_hamlet_verses()
    return 0


if __name__ == "__main__":
    import sys

    try:
        sys.exit(main())
    except BrokenPipeError:
        sys.exit(0)
