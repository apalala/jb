#!/usr/bin/env python3
# Copyright (c) 2017-2026 Juancarlo Añez (apalala@gmail.com)
# SPDX-License-Identifier: BSD-4-Clause
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
        ascii_payload[i : i + width] 
        for i in range(0, len(ascii_payload), width)
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
    except (IndexError, ValueError):
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
        "-d", "--decompress", 
        action="store_true", 
        help="unseal and decompress the input (auto-triggered if file ends in .bmx)"
    )
    parser.add_argument(
        "-w", "--width", 
        type=int, 
        default=80, 
        help="column width for text wrapping during sealing (default: 80)"
    )
    parser.add_argument(
        "input", 
        nargs="?", 
        help="path to the input file (reads from stdin if omitted)"
    )
    parser.add_argument(
        "output", 
        nargs="?", 
        help="path to the output file (writes to stdout or replaces input if omitted)"
    )
    
    args = parser.parse_args()

    # 1. Read input using strict newline-preservation rules
    input_path = Path(args.input) if args.input else None
    if input_path:
        with open(input_path, mode="r", encoding="utf-8", newline="") as f:
            content = f.read()
    else:
        content = sys.stdin.read()

    if not content.strip():
        return

    # 2. Smart Protocol Mode Selection
    should_decompress = args.decompress or (input_path and input_path.suffix == ".bmx") or content.startswith(HEADER)

    if not should_decompress and input_path and input_path.suffix == ".bmx":
        sys.stderr.write("Error: File is already a .bmx matrix. Aborting to avoid double compression.\n")
        sys.exit(1)

    # 3. Process content with raw byte structures preserved
    try:
        if should_decompress:
            result = unseal_text(content)
        else:
            result = seal_text(content, width=args.width)
    except Exception as err:
        sys.stderr.write(f"Error: {err}\n")
        sys.exit(1)

    # 4. Write out payload keeping text representation and line endings unchanged
    if args.output:
        with open(Path(args.output), mode="w", encoding="utf-8", newline="") as f:
            f.write(result)
    elif input_path:
        if should_decompress:
            out_path = input_path.with_suffix("") if input_path.suffix == ".bmx" else input_path.with_name(f"{input_path.name}.out")
        else:
            out_path = input_path.with_name(f"{input_path.name}.bmx")
            
        with open(out_path, mode="w", encoding="utf-8", newline="") as f:
            f.write(result)
        input_path.unlink()
    else:
        sys.stdout.write(result)
        if not result.endswith("\n"):
            sys.stdout.write("\n")


if __name__ == "__main__":
    main()
