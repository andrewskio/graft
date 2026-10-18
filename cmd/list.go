package cmd

import (
	"graft/internal/cli/render"
	"graft/internal/server"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "List all servers found in Graft directory",
	Long:    `List all servers found in Graft directory`,
	RunE: func(cmd *cobra.Command, args []string) error {
		srvConfs, err := server.EnumerateServers()
		if err != nil {
			return err
		}

		cmd.Println(render.Table(
			[]string{"Name", "Type", "Version", "Java"},
			srvConfs,
			func(s server.Server) []string {
				return []string{s.ServerName, s.ServerType, s.MinecraftVersion, s.JavaVersion}
			},
		))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
