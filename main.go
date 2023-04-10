package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

//go:embed redirect.json
var embeddedRedirectJson []byte

type Redirect struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Permanent   bool   `json:"permanent"`
}

func getRedirects() []Redirect {
	var redirects []Redirect

	err := json.Unmarshal(embeddedRedirectJson, &redirects)

	if err != nil {
		panic(err)
	}

	return redirects
}

func redirectHttpHandler(redirectUrl string, isPermanent bool) func(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	return func(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
		var httpStatusCode = http.StatusTemporaryRedirect

		if isPermanent {
			httpStatusCode = http.StatusPermanentRedirect
		}

		http.Redirect(httpResponseWriter, httpRequest, redirectUrl, httpStatusCode)
	}
}

// See: https://pkg.go.dev/os#example-LookupEnv
func getEnv(key string, fallback string) string {
	value, exists := os.LookupEnv(key)

	if exists {
		return value
	}

	return fallback
}

func main() {
	for _, redirect := range getRedirects() {
		http.HandleFunc(redirect.Source, redirectHttpHandler(redirect.Destination, redirect.Permanent))
	}

	port := getEnv("PORT", "80")

	log.Printf("listening on http://localhost:%s", port)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)

	if err != nil {
		panic(err)
	}
}
