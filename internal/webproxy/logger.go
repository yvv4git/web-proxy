package webproxy

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type WrapedZapLogger struct {
	log *zap.Logger
}

func NewWrapedZapLogger(log *zap.Logger) *WrapedZapLogger {
	return &WrapedZapLogger{log: log}
}

func (w *WrapedZapLogger) Printf(format string, v ...any) {
	w.log.Sugar().Logf(zapcore.InfoLevel, format, v...)
}

func (w *WrapedZapLogger) Write(p []byte) (n int, err error) {
	w.log.Sugar().Debug(string(p))
	return len(p), nil
}
