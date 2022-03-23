package releaser

import (
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
			ret.WriteString(fmt.Sprintf("%s %s %s\n", app.Name, rel.Name, rel.Status))
		}
	}
	return []byte(ret.String()), nil
}

var _ encoding.TextMarshaler = &ApplicationList{}

type ReleaseCandidate struct {
	Name   string                 `json:"name"`
	Status ReleaseCandidateStatus `json:"status"`
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

func GetAllPendingReleases(a Api) (*ApplicationList, error) {
	all, err := GetAllReleaseStatus(a)
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

func GetAllReleaseStatus(a Api) (*ApplicationList, error) {
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
			app.ReleaseCandidate = append(app.ReleaseCandidate, ReleaseCandidate{
				Name:   release,
				Status: getStatus(hasChange),
			})
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
