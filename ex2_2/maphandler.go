package main

import (
	"encoding/json"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"
)

type pathURL struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

func MapHandlerFromMappings(shortPath map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dest, ok := shortPath[r.URL.Path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

func YAMLHandler(filepath string, fallback http.Handler) (http.HandlerFunc, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	parsedYaml := make([]pathURL, 0)
	err = yaml.Unmarshal(data, &parsedYaml)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYaml)
	return MapHandlerFromMappings(pathMap, fallback), nil
}

func JSONHandler(filepath string, fallback http.Handler) (http.HandlerFunc, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var mappings []pathURL
	err = json.Unmarshal(data, &mappings)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(mappings)
	return MapHandlerFromMappings(pathMap, fallback), nil
}

func buildMap(pathURLs []pathURL) map[string]string {
	pathMap := make(map[string]string)
	for _, pu := range pathURLs {
		pathMap[pu.Path] = pu.URL
	}
	return pathMap
}
