package main

import (
	"fmt"
	"restgrpcproxy/internal/config"
	"restgrpcproxy/internal/lib/logger/setuplogger"
	"restgrpcproxy/internal/server"
	"time"

	"golang.org/x/exp/slog"
)

func main() {
	cfg := config.MustLoad()

	log := setuplogger.Setup(cfg.Env)
	log.Info("starting restgrpcproxy", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	shutdownCh := make(chan struct{})
	port := ":" + fmt.Sprintf("%d", cfg.Port)
	server.ListenPort(port, cfg.GRPCServer, cfg.TlsPath, shutdownCh, log, cfg.Timeout)
	server.ListenStopSig()
	close(shutdownCh)
	time.Sleep(cfg.Timeout)

}
