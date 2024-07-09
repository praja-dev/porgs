package main

import "net/http"

func handleRoot() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		// Below HTML is correct according to: https://validator.w3.org/
		_, _ = w.Write([]byte(`<!DOCTYPE HTML><html lang="en"><head><title>PORGS</title><p>PORGS`))
	})
}
