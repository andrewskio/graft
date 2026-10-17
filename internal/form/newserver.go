package form

import (
	"fmt"
	"graft/internal/util"
	"graft/pkg/fetch"
	"log/slog"
	"strings"

	"graft/internal/compose"
	"graft/internal/config"

	"charm.land/huh/v2"
)

func Run() (*compose.ServerConf, error) {
	fields := &compose.ServerConf{}
	fields.RecommendedJava = true

	slog.Debug("Starting NewServer huh form")

	mcManifest, err := fetch.FetchMinecraftManifest()
	if err != nil {
		return nil, err
	}
	mcVersions := mcManifest.GetStrList()

	javaManifest, err := fetch.FetchJavaManifest()
	if err != nil {
		return nil, err
	}
	javaVersions := javaManifest.GetStrList()

	err = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Server Name").
				Description("This will identify the folder name and container name. Must be 2-25 characters and alphanumeric, hyphen, and underscore only.").
				Value(&fields.ServerName).
				Validate(func(str string) error {
					return util.ValidateSlug(str, "Server name")
				}),
			huh.NewInput().
				Title("Minecraft Version").
				Description("Enter version (e.g., LATEST, SNAPSHOT, 26.2, 1.21.1)").
				Value(&fields.MinecraftVersion).
				Validate(func(str string) error {
					upperVal := strings.ToUpper(strings.TrimSpace(str))
					for _, v := range mcVersions {
						if strings.ToUpper(v) == upperVal {
							return nil
						}
					}
					return fmt.Errorf("Version '%s' not found.", str)
				}),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Server Type").
				Description("Type of server. For vanilla servers, Paper is recommended for performance and supports Bukkit/Spigot plugins.").
				Value(&fields.ServerType).
				Options(huh.NewOptions(config.GetServerTypes()...)...),
			huh.NewInput().
				Title("Memory").
				Description("Amount of RAM to assign to server. (1G, 500MB, etc.)").
				Value(&fields.Memory).
				Validate(func(s string) error {
					return util.ValidateJVMMemory(s)
				}),
		),

		huh.NewGroup(

			huh.NewConfirm().
				Title("Recommended Java").
				Description("Show only recommened Java versions").
				Value(&fields.RecommendedJava),
			huh.NewSelect[string]().
				Title("Java Version").
				Description("Select java version.").
				Value(&fields.JavaVersion).
				OptionsFunc(func() []huh.Option[string] {
					slog.Debug("Running Java version selection OptionsFunc", "McVersion", fields.MinecraftVersion)
					if !fields.RecommendedJava {
						return huh.NewOptions(javaVersions...)
					}

					versionDetails, err := mcManifest.FetchVersionDetails(fields.MinecraftVersion)
					if err != nil {
						return huh.NewOptions(javaVersions...)
					}
					matchedVersions := versionDetails.GetJavaVersions(javaManifest)
					if len(matchedVersions) == 0 {
						return huh.NewOptions(javaVersions...)
					}

					return huh.NewOptions(matchedVersions...)
				}, &fields),
		),
	).WithShowHelp(true).Run()
	if err != nil {
		return nil, err
	}
	return fields, nil
}
