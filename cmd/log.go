package cmd

import (
	"context"
	"graft/internal/tui"

	"github.com/spf13/cobra"
)

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Prints log for given server",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		srvName := args[0]

		service, err := newComposeService(cmd)
		if err != nil {
			return err
		}

		project, err := service.GetServerProject(srvName)
		if err != nil {
			return err
		}

		return tui.RunLogViewer(cmd.Context(), "Logs: "+srvName, func(ctx context.Context, consumer tui.LogConsumer) error {
			return project.Logs(ctx, consumer)
		})

	},
}

func init() {
	rootCmd.AddCommand(logCmd)
}
