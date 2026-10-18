package cmd

import (
	"graft/internal/cli/render"
	"graft/internal/dockerservice"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop server",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		srvName := args[0]

		service, err := dockerservice.NewComposeService(cmd.Context(), cmd.OutOrStdout(), cmd.OutOrStderr())
		if err != nil {
			return err
		}

		cmd.Println(render.Success("Executing Docker stop\n"))

		project, err := service.LoadServerProject(srvName)
		if err != nil {
			return err
		}

		err = service.Stop(project)
		if err != nil {
			return err
		}

		cmd.Println(render.Success("\nService stopped: %s", srvName))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
