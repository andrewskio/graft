package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create new Minecraft server",
	Long:  `Create new Minecraft server`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("new", args[0])
	},
}

func init() {
	rootCmd.AddCommand(newCmd)

}
