package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/sync/errgroup"
)

func run(ctx context.Context) error {
	s := &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello %s!", r.URL.Path[1:])
		}),
	}
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
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
	if err := run(context.Background()); err != nil {
		log.Fatalf("Error running server: %v\n", err)
	}
}
