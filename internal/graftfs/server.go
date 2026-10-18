package graftfs

import "path/filepath"

// GetServerPath(srvName) returns the absolute path for a server's directory, or an error if the home directory could not be resolved.
// It does not check if path exists or create the directory.
func GetServerPath(srvName string) (string, error) {
	graftDir, err := GetGraftDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(graftDir, "servers", srvName), nil
}

// GetServerComposePath(srvName) returns the absolute path for a server's compose.yml, or an error if the home directory could not be resolved.
// It does not check if path exists or create the directory.
func GetServerComposePath(srvName string) (string, error) {
	serverPath, err := GetServerPath(srvName)
	if err != nil {
		return "", nil
	}
	return filepath.Join(serverPath, "compose.yml"), nil
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
