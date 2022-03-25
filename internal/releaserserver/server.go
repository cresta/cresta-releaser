package releaserserver

import (
	"context"
	"fmt"
	"github.com/cresta/cresta-releaser/internal/managedgitrepo"
	"github.com/cresta/cresta-releaser/releaser"
	releaser_protobuf "github.com/cresta/cresta-releaser/rpc/releaser"
	"github.com/cresta/zapctx"
	"go.uber.org/zap"
)

type Server struct {
	Logger *zap.Logger
	Api    releaser.Api
	Repo   *managedgitrepo.Repo
}

func NewServer(ctx context.Context, logger *zapctx.Logger, api releaser.Api, repo *managedgitrepo.Repo) (*Server, error) {
	zapLogger := logger.Unwrap(ctx)
	return &Server{
		Logger: zapLogger,
		Api:    api,
		Repo:   repo,
	}, nil
}

func (s *Server) GetAllApplicationStatus(ctx context.Context, _ *releaser_protobuf.GetAllApplicationStatusRequest) (*releaser_protobuf.GetAllApplicationStatusResponse, error) {
	if err := s.Repo.UpdateCheckout(ctx); err != nil {
		return nil, fmt.Errorf("failed to update checkout: %w", err)
	}
	if err := s.Repo.G.ResetClean(ctx); err != nil {
		return nil, fmt.Errorf("failed to reset git repo: %w", err)
	}
	if err := s.Repo.G.ResetToOriginalBranch(ctx); err != nil {
		return nil, fmt.Errorf("failed to reset to original branch: %w", err)
	}
	releaseList, err := releaser.GetAllReleaseStatus(ctx, s.Api)
	if err != nil {
		return nil, fmt.Errorf("failed to get all release status: %w", err)
	}
	var ret releaser_protobuf.GetAllApplicationStatusResponse
	for _, app := range releaseList.Application {
		appStatus := &releaser_protobuf.ApplicationStatus{
			Name: app.Name,
		}
		for _, rc := range app.ReleaseCandidate {
			appStatus.ReleaseStatus = append(appStatus.ReleaseStatus, &releaser_protobuf.ReleaseStatus{
				Name:     rc.Name,
				PrNumber: rc.ExistingPR,
				Status:   statusAsProto(rc.Status),
			})
		}
		ret.ApplicationStatus = append(ret.ApplicationStatus, appStatus)
	}
	return &ret, nil
}

func statusAsProto(status releaser.ReleaseCandidateStatus) releaser_protobuf.ReleaseStatus_Status {
	switch status {
	case releaser.RC_STATUS_PENDING:
		return releaser_protobuf.ReleaseStatus_PENDING
	case releaser.RC_STATUS_RELEASED:
		return releaser_protobuf.ReleaseStatus_RELEASED
	default:
		return releaser_protobuf.ReleaseStatus_UNKNOWN
	}
}

var _ releaser_protobuf.Releaser = &Server{}
