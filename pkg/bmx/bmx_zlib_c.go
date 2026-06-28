//go:build czlib

package bmx

/*
#cgo LDFLAGS: ${SRCDIR}/../../lib/zlib/libz.a

#include "zlib.h"
#include <stdlib.h>

typedef struct {
	unsigned char *data;
	unsigned long  len;
	unsigned long  cap;
} cbuf;

static int c_deflate(const unsigned char *in, unsigned long in_len, cbuf *out) {
	z_stream strm;
	strm.zalloc = Z_NULL;
	strm.zfree  = Z_NULL;
	strm.opaque = Z_NULL;

	int ret = deflateInit2(&strm, 9, Z_DEFLATED, 15, 8, Z_DEFAULT_STRATEGY);
	if (ret != Z_OK) return ret;

	out->cap  = deflateBound(&strm, in_len);
	out->data = malloc(out->cap);
	if (!out->data) { deflateEnd(&strm); return Z_MEM_ERROR; }

	strm.next_in   = (unsigned char *)in;
	strm.avail_in  = in_len;
	strm.next_out  = out->data;
	strm.avail_out = out->cap;

	ret = deflate(&strm, Z_FINISH);
	if (ret != Z_STREAM_END) {
		free(out->data);
		out->data = NULL;
		deflateEnd(&strm);
		return ret == Z_OK ? Z_BUF_ERROR : ret;
	}

	out->len = strm.total_out;
	deflateEnd(&strm);
	return Z_OK;
}

static int c_inflate(const unsigned char *in, unsigned long in_len, cbuf *out) {
	z_stream strm;
	strm.zalloc  = Z_NULL;
	strm.zfree   = Z_NULL;
	strm.opaque  = Z_NULL;
	strm.avail_in  = in_len;
	strm.next_in   = (unsigned char *)in;

	int ret = inflateInit2(&strm, 15);
	if (ret != Z_OK) return ret;

	out->cap  = in_len * 4 + 1024;
	out->data = malloc(out->cap);
	if (!out->data) { inflateEnd(&strm); return Z_MEM_ERROR; }

	strm.avail_out = out->cap;
	strm.next_out  = out->data;

	while (1) {
		ret = inflate(&strm, Z_NO_FLUSH);
		if (ret == Z_STREAM_END) {
			out->len = strm.total_out;
			inflateEnd(&strm);
			return Z_OK;
		}
		if (ret != Z_OK) {
			free(out->data);
			out->data = NULL;
			inflateEnd(&strm);
			return ret;
		}

		out->cap *= 2;
		unsigned char *p = realloc(out->data, out->cap);
		if (!p) {
			free(out->data);
			out->data = NULL;
			inflateEnd(&strm);
			return Z_MEM_ERROR;
		}
		out->data = p;
		strm.avail_out = out->cap - strm.total_out;
		strm.next_out  = out->data + strm.total_out;
	}
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func init() {
	ZlibCompress = cDeflate
	ZlibDecompress = cInflate
}

func cDeflate(raw []byte) ([]byte, error) {
	if len(raw) == 0 {
		return nil, fmt.Errorf("bmx: compress: empty input")
	}

	var buf C.cbuf
	ret := C.c_deflate(ptr(raw), C.ulong(len(raw)), &buf)
	if ret != C.Z_OK {
		return nil, fmt.Errorf("bmx: compress: zlib error %d", int(ret))
	}
	defer C.free(unsafe.Pointer(buf.data))

	return C.GoBytes(unsafe.Pointer(buf.data), C.int(buf.len)), nil
}

func cInflate(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("bmx: decompress: empty input")
	}

	var buf C.cbuf
	ret := C.c_inflate(ptr(data), C.ulong(len(data)), &buf)
	if ret != C.Z_OK {
		return nil, fmt.Errorf("bmx: decompress: zlib error %d", int(ret))
	}
	defer C.free(unsafe.Pointer(buf.data))

	return C.GoBytes(unsafe.Pointer(buf.data), C.int(buf.len)), nil
}

func ptr(b []byte) *C.uchar {
	if len(b) == 0 {
		return nil
	}
	return (*C.uchar)(unsafe.Pointer(&b[0]))
}
