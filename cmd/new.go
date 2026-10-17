package cmd

import (
	"graft/internal/form"
	"graft/internal/manifest"
	"graft/internal/sys"
	"graft/internal/util"

	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create new Minecraft server",
	Long:  `Create new Minecraft server`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// fetch manifests
		javaManifest, err := manifest.FetchJavaManifest()
		if err != nil {
			return err
		}
		mcManifest, err := manifest.FetchMinecraftManifest()
		if err != nil {
			return err
		}

		// launch new server tui form
		serverConf, err := form.RunNewServerForm(mcManifest, javaManifest)
		if err != nil {
			return err
		}

		// composify from form configuration
		newCompose := serverConf.NewCompose()
		composeBytes, err := newCompose.CreateComposeBytes()
		if err != nil {
			return err
		}
		composePath, err := sys.WriteServerCompose(serverConf.ServerName, composeBytes, false)
		if err != nil {
			return err
		}

		util.PrintGreen("New server compose written to ", composePath)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
