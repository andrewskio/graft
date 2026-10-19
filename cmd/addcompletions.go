package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var addCompletionsCmd = &cobra.Command{
	Use:     "add-completions",
	Aliases: []string{"ac"},
	Short:   "Install Graft shell completions for the current user",
	Long: `Generate and install Graft command completions for the shell reported by $SHELL.
Supports Bash, Zsh, and Fish, and writes the completion script to the current user's
standard completion directory. For Zsh, ensure ~/.zsh/completions is included in fpath.

Use 'graft completion' to generate completion scripts manually for other shells.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		shell := os.Getenv("SHELL")
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		rootCmd := cmd.Root()

		switch {
		case strings.Contains(shell, "zsh"):
			zshDir := filepath.Join(home, ".zsh/completions")
			_ = os.MkdirAll(zshDir, 0755)
			targetPath := filepath.Join(zshDir, "_"+rootCmd.Name())

			// Generate and write completion file
			err := rootCmd.GenZshCompletionFile(targetPath)
			if err != nil {
				return err
			}
			fmt.Printf("Zsh completions installed to %s\n", targetPath)
			fmt.Println("Ensure fpath=(~/.zsh/completions $fpath) is in your ~/.zshrc")

		case strings.Contains(shell, "bash"):
			bashDir := filepath.Join(home, ".local/share/bash-completion/completions")
			_ = os.MkdirAll(bashDir, 0755)
			targetPath := filepath.Join(bashDir, rootCmd.Name())

			err := rootCmd.GenBashCompletionFileV2(targetPath, true)
			if err != nil {
				return err
			}
			fmt.Printf("Bash completions installed to %s\n", targetPath)

		case strings.Contains(shell, "fish"):
			fishDir := filepath.Join(home, ".config/fish/completions")
			_ = os.MkdirAll(fishDir, 0755)
			targetPath := filepath.Join(fishDir, rootCmd.Name()+".fish")

			err := rootCmd.GenFishCompletionFile(targetPath, true)
			if err != nil {
				return err
			}
			fmt.Printf("Fish completions installed to %s\n", targetPath)

		default:
			return fmt.Errorf("unsupported shell: %s. Please install manually using '%s completion'", shell, rootCmd.Name())
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCompletionsCmd)
}
