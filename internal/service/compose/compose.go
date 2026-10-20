// Package compose adapts Docker Compose operations to Graft's server store.
package compose

import (
	"context"
	"fmt"
	"io"

	"graft/internal/composefile"
	"graft/internal/server"
	"graft/internal/serverstore"

	"github.com/compose-spec/compose-go/v2/types"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/flags"
	"github.com/docker/compose/v5/cmd/display"
	"github.com/docker/compose/v5/pkg/api"
	composeapi "github.com/docker/compose/v5/pkg/compose"
)

type ComposeProject = types.Project

// Service wraps the Docker Compose SDK for Graft's server operations.
type Service struct {
	ctx   context.Context
	api   api.Compose
	store *serverstore.Store
}

// New creates a Compose service using the context, root store, and streams
func New(ctx context.Context, store *serverstore.Store, stdout, stderr io.Writer) (*Service, error) {
	dockerCLI, err := command.NewDockerCli()
	if err != nil {
		return nil, err
	}
	if err := dockerCLI.Initialize(&flags.ClientOptions{}); err != nil {
		return nil, err
	}
	composeService, err := composeapi.NewComposeService(
		dockerCLI,
		composeapi.WithPrompt(composeapi.AlwaysOkPrompt()),
		composeapi.WithOutputStream(stdout),
		composeapi.WithErrorStream(stderr),
		composeapi.WithEventProcessor(display.Full(stdout, stderr, false)),
	)
	if err != nil {
		return nil, err
	}
	return &Service{ctx: ctx, api: composeService, store: store}, nil
}

// GetServerProject loads the server model and its Compose project.
func (s *Service) GetServerProject(name string) (*Project, error) {
	data, err := s.store.LoadCompose(name)
	if err != nil {
		return nil, err
	}
	compose, err := composefile.Parse(data)
	if err != nil {
		return nil, err
	}
	configuration, err := compose.ToServer()
	if err != nil {
		return nil, err
	}
	composePath, err := s.store.ComposePath(name)
	if err != nil {
		return nil, err
	}
	project, err := s.api.LoadProject(s.ctx, api.ProjectLoadOptions{ConfigPaths: []string{composePath}})
	if err != nil {
		return nil, err
	}
	return &Project{service: s, project: project, name: name, Server: configuration}, nil
}

// Project is the Compose SDK project loaded for a managed server.
type Project struct {
	service *Service
	project *ComposeProject
	name    string
	Server  server.Server
}

type ContainerStatus = api.ContainerSummary

// ServerName returns the managed server name.
func (p *Project) ServerName() string { return p.name }

// Start starts the project.
func (p *Project) Start() error {
	return p.service.api.Up(p.service.ctx, p.project, api.UpOptions{Create: api.CreateOptions{}, Start: api.StartOptions{}})
}

// Stop stops the project.
func (p *Project) Stop() error {
	return p.service.api.Stop(p.service.ctx, p.project.Name, api.StopOptions{})
}

// Status returns the project's container summaries.
func (p *Project) Status() ([]ContainerStatus, error) {
	return p.service.api.Ps(p.service.ctx, p.project.Name, api.PsOptions{All: true})
}

func (p *Project) Logs(ctx context.Context, logConsumer api.LogConsumer) error {
	return p.service.api.Logs(ctx, p.project.Name, logConsumer, api.LogOptions{
		Follow: true,
		Tail:   "1000",
	})
}

// Attach attaches to the first service's console.
func (p *Project) Attach() error {
	services := p.project.ServiceNames()
	if len(services) == 0 {
		return fmt.Errorf("compose project %q has no services", p.project.Name)
	}
	return p.service.api.Attach(p.service.ctx, p.project.Name, api.AttachOptions{
		Service: services[0], DetachKeys: "ctrl-c", Proxy: false,
	})
}
