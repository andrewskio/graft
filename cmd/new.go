package cmd

import (
	"graft/internal/cli/render"
	"graft/internal/composefile"
	"graft/internal/manifest"
	"graft/internal/serverstore"
	"graft/internal/tui"

	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a Minecraft server configuration interactively",
	Long: `Open an interactive form to configure a new Minecraft server.
Available Minecraft versions, Java versions, and server types are loaded before
the form opens. After you complete it, Graft saves a Docker Compose configuration
in the Graft server directory; use 'graft start <server>' to launch it.`,
	Args: cobra.NoArgs,
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

		serverTypes := manifest.GetServerTypes()

		store, err := serverstore.NewDefault()
		if err != nil {
			return err
		}

		// launch new server tui form
		serverConf, err := tui.RunNewServerForm(mcManifest, javaManifest, serverTypes, store.ValidateNewName)
		if err != nil {
			return err
		}

		composeData, err := composefile.FromServer(*serverConf).Marshal()
		if err != nil {
			return err
		}
		composePath, err := store.SaveCompose(serverConf.ServerName, composeData, false)
		if err != nil {
			return err
		}

		cmd.Println(render.Success("New server compose written to %s", composePath))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
