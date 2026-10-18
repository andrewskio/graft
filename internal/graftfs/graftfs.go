/*
The graftfs package provides utility functions for getting paths, checking if
they exist, reading & writing compose files, and other filesystem action.
*/
package graftfs

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"
)

// doesPathExist checks if path (file/folder) already exists, error will be nil unless permission or other unkown error occured.
// Path is the absolute path, not relative to Graft dir.
func doesPathExist(path string) (bool, error) {
	slog.Debug("Checking path exists", "path", path)
	_, err := os.Stat(path)
	if err == nil {
		return true, nil // Path exists
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil // path doesn't exists
	}
	// Path might exist, permission error or something
	return false, err
}

// createNewDir creates a new directory at absolute path if it does not exist.
// Returns a non-nil error if the path could not be read.
func createNewDir(path string) error {
	slog.Debug("Creating new path", "path", path)
	pathExists, err := doesPathExist(path)
	if err != nil {
		return err
	}

	if pathExists {
		return nil
	}

	err = os.MkdirAll(path, 0755)
	if err != nil {
		slog.Error("Error creating new path", "path", path)
		return err
	}
	return nil
}

// GetGraftDir() returns the absolute path for Graft directory, or an error if the home directory could not be resolved.
// It does not check if path exists or create the directory.
func GetGraftDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".local", "share", "graft"), nil
}

// CreateGraftDir() creates the graft directory if it does not exist and returns the path,
// returning a non-nil error if any errors occured.
func CreateGraftDir() (string, error) {
	graftDir, err := GetGraftDir()
	if err != nil {
		return graftDir, err
	}
	return graftDir, createNewDir(graftDir)
}
