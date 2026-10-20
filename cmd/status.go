package cmd

import (
	"graft/internal/cli/render"
	"graft/internal/service/compose"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status <server>",
	Short: "Show the container status for a Minecraft server",
	Long: `Display the Docker Compose container name, status, and state for the named server.
The server name must match a configuration in the Graft server directory. If the
server has not been started, Graft reports that no container status is available.`,
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

		summaryList, err := project.Status()
		if err != nil {
			return err
		}

		if len(summaryList) == 0 {
			cmd.Println(render.Warning("%s found, but has no status. Has the server been started?", srvName))
			return nil
		}

		cmd.Println(render.Table(
			[]string{"Name", "Status", "State"},
			summaryList,
			func(s compose.ContainerStatus) []string {
				return []string{s.Project, s.Status, string(s.State)}
			},
		))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
