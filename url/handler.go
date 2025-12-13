package main

import (
	"gopkg.in/yaml.v3"
	"net/http"
)

type pathMap map[string]string

type pathUrl struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}

// MapHandler will return a http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(p pathMap, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestPath := r.URL.Path

		if redirectedUrl, exist := p[requestPath]; exist {
			http.Redirect(w, r, redirectedUrl, http.StatusFound)
			return
		}

		fallback.ServeHTTP(w, r)
	}
}

// YAMLHandler will parse the provided YAML and then return
// a http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}
	return MapHandler(parsedYaml, fallback), nil
}

// parseYAML will parse the provided YAML to pathMap
func parseYAML(yml []byte) (pathMap, error) {
	var paths []pathUrl
	err := yaml.Unmarshal(yml, &paths)
	if err != nil {
		return nil, err
	}
	return buildMap(paths), nil
}

// parseYAML will convert the provided []pathUrl to pathMap
func buildMap(paths []pathUrl) pathMap {
	pm := make(pathMap)
	for _, pu := range paths {
		if pu.Path != "" {
			pm[pu.Path] = pu.Url
		}
	}
	return pm
}
