/*
The dockerservice package provides helpers that wrap around the Docker
and Docker Compose SDKs to control containers. Use `NewComposeService(ctx)` and
`service.LoadProject(composePath)` to use helper methods such as
`service.Start(project)`.
*/
package dockerservice

import (
	"context"
	"graft/internal/server"
	"io"

	"github.com/compose-spec/compose-go/v2/types"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/flags"
	"github.com/docker/compose/v5/cmd/display"
	"github.com/docker/compose/v5/pkg/api"
	"github.com/docker/compose/v5/pkg/compose"
)

type ComposeService struct {
	ctx     context.Context
	service api.Compose
}

type ComposeProject = types.Project

/*
NewComposeService initializes a ComposeService from the root context (cobra cmd.Context())
and writes output it stdout and stderr (cobra cmd.OutOrStdout(), cmd.OutOrStderr())
*/
func NewComposeService(ctx context.Context, stdout io.Writer, stderr io.Writer) (*ComposeService, error) {
	dockerCLI, err := command.NewDockerCli()
	if err != nil {
		return nil, err
	}

	err = dockerCLI.Initialize(&flags.ClientOptions{})
	if err != nil {
		return nil, err
	}

	service, err := compose.NewComposeService(
		dockerCLI,
		compose.WithPrompt(compose.AlwaysOkPrompt()),
		compose.WithOutputStream(stdout),
		compose.WithErrorStream(stderr),
		compose.WithEventProcessor(display.Full(stdout, stderr, false)),
	)
	if err != nil {
		return nil, err
	}

	return &ComposeService{
		ctx:     ctx,
		service: service,
	}, nil
}

// LoadServerProject finds the compose for srvName and returns a ComposeProject.
// Returns non-nil error if server does not exist or could not be loaded.
func (c ComposeService) LoadServerProject(srvName string) (*ComposeProject, error) {
	srv, err := server.LoadServer(srvName)
	if err != nil {
		return nil, err
	}

	compPath, err := srv.GetComposePath()
	if err != nil {
		return nil, err
	}

	project, err := c.service.LoadProject(c.ctx, api.ProjectLoadOptions{
		ConfigPaths: []string{compPath},
	})
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (c ComposeService) Start(project *ComposeProject) error {
	return c.service.Up(c.ctx, project, api.UpOptions{
		Create: api.CreateOptions{},
		Start:  api.StartOptions{},
	})

}

func (c ComposeService) Stop(project *ComposeProject) error {
	return c.service.Stop(c.ctx, project.Name, api.StopOptions{})
}
