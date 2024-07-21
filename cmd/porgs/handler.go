package main

import (
	"github.com/praja-dev/porgs"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
)

func getHandlers() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/a/", getAssetHandler())
	mux.Handle("GET /{$}", handleRoot())
	mux.Handle("GET /login", handleLoginGet())
	mux.Handle("POST /login", handleLoginPost())
	mux.Handle("GET /logout", handleLogout())
	mux.Handle("GET /home", handleHome())

	for name, plugin := range porgs.Plugins {
		mux.Handle("/a/"+name+"/", getPluginAssetHandler(plugin))
		mux.Handle("/"+name+"/", http.StripPrefix("/"+name, plugin.GetHandler()))
	}

	return mux
}

func getAssetHandler() http.Handler {
	assetsDir, err := fs.Sub(embeddedFS, "assets")
	if err != nil {
		slog.Error("handlers: assets", "err", err)
		os.Exit(1)
	}
	return http.StripPrefix("/a", http.FileServer(http.FS(assetsDir)))
}

func getPluginAssetHandler(plugin porgs.Plugin) http.Handler {
	assetsDir, err := fs.Sub(plugin.GetFS(), "assets")
	if err != nil {
		slog.Error("handlers: plugin assets", "plugin", plugin.GetName(), "err", err)
		os.Exit(1)
	}
	return http.StripPrefix("/a/"+plugin.GetName(), http.FileServer(http.FS(assetsDir)))
}
