package cmd

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/spf13/cobra"
)

var allowedPayloadTypes = map[string]bool{
	"application/atom+xml":      true, // https://www.rfc-editor.org/rfc/rfc5023.html
	"application/json":          true, // https://www.rfc-editor.org/rfc/rfc8259.html
	"application/mbox":          true, // https://www.rfc-editor.org/rfc/rfc4155.html
	"application/msword":        true, // Microsoft Word Document or Document Template
	"application/pgp-signature": true,
	"application/rdf+xml":       true,
	"application/rss+xml":       true,
	"application/rtf":           true,
	"application/vnd.ms-excel":  true,
	"application/x-sh":          true,
	"application/xhtml+xml":     true,
	"application/xml":           true,
	"image/svg+xml":             true,
	"message/rfc822":            true,
	"text/css":                  true,
	"text/csv":                  true,
	"text/html":                 true,
	"text/plain":                true,
	"text/x-chdr":               true,
	"text/x-diff":               true,
	"text/x-log":                true,
	"text/x-perl":               true,
	"text/x-php":                true,
	"text/x-vcard":              true,
}

var githubAccessRegexp = regexp.MustCompile(`github_pat_[0-9a-zA-Z_]{82}|gho_[0-9a-zA-Z]{36}|ghp_[0-9a-zA-Z]{36}|ghr_[0-9a-zA-Z]{36}`)

// Buffer ...
type Buffer struct {
	TargetURI string
	Content   []byte
}

// CmdTest defines the CLI sub-command 'test'.
var CmdTest = &cobra.Command{
	Use:   "test [flags] [warc url]",
	Short: "...",
	Args:  cobra.ExactArgs(1),
	Run:   runTest,
}

// Initialize CLI options.
func init() {
}

// runTest is called when the "test" command is used.
func runTest(_ *cobra.Command, args []string) {
	// ...
	hc := &http.Client{
		Timeout: 60 * 60 * 4 * time.Second,
	}

	res, err := hc.Get(args[0])
	if err != nil {
		slog.Error("Unable to fetch WARC file", slog.Any("error", err))
		os.Exit(1) //nolint
	}

	defer res.Body.Close()

	// Decompress
	gr, err := gzip.NewReader(res.Body)
	if err != nil {
		slog.Error("Unable to decompress WARC body", slog.Any("error", err))
		os.Exit(1) //nolint
	}

	// Spawn go routines to check buffers for secrets
	bufferCh := make(chan *Buffer)
	wg := &sync.WaitGroup{}

	for j := 0; j < 8; j++ {
		wg.Add(1)
		go findSecret(wg, bufferCh)
	}

	// Buffered IO
	br := bufio.NewReaderSize(gr, 4*1024*1024)

	for {
		// Read version
		version, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}

		if err != nil {
			slog.Error("Error while reading WARC record version", slog.Any("error", err))
			os.Exit(1) //nolint
		}

		if !strings.HasPrefix(string(version), "WARC/") {
			slog.Error("Unknown WARC record version", slog.String("version", string(version)))
			os.Exit(1) //nolint
		}

		// Read headers
		headers := make(map[string]string)

		for {
			// Read header
			header, isPrefix, err := br.ReadLine()
			if err != nil {
				slog.Error("Error while processing WARC record header", slog.Any("error", err))
				os.Exit(1) //nolint
			}

			// Exit if the buffer is not big enough (32KiB)
			if isPrefix {
				slog.Error("WARC record header seems too big", slog.Any("error", err))
				os.Exit(1) //nolint
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
		length, err := strconv.Atoi(headers["content-length"])
		if err != nil {
			slog.Error("Unable to read record content length", slog.Any("error", err))
			os.Exit(1) //nolint
		}

		// ...
		cr := io.LimitReader(br, int64(length))

		if (headers["warc-type"] == "response") && allowedPayloadTypes[headers["warc-identified-payload-type"]] {
			// ...
			content, err := io.ReadAll(cr)
			if err != nil {
				slog.Error("Unable to read content block", slog.Any("error", err))
				os.Exit(1) //nolint
			}

			bufferCh <- &Buffer{
				TargetURI: headers["warc-target-uri"],
				Content:   content,
			}
		}

		// Discard remaining content block
		_, err = io.Copy(io.Discard, cr)
		if err != nil {
			slog.Error("Unable to discard remaining content block", slog.Any("error", err))
			os.Exit(1) //nolint
		}

		// Skip two empty lines
		for i := 0; i < 2; i++ {
			boundary, _, err := br.ReadLine()
			if (err != nil) && (err != io.EOF) {
				slog.Error("Unable to read WARC record boundary", slog.Any("error", err))
				os.Exit(1) //nolint
			}

			if len(boundary) != 0 {
				slog.Error("WARC record boundary not empty", slog.Any("error", err))
				os.Exit(1) //nolint
			}
		}
	}

	slog.Info("Done", slog.String("url", args[0]))
}

// ...
func findSecret(wg *sync.WaitGroup, bufferCh chan *Buffer) {
	// ...
	defer wg.Done()

	for buffer := range bufferCh {
		// ...
		indexes := githubAccessRegexp.FindAllIndex(buffer.Content, -1)

		for _, idx := range indexes {
			// ...
			if (idx[0] > 0) && isAlphaNum(buffer.Content[idx[0]-1]) {
				continue
			}

			if (idx[1] < len(buffer.Content)-1) && isAlphaNum(buffer.Content[idx[1]+1]) {
				continue
			}

			// ...
			fmt.Printf(
				"!!! Match [%s:%d]: \033[37m%s\033[1;33m%s\033[37m%s\033[0m\n",
				buffer.TargetURI,
				idx[0],
				cleanUpStrings(string(buffer.Content[max(0, idx[0]-20):idx[0]])),
				string(buffer.Content[idx[0]:idx[1]]),
				cleanUpStrings(string(buffer.Content[idx[1]:min(len(buffer.Content), idx[1]+20)])),
			)
		}
	}
}

// ...
func cleanUpStrings(in string) string {
	return strings.Map(
		func(r rune) rune {
			if unicode.IsPrint(r) {
				return r
			}
			return -1
		},
		in,
	)
}

// ...
func isAlphaNum(in byte) bool {
	return (in >= 48 && in <= 57) || (in >= 65 && in <= 90) || (in >= 97 && in <= 122)
}
