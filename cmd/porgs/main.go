package main

import (
	"context"
	"embed"
	"errors"
	"github.com/praja-dev/porgs"
	"github.com/praja-dev/porgs/task"
	"html/template"
	"io/fs"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"sync"
	"time"
)

func main() {
	porgs.BootConfig = getBootConfig()
	porgs.SiteConfig = getSiteConfig()
	porgs.Plugins = getPlugins()
	porgs.Layout = getLayoutTemplate()
	porgs.Templates = getTemplates()
	porgs.Handler = getHandlers()

	run(context.Background())
}

func getBootConfig() porgs.AppBootConfig {
	host := os.Getenv("HOST")
	portStr := os.Getenv("PORT")
	if portStr == "" {
		portStr = "8642"
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		slog.Error("config: port", "err", err)
		os.Exit(1)
	}

	return porgs.AppBootConfig{
		Host: host,
		Port: port,
	}
}

func getSiteConfig() porgs.AppSiteConfig {
	return porgs.AppSiteConfig{
		Title:       "PORGS",
		Description: "A website powered by Praja Organizations (PORGS)",
	}

}

func getPlugins() map[string]porgs.Plugin {
	plugins := make(map[string]porgs.Plugin)

	corePlugin := &Plugin{}
	plugins[corePlugin.GetName()] = corePlugin

	taskPlugin := &task.Plugin{}
	plugins[taskPlugin.GetName()] = taskPlugin

	return plugins
}

func getLayoutTemplate() *template.Template {
	fm := template.FuncMap{
		"cfg": func() porgs.AppSiteConfig {
			return porgs.SiteConfig
		},
	}

	layout, err := template.New("layout").Funcs(fm).ParseFS(embeddedFS, "layouts/default.go.html")
	if err != nil {
		slog.Error("templates: parse layouts", "err", err)
		os.Exit(1)
	}

	return layout
}

func getTemplates() map[string]*template.Template {
	tm := parseViewTemplates(embeddedFS, porgs.Layout)

	for _, plugin := range porgs.Plugins {
		pluginTemplates := parseViewTemplates(plugin.GetFS(), porgs.Layout)
		for k, v := range pluginTemplates {
			tm[k] = v
		}
	}

	return tm
}

func parseViewTemplates(embedFS embed.FS, layout *template.Template) map[string]*template.Template {
	tm := make(map[string]*template.Template)

	rgxpViewName := regexp.MustCompile(`views/(.+)\.go\.html`)

	viewFiles, err := fs.Glob(embedFS, "views/*.go.html")
	if err != nil {
		slog.Error("parse views", "err", err)
		os.Exit(1)
	}
	for _, viewFile := range viewFiles {
		viewNameMatches := rgxpViewName.FindStringSubmatch(viewFile)
		if viewNameMatches == nil {
			slog.Error("parse view: incorrect file name", "file", viewFile)
			os.Exit(1)
		}
		viewName := viewNameMatches[1]

		tp, err := layout.Clone()
		if err != nil {
			slog.Error("clone layout", "err", err)
			os.Exit(1)
		}
		tp, err = tp.ParseFS(embedFS, viewFile)
		if err != nil {
			slog.Error("parse view", "view", viewName, "err", err)
			os.Exit(1)
		}
		tm[viewName] = tp
	}

	return tm
}

func getHandlers() *http.ServeMux {
	mux := http.NewServeMux()

	for name, plugin := range porgs.Plugins {
		if name == "core" {
			mux.Handle("/", plugin.GetHandler())
		} else {
			mux.Handle("/"+name+"/", http.StripPrefix("/"+name, plugin.GetHandler()))
		}

		mux.Handle("/a/"+name+"/", getAssetHandler(plugin))
	}

	return mux
}

func getAssetHandler(plugin porgs.Plugin) http.Handler {
	assetsDir, err := fs.Sub(plugin.GetFS(), "assets")
	if err != nil {
		slog.Error("handlers: assets", "err", err)
		os.Exit(1)
	}
	assetsHandler := http.StripPrefix("/a/"+plugin.GetName(), http.FileServer(http.FS(assetsDir)))

	return assetsHandler
}

func run(ctx context.Context) {
	server := &http.Server{
		Addr:    net.JoinHostPort(porgs.BootConfig.Host, strconv.Itoa(porgs.BootConfig.Port)),
		Handler: porgs.Handler,
	}
	runServer := func() {
		slog.Info("run: server starting", "host", porgs.BootConfig.Host, "port", porgs.BootConfig.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("run: server failed", "err", err)
		}
	}
	go runServer()

	var wg sync.WaitGroup
	wg.Add(1)
	shutdownGracefully := func() {
		ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
		defer stop()
		defer wg.Done()

		<-ctx.Done()
		slog.Info("run: shutdown starting")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			slog.Error("run: shutdown failed", "err", err)
		} else {
			slog.Info("run: shutdown complete")
		}
	}
	go shutdownGracefully()
	wg.Wait()
}
