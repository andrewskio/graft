package cmd

import (
	"fmt"

	"graft/internal/serverstore"

	"github.com/spf13/cobra"
)

// dirCmd represents the dir command
var dirCmd = &cobra.Command{
	Use:   "dir",
	Short: "Print the directory where Graft stores server configurations",
	Long: `Print the default root path to Graft's server directory.
Server configurations and related server files are stored under this directory.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := serverstore.DefaultRoot()
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
