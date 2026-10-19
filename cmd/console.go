package cmd

import (
	"graft/internal/cli/render"

	"github.com/spf13/cobra"
)

var consoleCmd = &cobra.Command{
	Use:     "console <server>",
	Aliases: []string{"attach"},
	Short:   "Attach to a running Minecraft server console",
	Long: `Attach the current terminal to the console of the named running server.
Use Ctrl-C to detach from the console without issuing a server command. The server
name must match a configuration in the Graft server directory. The 'attach' alias
can be used instead of 'console'.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		srvName := args[0]

		service, err := newComposeService(cmd)
		if err != nil {
			return err
		}

		container, err := service.GetServerProject(srvName)
		if err != nil {
			return err
		}

		cmd.Println(render.Success("Attaching to console. Use ctrl-c to detach."))

		err = container.Attach()

		if err != nil {
			if err.Error() == "read escape sequence" {
				cmd.Println(render.Success("\n\nDetached from console."))
			} else {
				return err
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(consoleCmd)
}
