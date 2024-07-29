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
	mux.Handle("GET /{$}", idUser(handleRoot()))
	mux.Handle("GET /login", idUser(handleLoginGet()))
	mux.Handle("POST /login", idUser(handleLoginPost()))
	mux.Handle("GET /logout", idUser(rejectAnon(handleLogout())))
	mux.Handle("GET /home", idUser(rejectAnon(handleHome())))

	for name, plugin := range porgs.Plugins {
		mux.Handle("/a/"+name+"/", getPluginAssetHandler(plugin))
		mux.Handle("/"+name+"/", idUser(rejectAnon(http.StripPrefix("/"+name, plugin.GetHandler()))))
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
