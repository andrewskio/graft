package cmd

import (
	"graft/internal/form"
	"graft/internal/util"

	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create new Minecraft server",
	Long:  `Create new Minecraft server`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fields, err := form.Run()
		if err != nil {
			util.PrintRed(err.Error())
		}
		util.PrintGreen(fields.JavaVersion)

	},
}

func init() {
	rootCmd.AddCommand(newCmd)

}
