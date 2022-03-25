package main

import (
	"context"
	"github.com/cresta/zapctx"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func MustReturn[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func envWithDefault(s string, defaultVal string) string {
	if ret := os.Getenv(s); ret == "" {
		return defaultVal
	} else {
		return ret
	}
}

func killOnSigTerm(ctx context.Context, logger *zapctx.Logger, httpServer *http.Server) {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		<-sigc
		logger.Info(ctx, "shutting down", zap.String("reason", "signal"))
		if err := httpServer.Shutdown(ctx); err != nil {
			logger.Error(ctx, "failed to shutdown http server", zap.Error(err))
		}
	}()
}
