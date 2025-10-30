package commands

import (
	"errors"

	"github.com/homekit/homekit-cli/internal/core"
	"github.com/spf13/cobra"
)

func runtimeFrom(cmd *cobra.Command) (*core.Runtime, error) {
	rt, ok := core.FromContext(cmd.Context())
	if !ok {
		return nil, errors.New("runtime unavailable")
	}
	return rt, nil
}
