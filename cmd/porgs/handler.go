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

	mux.Handle("/a/", idLang(getAssetHandler()))
	mux.Handle("GET /{$}", idLang(idUser(handleRoot())))
	mux.Handle("GET /lang/{id}", idUser(handleLang()))
	mux.Handle("GET /login", idLang(idUser(handleLoginGet())))
	mux.Handle("POST /login", idLang(idUser(handleLoginPost())))
	mux.Handle("GET /logout", idLang(idUser(rejectAnon(handleLogout()))))
	mux.Handle("GET /home", idLang(idUser(rejectAnon(handleHome()))))

	for name, plugin := range porgs.Plugins {
		mux.Handle("/a/"+name+"/", getPluginAssetHandler(plugin))
		mux.Handle("/"+name+"/", idLang(idUser(rejectAnon(
			http.StripPrefix("/"+name, plugin.GetHandler())))))
	}

	return mux
}

func getAssetHandler() http.Handler {
	assetsDir, err := fs.Sub(embeddedFS, "assets")
	if err != nil {
		slog.Error("getAssetHandler", "err", err)
		os.Exit(1)
	}
	return http.StripPrefix("/a", http.FileServer(http.FS(assetsDir)))
}

func getPluginAssetHandler(plugin porgs.Plugin) http.Handler {
	assetsDir, err := fs.Sub(plugin.GetFS(), "assets")
	if err != nil {
		slog.Error("getPluginAssetHandler", "plugin", plugin.GetName(), "err", err)
		os.Exit(1)
	}
	return http.StripPrefix("/a/"+plugin.GetName(), http.FileServer(http.FS(assetsDir)))
}
