package cmd

import (
	"testing"

	"github.com/kube-peering/internal/kpeering"
	"github.com/kube-peering/internal/pkg/config"
	"github.com/kube-peering/internal/pkg/model"
	util_test "github.com/kube-peering/internal/pkg/util/test"
	"github.com/stretchr/testify/assert"
)

func TestStartCommand(t *testing.T) {
	startCmd.Run = util_test.MockRun

	c, out, err := util_test.ExecuteCommandC(rootCmd, []string{"start"}...)

	assert.Equal(t, "start", c.Name())
	assert.Empty(t, out)
	assert.NoError(t, err)
	assert.Equal(t,
		&kpeering.Kpeering{
			Frontdoor: model.Frontdoor{
				Endpoint: model.Endpoint{
					Name:       "frontdoor",
					Protocol:   "tcp",
					Host:       "localhost",
					ListenPort: config.DefautlFrontdoorPort,
				},
			},
			Backdoor: model.DefaultBackdoor,
		},
		instance,
	)
}
