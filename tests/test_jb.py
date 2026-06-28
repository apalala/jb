from __future__ import annotations

import pytest

from jb.jb import (
    WORKS_DATABASE,
    CLEANING_PATTERNS,
    _extract_id,
    _find_work_file,
    clean_work,
    load_work,
    print_work,
    stream_blue_signal,
    stream_blue_verses,
)

GUTENBERG_SAMPLE = """Some preamble text.
*** START OF THE PROJECT GUTENBERG EBOOK HAMLET ***
To be, or not to be, that is the question.
Whether 'tis nobler in the mind to suffer.
*** END OF THE PROJECT GUTENBERG EBOOK HAMLET ***
Some trailer text."""

THEATRE_SAMPLE = """HAMLET. To be, or not to be.
HORATIO. My lord!
[Exit Ghost]
The rest is silence."""

NOVEL_SAMPLE = """CHAPTER I. Loomings.
Call me Ishmael.
CHAPTER II. The Carpet-Bag.
Some more text."""


def test_clean_work_strips_gutenberg():
    result = clean_work("T", GUTENBERG_SAMPLE)
    assert "START OF THE PROJECT" not in result
    assert "END OF THE PROJECT" not in result
    assert "To be, or not to be" in result
    assert "Some trailer text" not in result


def test_clean_work_theatre_patterns():
    result = clean_work("T", THEATRE_SAMPLE)
    assert "HAMLET." not in result
    assert "HORATIO." not in result
    assert "[Exit Ghost]" not in result
    assert "The rest is silence" in result


def test_clean_work_novel_patterns():
    result = clean_work("N", NOVEL_SAMPLE)
    assert "CHAPTER" not in result
    assert "Call me Ishmael" in result


def test_clean_work_unknown_type():
    result = clean_work("X", THEATRE_SAMPLE)
    assert "HAMLET." in result
    assert "[Exit Ghost]" in result


def test_clean_work_empty_input():
    result = clean_work("T", "")
    assert result == ""


def test_find_work_file_exists():
    f = _find_work_file("pg001524", ".txt.bmx")
    assert f is not None
    assert f.name.endswith(".txt.bmx")


def test_find_work_file_unknown():
    f = _find_work_file("nonexistent", ".txt.bmx")
    assert f is None


def test_extract_id():
    assert _extract_id("pg1524") == "1524"
    assert _extract_id("pg2701") == "2701"
    assert _extract_id("pg84") == "84"


def test_extract_id_no_prefix():
    assert _extract_id("1524") == "1524"


def test_stream_blue_signal_statistics():
    samples = [next(stream_blue_signal(0.65)) for _ in range(10000)]
    mean = sum(samples) / len(samples)
    assert abs(mean) < 0.1


def test_stream_blue_verses_from_set():
    verses = ["a", "b", "c", "d", "e"]
    count = 0
    for v in stream_blue_verses(verses, 3):
        assert v in verses
        count += 1
        if count >= 50:
            break
    assert count > 0


def test_stream_blue_verses_no_immediate_repeat():
    verses = ["a", "b", "c", "d", "e", "f", "g", "h", "i", "j"]
    window_size = 4
    seen: list[str] = []
    for v in stream_blue_verses(verses, window_size):
        recent = seen[-window_size:]
        assert v not in recent
        seen.append(v)
        if len(seen) >= 100:
            break
    assert len(seen) > 0


def test_stream_blue_verses_empty_input():
    with pytest.raises(ValueError, match="cannot be empty"):
        for _ in stream_blue_verses([], 5):
            pass


def test_stream_blue_verses_single_verse():
    results = []
    for v in stream_blue_verses(["only"], 3):
        results.append(v)
        if len(results) >= 5:
            break
    assert len(results) == 5
    assert all(v == "only" for v in results)


def test_load_work_known():
    text = load_work("pg1524")
    assert isinstance(text, str)
    assert len(text) > 0
    assert "HAMLET" in text or "Hamlet" in text


def test_print_work_short_stream():
    verses = "\n".join(f"line {i}" for i in range(20))
    print_work(verses, window_size=5, stream_time=0.05)


def test_works_database_has_entries():
    assert len(WORKS_DATABASE) > 0


def test_works_database_types():
    for w in WORKS_DATABASE:
        assert w.type in ("T", "N", "P")


def test_cleaning_patterns_defined():
    assert "T" in CLEANING_PATTERNS
    assert "N" in CLEANING_PATTERNS
