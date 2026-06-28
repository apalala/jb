package jb

import (
	"regexp"
	"strings"
	"testing"
)

const gutenbergSample = `Some preamble text.
*** START OF THE PROJECT GUTENBERG EBOOK HAMLET ***
To be, or not to be, that is the question.
Whether 'tis nobler in the mind to suffer.
*** END OF THE PROJECT GUTENBERG EBOOK HAMLET ***
Some trailer text.`

const theatreSample = `HAMLET. To be, or not to be.
HORATIO. My lord!
[Exit Ghost]
The rest is silence.`

const novelSample = `CHAPTER I. Loomings.
Call me Ishmael.
CHAPTER II. The Carpet-Bag.
Some more text.`

func TestParseVersesGutenberg(t *testing.T) {
	verses := ParseVerses(gutenbergSample, nil)
	for _, v := range verses {
		if strings.Contains(v, "START OF THE PROJECT") {
			t.Errorf("verse contains Gutenberg header: %q", v)
		}
		if strings.Contains(v, "END OF THE PROJECT") {
			t.Errorf("verse contains Gutenberg footer: %q", v)
		}
	}
	if len(verses) == 0 {
		t.Fatal("expected verses after Gutenberg cleanup")
	}
}

func TestParseVersesTheatreCleaning(t *testing.T) {
	var pats []*regexp.Regexp
	for _, s := range []string{`(?m)^[A-Z0-9_\s]{2,15}[.:]\s*`, `(?m)[\[(].*?[\])]`} {
		pats = append(pats, regexp.MustCompile(s))
	}
	verses := ParseVerses(theatreSample, pats)

	for _, v := range verses {
		if strings.HasPrefix(v, "HAMLET.") || strings.HasPrefix(v, "HORATIO.") {
			t.Errorf("verse still has speaker tag: %q", v)
		}
		if strings.Contains(v, "[Exit Ghost]") {
			t.Errorf("verse still has stage direction: %q", v)
		}
	}
	if len(verses) == 0 {
		t.Fatal("expected verses after theatre cleaning")
	}
}

func TestParseVersesNovelCleaning(t *testing.T) {
	verses := ParseVerses(novelSample, NovelCleaningPatterns)
	for _, v := range verses {
		if strings.HasPrefix(v, "CHAPTER") {
			t.Errorf("verse still has chapter header: %q", v)
		}
	}
	if len(verses) == 0 {
		t.Fatal("expected verses after novel cleaning")
	}
}

func TestParseVersesEmptyInput(t *testing.T) {
	verses := ParseVerses("", nil)
	if len(verses) != 0 {
		t.Fatalf("expected empty result, got %d verses", len(verses))
	}
}

func TestStreamBlueVersesFromSet(t *testing.T) {
	verses := []string{"a", "b", "c", "d", "e"}
	var count int
	for v := range StreamBlueVerses(verses, 3) {
		found := false
		for _, orig := range verses {
			if v == orig {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("unexpected verse %q not in original set", v)
		}
		count++
		if count >= 50 {
			break
		}
	}
	if count == 0 {
		t.Fatal("expected at least one verse from stream")
	}
}

func TestStreamBlueVersesNoImmediateRepeat(t *testing.T) {
	verses := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	windowSize := 4
	var seen []string
	for v := range StreamBlueVerses(verses, windowSize) {
		for i := max(0, len(seen)-windowSize); i < len(seen); i++ {
			if seen[i] == v {
				t.Fatalf("repeat within window: %q at positions %d and %d", v, i, len(seen))
			}
		}
		seen = append(seen, v)
		if len(seen) >= 100 {
			break
		}
	}
	if len(seen) == 0 {
		t.Fatal("expected verses from stream")
	}
}

func TestStreamBlueVersesEmptyInput(t *testing.T) {
	count := 0
	for range StreamBlueVerses(nil, 5) {
		count++
	}
	if count != 0 {
		t.Fatalf("expected 0 verses from empty input, got %d", count)
	}
}

func TestStreamBlueVersesSingleVerse(t *testing.T) {
	verses := []string{"only"}
	var results []string
	for v := range StreamBlueVerses(verses, 3) {
		results = append(results, v)
		if len(results) >= 5 {
			break
		}
	}
	if len(results) != 5 {
		t.Fatalf("expected 5 verses, got %d", len(results))
	}
	for _, v := range results {
		if v != "only" {
			t.Fatalf("expected 'only', got %q", v)
		}
	}
}
