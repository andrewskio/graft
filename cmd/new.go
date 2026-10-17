package cmd

import (
	"graft/internal/compose"
	"graft/internal/form"

	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create new Minecraft server",
	Long:  `Create new Minecraft server`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		serverConf, err := form.Run()
		if err != nil {
			return err
		}
		newCompose := serverConf.NewCompose()
		return compose.CreateComposeFile(newCompose)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)

}
