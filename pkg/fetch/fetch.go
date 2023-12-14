package fetch

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
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

// URL will fetch address addr using the optional options.
func URL(addr string, opts ...func(*state)) (io.ReadCloser, error) {
	// Bootstrap state
	state := &state{
		timeout: DefaultTimeout,
		backOff: DefaultBackOff,
	}

	for _, o := range opts {
		o(state)
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
				hc := &http.Client{Timeout: state.timeout}

				res, err := hc.Get(u.String()) //nolint // res.Body will be closed by the decompression wrapper!
				if err != nil {
					return fmt.Errorf("HTTP fetch [url=%s]: %w", u.String(), err)
				}

				if res.StatusCode != http.StatusOK {
					return fmt.Errorf("unexpected HTTP status: %d", res.StatusCode)
				}

				rc = res.Body
				return nil

			case "s3":
				// Amazon S3
				hc := &http.Client{Timeout: state.timeout}

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

				rc = res.Body
				return nil

			default:
				// Not supported
				return backoff.Permanent(fmt.Errorf("schema not supported"))
			}
		},
		state.backOff,
	)

	if err != nil {
		return nil, err
	}

	// Decompress, if necessary
	dr, err := NewWrappedDecompressionReader(rc)
	if err != nil {
		return nil, fmt.Errorf("decompression: %w", err)
	}

	return dr, nil
}
