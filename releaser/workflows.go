package releaser

import (
	"context"
	"encoding"
	"fmt"
	"strings"
)

type ApplicationList struct {
	Application []Application `json:"applications"`
}

func (a *ApplicationList) MarshalText() (text []byte, err error) {
	var ret strings.Builder
	for _, app := range a.Application {
		for _, rel := range app.ReleaseCandidate {
			ret.WriteString(fmt.Sprintf("%s %s %s %d\n", app.Name, rel.Name, rel.Status, rel.ExistingPR))
		}
	}
	return []byte(ret.String()), nil
}

var _ encoding.TextMarshaler = &ApplicationList{}

type ReleaseCandidate struct {
	Name       string                 `json:"name"`
	Status     ReleaseCandidateStatus `json:"status"`
	ExistingPR int64                  `json:"existing_pr"`
}

type Application struct {
	Name             string
	ReleaseCandidate []ReleaseCandidate `json:"releaseCandidates"`
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

func GetAllReleaseStatus(ctx context.Context, a Api) (*ApplicationList, error) {
	apps, err := a.ListApplications()
	if err != nil {
		return nil, fmt.Errorf("failed to get application list: %w", err)
	}
	var ret ApplicationList
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
				app.ReleaseCandidate = append(app.ReleaseCandidate, ReleaseCandidate{
					Name:   release,
					Status: RC_STATUS_RELEASED,
				})
				continue
			}
			old, newRelease, err := a.PreviewRelease(app.Name, release)
			if err != nil {
				return nil, fmt.Errorf("failed to get preview for %s:%s: %w", app.Name, release, err)
			}
			hasChange := old.Yaml() != newRelease.Yaml()
			rc := ReleaseCandidate{
				Name:   release,
				Status: getStatus(hasChange),
			}
			if rc.Status == RC_STATUS_PENDING {
				prNum, err := CheckForPRForRelease(ctx, a, app.Name, release)
				if err != nil {
					return nil, fmt.Errorf("failed to check for PR for %s:%s: %w", app.Name, release, err)
				}
				rc.ExistingPR = prNum
			}
			app.ReleaseCandidate = append(app.ReleaseCandidate, rc)
		}
		ret.Application = append(ret.Application, app)
	}
	return &ret, nil
}

func getStatus(change bool) ReleaseCandidateStatus {
	if change {
		return RC_STATUS_PENDING
	}
	return RC_STATUS_RELEASED
}
