package composer

import (
	"log/slog"

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

/*
CreateComposeBytes method returns a formatted compose.yaml []byte.
Compose can be generated from ServerConf.NewCompose() method.
*/
func (compose Compose) CreateComposeBytes() ([]byte, error) {
	composeYaml, err := yaml.Marshal(&compose)
	if err != nil {
		return []byte{}, err
	}

	slog.Debug("--- Generated Compose ---")
	slog.Debug(string(composeYaml))

	return composeYaml, nil
}
