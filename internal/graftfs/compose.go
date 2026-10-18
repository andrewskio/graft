package graftfs

import (
	"fmt"
	"os"
	"path/filepath"
)

// WriteServerCompose creates a server path (if non-existent) writes the yamlData to compose.yaml for a server srvName.
// Returns the path of the compose file and an error.
// When allowOverwrite is false, it will return a non-nil error if the file already exists.
func WriteServerCompose(srvName string, yamlData []byte, allowOverwrite bool) (string, error) {
	_, err := CreateServerDir(srvName)
	if err != nil {
		return "", err
	}

	composePath, err := GetServerComposePath(srvName)
	if err != nil {
		return "", err
	}

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

// ReadServerCompose(srvName) returns a []byte containing the contents of a server's compose.yml
// Returns non-nil error if server not found or compose could not be loaded
func ReadServerCompose(srvName string) ([]byte, error) {
	composePath, err := GetServerComposePath(srvName)
	if err != nil {
		return []byte{}, err
	}

	exists, err := doesPathExist(composePath)
	if err != nil {
		return []byte{}, err
	}
	if !exists {
		return []byte{}, fmt.Errorf("compose.yml not found or could not be read: %s", composePath)
	}

	data, err := os.ReadFile(composePath)
	if err != nil {
		return []byte{}, err
	}

	return data, nil
}

// DoesServerComposeExist(srvName) checks if a compose.yml path exists for given server name.
// It will return false and a non-nil error if any errors occured.
func DoesServerComposeExist(srvName string) (bool, error) {
	compPath, err := GetServerComposePath(srvName)
	if err != nil {
		return false, err
	}
	pathExists, err := doesPathExist(compPath)
	if err != nil {
		return false, err
	}
	return pathExists, nil
}

// EnumerateComposePaths() finds all server paths within the Graft directory
// and returns the names of servers with a compose.yml.
func EnumerateComposePaths() ([]string, error) {
	serverNames := []string{}

	rootDir, err := GetGraftDir()
	if err != nil {
		return serverNames, err
	}

	entries, err := os.ReadDir(filepath.Join(rootDir, "servers"))
	if err != nil {
		return serverNames, err
	}
	if len(entries) == 0 {
		return serverNames, nil
	}

	for _, entry := range entries {
		compPath, err := GetServerComposePath(entry.Name())
		if err != nil {
			return serverNames, err
		}
		composeExists, err := doesPathExist(compPath)
		if err != nil {
			return serverNames, err
		}
		if composeExists {
			serverNames = append(serverNames, entry.Name())
		}
	}

	return serverNames, nil
}
