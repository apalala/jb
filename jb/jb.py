# Copyright (c) 2017-2026 Juancarlo Añez (apalala@gmail.com)
# SPDX-License-Identifier: BSD-4-Clause
from __future__ import annotations

import random
import re
import urllib.request
from typing import Iterator

# Project Gutenberg Stable UTF-8 Raw Text URLs
HAMLET_URL = "https://www.gutenberg.org/cache/epub/1524/pg1524.txt"
MOBY_DICK_URL = "https://www.gutenberg.org/cache/epub/2701/pg2701.txt"


def fetch_and_parse_verses(url: str) -> list[str]:
    """Downloads a text file from Gutenberg and parses it into cleaned lines."""
    print(f"Fetching source text from {url}...")
    with urllib.request.urlopen(url) as response:
        raw_text = response.read().decode("utf-8")

    # Split into lines and strip trailing/leading whitespace tokens
    raw_lines = [line.strip() for line in raw_text.splitlines()]

    # Clean out empty lines, Project Gutenberg header indicators, or pure structural markup
    # (A simple line filter looking for lines with alphanumeric content)
    lines = [line for line in raw_lines if line and re.search(r"\w", line)]

    # Slice out metadata blocks if known, or return the active body
    # Gutenberg files typically contain indicators like '*** START OF THE PROJECT...'
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

    return lines[start_idx:end_idx]


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
    current_idx = total_lines // 2

    while True:
        # 1. Pull next raw blue noise step (-3.0 to +3.0 scale roughly)
        step = next(noise_stream)

        # 2. Translate the noise value into a bounded structural jump index
        # Window size dictates how large an instantaneous line jump can be
        jump = int(step * window_size)

        # 3. Apply the shift and bounce smoothly off borders
        current_idx += jump
        if current_idx < 0:
            current_idx = abs(current_idx) % total_lines
        elif current_idx >= total_lines:
            current_idx = total_lines - (current_idx % total_lines) - 1

        yield verses[current_idx]


if __name__ == "__main__":
    # --- Run Hamlet Mode ---
    hamlet_lines = fetch_and_parse_verses(HAMLET_URL)
    print(f"Loaded {len(hamlet_lines)} lines into memory footprint.\n")

    hamlet_stream = stream_blue_verses(hamlet_lines, window_size=15)

    print("--- Streaming Blue Noise Hamlet Lines ---")
    for _ in range(5):
        print(f"> {next(hamlet_stream)}")

    print("\n" + "=" * 40 + "\n")

    # --- Run Moby Dick Mode ---
    moby_lines = fetch_and_parse_verses(MOBY_DICK_URL)
    print(f"Loaded {len(moby_lines)} lines into memory footprint.\n")

    moby_stream = stream_blue_verses(moby_lines, window_size=25)

    print("--- Streaming Blue Noise Moby Dick Lines ---")
    for _ in range(5):
        print(f"> {next(moby_stream)}")
