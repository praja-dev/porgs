package main

import (
	"net/http"
)

func (p *Plugin) GetHandler() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("GET /{$}", handleRoot())

	return mux
}
