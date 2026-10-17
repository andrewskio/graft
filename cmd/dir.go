package cmd

import (
	"fmt"

	"graft/internal/sys"

	"github.com/spf13/cobra"
)

// dirCmd represents the dir command
var dirCmd = &cobra.Command{
	Use:   "dir",
	Short: "Simple utility command to print path that Graft saves servers.",
	Long:  "Simple utility command to print path that Graft saves servers.",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := sys.GetGraftDir()
		if err != nil {
			return err
		}
		fmt.Println(path)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(dirCmd)
}
