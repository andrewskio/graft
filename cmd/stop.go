package cmd

import (
	"graft/internal/cli/render"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop <server>",
	Short: "Stop a running Minecraft server",
	Long: `Stop the Docker Compose project for the named server.
The server name must match a configuration in the Graft server directory.
The 'down' alias can be used instead of 'stop'.`,
	Aliases: []string{"down"},

	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		srvName := args[0]

		service, err := newComposeService(cmd)
		if err != nil {
			return err
		}

		cmd.Println(render.Success("Executing Compose project down...\n"))

		project, err := service.GetServerProject(srvName)
		if err != nil {
			return err
		}

		err = project.Stop()
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
