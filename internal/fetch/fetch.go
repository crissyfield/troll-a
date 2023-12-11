package fetch

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	// DefaultTimeout is the default timeout for fetching the URL.
	DefaultTimeout = 60 * time.Second
)

// URL with fetch address addr using the optional options.
func URL(addr string, opts ...Option) (io.ReadCloser, error) {
	// Build config
	c := &config{
		timeout: DefaultTimeout,
	}

	for _, o := range opts {
		o(c)
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
		hc := &http.Client{
			Timeout: c.timeout,
		}

		res, err := hc.Get(u.String()) //nolint // res.Body will be closed by the decompression wrapper!
		if err != nil {
			return nil, fmt.Errorf("HTTP fetch: %w", err)
		}

		r, err := NewWrappedDecompressionReader(res.Body)
		if err != nil {
			return nil, fmt.Errorf("decompression: %w", err)
		}

		return r, nil

	case "s3":
		// Amazon S3
		return nil, fmt.Errorf("schema 's3' not yet supported")

	default:
		// Not supported
		return nil, fmt.Errorf("schema '%s' not supported", u.Scheme)
	}
}
