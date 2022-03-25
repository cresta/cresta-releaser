package main

import (
	"context"
	// mage:import go
	_ "github.com/cresta/magehelper/gobuild"
	"github.com/cresta/magehelper/pipe"
)

func BuildTwirp(ctx context.Context) error {
	return pipe.Shell("protoc --proto_path=. --go_out=. --go_opt=module=github.com/cresta/cresta-releaser --twirp_out=. --twirp_opt=module=github.com/cresta/cresta-releaser rpc/releaser/Releaser.proto").Run(ctx)
}

func InstallTwirpDeps(ctx context.Context) error {
	if err := pipe.Shell("go install github.com/twitchtv/twirp/protoc-gen-twirp").Run(ctx); err != nil {
		return err
	}
	if err := pipe.Shell("go install google.golang.org/protobuf/cmd/protoc-gen-go").Run(ctx); err != nil {
		return err
	}
	return nil
}
