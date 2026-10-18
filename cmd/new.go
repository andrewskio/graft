package cmd

import (
	"graft/internal/cli"
	"graft/internal/graftfs"
	"graft/internal/manifest"

	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Launch TUI to create new Minecraft server",
	Long:  `Launch TUI to create new Minecraft server`,
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
		serverConf, err := cli.RunNewServerForm(mcManifest, javaManifest)
		if err != nil {
			return err
		}

		// composify from form configuration
		newCompose := serverConf.ToCompose()
		composeBytes, err := newCompose.Marshal()
		if err != nil {
			return err
		}
		composePath, err := graftfs.WriteServerCompose(serverConf.ServerName, composeBytes, false)
		if err != nil {
			return err
		}

		cli.PrintGreen("New server compose written to", composePath)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
