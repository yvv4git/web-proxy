package application

import (
	stdLog "log"

	"github.com/davecgh/go-spew/spew"
	"github.com/yvv4git/web-proxy/internal/config"
	"github.com/yvv4git/web-proxy/internal/infra"
	"github.com/yvv4git/web-proxy/internal/webproxy"
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

	wp := webproxy.NewWebProxy()
	if err := wp.Start(); err != nil {
		stdLog.Fatalf("Failed to start web proxy: %v", err)
	}
}
