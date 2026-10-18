package server

import (
	"graft/internal/compose"
	"graft/internal/graftfs"

	"fmt"
	"graft/internal/config"
	"strings"
)

// EnumerateServerConfs() fetches all found servers, unmarshals the Compose, and returns
// the server configurations
func EnumerateServerConfs() ([]ServerConf, error) {
	serverConfs := []ServerConf{}
	names, err := graftfs.EnumerateServerNames()
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

// ServerConf contains configuration for a single Minecraft server.
// Used with the ToCompose and FromCompose method to convert to and from Compose YAML.
type ServerConf struct {
	ServerName       string
	ServerType       string
	MinecraftVersion string
	JavaVersion      string
	RecommendedJava  bool
	Memory           string
}

// ToString() returns the values of ServerConf in a formatted, readable string,
// used primarily for users view/confirm server configuration.
func (s ServerConf) ToString() string {
	str := fmt.Sprintf(`
	Server Name: %s
	Server Type: %s
	Minecraft Version: %s
	Java Version: %s
	Memory: %s
	`, s.ServerName, s.ServerType, s.MinecraftVersion, s.JavaVersion, s.Memory)
	return strings.TrimLeft(str, "\r\n")
}

// NewCompose generates a Compose from a ServerConf.
// Compose is the type used to marshalling/unmarshalling compose YAML
func (srvConf ServerConf) ToCompose() compose.Compose {
	service := compose.Service{
		Image:     fmt.Sprintf("%s:%s", config.DockerImage, srvConf.JavaVersion),
		Tty:       true,
		StdinOpen: true,
		Ports:     []string{"25565:25565"},
		Volumes:   []string{"./data:/data"},
	}
	service.Environment.EULA = "TRUE"
	service.Environment.TYPE = strings.ToUpper(srvConf.ServerType)
	service.Environment.VERSION = srvConf.MinecraftVersion
	service.Environment.MEMORY = strings.ToUpper(srvConf.Memory)

	return compose.Compose{
		Services: map[string]compose.Service{
			srvConf.ServerName: service,
		},
	}
}

// FromCompose() returns a ServerConf constructed from an unmarshalled Compose
func FromCompose(comp compose.Compose) (ServerConf, error) {
	if len(comp.Services) != 1 {
		return ServerConf{}, fmt.Errorf("expected exactly 1 service in compose, found %d", len(comp.Services))
	}

	var serverName string
	var service compose.Service
	for name, svc := range comp.Services {
		serverName = name
		service = svc
	}

	// Extract Java version from the image string (format: image:tag)
	parts := strings.Split(service.Image, ":")
	javaVersion := ""
	if len(parts) == 2 {
		javaVersion = parts[1]
	}

	return ServerConf{
		ServerName:       serverName,
		ServerType:       strings.ToLower(service.Environment.TYPE),
		MinecraftVersion: service.Environment.VERSION,
		JavaVersion:      javaVersion,
		RecommendedJava:  true, // Cannot be inferred, defaults to true
		Memory:           service.Environment.MEMORY,
	}, nil
}
