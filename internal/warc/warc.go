package warc

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	// bufferSize defines the size of the read buffer.
	bufferSize = 4 * 1024 * 1024

	// Headers
	contentLengthHeader             = "content-length"
	warcTypeHeader                  = "warc-type"
	warcIdentifiedPayloadTypeHeader = "warc-identified-payload-type"
	warchTargetURIHeader            = "warc-target-uri"
)

const (
	// RecordTypeRequest is used for requests.
	RecordTypeRequest = "request"

	// RecordTypeResponse is used for responses.
	RecordTypeResponse = "response"

	// RecordTypeMetadata is used for metadata.
	RecordTypeMetadata = "metadata"
)

// Record contains all information about a record.
type Record struct {
	Type                  string    // Type of record ("request", "response", or "metadata")
	TargetURI             string    // Target URI of the record
	IdentifiedPayloadType string    // Identified MIME type of the payload
	Content               io.Reader // Reader for the content
}

// RecordCallback is called for each record during traversal.
type RecordCallback func(r *Record) error

// Traverse will traverse the stream via r, calling fn for each record.
func Traverse(r io.Reader, fn RecordCallback) error {
	// Buffered IO
	br := bufio.NewReaderSize(r, bufferSize)

	for {
		// Read and validate version
		version, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("reading record version: %w", err)
		}

		if !strings.HasPrefix(string(version), "WARC/") {
			return fmt.Errorf("unknown record version [version=%s]", string(version))
		}

		// Read headers
		headers := make(map[string]string)

		for {
			// Read header
			header, isPrefix, err := br.ReadLine()
			if err != nil {
				return fmt.Errorf("reading record header: %w", err)
			}

			// Exit if the buffer is not big enough (32KiB)
			if isPrefix {
				return fmt.Errorf("record header too big")
			}

			// Stop reading headers on empty line
			if len(header) == 0 {
				break
			}

			// Split header into key and value
			parts := strings.SplitN(string(header), ":", 2)

			if len(parts) == 2 {
				key := strings.ToLower(parts[0])
				value := strings.TrimSpace(parts[1])

				headers[key] = value
			}
		}

		// Extract length of record content
		length, err := strconv.Atoi(headers[contentLengthHeader])
		if err != nil {
			return fmt.Errorf("read record content length: %w", err)
		}

		// Call record
		cr := io.LimitReader(br, int64(length))

		err = fn(&Record{
			Type:                  headers[warcTypeHeader],
			TargetURI:             headers[warchTargetURIHeader],
			IdentifiedPayloadType: headers[warcIdentifiedPayloadTypeHeader],
			Content:               cr,
		})

		if err != nil {
			return fmt.Errorf("call back: %w", err)
		}

		// Discard remaining record content
		_, err = io.Copy(io.Discard, cr)
		if err != nil {
			return fmt.Errorf("discard remaining record content: %w", err)
		}

		// Skip two empty lines
		for i := 0; i < 2; i++ {
			boundary, _, err := br.ReadLine()
			if (err != nil) && (err != io.EOF) {
				return fmt.Errorf("read record boundary: %w", err)
			}

			if len(boundary) != 0 {
				return fmt.Errorf("non-empty record boundary [boundary: %s]", boundary)
			}
		}
	}

	return nil
}
