package server

import (
	"fmt"
	"graft/internal/compose"
	"graft/internal/graftfs"
	"strings"
)

// EnumerateServerConfs() fetches all found servers, unmarshals the Compose, and returns
// the Server
func EnumerateServers() ([]Server, error) {
	serverConfs := []Server{}
	names, err := graftfs.EnumerateComposePaths()
	if err != nil {
		return serverConfs, err
	}
	for _, name := range names {
		composeBytes, err := graftfs.ReadServerCompose(name)
		if err != nil {
			return serverConfs, err
		}
		newCompose, err := compose.Unmarshal(composeBytes)
		if err != nil {
			return serverConfs, err
		}
		serverConf, err := FromCompose(newCompose)
		if err != nil {
			return serverConfs, err
		}
		serverConfs = append(serverConfs, serverConf)
	}
	return serverConfs, nil
}

// LoadServer(srvName) loads the given server compose.yml, unmarshals it, and
// returns the Server created from the Compose.
// Returns non-nil error if server not found or unkown issues occured.
func LoadServer(srvName string) (Server, error) {
	compFile, err := graftfs.ReadServerCompose(srvName)
	if err != nil {
		return Server{}, err
	}

	comp, err := compose.Unmarshal(compFile)
	if err != nil {
		return Server{}, err
	}

	return FromCompose(comp)
}

// FromCompose() returns a Server constructed from an unmarshalled Compose
func FromCompose(comp compose.Compose) (Server, error) {
	if len(comp.Services) != 1 {
		return Server{}, fmt.Errorf("expected exactly 1 service in compose, found %d", len(comp.Services))
	}

	var service compose.Service
	for _, svc := range comp.Services {
		service = svc
	}

	// Extract Java version from the image string (image:tag)
	parts := strings.Split(service.Image, ":")
	javaVersion := ""
	if len(parts) == 2 {
		javaVersion = parts[1]
	}

	return Server{
		ServerName:       comp.Name,
		ServerType:       strings.ToLower(service.Environment.TYPE),
		MinecraftVersion: service.Environment.VERSION,
		JavaVersion:      javaVersion,
		RecommendedJava:  true, // Cannot be inferred, defaults to true
		Memory:           service.Environment.MEMORY,
	}, nil
}
