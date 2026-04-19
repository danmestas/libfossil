package libfossil

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
)

// Transport delivers sync payloads between peers.
// Implementations handle the network layer (HTTP, NATS, etc.).
// Payloads are opaque zlib-compressed xfer card streams.
type Transport interface {
	RoundTrip(ctx context.Context, payload []byte) ([]byte, error)
}

// NewHTTPTransport creates a Transport that speaks Fossil's HTTP /xfer protocol.
func NewHTTPTransport(url string, opts ...HTTPOption) Transport {
	t := &httpTransport{url: url, client: http.DefaultClient}
	for _, o := range opts {
		o(t)
	}
	return t
}

// HTTPOption configures an HTTP transport.
type HTTPOption func(*httpTransport)

// WithHTTPClient sets a custom http.Client for the transport.
func WithHTTPClient(c *http.Client) HTTPOption {
	return func(t *httpTransport) { t.client = c }
}

type httpTransport struct {
	url    string
	client *http.Client
}

func (t *httpTransport) RoundTrip(ctx context.Context, payload []byte) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", t.url, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("libfossil: http transport: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-fossil")
	resp, err := t.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("libfossil: http transport: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("libfossil: http transport read: %w", err)
	}
	return body, nil
}

// MockTransport is a test double that delegates to a handler function.
type MockTransport struct {
	Handler func(req []byte) []byte
}

// RoundTrip calls the Handler function if set, otherwise returns empty bytes.
func (t *MockTransport) RoundTrip(_ context.Context, payload []byte) ([]byte, error) {
	if t.Handler == nil {
		return []byte{}, nil
	}
	return t.Handler(payload), nil
}

// TransportFunc adapts a plain function to the Transport interface.
// This is the Transport equivalent of http.HandlerFunc.
type TransportFunc func(ctx context.Context, payload []byte) ([]byte, error)

// RoundTrip calls the function.
func (f TransportFunc) RoundTrip(ctx context.Context, payload []byte) ([]byte, error) {
	return f(ctx, payload)
}
