package main

import (
	"github.com/praja-dev/porgs"
	"log/slog"
	"os"
	"strconv"
)

func main() {
	porgs.BootConfig = getBootConfig()
	slog.Info("porgs", "host", porgs.BootConfig.Host,
		"port", porgs.BootConfig.Port)
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
