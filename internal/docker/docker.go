package docker

import (
	"context"

	"github.com/compose-spec/compose-go/v2/types"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/flags"
	"github.com/docker/compose/v5/pkg/api"
	"github.com/docker/compose/v5/pkg/compose"
)

func initComposeService(ctx context.Context, composePath string, projectName string) (api.Compose, *types.Project, error) {
	dockerCLI, err := command.NewDockerCli()
	if err != nil {
		return nil, nil, err
	}

	err = dockerCLI.Initialize(&flags.ClientOptions{})
	if err != nil {
		return nil, nil, err
	}

	service, err := compose.NewComposeService(dockerCLI)
	if err != nil {
		return nil, nil, err
	}

	project, err := service.LoadProject(ctx, api.ProjectLoadOptions{
		ConfigPaths: []string{composePath},
		ProjectName: projectName,
	})
	if err != nil {
		return nil, nil, err
	}

	return service, project, nil
}

func StartService(ctx context.Context, composePath string, projectName string) error {
	service, project, err := initComposeService(ctx, composePath, projectName)
	if err != nil {
		return err
	}

	return service.Up(ctx, project, api.UpOptions{
		Create: api.CreateOptions{},
		Start:  api.StartOptions{},
	})

}

func StopService(ctx context.Context, composePath string, projectName string) error {
	service, project, err := initComposeService(ctx, composePath, projectName)
	if err != nil {
		return err
	}

	return service.Stop(ctx, project.Name, api.StopOptions{})
}
