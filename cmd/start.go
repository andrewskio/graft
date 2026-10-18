package cmd

import (
	"graft/internal/cli/render"
	"graft/internal/dockerservice"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start server",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		srvName := args[0]

		service, err := dockerservice.NewComposeService(cmd.Context(), cmd.OutOrStdout(), cmd.OutOrStderr())
		if err != nil {
			return err
		}

		project, err := service.LoadServerProject(srvName)
		if err != nil {
			return err
		}

		cmd.Println(render.Success("Executing Docker start\n"))

		err = service.Start(project)
		if err != nil {
			return err
		}

		cmd.Println(render.Success("\nService started: %s", srvName))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
