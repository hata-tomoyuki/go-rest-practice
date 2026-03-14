package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"golang.org/x/sync/errgroup"
)

func run(ctx context.Context, l net.Listener) error {
	s := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello %s!", r.URL.Path[1:])
		}),
	}
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		if err := s.Serve(l); err != nil && err != http.ErrServerClosed {
			log.Printf("Error starting server: %v\n", err)
			return err
		}
		return nil
	})
	<-ctx.Done()
	if err := s.Shutdown(context.Background()); err != nil {
		log.Printf("Error shutting down server: %v\n", err)
		return err
	}
	return eg.Wait()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <address>\n", os.Args[0])
		os.Exit(1)
	}
	p := os.Args[1]
	l, err := net.Listen("tcp", p)
	if err != nil {
		log.Fatalf("Error listening on %s: %v\n", p, err)
	}
	if err := run(context.Background(), l); err != nil {
		log.Fatalf("Error running server: %v\n", err)
		os.Exit(1)
	}
}
