package webproxy

import (
	"context"
	"net/http"
	"time"

	"github.com/elazarl/goproxy"
	"go.uber.org/zap"
)

type WebProxy struct {
	log    *zap.Logger
	webSrv *http.Server
}

func NewWebProxy(log *zap.Logger) *WebProxy {
	const (
		defaultAddr = "127.0.0.1:8080"
	)

	proxy := goproxy.NewProxyHttpServer()

	webSrv := &http.Server{
		Addr:    defaultAddr,
		Handler: proxy,
	}

	wp := &WebProxy{
		log:    log,
		webSrv: webSrv,
	}

	// For all standard requests
	proxy.OnRequest().Do(goproxy.FuncReqHandler(wp.authMiddleware))

	// For connect requests in https
	proxy.OnRequest().HandleConnect(goproxy.FuncHttpsHandler(func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
		_, resp := wp.authMiddleware(ctx.Req, ctx)
		if resp != nil {
			ctx.Resp = resp
			return goproxy.RejectConnect, host
		}

		return goproxy.OkConnect, host
	}))

	return wp
}

func (wp *WebProxy) Start(ctx context.Context) error {
	const shutdownTimeout = 5 * time.Second // TODO: вынести в конфиг

	go func() {
		wp.log.Info("starting web proxy", zap.String("address", wp.webSrv.Addr))
		if err := wp.webSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			wp.log.Error("web proxy failed", zap.Error(err))
		}
	}()

	<-ctx.Done()

	wp.log.Info("shutting down web proxy")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := wp.webSrv.Shutdown(shutdownCtx); err != nil {
		wp.log.Error("failed to gracefully shutdown web proxy", zap.Error(err))
		return err
	}

	wp.log.Info("web proxy stopped successfully")

	return nil
}
