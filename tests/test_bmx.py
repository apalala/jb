from __future__ import annotations

from importlib.resources import files

import pytest

from jb.bmx import seal_text, unseal_text, HEADER


def test_seal_unseal_roundtrip():
    input_text = "Hello, this is a test of the BMX roundtrip."
    sealed = seal_text(input_text, 80)
    result = unseal_text(sealed)
    assert result == input_text
    assert sealed.startswith(HEADER)
    assert sealed.strip().endswith("---")


def test_seal_width_wrapping():
    input_text = "x" * 1000
    sealed = seal_text(input_text, 40)
    lines = sealed.strip().split("\n")
    for i, line in enumerate(lines[1:-1]):
        assert len(line) <= 40, f"line {i + 1} exceeds width 40: len={len(line)}"


def test_crc_integrity():
    input_text = "Data integrity check test."
    sealed = seal_text(input_text, 80)
    corrupted = sealed.replace("[CRC:", "[CRC:DEADBEEF", 1)
    with pytest.raises(ValueError, match="Validation Failed"):
        unseal_text(corrupted)


def test_bad_header():
    with pytest.raises(ValueError, match="Missing or corrupted header"):
        unseal_text("garbage input\n")


def test_empty_envelope():
    with pytest.raises(ValueError, match="Missing or corrupted header"):
        unseal_text("")


def _test_decode(filename: str, expected: str) -> None:
    path = files("jb") / "works" / filename
    data = path.read_text()
    result = unseal_text(data)
    assert expected in result


def test_decode_hamlet():
    _test_decode("pg001524-T15-Hamlet.txt.bmx", "THE TRAGEDY OF HAMLET")


def test_decode_moby_dick():
    _test_decode("pg002701-N25-Moby_Dick.txt.bmx", "MOBY-DICK")


def test_unicode_roundtrip():
    input_text = "Hello, 世界! ñoño — émoji 🎉 test ✓"
    sealed = seal_text(input_text, 80)
    result = unseal_text(sealed)
    assert result == input_text


def test_empty_string_roundtrip():
    sealed = seal_text("", 80)
    result = unseal_text(sealed)
    assert result == ""


def test_multi_line_input():
    input_text = "Line one\nLine two\n\nLine four after blank\n  indented  "
    sealed = seal_text(input_text, 80)
    result = unseal_text(sealed)
    assert result == input_text


def test_trailing_whitespace_in_envelope():
    input_text = "Whitespace tolerance test."
    sealed = seal_text(input_text, 80)
    with_whitespace = "\n\n\n" + sealed + "\n\n\n"
    result = unseal_text(with_whitespace)
    assert result == input_text
