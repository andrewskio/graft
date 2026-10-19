package cmd

import (
	"graft/internal/cli/render"
	"graft/internal/server"
	"graft/internal/serverstore"

	"github.com/docker/compose/v5/pkg/api"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "List saved Minecraft servers and their configuration details",
	Long: `List every server configuration in the Graft server directory, including its
server type, Minecraft version, Java version, and container status. Servers that
have not been started are listed with that status. The alias 'list' can also be used.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		service, err := newComposeService(cmd)
		if err != nil {
			return err
		}

		store, err := serverstore.NewDefault()
		if err != nil {
			return err
		}
		names, err := store.ListNames()
		if err != nil {
			return err
		}

		type serverStatus struct {
			status api.ContainerSummary
			server server.Server
		}

		serverStatuses := []serverStatus{}

		for _, srvName := range names {
			container, err := service.GetServerProject(srvName)
			if err != nil {
				return err
			}

			summaries, err := container.Status()
			if err != nil {
				return err
			}

			if len(summaries) == 0 {
				serverStatuses = append(serverStatuses, serverStatus{
					status: api.ContainerSummary{
						Status: "Server has never been started",
					},
					server: container.Server,
				})
			} else {
				serverStatuses = append(serverStatuses, serverStatus{
					status: summaries[0],
					server: container.Server,
				})
			}
		}

		cmd.Println(render.Table(
			[]string{"Name", "Type", "Version", "Java", "Container Status"},
			serverStatuses,
			func(s serverStatus) []string {
				return []string{s.server.ServerName, s.server.ServerType, s.server.MinecraftVersion, s.server.JavaVersion, s.status.Status}
			},
		))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
