package main

import (
	"net/http"
)

func (p *Plugin) GetHandler() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("GET /{$}", handleRoot())
	mux.Handle("GET /login", handleLoginGet())
	mux.Handle("POST /login", handleLoginPost())

	return mux
}
