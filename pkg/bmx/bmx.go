package bmx

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"hash/crc32"
	"strings"

	"github.com/tortxof/z85"
)

const (
	Header        = "--- BEGIN THE BARD'S MATRIX ---"
	FooterPattern = "--- END OF THE BARD'S MATRIX [CRC:%08X] ---"
)

func SealText(text string, width int) (string, error) {
	raw := []byte(text)
	crc := crc32.ChecksumIEEE(raw)

	var buf bytes.Buffer
	w, err := zlib.NewWriterLevel(&buf, 9)
	if err != nil {
		return "", fmt.Errorf("bmx: create zlib writer: %w", err)
	}
	if _, err := w.Write(raw); err != nil {
		return "", fmt.Errorf("bmx: compress: %w", err)
	}
	if err := w.Close(); err != nil {
		return "", fmt.Errorf("bmx: close zlib: %w", err)
	}

	encoded := z85.Encode(buf.Bytes())

	var sb strings.Builder
	sb.WriteString(Header)
	sb.WriteByte('\n')

	for i := 0; i < len(encoded); i += width {
		end := i + width
		if end > len(encoded) {
			end = len(encoded)
		}
		sb.Write(encoded[i:end])
		sb.WriteByte('\n')
	}

	sb.WriteString(fmt.Sprintf(FooterPattern, crc))

	return sb.String(), nil
}

func UnsealText(envelope string) (string, error) {
	lines := splitLines(envelope)

	if len(lines) == 0 || lines[0] != Header {
		return "", fmt.Errorf("Invalid matrix: Missing or corrupted header marker.")
	}

	last := lines[len(lines)-1]
	if !strings.HasPrefix(last, "--- END OF THE BARD'S MATRIX") {
		return "", fmt.Errorf("Invalid matrix: Missing or cut-off footer marker.")
	}

	footerLine := lines[len(lines)-1]
	expectedCRC, err := extractCRC(footerLine)
	if err != nil {
		return "", fmt.Errorf("Corrupted matrix: Validation signature format is unreadable.")
	}

	var payload strings.Builder
	for _, line := range lines[1 : len(lines)-1] {
		payload.WriteString(line)
	}

	decoded := z85.Decode([]byte(payload.String()))

	uncompressed, err := zlibDecode(decoded)
	if err != nil {
		return "", fmt.Errorf("Corrupted matrix: %s", err)
	}

	actualCRC := crc32.ChecksumIEEE(uncompressed)
	if actualCRC != expectedCRC {
		return "", fmt.Errorf(
			"Validation Failed! Data corruption detected.\nExpected CRC: %08X vs Actual CRC: %08X",
			expectedCRC, actualCRC,
		)
	}

	return string(uncompressed), nil
}

func splitLines(s string) []string {
	var lines []string
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines
}

func extractCRC(line string) (uint32, error) {
	_, after, ok := strings.Cut(line, "[CRC:")
	if !ok {
		return 0, fmt.Errorf("CRC marker not found")
	}
	hexStr, _, ok := strings.Cut(after, "]")
	if !ok {
		return 0, fmt.Errorf("CRC closing bracket not found")
	}
	var val uint32
	if _, err := fmt.Sscanf(hexStr, "%X", &val); err != nil {
		return 0, fmt.Errorf("invalid CRC hex: %w", err)
	}
	return val, nil
}

func zlibDecode(data []byte) ([]byte, error) {
	r, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var result bytes.Buffer
	if _, err := result.ReadFrom(r); err != nil {
		return nil, err
	}
	return result.Bytes(), nil
}
