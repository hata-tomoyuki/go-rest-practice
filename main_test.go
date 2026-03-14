package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"

	"golang.org/x/sync/errgroup"
)

func TestRun(t *testing.T) {
	l, err := net.Listen("tcp", "localhost:0")
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return run(ctx, l)
	})
	in := "message"
	url := fmt.Sprintf("http://%s/%s", l.Addr().String(), in)

	t.Logf("Making GET request to http://%s/%s", l.Addr(), in)
	rsp, err := http.Get(url)

	if err != nil {
		t.Error("Error making GET request:", err)
	}
	defer rsp.Body.Close()
	got, err := io.ReadAll(rsp.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %v", err)
	}

	want := fmt.Sprintf("Hello %s!", in)
	if string(got) != want {
		t.Errorf("Expected response %q, got %q", want, string(got))
	}
	cancel()
	if err := eg.Wait(); err != nil {
		t.Fatalf("Error waiting for server to shut down: %v", err)
	}
}
