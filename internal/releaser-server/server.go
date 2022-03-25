package releaser_server

import (
	"context"
	"fmt"
	"github.com/cresta/cresta-releaser/releaser"
	releaser_server "github.com/cresta/cresta-releaser/rpc/releaser-server"
	"github.com/cresta/zapctx"
	"go.uber.org/zap"
)

type Server struct {
	Logger *zap.Logger
	Api    releaser.Api
}

func NewServer(ctx context.Context, logger *zapctx.Logger) (*Server, error) {
	zapLogger := logger.Unwrap(ctx)
	api, err := releaser.NewFromCommandLine(ctx, zapLogger, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create releaser api: %w", err)
	}
	return &Server{
		Logger: zapLogger,
		Api:    api,
	}, nil
}

func (s *Server) GetAllApplicationStatus(ctx context.Context, _ *releaser_server.GetAllApplicationStatusRequest) (*releaser_server.GetAllApplicationStatusResponse, error) {
	releaseList, err := releaser.GetAllReleaseStatus(ctx, s.Api)
	if err != nil {
		return nil, fmt.Errorf("failed to get all release status: %w", err)
	}
	var ret releaser_server.GetAllApplicationStatusResponse
	for _, app := range releaseList.Application {
		appStatus := &releaser_server.ApplicationStatus{
			Name: app.Name,
		}
		for _, rc := range app.ReleaseCandidate {
			appStatus.ReleaseStatus = append(appStatus.ReleaseStatus, &releaser_server.ReleaseStatus{
				Name:     rc.Name,
				PrNumber: rc.ExistingPR,
				Status:   statusAsProto(rc.Status),
			})
		}
		ret.ApplicationStatus = append(ret.ApplicationStatus, appStatus)
	}
	return &ret, nil
}

func statusAsProto(status releaser.ReleaseCandidateStatus) releaser_server.ReleaseStatus_Status {
	switch status {
	case releaser.RC_STATUS_PENDING:
		return releaser_server.ReleaseStatus_PENDING
	case releaser.RC_STATUS_RELEASED:
		return releaser_server.ReleaseStatus_RELEASED
	default:
		return releaser_server.ReleaseStatus_UNKNOWN
	}
}

var _ releaser_server.ReleaserServer = &Server{}
