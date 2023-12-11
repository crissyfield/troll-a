package cmd

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/spf13/cobra"

	"github.com/crissyfield/troll-a/internal/detect"
	"github.com/crissyfield/troll-a/internal/fetch"
	"github.com/crissyfield/troll-a/internal/warc"
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
	// Open reader for URL
	r, err := fetch.URL(args[0], fetch.WithTimeout(4*time.Hour))
	if err != nil {
		slog.Error("Failed to fetch WARC file", slog.Any("error", err))
		os.Exit(1) //nolint
	}

	defer r.Close()

	// Spawn go routines to check buffers for secrets
	var wg sync.WaitGroup

	bufferCh := make(chan *Buffer)

	for j := 0; j < 8; j++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for buffer := range bufferCh {
				// Detect secrets
				findings, err := detect.Detect(bytes.NewBuffer(buffer.Content))
				if err != nil {
					slog.Error("Unable to read WARC content block", slog.Any("error", err))
					continue
				}

				// Print findings
				for _, f := range findings {
					fmt.Printf(
						"\033[96m%s:%d:%d\033[0m: \033[91m%s\033[0m: \033[93m%s\033[0m\n",
						buffer.TargetURI,
						f.StartLine,
						f.StartColumn,
						f.ID,
						f.Secret,
					)
				}
			}
		}()
	}

	// Traverse WARC file
	err = warc.Traverse(r, func(r *warc.Record) error {
		// Bail if wrong type or payload
		if (r.Type != warc.RecordTypeResponse) || !allowedPayloadTypes[r.IdentifiedPayloadType] {
			return nil
		}

		// Read full record content
		content, err := io.ReadAll(r.Content)
		if err != nil {
			return fmt.Errorf("read record content: %w", err)
		}

		// Hand over to processing
		bufferCh <- &Buffer{TargetURI: r.TargetURI, Content: content}

		return nil
	})

	if err != nil {
		slog.Error("Failed to process WARC file", slog.Any("error", err))
		os.Exit(1) //nolint
	}

	// Clean up
	close(bufferCh)
	wg.Wait()

	slog.Info("Done", slog.String("url", args[0]))
}
