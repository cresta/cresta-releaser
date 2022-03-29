package main

import (
	"context"

	"github.com/cresta/magehelper/docker/registry"
	"github.com/cresta/magehelper/docker/registry/ghcr"
	"github.com/cresta/magehelper/env"

	_ "github.com/cresta/magehelper/cicd/githubactions"
	// mage:import go
	_ "github.com/cresta/magehelper/gobuild"
	// mage:import docker
	_ "github.com/cresta/magehelper/docker"
	// mage:import ghcr
	_ "github.com/cresta/magehelper/docker/registry/ghcr"
	"github.com/cresta/magehelper/pipe"
)

func init() {
	// Install ECR as my registry
	registry.Instance = ghcr.Instance
	env.Default("DOCKER_MUTABLE_TAGS", "true")
}

func BuildTwirp(ctx context.Context) error { //nolint:golint,deadcode
	return pipe.Shell("protoc --proto_path=. --go_out=. --go_opt=module=github.com/cresta/cresta-releaser --twirp_out=. --twirp_opt=module=github.com/cresta/cresta-releaser rpc/releaser/Releaser.proto").Run(ctx)
}

func InstallTwirpDeps(ctx context.Context) error { //nolint:golint,deadcode
	if err := pipe.Shell("go install github.com/twitchtv/twirp/protoc-gen-twirp").Run(ctx); err != nil {
		return err
	}
	if err := pipe.Shell("go install google.golang.org/protobuf/cmd/protoc-gen-go").Run(ctx); err != nil {
		return err
	}
	return nil
}
