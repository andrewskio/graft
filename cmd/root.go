package cmd

import (
	"context"
	"os"

	"charm.land/fang/v2"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "graft",
	Short: "Simple Minecraft server management CLI built around docker-compose",
	Long:  `Simple Minecraft server management CLI built around docker-compose`,
}

func Execute() {
	err := fang.Execute(context.Background(), rootCmd)
	if err != nil {
		os.Exit(1)
	}
}

func init() {

}
