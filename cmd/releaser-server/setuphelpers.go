package main

import (
	"context"
	mux2 "github.com/gorilla/mux"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/cresta/zapctx"
	"go.uber.org/zap"
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

func muxWithHealthCheck() *mux2.Router {
	mux := mux2.NewRouter()
	mux.Handle("/healthz", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	return mux
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
