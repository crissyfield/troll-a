package warc

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	// RecordTypeRequest is used for requests.
	RecordTypeRequest = "request"

	// RecordTypeResponse is used for responses.
	RecordTypeResponse = "response"

	// RecordTypeMetadata is used for metadata.
	RecordTypeMetadata = "metadata"
)

var (
	// ErrBreakTraversal should be returned from the callback to break traversal.
	ErrBreakTraversal = errors.New("stop traversal")
)

const (
	// bufferSize defines the size of the read buffer.
	bufferSize = 4 * 1024 * 1024

	// Headers
	contentLengthHeader             = "content-length"
	contentTypeHeader               = "content-type"
	warcTypeHeader                  = "warc-type"
	warcIdentifiedPayloadTypeHeader = "warc-identified-payload-type"
	warcTargetURIHeader             = "warc-target-uri"
	httpContentTypeHeader           = "content-type"
)

// Record contains all information about a record.
type Record struct {
	Type                  string    // Type of record ("request", "response", or "metadata")
	TargetURI             string    // Target URI of the record
	IdentifiedPayloadType string    // Identified MIME type of the payload
	HTTPContentType       string    // Content type defined by HTTP header
	Content               io.Reader // Reader for the content
}

// Traverse will traverse the stream via r, calling fn for each record.
func Traverse(r io.Reader, fn func(r *Record) error) error {
	// Buffered IO
	br := bufio.NewReaderSize(r, bufferSize)

	for {
		// Parse WARC header
		warcHeader, err := parseWARCHeader(br)
		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("parse WARC header: %w", err)
		}

		// Extract length of record content
		length, err := strconv.Atoi(warcHeader[contentLengthHeader])
		if err != nil {
			return fmt.Errorf("read record content length: %w", err)
		}

		// Extract HTTP headers
		lr := io.LimitReader(br, int64(length))

		if strings.HasPrefix(warcHeader[contentTypeHeader], "application/http") {
			// We want to read the HTTP header, but also want to pass a reader of the full record (including
			// the HTTP header) into the callback. To achieve this, we create TeeReader tr, which reads from lr
			// but also writes everything that was read into a buffer buf. Then we create MultiReader mr that
			// concatenates whatever is in the buffer (= what we already read from lr) with whatever is left in
			// lr, to re-create the full content reader again.
			buf := &bytes.Buffer{}

			tr := io.TeeReader(lr, buf)
			mr := io.MultiReader(buf, lr)

			// Parse HTTP header
			httpHeader, err := parseHTTPHeader(tr)
			if err != nil {
				return fmt.Errorf("parse HTTP header: %w", err)
			}

			// Call record
			err = fn(&Record{
				Type:                  warcHeader[warcTypeHeader],
				TargetURI:             warcHeader[warcTargetURIHeader],
				IdentifiedPayloadType: warcHeader[warcIdentifiedPayloadTypeHeader],
				HTTPContentType:       httpHeader[httpContentTypeHeader],
				Content:               mr,
			})

			if err != nil {
				if errors.Is(err, ErrBreakTraversal) {
					// Don't report an error if break was requested
					return nil
				}

				return fmt.Errorf("callback: %w", err)
			}
		}

		// Discard remaining record content
		_, err = io.Copy(io.Discard, lr)
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

// parseWARCHeader parses WARC header from incoming stream.
func parseWARCHeader(br *bufio.Reader) (map[string]string, error) {
	// Read and validate version
	version, _, err := br.ReadLine()
	if err == io.EOF {
		return nil, err
	}

	if err != nil {
		return nil, fmt.Errorf("reading record version: %w", err)
	}

	if !strings.HasPrefix(string(version), "WARC/") {
		return nil, fmt.Errorf("unknown record version [version=%s]", string(version))
	}

	// Read warc header
	header := make(map[string]string)

	for {
		// Read header key
		hk, isPrefix, err := br.ReadLine()
		if err != nil {
			return nil, fmt.Errorf("reading record header: %w", err)
		}

		// Exit if the buffer is not big enough (32KiB)
		if isPrefix {
			return nil, fmt.Errorf("record header too big")
		}

		// Stop reading headers on empty line
		if len(hk) == 0 {
			break
		}

		// Split header into key and value
		parts := strings.SplitN(string(hk), ":", 2)

		if len(parts) == 2 {
			key := strings.TrimSpace(strings.ToLower(parts[0]))
			value := strings.TrimSpace(parts[1])

			header[key] = value
		}
	}

	return header, nil
}

// parseHTTPHeader parses HTTP header from incoming stream.
func parseHTTPHeader(r io.Reader) (map[string]string, error) {
	// Go line by line until we hit an empty one
	header := make(map[string]string)

	scanner := bufio.NewScanner(r)
	for i := 0; scanner.Scan(); i++ {
		// Skip the first line
		if i == 0 {
			continue
		}

		// Break on empty line
		line := scanner.Text()
		if line == "" {
			break
		}

		// Split header into key and value
		parts := strings.SplitN(line, ":", 2)

		if len(parts) == 2 {
			key := strings.TrimSpace(strings.ToLower(parts[0]))
			value := strings.TrimSpace(parts[1])

			header[key] = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading content header: %w", err)
	}

	return header, nil
}
