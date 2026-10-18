package cmd

import (
	"graft/internal/server"

	"fmt"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "List all servers found in Graft directory",
	Long:    `List all servers found in Graft directory`,
	RunE: func(cmd *cobra.Command, args []string) error {
		srvConfs, err := server.EnumerateServerConfs()
		if err != nil {
			return err
		}

		t := table.New().
			Border(lipgloss.NormalBorder()).
			BorderRow(true).
			BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("238"))).
			Headers("Name", "Type", "Version", "Java")

		for _, c := range srvConfs {
			t.Row(c.ServerName, c.ServerType, c.MinecraftVersion, c.JavaVersion)
		}

		fmt.Println(t.Render())

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
