package main

import (
	"context"
	"embed"
	"errors"
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
	"zombiezen.com/go/sqlite/sqlitex"
)

//go:embed assets/*
//go:embed layouts/*.go.html
//go:embed views/*.go.html
//go:embed schema.sql
//go:embed seed.sql
var embeddedFS embed.FS

func main() {
	porgs.BootConfig = getBootConfig()
	porgs.DbConnPool = getDbConnPool()
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

	slog.Info("config: boot config", "host", host, "port", port, "dsn", dsn)

	return porgs.AppBootConfig{
		Host: host,
		Port: port,
		DSN:  dsn,
	}
}

func getDbConnPool() *sqlitex.Pool {
	cpl, err := sqlitex.NewPool(porgs.BootConfig.DSN, sqlitex.PoolOptions{
		PoolSize: -1,
	})
	if err != nil {
		slog.Error("db: conn pool", "err", err)
		os.Exit(1)
	}
	slog.Info("db: conn pool ready")

	return cpl
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
	conn, err := porgs.DbConnPool.Take(context.Background())
	if err != nil {
		slog.Error("init-db: take conn", "err", err)
		os.Exit(1)
	}
	defer porgs.DbConnPool.Put(conn)

	// # Check if user table exists - a negative implies this is a fresh database
	qryUserTbl := "SELECT name FROM sqlite_master WHERE type='table' AND name='user';"
	stmt, _, err := conn.PrepareTransient(qryUserTbl)
	if err != nil {
		slog.Error("init-db: is fresh: prepare", "err", err)
		os.Exit(1)
	}
	defer func() { _ = stmt.Finalize() }()
	hasRow, err := stmt.Step()
	if err != nil {
		slog.Error("init-db: is fresh: step", "err", err)
		os.Exit(1)
	}
	if hasRow {
		slog.Info("init-db: is fresh: no")
		return
	}
	slog.Info("init-db: is fresh: yes")

	// # Run schema.sql and seed.sql scripts
	err = sqlitex.ExecuteScriptFS(conn, embeddedFS, "schema.sql", &sqlitex.ExecOptions{})
	if err != nil {
		slog.Error("init-db: schema", "err", err)
		os.Exit(1)
	}
	slog.Info("init-db: schema created")
	err = sqlitex.ExecuteScriptFS(conn, embeddedFS, "seed.sql", &sqlitex.ExecOptions{})
	if err != nil {
		slog.Error("init-db: seed", "err", err)
		os.Exit(1)
	}
	slog.Info("init-db: seed ok")
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
