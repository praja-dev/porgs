package main

import "net/http"

func getHandler() (http.Handler, error) {
	handler := http.NewServeMux()
	handler.Handle("/{$}", handleRoot())
	return handler, nil
}
