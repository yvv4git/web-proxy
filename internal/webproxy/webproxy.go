package webproxy

import (
	"context"
	"net/http"
	"time"

	"github.com/elazarl/goproxy"
	"go.uber.org/zap"
)

type Option func(*WebProxy)

func WithAddr(addr string) Option {
	return func(wp *WebProxy) {
		wp.webSrv.Addr = addr
	}
}

func WithShutdownTimeout(timeout time.Duration) Option {
	return func(wp *WebProxy) {
		wp.shutdownTimeout = timeout
	}
}

type WebProxy struct {
	log             *zap.Logger
	shutdownTimeout time.Duration
	webSrv          *http.Server
}

func NewWebProxy(log *zap.Logger, opts ...Option) *WebProxy {
	const (
		defaultAddr = "127.0.0.1:8080"
	)

	proxy := goproxy.NewProxyHttpServer()

	// Create entity with default options
	entity := &WebProxy{}
	entity.log = log
	entity.webSrv = &http.Server{
		Addr:    defaultAddr,
		Handler: proxy,
	}

	// Update entity with custom options
	for _, opt := range opts {
		opt(entity)
	}

	// For all standard requests
	proxy.OnRequest().Do(goproxy.FuncReqHandler(entity.authMiddleware))

	// For connect requests in https
	proxy.OnRequest().HandleConnect(goproxy.FuncHttpsHandler(func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
		_, resp := entity.authMiddleware(ctx.Req, ctx)
		if resp != nil {
			ctx.Resp = resp
			return goproxy.RejectConnect, host
		}

		return goproxy.OkConnect, host
	}))

	return entity
}

func (wp *WebProxy) Start(ctx context.Context) error {
	go func() {
		wp.log.Info("starting web proxy", zap.String("address", wp.webSrv.Addr))
		if err := wp.webSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			wp.log.Error("web proxy failed", zap.Error(err))
		}
	}()

	<-ctx.Done()

	wp.log.Info("shutting down web proxy")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), wp.shutdownTimeout)
	defer cancel()

	if err := wp.webSrv.Shutdown(shutdownCtx); err != nil {
		wp.log.Error("failed to gracefully shutdown web proxy", zap.Error(err))
		return err
	}

	wp.log.Info("web proxy stopped successfully")

	return nil
}
