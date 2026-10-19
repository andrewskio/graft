// Package server defines the configuration for a managed Minecraft server within the application.
package server

// Server is the internal representation for a Minecraft server's configuration.
type Server struct {
	ServerName       string
	ServerType       string
	MinecraftVersion string
	JavaVersion      string
	Memory           string
}
