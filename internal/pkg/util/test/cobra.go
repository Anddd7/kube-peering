package util

import (
	"bytes"

	"github.com/spf13/cobra"
)

// execute the root command with args, return the target command and output or error
func ExecuteCommandC(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	c, err = root.ExecuteC()

	return c, buf.String(), err
}

func MockRun(cmd *cobra.Command, args []string) {
	// empty
}
