package main

import (
	"net/http"
)

func main() {
	yamlData := `
- path: /ex
  url: https://example.com/an-example-url
`

	fallback := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	yamlHandler, err := YamlHandler([]byte(yamlData), fallback)
	if err != nil {
		panic(err)
	}
	http.Handle("/", yamlHandler)
	http.ListenAndServe(":8080", nil)
}
