package bmx

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSealUnsealRoundtrip(t *testing.T) {
	input := "Hello, this is a test of the BMX roundtrip."
	sealed, err := SealText(input, 80)
	if err != nil {
		t.Fatalf("SealText: %v", err)
	}
	result, err := UnsealText(sealed)
	if err != nil {
		t.Fatalf("UnsealText: %v", err)
	}

	if result != input {
		t.Fatalf("roundtrip mismatch:\nwant: %q\ngot:  %q", input, result)
	}

	if !strings.HasPrefix(sealed, Header) {
		t.Fatal("sealed text missing header")
	}
	if !strings.HasSuffix(strings.TrimSpace(sealed), "---") {
		t.Fatal("sealed text missing footer")
	}
}

func TestSealWidthWrapping(t *testing.T) {
	input := strings.Repeat("x", 1000)

	sealed, err := SealText(input, 40)
	if err != nil {
		t.Fatalf("SealText: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(sealed), "\n")
	for i, line := range lines[1 : len(lines)-1] {
		if len(line) > 40 {
			t.Fatalf("line %d exceeds width 40: len=%d", i+1, len(line))
		}
	}
}

func TestCRCIntegrity(t *testing.T) {
	input := "Data integrity check test."
	sealed, err := SealText(input, 80)
	if err != nil {
		t.Fatalf("SealText: %v", err)
	}

	lines := strings.SplitN(sealed, "\n", -1)
	if len(lines) >= 2 {
		payload := []byte(lines[1])
		if len(payload) > 1 {
			payload[0] ^= 0xFF
			lines[1] = string(payload)
		}
	}
	corrupted := strings.Join(lines, "\n")

	_, err = UnsealText(corrupted)
	if err == nil {
		t.Fatal("expected error for corrupted data, got nil")
	}
}

func TestBadHeader(t *testing.T) {
	_, err := UnsealText("garbage input\n")
	if err == nil {
		t.Fatal("expected error for bad header, got nil")
	}
}

func TestEmptyEnvelope(t *testing.T) {
	_, err := UnsealText("")
	if err == nil {
		t.Fatal("expected error for empty envelope, got nil")
	}
}

func projectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

func testFile(t *testing.T, path, expected string) {
	t.Helper()

	full := filepath.Join(projectRoot(), path)
	data, err := os.ReadFile(full)
	if err != nil {
		t.Fatalf("read %s: %v", full, err)
	}

	result, err := UnsealText(string(data))
	if err != nil {
		t.Fatalf("UnsealText %s: %v", path, err)
	}

	if !strings.Contains(result, expected) {
		t.Fatalf("decoded %s does not contain expected text %q", path, expected)
	}
}

func TestDecodeHamlet(t *testing.T) {
	testFile(t, "works/pg1524.txt.bmx", "THE TRAGEDY OF HAMLET")
}

func TestDecodeMobyDick(t *testing.T) {
	testFile(t, "works/pg2701.txt.bmx", "MOBY-DICK")
}
