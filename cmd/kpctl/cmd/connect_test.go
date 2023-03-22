package cmd

import (
	"testing"

	"github.com/kube-peering/internal/kpctl"
	"github.com/kube-peering/internal/pkg/model"
	"github.com/kube-peering/internal/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestConnectCommand(t *testing.T) {
	connectCmd.Run = util.MockRun

	c, out, err := util.ExecuteCommandC(rootCmd, []string{"connect"}...)

	assert.Equal(t, "connect", c.Name())
	assert.Empty(t, out)
	assert.NoError(t, err)
	assert.Equal(t,
		&kpctl.Kpctl{
			Backdoor:    model.DefaultBackdoor,
			Application: model.CreateApplication("localhost", 8080),
		},
		instance,
	)
}
