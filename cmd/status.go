package cmd

import (
	"graft/internal/cli/render"
	"graft/internal/dockerservice"

	"github.com/docker/compose/v5/pkg/api"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Print the status of given server.",
	Long:  `Print the status of given server.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		srvName := args[0]

		service, err := dockerservice.NewComposeService(cmd.Context(), cmd.OutOrStdout(), cmd.OutOrStderr())
		if err != nil {
			return err
		}

		container, err := service.GetServerContainer(srvName)
		if err != nil {
			return err
		}

		summaryList, err := container.Status()
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
			func(s api.ContainerSummary) []string {
				return []string{s.Project, s.Status, string(s.State)}
			},
		))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
