package main

// Structs Compose and Service for constructing compose ymls
type Compose struct {
	Services map[string]Service `yaml:"services"`
}

type Service struct {
	Image       string            `yaml:"image"`
	PullPolicy  string            `yaml:"pull_policy"`
	Tty         bool              `yaml:"tty"`
	StdinOpen   bool              `yaml:"stdin_open"`
	Ports       []string          `yaml:"ports"`
	Environment map[string]string `yaml:"environment"`
	Volumes     []string          `yaml:"volumes"`
}
