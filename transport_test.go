package libfossil

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPTransportRoundTrip(t *testing.T) {
	echo := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ct := r.Header.Get("Content-Type"); ct != "application/x-fossil" {
			t.Errorf("Content-Type = %q, want application/x-fossil", ct)
		}
		body, _ := io.ReadAll(r.Body)
		w.Write(body)
	}))
	defer echo.Close()

	tr := NewHTTPTransport(echo.URL)
	resp, err := tr.RoundTrip(context.Background(), []byte("hello"))
	if err != nil {
		t.Fatal(err)
	}
	if string(resp) != "hello" {
		t.Errorf("got %q, want %q", resp, "hello")
	}
}

func TestMockTransportRoundTrip(t *testing.T) {
	tr := &MockTransport{
		Handler: func(req []byte) []byte { return []byte("mock") },
	}
	resp, err := tr.RoundTrip(context.Background(), []byte("req"))
	if err != nil {
		t.Fatal(err)
	}
	if string(resp) != "mock" {
		t.Errorf("got %q, want %q", resp, "mock")
	}
}

func TestTransportFunc(t *testing.T) {
	called := false
	fn := TransportFunc(func(ctx context.Context, payload []byte) ([]byte, error) {
		called = true
		return []byte("response"), nil
	})
	resp, err := fn.RoundTrip(context.Background(), []byte("request"))
	if err != nil {
		t.Fatalf("RoundTrip: %v", err)
	}
	if !called {
		t.Error("handler function was not called")
	}
	if string(resp) != "response" {
		t.Errorf("got %q, want %q", resp, "response")
	}
}
