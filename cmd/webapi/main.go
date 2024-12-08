package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Hello, %q"+name, r.URL.Path)

}
func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error", err)
		os.Exit(1)
	}
}

func run() error {
	http.HandleFunc("/", Hello)
	log.Println("Starting server...")
	return http.ListenAndServe(":8080", nil)
}
