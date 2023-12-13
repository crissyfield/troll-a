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

	"github.com/crissyfield/troll-a/internal/cli"
)

var (
	// Version will be set during build.
	Version = "(unknown)"

	// Configuration
	configLogLevel    = cli.LogLevel{Val: slog.LevelInfo}
	configJSON        = false
	configJobs        = uint(8)
	configBackoff     = cli.BackoffStrategy{Val: cli.BackoffStrategyValNone}
	configRulesPreset = cli.RulesPreset{Val: preset.Secret}
	configEnclosed    = false
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
		Run:               runCommand,
	}

	// Settings
	cmd.Flags().VarP(&configLogLevel, "verbosity", "V", `verbosity of logging output, allowed: "debug", "info", "warn", "error"`)
	cmd.Flags().BoolVarP(&configJSON, "json", "s", configJSON, `change output format to JSON`)
	cmd.Flags().UintVarP(&configJobs, "jobs", "j", configJobs, `number of concurrent jobs to detect secrets`)
	cmd.Flags().VarP(&configBackoff, "backoff", "b", `backoff strategy for fetching, allowed: "none", "constant", "exponential", "zero"`)
	cmd.Flags().VarP(&configRulesPreset, "preset", "p", `rules preset to use, allowed: "all", "most", "secret"`)
	cmd.Flags().BoolVarP(&configEnclosed, "enclosed", "e", configEnclosed, `only report secrets that are enclosed`)

	// Execute
	if err := cmd.Execute(); err != nil {
		// Error has already been printed, just exit
		os.Exit(1)
	}
}

// runCommand is called when the command is used.
func runCommand(_ *cobra.Command, args []string) {
	// Logging
	var slogErr, slogOut *slog.Logger

	if configJSON {
		slogErr = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: configLogLevel.Val}))
		slogOut = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	} else {
		slogErr = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: configLogLevel.Val}))
		slogOut = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}

	// Create detector on given rules preset
	detector := detect.NewDetector(configRulesPreset.Val, configEnclosed)

	// Open reader for URL
	r, err := fetch.URL(
		args[0],
		fetch.WithTimeout(4*time.Hour),
		fetch.WithBackoff(configBackoff.Val),
	)

	if err != nil {
		slogErr.Error("Failed to fetch WARC file", slog.Any("error", err))
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

	for j := uint(0); j < configJobs; j++ {
		eg.Go(func() error {
			for buf := range bufferCh {
				// Detect secrets
				findings, err := detector.Detect(bytes.NewBuffer(buf.Content))
				if err != nil {
					return fmt.Errorf("detect secrets: %w", err)
				}

				// Print findings
				for _, f := range findings {
					slogOut.Info(
						"Matched",
						slog.String("uri", buf.TargetURI),
						slog.Int("line", f.Location.StartLine),
						slog.Int("column", f.Location.StartColumn),
						slog.String("rule", f.ID),
						slog.String("secret", f.Secret),
						slog.String("full", f.Location.Line(string(buf.Content))),
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
		slogErr.Error("Failed to process WARC file", slog.Any("error", err))
		os.Exit(1) //nolint
	}

	// Clean up
	close(bufferCh)

	err = eg.Wait()
	if err != nil {
		slogErr.Error("Failed to detect secrets", slog.Any("error", err))
		os.Exit(1) //nolint
	}

	slogErr.Info("Done", slog.String("url", args[0]))
}
