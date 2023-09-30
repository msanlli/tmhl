package main

import (
	"net/http"
	"gopkg.in/yaml.v3"
)

// The MapHandler function maps (duh...) the URLs with their shortened URL
func MapHandler(shortPath map[string]string, fallback http.Handler) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if dest, ok := shortPath[r.URL.Path]; ok {
            http.Redirect(w, r, dest, http.StatusFound)
            return
        }
        fallback.ServeHTTP(w, r)
    }
}


// Setting structure for the yaml version
type pathURL struct {
	Path string `yaml:"path"`
	URL string `yaml:"url"`
}

// Yaml parsing the URL
func parseYaml(yamlBytes []byte) ([]pathURL, error) {
	var pathURLs []pathURL
	err := yaml.Unmarshal(yamlBytes, &pathURLs)
	if err != nil {
		return nil, err
	}
	return pathURLs, nil
}

// Yaml to map
func buildMap(pathURLs []pathURL) map[string]string {
	shortPath := make(map[string]string)
	for _, pu := range pathURLs {
		shortPath[pu.Path] = pu.URL
	}
	return shortPath
}

func YamlHandler(yamlBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err:= parseYaml(yamlBytes)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}
