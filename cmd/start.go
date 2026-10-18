package cmd

import (
	"context"
	"fmt"
	"graft/internal/cli"
	"graft/internal/docker"
	"graft/internal/graftfs"
	"path/filepath"

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

		serverPathExists, err := graftfs.DoesServerPathExist(srvName)
		if err != nil {
			return err
		}
		if !serverPathExists {
			return fmt.Errorf("Server %s not found.", srvName)
		}

		srvPath, err := graftfs.GetServerPath(srvName)
		if err != nil {
			return err
		}

		ctx := context.Background()

		err = docker.StartService(ctx, filepath.Join(srvPath, "compose.yml"), srvName)
		if err != nil {
			return err
		}

		cli.PrintGreen("Service started: ", srvName)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
