package fetch

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/cenkalti/backoff/v4"
)

var (
	// DefaultTimeout is the default timeout for fetching the URL.
	DefaultTimeout = 60 * time.Second

	// DefaultBackOff is the default backoff strategy for fetching the URL.
	DefaultBackOff = &backoff.StopBackOff{}
)

// Open will fetch address addr using the given options.
func Open(addr string, opts ...Option) (io.ReadCloser, error) {
	// Bootstrap params
	params := &params{
		timeout: DefaultTimeout,
		backOff: DefaultBackOff,
	}

	for _, o := range opts {
		o(params)
	}

	// Special case: STDIN
	if addr == "-" {
		return io.NopCloser(os.Stdin), nil
	}

	// Parse URL
	u, err := url.Parse(addr)
	if err != nil {
		return nil, fmt.Errorf("parse URL: %w", err)
	}

	// Pick proper fetch strategy
	var rc io.ReadCloser

	err = backoff.Retry(
		func() error {
			switch u.Scheme {
			case "http", "https":
				// HTTP/HTTPS
				return openHTTPURL(u, params, &rc)

			case "s3":
				// Amazon S3
				return openS3URL(u, params, &rc)

			case "file", "":
				// File URL
				return openFileURL(u, params, &rc)

			default:
				// Unknown schema
				return backoff.Permanent(fmt.Errorf("schema not supported"))
			}
		},
		params.backOff,
	)

	if err != nil {
		return nil, err
	}

	return rc, nil
}

// openHTTPURL returns a reader for the given HTTP/HTTPS URL.
func openHTTPURL(u *url.URL, params *params, rc *io.ReadCloser) error {
	// HTTP/HTTPS
	hc := &http.Client{Timeout: params.timeout}

	res, err := hc.Get(u.String()) //nolint // res.Body will be closed by the decompression wrapper!
	if err != nil {
		return fmt.Errorf("HTTP fetch [url=%s]: %w", u.String(), err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected HTTP status: %d", res.StatusCode)
	}

	*rc = res.Body
	return nil
}

// openS3URL returns a reader for the given Amazon S3 URL.
func openS3URL(u *url.URL, params *params, rc *io.ReadCloser) error {
	hc := &http.Client{Timeout: params.timeout}

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithHTTPClient(hc))
	if err != nil {
		return backoff.Permanent(fmt.Errorf("load default AWS config: %w", err))
	}

	s3c := s3.NewFromConfig(cfg)

	res, err := s3c.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(u.Host),
		Key:    aws.String(strings.TrimPrefix(u.Path, "/")),
	})

	if err != nil {
		return fmt.Errorf("S3 fetch [url=%s]: %w", u.String(), err)
	}

	*rc = res.Body
	return nil
}

// openFileURL returns a reader for the given file URL.
func openFileURL(u *url.URL, _ *params, rc *io.ReadCloser) error {
	// Get path from URL
	path, err := pathFromURL(u)
	if err != nil {
		return fmt.Errorf("file fetch [url=%s]: %w", u.String(), err)
	}

	// Open file
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("file open [url=%s]: %w", u.String(), err)
	}

	*rc = f
	return nil
}

// Get the file path from the URL.
func pathFromURL(u *url.URL) (string, error) {
	// Special case windows
	if runtime.GOOS == "windows" {
		return pathFromURLWindows(u)
	}

	// For all others, only an empty or 'localhost' hosts are allowed
	if (u.Host != "") && (u.Host != "localhost") {
		return "", errors.New("file URL specifies non-local host")
	}

	return filepath.FromSlash(u.Path), nil
}

// pathFromURLWindows converts the URL u into a windows filesystem path.
//
// Copied from Golang's src/cmd/go/internal/web/url_windows.go
func pathFromURLWindows(u *url.URL) (string, error) {
	// Bail if path is not absolute
	if (len(u.Path) == 0) || (u.Path[0] != '/') {
		return "", errors.New("path is not absolute")
	}

	// Convert to native slashes
	path := filepath.FromSlash(u.Path)

	// We interpret Windows file URLs per the description in
	// https://blogs.msdn.microsoft.com/ie/2006/12/06/file-uris-in-windows/.
	// The host part of a file URL (if any) is the UNC volume name, but RFC 8089 reserves the authority
	// "localhost" for the local machine.
	if (u.Host != "") && (u.Host != "localhost") {
		// A common "legacy" format omits the leading slash before a drive letter, encoding the drive letter as
		// the host instead of part of the path. (See
		// https://blogs.msdn.microsoft.com/freeassociations/2005/05/19/the-bizarre-and-unhappy-story-of-file-urls/.)
		// We do not support that format, but we should at least emit a more helpful error message for it.
		if filepath.VolumeName(u.Host) != "" {
			return "", errors.New("file URL encodes volume in host field: too few slashes?")
		}

		return `\\` + u.Host + path, nil
	}

	// If host is empty, path must contain an initial slash followed by a drive letter and path. Remove the
	// slash and verify that the path is valid.
	if vol := filepath.VolumeName(path[1:]); vol == "" || strings.HasPrefix(vol, `\\`) {
		return "", errors.New("file URL missing drive letter")
	}

	return path[1:], nil
}
