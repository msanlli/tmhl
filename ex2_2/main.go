package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func main() {
	var input = flag.String("input", "", "URL or path to file")
	flag.Parse()

	// Check if the input flag is provided.
	if *input == "" {
		fmt.Println("Please provide an input using the -input flag.")
		return
	}

	// Fallback handler for URLs that are not found.
	fallback := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "URL not found!")
	})

	// Check if the input is a URL.
	if strings.Contains(*input, "=") {
		parts := strings.SplitN(*input, "=", 2)
		if len(parts) == 2 {
			http.HandleFunc(parts[0], func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, parts[1], http.StatusFound)
			})
			http.ListenAndServe(":8080", nil)
			return
		}
	}

	// Map of handlers for different input types.
	handlersMap := map[string]func(string, http.Handler) (http.HandlerFunc, error){
		".json": JSONHandler,
		".yml":  YAMLHandler,
		".yaml": YAMLHandler,
	}

	var handler http.HandlerFunc
	var err error

	// Find the handler for the input type.
	for prefixSuffix, handlerFunc := range handlersMap {
		if strings.HasSuffix(*input, prefixSuffix) {
			handler, err = handlerFunc(*input, fallback)
			break
		}
	}

	if err != nil {
		fmt.Println("Error setting up handler:", err)
		os.Exit(1)
	}

	if handler == nil {
		fmt.Println("Unknown input type")
		return
	}

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
