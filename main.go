package main

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

	"github.com/crissyfield/troll-a/pkg/detect"
	"github.com/crissyfield/troll-a/pkg/detect/preset"
	"github.com/crissyfield/troll-a/pkg/fetch"
	"github.com/crissyfield/troll-a/pkg/mime"
	"github.com/crissyfield/troll-a/pkg/warc"
)

var (
	// Version will be set during build.
	Version = "(unknown)"

	// Settings
	flagLogAsJSON bool
	flagLogLevel  string
)

// main is the main entry point of the command.
func main() {
	// Define command
	var cmd = &cobra.Command{
		Use:               "troll-a [flags] url",
		Short:             "Drill into WARC web archives",
		Args:              cobra.ExactArgs(1),
		Version:           Version,
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		RunE:              runCommand,
	}

	// Settings
	cmd.Flags().StringVar(&flagLogLevel, "log-level", "info", "verbosity of logging output")
	cmd.Flags().BoolVar(&flagLogAsJSON, "log-as-json", false, "change logging format to JSON")

	// Execute
	if err := cmd.Execute(); err != nil {
		// Error has already been printed, just exit
		os.Exit(1)
	}
}

// runCommand is called when the command is used.
func runCommand(_ *cobra.Command, args []string) error {
	// Logging
	var level slog.Level
	if err := level.UnmarshalText([]byte(flagLogLevel)); err != nil {
		return fmt.Errorf("invalid argument \"%s\" for \"--log-level\" flag: %w", flagLogLevel, err)
	}

	var handler slog.Handler
	if flagLogAsJSON {
		handler = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: level})
	} else {
		handler = slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})
	}

	slog.SetDefault(slog.New(handler))

	// Create detector on given rules preset
	detector, err := detect.NewDetector(preset.All)
	if err != nil {
		slog.Error("Failed to create detector", slog.Any("error", err))
		os.Exit(1) //nolint
	}

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
				findings, err := detector.Detect(bytes.NewBuffer(buf.Content))
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

	return nil
}
