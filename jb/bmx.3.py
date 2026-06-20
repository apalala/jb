#!/usr/bin/env python3
# Copyright (c) 2017-2026 Juancarlo Añez (apalala@gmail.com)
# SPDX-License-Identifier: MIT
from __future__ import annotations

import argparse
import base64
import sys
from pathlib import Path
import zlib

HEADER = "--- BEGIN THE BARD'S MATRIX ---"
FOOTER_TEMPLATE = "--- END OF THE BARD'S MATRIX [CRC:{:08X}] ---"


def seal_text(text: str, width: int = 80) -> str:
    """Compresses, Z85-encodes, and frames text with a marker and CRC validation."""
    raw_bytes = text.encode("utf-8")
    crc_checksum = zlib.crc32(raw_bytes)

    compressed = zlib.compress(raw_bytes, level=9)
    ascii_payload = base64.z85encode(compressed).decode("ascii")

    wrapped_lines = [
        ascii_payload[i : i + width] for i in range(0, len(ascii_payload), width)
    ]

    footer = FOOTER_TEMPLATE.format(crc_checksum)
    return "\n".join([HEADER, *wrapped_lines, footer])


def unseal_text(envelope_block: str) -> str:
    """Validates the envelope structure, checks the CRC signature, and decompresses."""
    lines = [line.strip() for line in envelope_block.splitlines() if line.strip()]

    if not lines or lines[0] != HEADER:
        raise ValueError("Invalid matrix: Missing or corrupted header marker.")

    if not lines[-1].startswith("--- END OF THE BARD'S MATRIX"):
        raise ValueError("Invalid matrix: Missing or cut-off footer marker.")

    footer_line = lines[-1]
    try:
        expected_crc = int(footer_line.split("[CRC:")[1].split("]")[0], 16)
    except IndexError, ValueError:
        raise ValueError("Corrupted matrix: Validation signature format is unreadable.")

    payload_str = "".join(lines[1:-1])

    compressed_bytes = base64.z85decode(payload_str.encode("ascii"))
    uncompressed_bytes = zlib.decompress(compressed_bytes)

    actual_crc = zlib.crc32(uncompressed_bytes)
    if actual_crc != expected_crc:
        raise ValueError(
            f"Validation Failed! Data corruption detected.\n"
            f"Expected CRC: {expected_crc:08X} vs Actual CRC: {actual_crc:08X}"
        )

    return uncompressed_bytes.decode("utf-8")


def main() -> None:
    parser = argparse.ArgumentParser(
        description="Pack or unpack text files using a compressed Z85 validation envelope."
    )
    parser.add_argument(
        "-d",
        "--decompress",
        action="store_true",
        help="unseal and decompress the input instead of sealing it",
    )
    parser.add_argument(
        "-w",
        "--width",
        type=int,
        default=80,
        help="column width for text wrapping during sealing (default: 80)",
    )
    parser.add_argument(
        "input", nargs="?", help="path to the input file (reads from stdin if omitted)"
    )
    parser.add_argument(
        "output",
        nargs="?",
        help="path to the output file (writes to stdout or replaces input if omitted)",
    )

    args = parser.parse_args()

    # 1. Read input from file or stream
    input_path = Path(args.input) if args.input else None
    if input_path:
        content = input_path.read_text(encoding="utf-8")
    else:
        content = sys.stdin.read()

    if not content.strip():
        return

    # 2. Process the content through the encoder or decoder
    try:
        if args.decompress:
            result = unseal_text(content)
        else:
            result = seal_text(content, width=args.width)
    except Exception as err:
        sys.stderr.write(f"Error: {err}\n")
        sys.exit(1)

    # 3. Determine output target matching default gzip semantics
    if args.output:
        # Explicit output file given
        Path(args.output).write_text(result, encoding="utf-8")
    elif input_path:
        # Gzip behavior: replace the original file
        if args.decompress:
            # Drop the .bmx suffix if present
            out_path = (
                input_path.with_suffix("")
                if input_path.suffix == ".bmx"
                else input_path.with_name(f"{input_path.name}.out")
            )
        else:
            # Append our .bmx suffix
            out_path = input_path.with_name(f"{input_path.name}.bmx")

        out_path.write_text(result, encoding="utf-8")
        input_path.unlink()  # Safely clean up original source file
    else:
        # Standard fallback to stdout stream
        sys.stdout.write(result)
        if not result.endswith("\n"):
            sys.stdout.write("\n")


if __name__ == "__main__":
    main()
