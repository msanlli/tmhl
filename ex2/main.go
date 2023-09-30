package main

import (
	"net/http"
	"fmt"
	"strings"
	"flag"
)

func main() {
	shortPath := map[string]string {
		"/ex": "https://www.example.com/an-example-url",
	}
	
	fallbackHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    	http.NotFound(w, r) // 404 page not found
	})
	
	// Flag for custom url
	cURL := flag.String("url", "", "Provide URL to shorten")
	flag.Parse()

	// Splitting the string and updating URL map
	if *cURL != "" {
		split := strings.Split(*cURL, ":")
		if len(split) == 2 {
			shortPath[split[0]] = split[1]
			//yaml += fmt.Sprintf("\nn- path: %s\n  url: %s", split[0], split[1])
		} else { 
			fmt.Println("Invalid URL")
			return
		}
	}
	handler := MapHandler(shortPath, fallbackHandler)
	http.Handle("/", handler)

	http.ListenAndServe("localhost:8080", nil)
}

// Messages
func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
