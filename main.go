package main

import (
	"fmt"
	"net/http"
)

func main() {
	err := http.ListenAndServe(
		":8080",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello %s!", r.URL.Path[1:])
		}),
	)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
