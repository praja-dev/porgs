package main

import (
	"io/fs"
	"net/http"
)

func getHandler() (http.Handler, error) {
	handler := http.NewServeMux()

	// # Serve assets at /a/
	assetsDir, err := fs.Sub(embeddedFS, "assets")
	if err != nil {
		return nil, err
	}
	assetsHandler := http.StripPrefix("/a/", http.FileServer(http.FS(assetsDir)))
	handler.Handle("/a/", assetsHandler)

	handler.Handle("/{$}", handleRoot())

	return handler, nil
}
