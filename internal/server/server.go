/*
The server package defines the Server type, the internal representation of a Server.
A Server can be created from a Compose (see FromCompose helper), which is loaded from the server's compose.yml.
*/
package server

import (
	"graft/internal/compose"
	"graft/internal/graftfs"

	"fmt"
	"graft/internal/config"
	"strings"
)

// ServerConf contains configuration for a single Minecraft server.
// Used with the ToCompose and FromCompose method to convert to and from Compose YAML.
type Server struct {
	ServerName       string
	ServerType       string
	MinecraftVersion string
	JavaVersion      string
	// recommended java is ui only, should be refractored out of here at some point
	// used in ServerConf so state change properly updates form
	RecommendedJava bool
	Memory          string
}

// ToString() returns the values of ServerConf in a formatted, readable string,
// used primarily for users view/confirm server configuration.
func (s Server) ToString() string {
	str := fmt.Sprintf(`
	Server Name: %s
	Server Type: %s
	Minecraft Version: %s
	Java Version: %s
	Memory: %s
	`, s.ServerName, s.ServerType, s.MinecraftVersion, s.JavaVersion, s.Memory)
	return strings.TrimLeft(str, "\r\n")
}

// GetPath() method returns the absolute path of Server directory
func (s Server) GetPath() (string, error) {
	return graftfs.GetServerPath(s.ServerName)
}

// GetComposePath() method returns the absolute path of Server's compose.yml file
func (s Server) GetComposePath() (string, error) {
	return graftfs.GetServerComposePath(s.ServerName)
}

// NewCompose generates a Compose from a ServerConf.
// Compose is the type used to marshalling/unmarshalling compose YAML
func (s Server) ToCompose() compose.Compose {
	service := compose.Service{
		Image:     fmt.Sprintf("%s:%s", config.DockerImage, s.JavaVersion),
		Tty:       true,
		StdinOpen: true,
		Ports:     []string{"25565:25565"},
		Volumes:   []string{"./data:/data"},
	}
	service.Environment.EULA = "TRUE"
	service.Environment.TYPE = strings.ToUpper(s.ServerType)
	service.Environment.VERSION = s.MinecraftVersion
	service.Environment.MEMORY = strings.ToUpper(s.Memory)

	return compose.Compose{
		Name: s.ServerName,
		Services: map[string]compose.Service{
			"mc": service,
		},
	}
}
