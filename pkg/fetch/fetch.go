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
)

const (
	// DefaultTimeout is the default timeout for fetching the URL.
	DefaultTimeout = 60 * time.Second
)

// URL with fetch address addr using the optional options.
func URL(addr string, opts ...Option) (io.ReadCloser, error) {
	// Bootstrap settings
	settings := &settings{
		timeout: DefaultTimeout,
	}

	for _, o := range opts {
		o(settings)
	}

	// Parse URL
	u, err := url.Parse(addr)
	if err != nil {
		return nil, fmt.Errorf("parse URL: %w", err)
	}

	// Pick proper code path
	switch u.Scheme {
	case "http", "https":
		// HTTP/HTTPS
		hc := &http.Client{Timeout: settings.timeout}

		res, err := hc.Get(u.String()) //nolint // res.Body will be closed by the decompression wrapper!
		if err != nil {
			return nil, fmt.Errorf("HTTP fetch [url=%s]: %w", u.String(), err)
		}

		r, err := NewWrappedDecompressionReader(res.Body)
		if err != nil {
			return nil, fmt.Errorf("decompression: %w", err)
		}

		return r, nil

	case "s3":
		// Amazon S3
		hc := &http.Client{Timeout: settings.timeout}

		cfg, err := config.LoadDefaultConfig(context.Background(), config.WithHTTPClient(hc))
		if err != nil {
			return nil, fmt.Errorf("load default AWS config: %w", err)
		}

		s3c := s3.NewFromConfig(cfg)

		res, err := s3c.GetObject(context.Background(), &s3.GetObjectInput{
			Bucket: aws.String(u.Host),
			Key:    aws.String(strings.TrimPrefix(u.Path, "/")),
		})

		if err != nil {
			return nil, fmt.Errorf("S3 fetch [url=%s]: %w", u.String(), err)
		}

		r, err := NewWrappedDecompressionReader(res.Body)
		if err != nil {
			return nil, fmt.Errorf("decompression: %w", err)
		}

		return r, nil

	default:
		// Not supported
		return nil, fmt.Errorf("schema '%s' not supported", u.Scheme)
	}
}
