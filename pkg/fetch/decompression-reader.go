package fetch

import (
	"bufio"
	"compress/bzip2"
	"fmt"
	"io"

	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/compress/zstd"
)

const (
	magicGZip               = "\x1f\x8b"         // Magic number of the Gzip format (RFC 1952, section 2.3.1)
	magicBZip2              = "\x42\x5a"         // Magic number of the BZip2 format (no formal spec exists)
	magicZStdFrame          = "\x28\xb5\x2f\xfd" // Magic number of the ZStandard frame format (RFC 8478, section 3.1.1)
	magicZStdSkippableFrame = "\x2a\x4d\x18"     // Magic number of the ZStandard skippable frame format (RFC 8478, section 3.1.2)
)

// NewDecompressionReader will return a new reader transparently doing decompression of GZip, BZip2, and ZStd.
func NewDecompressionReader(r io.ReadCloser) (io.ReadCloser, error) {
	// Read magic bytes
	br := bufio.NewReader(r)

	magic, err := br.Peek(4)
	if err != nil {
		if err == io.EOF {
			return io.NopCloser(br), nil
		}

		return nil, fmt.Errorf("read magic bytes: %w", err)
	}

	switch {
	case string(magic[0:2]) == magicGZip:
		// GZIP decompression
		gzr, err := gzip.NewReader(br)
		if err != nil {
			return nil, fmt.Errorf("create GZIP reader: %w", err)
		}

		return gzr, nil

	case string(magic[0:2]) == magicBZip2:
		// Use BZIP2 decompression
		bzr := bzip2.NewReader(br)

		return io.NopCloser(bzr), nil

	case (string(magic[0:4]) == magicZStdFrame) || ((magic[0]&0xf0 == 0x50) && (string(magic[1:4]) == magicZStdSkippableFrame)):
		// ZStandard decompression
		zsr, err := zstd.NewReader(br)
		if err != nil {
			_ = r.Close()
			return nil, fmt.Errorf("create ZStandard reader: %w", err)
		}

		return zsr.IOReadCloser(), nil

	default:
		// Use no decompression
		return io.NopCloser(br), nil
	}
}
