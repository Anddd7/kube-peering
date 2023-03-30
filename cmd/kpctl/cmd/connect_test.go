package cmd

import (
	"testing"

	"github.com/kube-peering/internal/kpctl"
	"github.com/kube-peering/internal/pkg/model"
	util_test "github.com/kube-peering/internal/pkg/util/test"
	"github.com/stretchr/testify/assert"
)

func TestConnectCommand(t *testing.T) {
	connectCmd.Run = util_test.MockRun

	c, out, err := util_test.ExecuteCommandC(rootCmd, []string{"connect", "-p", "8080"}...)

	assert.Equal(t, "connect", c.Name())
	assert.Empty(t, out)
	assert.NoError(t, err)
	assert.Equal(t,
		&kpctl.Kpctl{
			Tunnel:    model.DefaultTunnel,
			Forwarder: model.CreateForwarder("localhost", 8080),
		},
		instance,
	)
}
