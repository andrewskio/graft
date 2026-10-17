package config

const (
	JavaManifestURL = "https://raw.githubusercontent.com/itzg/docker-minecraft-server/refs/heads/master/images.json"
	McManifestURL   = "https://launchermeta.mojang.com/mc/game/version_manifest.json"
	DockerImage     = "itzg/minecraft-server"
)

// would love to find a way to programmatically fetch these from minecraft-docker-compose (TODO)
func GetServerTypes() []string {
	return []string{"Paper", "Fabric", "Neoforge", "Forge", "Vanilla"}
}
