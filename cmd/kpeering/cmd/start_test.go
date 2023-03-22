package cmd

import (
	"testing"

	"github.com/kube-peering/internal/kpeering"
	"github.com/kube-peering/internal/pkg/model"
	"github.com/kube-peering/internal/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestStartCommand(t *testing.T) {
	startCmd.Run = util.MockRun

	c, out, err := util.ExecuteCommandC(rootCmd, []string{"start"}...)

	assert.Equal(t, "start", c.Name())
	assert.Empty(t, out)
	assert.NoError(t, err)
	assert.Equal(t,
		&kpeering.Kpeering{
			Frontdoor: model.DefaultFrontdoor,
			Backdoor:  model.DefaultBackdoor,
		},
		instance,
	)
}
