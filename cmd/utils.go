package cmd

import (
	"graft/internal/serverstore"
	"graft/internal/service/compose"

	"github.com/spf13/cobra"
)

func newComposeService(cmd *cobra.Command) (*compose.Service, error) {
	store, err := serverstore.NewDefault()
	if err != nil {
		return nil, err
	}
	return compose.New(cmd.Context(), store, cmd.OutOrStdout(), cmd.OutOrStderr())
}
