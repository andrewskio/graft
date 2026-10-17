package compose

import (
	"fmt"
	"graft/internal/config"
	"graft/internal/util"
	"strings"

	"gopkg.in/yaml.v3"
)

// Compose is the primary struct that represents a compose.yaml file.
// Used for unmarshalling/marshalling to YAML
type Compose struct {
	Services map[string]Service `yaml:"services"`
}

// Service is the type used by Compose to represent a service in compose.yaml
type Service struct {
	Image       string            `yaml:"image"`
	PullPolicy  string            `yaml:"pull_policy"`
	Tty         bool              `yaml:"tty"`
	StdinOpen   bool              `yaml:"stdin_open"`
	Ports       []string          `yaml:"ports"`
	Environment map[string]string `yaml:"environment"`
	Volumes     []string          `yaml:"volumes"`
}

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

/*
CreateComposeFile takes a Compose (generated from ServerConf.NewCompose() method)
and prints a formatted compose.yaml file.

This is a WIP and will soon handle writing the file.
*/
func CreateComposeFile(compose Compose) error {
	composeYaml, err := yaml.Marshal(&compose)
	if err != nil {
		return err
	}

	util.PrintGreen("--- Generated Compose ---")
	fmt.Println(string(composeYaml))

	return nil
}
