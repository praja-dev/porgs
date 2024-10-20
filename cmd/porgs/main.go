package main

import (
	"context"
	"embed"
	"errors"
	"github.com/praja-dev/porgs"
	"github.com/praja-dev/porgs/core"
	"github.com/praja-dev/porgs/task"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
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
	porgs.Context = context.Background()
	porgs.Args = parseArgs()
	porgs.BootConfig = getBootConfig()
	porgs.DbConnPool = getDbConnPool()
	porgs.SiteConfig = getSiteConfig()
	porgs.Plugins = getPlugins()
	porgs.Layout = getLayoutTemplate()
	porgs.Templates = getTemplates()
	porgs.Handler = getHandlers()

	initDB()
	initPlugins()
	run(porgs.Context)
}

func parseArgs() map[string]string {
	args := make(map[string]string)
	for _, arg := range os.Args[1:] {
		parts := strings.Split(arg, "=")
		if len(parts) != 2 {
			slog.Error("parseArgs", "arg", arg, "err", "invalid argument")
			os.Exit(1)
		}
		name := parts[0]
		value := parts[1]

		// # Handle ~ and %USERPROFILE% in --load arg value
		if name == "--load" {
			if strings.HasPrefix(value, "~") || strings.HasPrefix(value, "%USERPROFILE%") {
				home, err := os.UserHomeDir()
				if err != nil {
					slog.Error("parseArgs", "arg", "--load", "err", err)
					os.Exit(1)
				}

				if strings.HasPrefix(value, "~") {
					value = filepath.Join(home, value[1:])
				} else {
					value = filepath.Join(home, value[13:])
				}
			}
		}

		args[name] = value
	}
	slog.Info("parseArgs: ok")

	return args
}

func getBootConfig() porgs.AppBootConfig {
	host := os.Getenv("PORGS_HOST")
	if host == "" {
		slog.Info("getBootConfig: host", "default", "", "msg", "PORGS_HOST not set or is empty")
	}

	portStr := os.Getenv("PORGS_PORT")
	if portStr == "" {
		slog.Info("getBootConfig: port", "default", "8642", "msg", "PORGS_PORT not set")
		portStr = "8642"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		slog.Error("getBootConfig: port", "err", err)
		os.Exit(1)
	}

	dsn := os.Getenv("PORGS_DSN")
	if dsn == "" {
		slog.Info("getBootConfig: dsn", "default", "porgs.db", "msg", "PORGS_DSN not set")
		dsn = "porgs.db"
	}

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
		slog.Error("getDbConnPool", "err", err)
		os.Exit(1)
	}
	slog.Info("getDbConnPool: ok")

	return cpl
}

func getSiteConfig() porgs.AppSiteConfig {
	return porgs.AppSiteConfig{
		Title:         "Ourville",
		Description:   "A website powered by Praja Organizations (PORGS)",
		LangSupported: []string{"en", "si", "ta"},
		LangDefault:   "en",
		Text:          text,
	}

}

func getPlugins() map[string]porgs.Plugin {
	plugins := make(map[string]porgs.Plugin)

	corePlugin := &core.Plugin{}
	plugins[corePlugin.GetName()] = corePlugin

	taskPlugin := &task.Plugin{}
	plugins[taskPlugin.GetName()] = taskPlugin

	return plugins
}

func initDB() {
	conn, err := porgs.DbConnPool.Take(context.Background())
	if err != nil {
		slog.Error("main.initDB: connect", "err", err)
		os.Exit(1)
	}
	slog.Info("main.initDB: connect: ok")
	defer porgs.DbConnPool.Put(conn)

	// # Check if user table exists - a negative implies this is a fresh database
	qryUserTbl := "SELECT name FROM sqlite_master WHERE type='table' AND name='user';"
	stmt, _, err := conn.PrepareTransient(qryUserTbl)
	if err != nil {
		slog.Error("main.initDB: check if db is fresh: stmt prepare", "err", err)
		os.Exit(1)
	}
	defer func() { _ = stmt.Finalize() }()
	hasRow, err := stmt.Step()
	if err != nil {
		slog.Error("main.initDB: check if db is fresh: stmt step", "err", err)
		os.Exit(1)
	}
	if hasRow {
		slog.Info("main.initDB: check if db is fresh: ok", "msg", "not fresh")
		return
	}
	slog.Info("main.initDB: check if db is fresh: ok", "msg", "fresh")

	// # Run schema.sql and seed.sql scripts in main
	err = sqlitex.ExecuteScriptFS(conn, embeddedFS, "schema.sql", &sqlitex.ExecOptions{})
	if err != nil {
		slog.Error("main.initDB: exec schema.sql in main", "err", err)
		os.Exit(1)
	}
	slog.Info("main.initDB: exec schema.sql in main: ok")

	err = sqlitex.ExecuteScriptFS(conn, embeddedFS, "seed.sql", &sqlitex.ExecOptions{})
	if err != nil {
		slog.Error("main.initDB: exec seed.sql in main", "err", err)
		os.Exit(1)
	}
	slog.Info("main.initDB: exec seed.sql in main: ok")

	// # Run schema.sql and seed.sql scripts in plugins
	for _, plugin := range porgs.Plugins {
		err = sqlitex.ExecuteScriptFS(conn, plugin.GetFS(), "schema.sql", &sqlitex.ExecOptions{})
		if err != nil {
			slog.Error("main.initDB: exec schema.sql in plugin", "plugin", plugin.GetName(), "err", err)
			os.Exit(2)
		}
		slog.Info("main.initDB: exec schema.sql in plugin: ok", "plugin", plugin.GetName())

		err = sqlitex.ExecuteScriptFS(conn, plugin.GetFS(), "seed.sql", &sqlitex.ExecOptions{})
		if err != nil {
			slog.Error("main.initDB: exec seed.sql in plugin", "plugin", plugin.GetName(), "err", err)
			os.Exit(2)
		}
		slog.Info("main.initDB: exec seed.sql in plugin: ok", "plugin", plugin.GetName())
	}
}

func initPlugins() {
	for _, plugin := range porgs.Plugins {
		err := plugin.GetInit()
		if err != nil {
			slog.Error("initPlugins", "plugin", plugin.GetName(), "err", err)
			os.Exit(2)
		}
		slog.Info("initPlugins: ok", "plugin", plugin.GetName())
	}
}

func run(ctx context.Context) {
	server := &http.Server{
		Addr:    net.JoinHostPort(porgs.BootConfig.Host, strconv.Itoa(porgs.BootConfig.Port)),
		Handler: porgs.Handler,
	}
	runServer := func() {
		slog.Info("run: listen and serve", "host", porgs.BootConfig.Host, "port", porgs.BootConfig.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("run: listen and serve", "err", err)
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
		slog.Info("run: shutdown")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			slog.Error("run: shutdown", "err", err)
		} else {
			slog.Info("run: shutdown: ok")
		}
	}
	go shutdownGracefully()
	wg.Wait()
}
