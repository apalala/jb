package bmx

import (
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

	compressed, err := ZlibCompress(raw)
	if err != nil {
		return "", err
	}

	encoded := z85.Encode(compressed)

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
	newlineIdx := strings.IndexByte(envelope, '\n')
	if newlineIdx < 0 {
		return "", fmt.Errorf("Invalid matrix: Missing or corrupted header marker.")
	}

	if envelope[:newlineIdx] != Header {
		return "", fmt.Errorf("Invalid matrix: Missing or corrupted header marker.")
	}

	rest := envelope[newlineIdx+1:]

	const footerPrefix = "--- END OF THE BARD'S MATRIX"
	footerIdx := strings.LastIndex(rest, footerPrefix)
	if footerIdx < 0 {
		return "", fmt.Errorf("Invalid matrix: Missing or cut-off footer marker.")
	}

	payloadStr := rest[:footerIdx]
	payloadStr = strings.ReplaceAll(payloadStr, "\n", "")

	footerLine := rest[footerIdx:]
	expectedCRC, err := extractCRC(footerLine)
	if err != nil {
		return "", fmt.Errorf("Corrupted matrix: Validation signature format is unreadable.")
	}

	decoded := z85.Decode([]byte(payloadStr))

	uncompressed, err := ZlibDecompress(decoded)
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


