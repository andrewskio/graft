package sys

import (
	"errors"
	"fmt"
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
// Returns a non-nil error if the path could not be read or already exists.
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

// CreateGraftDir() creates the graft directory if it does not exist,
// returning a non-nil error if any errors occured.
func CreateGraftDir() error {
	graftDir, err := GetGraftDir()
	if err != nil {
		return err
	}
	return createNewDir(graftDir)
}

// GetServerPath(srvName) returns the absolute path for a server's directory, or an error if the home directory could not be resolved.
// It does not check if path exists or create the directory.
func GetServerPath(srvName string) (string, error) {
	graftDir, err := GetGraftDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(graftDir, srvName), nil
}

// DoesServerPath(srvName) checks if a path exists for given server name.
// It will return false and a non-nil error if any errors occured.
func DoesServerPathExist(srvName string) (bool, error) {
	serverPath, err := GetServerPath(srvName)
	if err != nil {
		return false, err
	}
	pathExists, err := doesPathExist(serverPath)
	if err != nil {
		return false, err
	}
	return pathExists, nil
}

// CreateServerPath(srvName) creates a directory for the given server name if the path does not already exist.
// Returns the server path, or a non-nil error if any issues occur.
func CreateServerDir(srvName string) (string, error) {
	serverPath, err := GetServerPath(srvName)
	if err != nil {
		return "", err
	}
	pathExists, err := doesPathExist(serverPath)
	if err != nil {
		return "", err
	}
	if pathExists {
		return serverPath, nil
	}

	err = createNewDir(serverPath)
	if err != nil {
		return "", err
	}
	return serverPath, nil
}

// WriteServerCompose creates a server path (if non-existent) writes the yamlData to compose.yaml for a server srvName.
// Returns the path of the compose file and an error.
// When allowOverwrite is false, it will return a non-nil error if the file already exists.
func WriteServerCompose(srvName string, yamlData []byte, allowOverwrite bool) (string, error) {
	srvDir, err := CreateServerDir(srvName)
	if err != nil {
		return "", err
	}

	composePath := filepath.Join(srvDir, "compose.yml")

	if !allowOverwrite {
		pathExists, err := doesPathExist(composePath)
		if err != nil {
			return "", err
		}
		if pathExists {
			return "", fmt.Errorf("%s already exists and overwriting is disabled", composePath)
		}
	}

	err = os.WriteFile(composePath, yamlData, 0644)
	if err != nil {
		return "", err
	}
	return composePath, nil

}
