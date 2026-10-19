// Package serverstore owns filesystem operations and persistence of server
// configurations.
package serverstore

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"graft/internal/util"
)

const serverDirectory = "servers"
const composeFilename = "compose.yml"

// Store persists server configurations under directory root.
type Store struct {
	root string
}

// New creates a store rooted at root. Call Ensure to create its server directory.
func New(root string) *Store {
	return &Store{root: root}
}

// DefaultRoot returns Graft's default per-user data directory.
func DefaultRoot() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".local", "share", "graft"), nil
}

// NewDefault creates a store using the current user's default Graft directory.
func NewDefault() (*Store, error) {
	root, err := DefaultRoot()
	if err != nil {
		return nil, err
	}
	return New(root), nil
}

// Root returns the configured store root.
func (s *Store) Root() string { return s.root }

// Ensure creates the server storage directory if it does not exist.
func (s *Store) Ensure() error {
	return os.MkdirAll(filepath.Join(s.root, serverDirectory), 0755)
}

// ServerPath returns the path for a server name, if the server name is valid
func (s *Store) ServerPath(name string) (string, error) {
	if err := util.ValidateServerName(name); err != nil {
		return "", err
	}
	return filepath.Join(s.root, serverDirectory, name), nil
}

// ComposePath returns the Compose file path for a server name.
func (s *Store) ComposePath(name string) (string, error) {
	serverPath, err := s.ServerPath(name)
	if err != nil {
		return "", err
	}
	return filepath.Join(serverPath, composeFilename), nil
}

// ValidateNewName returns nil if name is valid and is not already in use.
func (s *Store) ValidateNewName(name string) error {
	if err := util.ValidateServerName(name); err != nil {
		return err
	}
	serverPath, err := s.ServerPath(name)
	if err != nil {
		return err
	}
	_, err = os.Stat(serverPath)
	if err == nil {
		return fmt.Errorf("server name already in use")
	}
	if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

// SaveCompose writes Compose bytes for server name. Existing files are preserved
// unless allowOverwrite is true.
func (s *Store) SaveCompose(name string, data []byte, allowOverwrite bool) (string, error) {
	composePath, err := s.ComposePath(name)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Dir(composePath), 0755); err != nil {
		return "", err
	}
	if !allowOverwrite {
		if _, err := os.Stat(composePath); err == nil {
			return "", fmt.Errorf("%s already exists and overwriting is disabled", composePath)
		} else if !errors.Is(err, os.ErrNotExist) {
			return "", err
		}
	}
	if err := os.WriteFile(composePath, data, 0644); err != nil {
		return "", err
	}
	return composePath, nil
}

// LoadCompose reads a server's Compose file without parsing it
func (s *Store) LoadCompose(name string) ([]byte, error) {
	composePath, err := s.ComposePath(name)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(composePath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// ListNames returns directories containing a Compose file
func (s *Store) ListNames() ([]string, error) {
	serversPath := filepath.Join(s.root, serverDirectory)
	entries, err := os.ReadDir(serversPath)
	if errors.Is(err, os.ErrNotExist) {
		return []string{}, nil
	}
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() || util.ValidateServerName(entry.Name()) != nil {
			continue
		}
		composePath := filepath.Join(serversPath, entry.Name(), composeFilename)
		info, err := os.Stat(composePath)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, err
		}
		if info.Mode().IsRegular() {
			names = append(names, entry.Name())
		}
	}
	return names, nil
}
