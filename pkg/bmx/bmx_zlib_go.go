//go:build !czlib

package bmx

import (
	"bytes"
	"compress/zlib"
	"fmt"
)

func init() {
	ZlibCompress = goDeflate
	ZlibDecompress = goInflate
}

func goDeflate(raw []byte) ([]byte, error) {
	var buf bytes.Buffer
	w, err := zlib.NewWriterLevel(&buf, 9)
	if err != nil {
		return nil, fmt.Errorf("bmx: create zlib writer: %w", err)
	}
	if _, err := w.Write(raw); err != nil {
		return nil, fmt.Errorf("bmx: compress: %w", err)
	}
	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("bmx: close zlib: %w", err)
	}
	return buf.Bytes(), nil
}

func goInflate(data []byte) ([]byte, error) {
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
