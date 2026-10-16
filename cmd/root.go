package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "graft",
	Short: "Simple Minecraft server management CLI built around docker-compose",
	Long:  `Simple Minecraft server management CLI built around docker-compose`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

}
