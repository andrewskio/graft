package cmd

import (
	"context"
	"graft/internal/sys"
	"log/slog"
	"os"

	"charm.land/fang/v2"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "graft",
	Short: "Simple Minecraft server management CLI built around docker-compose",
	Long:  `Simple Minecraft server management CLI built around docker-compose`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Create graft dir if non existent
		err := sys.CreateGraftDir()
		if err != nil {
			return err
		}

		// Init logging if enabled
		if Log {
			file, err := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err == nil {
				handler := slog.NewTextHandler(file, &slog.HandlerOptions{
					Level: slog.LevelDebug,
				})
				slog.SetDefault(slog.New(handler))
				slog.Info("LOGGING INITIALIZED")
			}
			return err
		}

		return nil
	},
}

var Log bool

func Execute() {
	err := fang.Execute(context.Background(), rootCmd)
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Log, "log", "l", false, "Enable logging to debug.log file in Graft directory")
}
