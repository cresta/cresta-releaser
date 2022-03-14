package releaser

import (
	"context"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

func TestCurrentBranchName(t *testing.T) {
	g := GitCli{
		Logger: zap.NewNop(),
	}
	branchName, err := g.CurrentBranchName(context.Background())
	require.NoError(t, err)
	require.NotEqual(t, "", branchName)
}
