package fetch

import (
	"bufio"
	"compress/bzip2"
	"compress/gzip"
	"fmt"
	"io"
)

// decompressionReader will transparently decompress a data stream.
type decompressionReader struct {
	o io.ReadCloser // Original reader-closer
	w io.ReadCloser // Wrapped reader-closer
}

// NewWrappedDecompressionReader will return a new reader transparently doing decompression of GZIP or BZIP2.
func NewWrappedDecompressionReader(r io.ReadCloser) (io.ReadCloser, error) {
	// Read magic bytes
	br := bufio.NewReader(r)

	magic, err := br.Peek(2)
	if err != nil {
		if err == io.EOF {
			return &decompressionReader{o: r, w: io.NopCloser(br)}, nil
		}

		_ = r.Close()
		return nil, fmt.Errorf("read magic bytes: %w", err)
	}

	switch {
	case (magic[0] == 0x1f) && (magic[1] == 0x8b):
		// Use GZIP decompression
		gzr, err := gzip.NewReader(br)
		if err != nil {
			_ = r.Close()
			return nil, fmt.Errorf("create GZIP reader: %w", err)
		}

		return &decompressionReader{o: r, w: gzr}, nil

	case (magic[0] == 0x42) && (magic[1] == 0x5a):
		// Use BZIP2 decompression
		bzr := bzip2.NewReader(br)

		return &decompressionReader{o: r, w: io.NopCloser(bzr)}, nil

	default:
		// Use no decompression
		return &decompressionReader{o: r, w: io.NopCloser(br)}, nil
	}
}

// Read reads up to len(b) bytes from the stream and stores them in b. It returns the number of bytes read and
// any error encountered. At end of file, Read returns 0, io.EOF.
func (r *decompressionReader) Read(b []byte) (int, error) {
	// Forward to wrapped reader
	return r.w.Read(b)
}

// Close closes the reader. Close will return an error if it has already been called.
func (r *decompressionReader) Close() error {
	// Close wrapped reader
	_ = r.w.Close()

	// Close original reader
	return r.o.Close()
}
