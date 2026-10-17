package composer

import (
	"fmt"
	"graft/internal/config"
	"strings"
)

// ServerConf contains configuration for a single Minecraft server.
// Used with the NewCompose method to generate Compose YAML.
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
func (srvConf ServerConf) NewCompose() Compose {
	service := Service{
		Image:      fmt.Sprintf("%s:%s", config.DockerImage, srvConf.JavaVersion),
		PullPolicy: "daily",
		Tty:        true,
		StdinOpen:  true,
		Ports:      []string{"25565:25565"},
		Environment: map[string]string{
			"EULA":    "TRUE",
			"TYPE":    strings.ToUpper(srvConf.ServerType),
			"VERSION": srvConf.MinecraftVersion,
			"MEMORY":  strings.ToUpper(srvConf.Memory),
		},
		Volumes: []string{"./data:/data"},
	}

	return Compose{map[string]Service{
		srvConf.ServerName: service,
	}}
}
