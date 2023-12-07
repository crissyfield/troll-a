package cmd

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// CmdTest defines the CLI sub-command 'test'.
var CmdTest = &cobra.Command{
	Use:   "test [flags]",
	Short: "...",
	Args:  cobra.NoArgs,
	Run:   runTest,
}

// Initialize CLI options.
func init() {
}

// runTest is called when the "test" command is used.
func runTest(_ *cobra.Command, _ []string) {
	// ...
	hc := &http.Client{
		Timeout: 60 * time.Second,
	}

	res, err := hc.Get("https://data.commoncrawl.org/crawl-data/CC-MAIN-2023-40/segments/1695233510326.82/warc/CC-MAIN-20230927203115-20230927233115-00771.warc.gz")
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

	// Buffered IO
	br := bufio.NewReader(gr)

	for {
		// Read version
		version, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}

		if err != nil {
			slog.Error("Error while processing WARC body [header]", slog.Any("error", err))
			os.Exit(1) //nolint
		}

		_ = version

		// Read header
		headers := make(map[string]string)

		for {
			// ...
			data, isPrefix, err := br.ReadLine()
			if err != nil {
				slog.Error("Error while processing WARC body [key-value]", slog.Any("error", err))
				os.Exit(1) //nolint
			}

			// ...
			if isPrefix {
				slog.Error("Error while processing WARC body [prefix]", slog.Any("error", err))
				os.Exit(1) //nolint
			}

			// ...
			header := string(data)

			if header == "" {
				break
			}

			// ...
			parts := strings.SplitN(header, ":", 2)
			if len(parts) == 2 {
				key, value := strings.ToLower(parts[0]), strings.TrimSpace(parts[1])

				headers[key] = value
			}
		}

		// ...
		length, err := strconv.Atoi(headers["content-length"])
		if err != nil {
			slog.Error("Error while processing WARC body [content-length]", slog.Any("error", err))
			os.Exit(1) //nolint
		}

		// ...
		cr := io.LimitReader(br, int64(length))

		if (headers["warc-type"] != "response") || (headers["warc-identified-payload-type"] != "text/html") {
			// ...
			_, _ = io.Copy(io.Discard, cr)
		} else {
			//
			fmt.Println(
				headers["warc-target-uri"],
				headers["warc-type"],
				headers["warc-identified-payload-type"],
			)

			_, _ = io.Copy(io.Discard, cr)
		}

		// ...
		_, _, _ = br.ReadLine()
		_, _, _ = br.ReadLine()
	}
	slog.Info("Done")
}
