package main

import (
	"context"
	"errors"
	"github.com/praja-dev/porgs"
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
	porgs.Templates = getTemplates()
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

func getTemplates() map[string]*template.Template {
	tm := make(map[string]*template.Template)

	// # Parse the default layout
	layout, err := template.ParseFS(embeddedFS, "layouts/default.go.html")
	if err != nil {
		slog.Error("templates: parse layouts", "err", err)
		os.Exit(1)
	}

	rgxpViewName := regexp.MustCompile(`views/(.+)\.go\.html`)

	// # Parse all views in the main package
	viewFiles, err := fs.Glob(embeddedFS, "views/*.go.html")
	if err != nil {
		slog.Error("templates: parse views", "err", err)
		os.Exit(1)
	}
	for _, viewFile := range viewFiles {
		viewNameMatches := rgxpViewName.FindStringSubmatch(viewFile)
		if viewNameMatches == nil {
			slog.Error("templates: parse view: incorrect file name", "file", viewFile)
			os.Exit(1)
		}
		viewName := viewNameMatches[1]

		tp, err := layout.Clone()
		if err != nil {
			slog.Error("templates: clone layout", "err", err)
			os.Exit(1)
		}
		tp, err = tp.ParseFS(embeddedFS, viewFile)
		if err != nil {
			slog.Error("templates: parse view", "view", viewName, "err", err)
			os.Exit(1)
		}
		tm[viewName] = tp
	}

	return tm
}

func run(ctx context.Context) {
	handler, err := getHandler()
	if err != nil {
		slog.Error("run: getting handler", "err", err)
		os.Exit(1)
	}
	server := &http.Server{
		Addr:    net.JoinHostPort(porgs.BootConfig.Host, strconv.Itoa(porgs.BootConfig.Port)),
		Handler: handler,
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
