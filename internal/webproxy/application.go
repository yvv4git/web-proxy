package webproxy

import (
	"context"
	stdLog "log"
	"os"
	"os/signal"
	"syscall"

	"github.com/yvv4git/web-proxy/internal/config"
	"github.com/yvv4git/web-proxy/internal/infra"
	"go.uber.org/zap"
)

func RunWebProxy(configPath string) {
	var cfg config.Config

	if err := config.Load(configPath, &cfg); err != nil {
		stdLog.Fatalf("failed to load config: %v", err)
	}

	log, err := infra.NewWithLogLevel(cfg.LogLevel)
	if err != nil {
		stdLog.Fatalf("failed to create logger: %v", err)
	}
	defer log.Sync()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	webProxy := NewWebProxy(log)

	for _, v := range cfg.Auth.PredifinedAuth.Accounts {
		webProxy.authManager.AddAccount(v.Username, v.Password)
	}

	if err := webProxy.Start(ctx); err != nil {
		log.Fatal("failed to start web proxy: %v", zap.Error(err))
	}
}
