// Package composefile defines the Compose YAML representation and its
// conversion to and from the application server model.
package composefile

import (
	"fmt"
	"log/slog"
	"strings"

	"graft/internal/server"

	"gopkg.in/yaml.v3"
)

const DockerImage = "itzg/minecraft-server"

type Environment struct {
	EULA    string `yaml:"EULA"`
	TYPE    string `yaml:"TYPE"`
	VERSION string `yaml:"VERSION"`
	MEMORY  string `yaml:"MEMORY"`
}

type Service struct {
	Image       string      `yaml:"image"`
	Tty         bool        `yaml:"tty"`
	StdinOpen   bool        `yaml:"stdin_open"`
	Ports       []string    `yaml:"ports"`
	Environment Environment `yaml:"environment"`
	Volumes     []string    `yaml:"volumes"`
}

// Compose represents a parsed compose.yaml file.
type Compose struct {
	Name     string             `yaml:"name"`
	Services map[string]Service `yaml:"services"`
}

// Marshal encodes Compose to Compose YAML
func (compose Compose) Marshal() ([]byte, error) {
	composeYaml, err := yaml.Marshal(&compose)
	if err != nil {
		return []byte{}, err
	}

	slog.Debug("--- Generated Compose ---")
	slog.Debug(string(composeYaml))

	return composeYaml, nil
}

// Parse decodes Compose YAML.
func Parse(composeFile []byte) (Compose, error) {
	compose := Compose{}
	err := yaml.Unmarshal(composeFile, &compose)
	if err != nil {
		return compose, err
	}

	return compose, nil
}

// FromServer projects the server configuration into Compose YAML.
func FromServer(s server.Server) Compose {
	service := Service{
		Image:     fmt.Sprintf("%s:%s", DockerImage, s.JavaVersion),
		Tty:       true,
		StdinOpen: true,
		Ports:     []string{"25565:25565"},
		Volumes:   []string{"./data:/data"},
		Environment: Environment{
			EULA:    "TRUE",
			TYPE:    strings.ToUpper(s.ServerType),
			VERSION: s.MinecraftVersion,
			MEMORY:  strings.ToUpper(s.Memory),
		},
	}
	return Compose{Name: s.ServerName, Services: map[string]Service{"mc": service}}
}

// ToServer maps a single service Minecraft Compose project to the Server model.
func (compose Compose) ToServer() (server.Server, error) {
	if len(compose.Services) != 1 {
		return server.Server{}, fmt.Errorf("expected exactly 1 service in compose, found %d", len(compose.Services))
	}
	var service Service
	for _, candidate := range compose.Services {
		service = candidate
	}
	javaVersion := ""
	if _, tag, found := strings.Cut(service.Image, ":"); found {
		javaVersion = tag
	}
	return server.Server{
		ServerName:       compose.Name,
		ServerType:       strings.ToLower(service.Environment.TYPE),
		MinecraftVersion: service.Environment.VERSION,
		JavaVersion:      javaVersion,
		Memory:           service.Environment.MEMORY,
	}, nil
}
