package main

import (
	"context"
	"embed"
	"errors"
	"github.com/eatonphil/gosqlite"
	"github.com/praja-dev/porgs"
	"github.com/praja-dev/porgs/task"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"
)

//go:embed assets/*
//go:embed layouts/*.go.html
//go:embed views/*.go.html
var embeddedFS embed.FS

func main() {
	porgs.BootConfig = getBootConfig()
	porgs.SiteConfig = getSiteConfig()
	porgs.Plugins = getPlugins()
	porgs.Layout = getLayoutTemplate()
	porgs.Templates = getTemplates()
	porgs.Handler = getHandlers()

	initDB()
	run(context.Background())
}

func getBootConfig() porgs.AppBootConfig {
	host := os.Getenv("PORGS_HOST")

	portStr := os.Getenv("PORGS_PORT")
	if portStr == "" {
		portStr = "8642"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		slog.Error("config: port", "err", err)
		os.Exit(1)
	}

	dsn := os.Getenv("PORGS_DSN")
	if dsn == "" {
		dsn = "porgs.db"
	}

	return porgs.AppBootConfig{
		Host: host,
		Port: port,
		DSN:  dsn,
	}
}

func getSiteConfig() porgs.AppSiteConfig {
	return porgs.AppSiteConfig{
		Title:       "Ourville",
		Description: "A website powered by Praja Organizations (PORGS)",
	}

}

func getPlugins() map[string]porgs.Plugin {
	plugins := make(map[string]porgs.Plugin)

	taskPlugin := &task.Plugin{}
	plugins[taskPlugin.GetName()] = taskPlugin

	return plugins
}

func initDB() {
	slog.Info("db: initializing")
	conn, err := gosqlite.Open(porgs.BootConfig.DSN)
	if err != nil {
		slog.Error("db: open", "err", err)
		os.Exit(1)
	}
	defer func() { _ = conn.Close() }()
	conn.BusyTimeout(3 * time.Second)
	slog.Info("db: connected", "dsn", porgs.BootConfig.DSN)

	// # Check if user table exists
	qryUserTbl := "SELECT name FROM sqlite_master WHERE type='table' AND name='user';"
	stmt, err := conn.Prepare(qryUserTbl)
	if err != nil {
		slog.Error("db statement: prepare", "err", err)
		os.Exit(1)
	}
	defer func() { _ = stmt.Close() }()
	hasRow, err := stmt.Step()
	if err != nil {
		slog.Error("db statement: step", "err", err)
		os.Exit(1)
	}
	if hasRow {
		slog.Info("db: ready")
		return
	}

	// # Create user table
	slog.Info("db: preparing for first use")
	qryCreateUserTbl := "CREATE TABLE user (id INTEGER PRIMARY KEY, name TEXT);"
	err = conn.Exec(qryCreateUserTbl)
	if err != nil {
		slog.Error("db: create user table", "err", err)
		os.Exit(1)
	}
	slog.Info("db: ready")
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
