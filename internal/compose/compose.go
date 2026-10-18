/*
The compose package contains the Compose type, and helpers to
Marshal and Unmarshal to and from a compose.yml. The Compose type
internally represents a compose.yml file, and should be converted to and from
a Server type for use. (see server.FromCompose(compose))
*/
package compose

import (
	"log/slog"

	"gopkg.in/yaml.v3"
)

type Service struct {
	Image       string   `yaml:"image"`
	Tty         bool     `yaml:"tty"`
	StdinOpen   bool     `yaml:"stdin_open"`
	Ports       []string `yaml:"ports"`
	Environment struct {
		EULA    string `yaml:"EULA"`
		TYPE    string `yaml:"TYPE"`
		VERSION string `yaml:"VERSION"`
		MEMORY  string `yaml:"MEMORY"`
	} `yaml:"environment"`
	Volumes []string `yaml:"volumes"`
}

// Compose is the primary struct that represents a compose.yaml file.
// Used for unmarshalling/marshalling to YAML
type Compose struct {
	Name     string             `yaml:"name"`
	Services map[string]Service `yaml:"services"`
}

/*
Marshal method returns a formatted compose.yaml []byte from a Compose.
Compose can be generated from ServerConf.NewCompose() method.
*/
func (compose Compose) Marshal() ([]byte, error) {
	composeYaml, err := yaml.Marshal(&compose)
	if err != nil {
		return []byte{}, err
	}

	slog.Debug("--- Generated Compose ---")
	slog.Debug(string(composeYaml))

	return composeYaml, nil
}

// Unmarshal([]byte) takes []byte from a read compose.yml and
// returns a Compose
func Unmarshal(composeFile []byte) (Compose, error) {
	compose := Compose{}
	err := yaml.Unmarshal(composeFile, &compose)
	if err != nil {
		return compose, err
	}

	return compose, nil
}
