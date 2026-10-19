package cmd

import (
	"context"
	"graft/internal/serverstore"
	"log/slog"
	"os"
	"path/filepath"

	"charm.land/fang/v2"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "graft",
	Short: "Manage Minecraft server configurations and containers with Docker Compose",
	Long: `Graft creates and manages Minecraft server projects stored in its server directory.
Use the commands below to create a server configuration, start or stop a server,
inspect server status, list saved servers, or attach to a running server console.

Server containers are managed through Docker Compose. Use --log to write diagnostic
logs to debug.log in the Graft server directory.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Create graft dir if non existent
		store, err := serverstore.NewDefault()
		if err != nil {
			return err
		}
		if err := store.Ensure(); err != nil {
			return err
		}

		// Init logging if enabled
		if Log {
			logFilePath := filepath.Join(store.Root(), "debug.log")
			file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
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
