package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"url/urlshort"
)

func main() {
	pathFile := flag.String("file", "", "file with path redirection info")
	flag.Parse()
	if *pathFile == "" {
		fmt.Fprintln(os.Stderr, "Need to set flag, supported format: yaml")
		flag.PrintDefaults()
		os.Exit(1)
	}

	f, err := os.Open(*pathFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	pathInfo, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	var pathsToUrls urlshort.PathMap
	// pathsToUrls := urlshort.PathMap{
	// 	"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
	// 	"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	// }
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	// 	pathInfo := `
	// - path: /urlshort
	//   url: https://github.com/gophercises/urlshort
	// - path: /urlshort-final
	//   url: https://github.com/gophercises/urlshort/tree/solution
	// `
	yamlHandler, err := urlshort.YAMLHandler(pathInfo, mapHandler)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	log.Fatal(http.ListenAndServe(":8080", yamlHandler))
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
