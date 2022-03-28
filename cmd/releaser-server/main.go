package main

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/cresta/cresta-releaser/internal/logging"
	"github.com/cresta/cresta-releaser/internal/managedgitrepo"
	"github.com/cresta/cresta-releaser/internal/releaserserver"
	"github.com/cresta/cresta-releaser/releaser"
	releaser_protobuf "github.com/cresta/cresta-releaser/rpc/releaser"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()
	logger := MustReturn(logging.SetupLogging(envWithDefault("LOG_LEVEL", "info")))
	logger.Info(ctx, "Starting application")
	api := MustReturn(releaser.NewFromCommandLine(ctx, logger.Unwrap(ctx), nil))
	repo := MustReturn(managedgitrepo.NewRepo(ctx, envWithDefault("REPO_DISK_LOCATION", "/tmp/repo"), os.Getenv("REPO_URL"), api.Fs, api.Github, api.Git))
	serverImpl := MustReturn(releaserserver.NewServer(ctx, logger, api, repo))
	twirpServer := releaser_protobuf.NewReleaserServer(serverImpl)
	Must(os.Chdir(envWithDefault("REPO_DISK_LOCATION", "/tmp/repo")))
	mux := muxWithHealthCheck()
	mux.Handle("/", twirpServer)
	httpServer := http.Server{
		Addr:    envWithDefault("LISTEN_ADDR", ":8080"),
		Handler: mux,
	}
	killOnSigTerm(ctx, logger, &httpServer)
	logger.Info(ctx, "starting server", zap.String("addr", httpServer.Addr))
	err := httpServer.ListenAndServe()
	logger.Info(ctx, "server stopped", zap.Error(err))
	if !errors.Is(err, http.ErrServerClosed) {
		logger.Error(ctx, "http server error", zap.Error(err))
	}
}
