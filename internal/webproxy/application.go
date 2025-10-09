package webproxy

import (
	"context"
	stdLog "log"
	"os"
	"os/signal"
	"syscall"

	"github.com/davecgh/go-spew/spew"
	"github.com/yvv4git/web-proxy/internal/config"
	"github.com/yvv4git/web-proxy/internal/infra"
)

func RunWebProxy(configPath string) {
	var cfg config.Config

	if err := config.Load(configPath, &cfg); err != nil {
		stdLog.Fatalf("Failed to load config: %v", err)
	}

	spew.Dump(cfg)

	log, err := infra.NewWithLogLevel(cfg.LogLevel)
	if err != nil {
		stdLog.Fatalf("Failed to create logger: %v", err)
	}
	defer log.Sync()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	wp := NewWebProxy(log)
	if err := wp.Start(ctx); err != nil {
		stdLog.Fatalf("Failed to start web proxy: %v", err)
	}
}
