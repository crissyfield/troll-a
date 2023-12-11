package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"

	"github.com/crissyfield/troll-a/internal/detect"
	"github.com/crissyfield/troll-a/internal/fetch"
	"github.com/crissyfield/troll-a/internal/mime"
	"github.com/crissyfield/troll-a/internal/warc"
)

// CmdTest defines the CLI sub-command 'test'.
var CmdTest = &cobra.Command{
	Use:   "test [flags] [warc url]",
	Short: "...",
	Args:  cobra.ExactArgs(1),
	Run:   runTest,
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

	// Create buffer channel
	type buffer struct {
		TargetURI string
		Content   []byte
	}

	bufferCh := make(chan *buffer)

	// Spawn go routines to check buffers for secrets
	eg, ctx := errgroup.WithContext(context.Background())

	for j := 0; j < 8; j++ {
		eg.Go(func() error {
			for buf := range bufferCh {
				// Detect secrets
				findings, err := detect.Detect(bytes.NewBuffer(buf.Content))
				if err != nil {
					return fmt.Errorf("detect secrets: %w", err)
				}

				// Print findings
				for _, f := range findings {
					fmt.Printf(
						"\033[96m%s:%d:%d\033[0m: \033[91m%s\033[0m: \033[93m%s\033[0m\n",
						buf.TargetURI,
						f.StartLine,
						f.StartColumn,
						f.ID,
						f.Secret,
					)
				}
			}

			return nil
		})
	}

	// Traverse WARC file
	err = warc.Traverse(r, func(r *warc.Record) error {
		select {
		case <-ctx.Done():
			// Break traversal if jobs have stopped
			return warc.ErrBreakTraversal

		default:
			// Bail if wrong type or payload
			if (r.Type != warc.RecordTypeResponse) || !mime.IsText(r.IdentifiedPayloadType) {
				return nil
			}

			// Read full record content
			content, err := io.ReadAll(r.Content)
			if err != nil {
				return fmt.Errorf("read record content: %w", err)
			}

			// Hand over to processing
			bufferCh <- &buffer{
				TargetURI: r.TargetURI,
				Content:   content,
			}
		}

		return nil
	})

	if err != nil {
		slog.Error("Failed to process WARC file", slog.Any("error", err))
		os.Exit(1) //nolint
	}

	// Clean up
	close(bufferCh)

	err = eg.Wait()
	if err != nil {
		slog.Error("Failed to detect secrets", slog.Any("error", err))
		os.Exit(1) //nolint
	}

	slog.Info("Done", slog.String("url", args[0]))
}
