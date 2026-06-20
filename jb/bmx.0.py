import base64
import textwrap
import zlib

HEADER = "--- BEGIN THE BARD'S MATRIX ---"
FOOTER_TEMPLATE = "--- END OF THE BARD'S MATRIX [CRC:{:08X}] ---"


def seal_text(text: str, width: int = 80) -> str:
    """Compresses, encodes, and frames a text string with a funny marker and CRC validation."""
    raw_bytes = text.encode("utf-8")

    # 1. Compute a 32-bit checksum of the original source text for signature verification
    crc_checksum = zlib.crc32(raw_bytes)

    # 2. Compress and Base85 encode
    compressed = zlib.compress(raw_bytes, level=9)
    ascii_payload = base64.b85encode(compressed).decode("ascii")

    # 3. Format the lines inside our fancy custom envelope boundaries
    wrapped_lines = textwrap.wrap(ascii_payload, width=width)
    footer = FOOTER_TEMPLATE.format(crc_checksum)

    return "\n".join([HEADER, *wrapped_lines, footer])


def unseal_text(envelope_block: str) -> str:
    """Validates the envelope structure, checks the CRC signature, and decompresses."""
    lines = [line.strip() for line in envelope_block.splitlines() if line.strip()]

    if not lines or lines[0] != HEADER:
        raise ValueError("Invalid matrix: Missing or corrupted header marker.")

    if not lines[-1].startswith("--- END OF THE BARD'S MATRIX"):
        raise ValueError("Invalid matrix: Missing or cut-off footer marker.")

    # Extract the embedded hex validation signature from the footer line
    footer_line = lines[-1]
    try:
        expected_crc = int(footer_line.split("[CRC:")[1].split("]")[0], 16)
    except IndexError, ValueError:
        raise ValueError("Corrupted matrix: Validation signature format is unreadable.")

    # Extract and merge the raw body payload lines
    payload_str = "".join(lines[1:-1])

    # Decode and decompress back to original text form
    compressed_bytes = base64.b85decode(payload_str.encode("ascii"))
    uncompressed_bytes = zlib.decompress(compressed_bytes)

    # Run the validation check!
    actual_crc = zlib.crc32(uncompressed_bytes)
    if actual_crc != expected_crc:
        raise ValueError(
            f"Validation Failed! Data corruption detected.\n"
            f"Expected CRC: {expected_crc:08X} vs Actual CRC: {actual_crc:08X}"
        )

    return uncompressed_bytes.decode("utf-8")
