package releaser

import (
	"context"
	"encoding"
	"fmt"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"
)

type ApplicationList struct {
	Application []Application `json:"applications"`
}

func (a *ApplicationList) MarshalText() (text []byte, err error) {
	var ret strings.Builder
	for _, app := range a.Application {
		for _, rel := range app.ReleaseCandidate {
			txt, err := rel.MarshalText()
			if err != nil {
				return nil, fmt.Errorf("failed to marshal release candidate: %w", err)
			}
			if _, err := fmt.Fprintf(&ret, "%s %s\n", app.Name, txt); err != nil {
				return nil, fmt.Errorf("failed to write to string: %w", err)
			}
		}
	}
	return []byte(ret.String()), nil
}

var _ encoding.TextMarshaler = &ApplicationList{}

type ReleaseCandidate struct {
	Name        string                 `json:"name"`
	Status      ReleaseCandidateStatus `json:"status"`
	ExistingPR  int64                  `json:"existing_pr"`
	OriginalSHA string                 `json:"original_sha"`
	Age         time.Duration          `json:"age"`
}

func (r *ReleaseCandidate) MarshalText() (text []byte, err error) {
	return []byte(fmt.Sprintf("%s %s %d %s %s", r.Name, r.Status, r.ExistingPR, r.OriginalSHA, r.Age)), nil
}

type Application struct {
	Name             string
	ReleaseCandidate []*ReleaseCandidate `json:"releaseCandidates"`
}

type ReleaseCandidateStatus int

func (r ReleaseCandidateStatus) String() string {
	switch r {
	case RC_STATUS_PENDING:
		return "pending"
	case RC_STATUS_RELEASED:
		return "released"
	default:
		return "unknown"
	}
}

const (
	RC_STATUS_UNKNOWN ReleaseCandidateStatus = iota
	RC_STATUS_PENDING
	RC_STATUS_RELEASED
)

func GetAllPendingReleases(ctx context.Context, a Api) (*ApplicationList, error) {
	all, err := GetAllReleaseStatus(ctx, a)
	if err != nil {
		return nil, fmt.Errorf("failed to get all release status: %w", err)
	}
	var ret ApplicationList
	for _, app := range all.Application {
		newApp := Application{Name: app.Name}
		for _, rel := range app.ReleaseCandidate {
			if rel.Status == RC_STATUS_PENDING {
				newApp.ReleaseCandidate = append(newApp.ReleaseCandidate, rel)
			}
		}
		if len(newApp.ReleaseCandidate) > 0 {
			ret.Application = append(ret.Application, newApp)
		}
	}
	return &ret, nil
}

func FindKustomizationForRelease(a Api, application string, release string) (string, error) {
	rel, err := a.GetRelease(application, release)
	if err != nil {
		return "", fmt.Errorf("failed to get release %s/%s: %w", application, release, err)
	}
	var validKustomizationNames = []string{"kustomization.yaml", "kustomization.yml", "Kustomization"}
	for _, r := range rel.Files {
		if r.Directory != "." {
			continue
		}
		for _, n := range validKustomizationNames {
			if r.Name == n {
				return n, nil
			}
		}
	}
	return "", nil
}

func DoesApplicationExist(a Api, application string) (bool, error) {
	apps, err := a.ListApplications()
	if err != nil {
		return false, fmt.Errorf("failed to list applications: %w", err)
	}
	for _, app := range apps {
		if app == application {
			return true, nil
		}
	}
	return false, nil
}

func NeedsPromotion(ctx context.Context, a Api, application string, release string) (bool, error) {
	old, newRelease, err := a.PreviewRelease(ctx, application, release, true)
	if err != nil {
		return false, fmt.Errorf("failed to get preview for %s:%s: %w", application, release, err)
	}
	hasChange := old.Yaml() != newRelease.Yaml()
	return hasChange, nil
}

func GetAllReleaseStatus(ctx context.Context, a Api) (*ApplicationList, error) {
	apps, err := a.ListApplications()
	if err != nil {
		return nil, fmt.Errorf("failed to get application list: %w", err)
	}
	var ret ApplicationList
	eg, egCtx := errgroup.WithContext(ctx)
	for _, app := range apps {
		releases, err := a.ListReleases(app)
		if err != nil {
			return nil, fmt.Errorf("failed to get release list for %s: %w", app, err)
		}
		app := Application{
			Name: app,
		}
		for idx, release := range releases {
			if idx == 0 {
				app.ReleaseCandidate = append(app.ReleaseCandidate, &ReleaseCandidate{
					Name:   release,
					Status: RC_STATUS_RELEASED,
				})
				continue
			}
			hasChange, err := NeedsPromotion(egCtx, a, app.Name, release)
			if err != nil {
				return nil, fmt.Errorf("failed to get preview for %s:%s: %w", app.Name, release, err)
			}
			existingRelease, err := a.GetRelease(app.Name, release)
			if err != nil {
				return nil, fmt.Errorf("failed to get preview for %s:%s: %w", app.Name, release, err)
			}
			releaseConfig, err := existingRelease.loadReleaseConfig()
			if err != nil {
				return nil, fmt.Errorf("failed to load release config for %s:%s: %w", app.Name, release, err)
			}
			rc := &ReleaseCandidate{
				Name:        release,
				Status:      getStatus(hasChange),
				OriginalSHA: releaseConfig.Metadata.OriginalRelease.GitSha,
			}
			if rc.Status == RC_STATUS_PENDING {
				a := a
				release := release
				rc := rc
				eg.Go(func() error {
					prNum, err := CheckForPRForRelease(egCtx, a, app.Name, release)
					if err != nil {
						return fmt.Errorf("failed to check for PR for %s:%s: %w", app.Name, release, err)
					}
					rc.ExistingPR = prNum
					return nil
				})
			}
			app.ReleaseCandidate = append(app.ReleaseCandidate, rc)
		}
		ret.Application = append(ret.Application, app)
	}
	if err := eg.Wait(); err != nil {
		return nil, fmt.Errorf("failed to wait for all PRs: %w", err)
	}
	return &ret, nil
}

func getStatus(change bool) ReleaseCandidateStatus {
	if change {
		return RC_STATUS_PENDING
	}
	return RC_STATUS_RELEASED
}
