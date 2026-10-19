package cmd

import (
	"graft/internal/cli/render"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:     "start <server>",
	Aliases: []string{"up"},

	Short: "Start a saved Minecraft server",
	Long: `Start the Docker Compose project for the named server.
The server name must match a configuration in the Graft server directory.
The 'up' alias can be used instead of 'start'.`,
	Args: cobra.ExactArgs(1),
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

		cmd.Println(render.Success("Executing Compose project up...\n"))

		err = project.Start()
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
