package releaserserver

import (
	"context"
	"fmt"
	"sync"

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
	mu     sync.RWMutex
}

func (s *Server) PushPromotion(ctx context.Context, request *releaser_protobuf.PushPromotionRequest) (*releaser_protobuf.PushPromotionResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	branchName := releaser.DefaultBranchNameForRelease(request.ApplicationName, request.ReleaseName)
	if pr, err := s.Api.CheckForPRForBranch(ctx, branchName); err != nil {
		return nil, fmt.Errorf("failed to check for existing PR for branch %s: %w", branchName, err)
	} else if pr != 0 {
		return &releaser_protobuf.PushPromotionResponse{
			Status:        releaser_protobuf.PushPromotionResponse_EXISTING_PULL_REQUEST,
			PullRequestId: pr,
		}, nil
	}
	if err := s.Repo.G.ResetToOriginalBranch(ctx); err != nil {
		return nil, fmt.Errorf("failed to reset to original branch: %w", err)
	}
	if exists, err := s.Repo.G.DoesBranchExist(ctx, branchName); err != nil {
		return nil, fmt.Errorf("failed to check if branch %s exists: %w", branchName, err)
	} else if exists {
		if err := s.Repo.G.ForceDeleteLocalBranch(ctx, branchName); err != nil {
			return nil, fmt.Errorf("failed to delete branch %s: %w", branchName, err)
		}
	}
	if err := s.Api.FreshGitBranch(ctx, request.ApplicationName, request.ReleaseName, ""); err != nil {
		return nil, fmt.Errorf("failed to create branch %s: %w", branchName, err)
	}
	oldRelease, newRelease, err := s.Api.PreviewRelease(request.ApplicationName, request.ReleaseName)
	if err != nil {
		return nil, fmt.Errorf("failed to preview release: %w", err)
	}
	if err := s.Api.ApplyRelease(request.ApplicationName, request.ReleaseName, oldRelease, newRelease); err != nil {
		return nil, fmt.Errorf("failed to apply release: %w", err)
	}
	if changes, err := s.Api.AreThereUncommittedChanges(ctx); err != nil {
		return nil, fmt.Errorf("failed to check for uncommitted changes: %w", err)
	} else if !changes {
		return &releaser_protobuf.PushPromotionResponse{
			Status: releaser_protobuf.PushPromotionResponse_NO_CHANGES,
		}, nil
	}
	if err := s.Api.CommitForRelease(ctx, request.ApplicationName, request.ReleaseName); err != nil {
		return nil, fmt.Errorf("failed to commit release: %w", err)
	}
	if err := s.Api.ForcePushCurrentBranch(ctx); err != nil {
		return nil, fmt.Errorf("failed to push release: %w", err)
	}
	if prNum, err := s.Api.PullRequestCurrent(ctx); err != nil {
		return nil, fmt.Errorf("failed to create pull request: %w", err)
	} else {
		return &releaser_protobuf.PushPromotionResponse{
			Status:        releaser_protobuf.PushPromotionResponse_NEW_PULL_REQUEST,
			PullRequestId: prNum,
		}, nil
	}
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
	s.mu.Lock()
	defer s.mu.Unlock()
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
