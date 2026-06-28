#!/usr/bin/env python3
# Copyright (c) 2017-2026 Juancarlo Añez (apalala@gmail.com)
# SPDX-License-Identifier: MIT
from __future__ import annotations

import argparse
import base64
import sys
import zlib
from pathlib import Path

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
        help="unseal and decompress the input (auto-triggered if file ends in .bmx)",
    )
    parser.add_argument(
        "-w",
        "--width",
        type=int,
        default=80,
        help="column width for text wrapping during sealing (default: 80)",
    )
    parser.add_argument(
        "inputs", nargs="*", help="paths to input files (reads from stdin if omitted)"
    )

    args = parser.parse_args()

    if not args.inputs:
        content = sys.stdin.read()
        if not content.strip():
            return

        should_decompress = args.decompress or content.startswith(HEADER)
        try:
            result = (
                unseal_text(content)
                if should_decompress
                else seal_text(content, width=args.width)
            )
        except Exception as err:
            sys.stderr.write(f"Error: {err}\n")
            sys.exit(1)

        sys.stdout.write(result)
        if not result.endswith("\n"):
            sys.stdout.write("\n")
        return

    for input_name in args.inputs:
        input_path = Path(input_name)
        with open(input_path, mode="r", encoding="utf-8", newline="") as f:
            content = f.read()

        if not content.strip():
            continue

        should_decompress = (
            args.decompress
            or input_path.suffix == ".bmx"
            or content.startswith(HEADER)
        )

        if not should_decompress and input_path.suffix == ".bmx":
            sys.stderr.write(
                "Error: File is already a .bmx matrix. Aborting to avoid double compression.\n"
            )
            sys.exit(1)

        try:
            result = (
                unseal_text(content)
                if should_decompress
                else seal_text(content, width=args.width)
            )
        except Exception as err:
            sys.stderr.write(f"Error: {err}\n")
            sys.exit(1)

        if should_decompress:
            out_path = (
                input_path.with_suffix("")
                if input_path.suffix == ".bmx"
                else input_path.with_name(f"{input_path.name}.out")
            )
        else:
            out_path = input_path.with_name(f"{input_path.name}.bmx")

        with open(out_path, mode="w", encoding="utf-8", newline="") as f:
            f.write(result)
        input_path.unlink()


if __name__ == "__main__":
    main()
