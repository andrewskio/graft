package dockerservice

import "github.com/docker/compose/v5/pkg/api"

/*
Container contains both a ComposeService and ComposeProject and has
methods to perform container operations (start, stop, etc.).
It is created with LoadServerProject() method on ComposeService.
*/
type Container struct {
	*ComposeService
	project *ComposeProject
}

/*
Start the container
*/
func (c Container) Start() error {
	return c.service.Up(c.ctx, c.project, api.UpOptions{
		Create: api.CreateOptions{},
		Start:  api.StartOptions{},
	})

}

/*
Stop the container
*/
func (c Container) Stop() error {
	return c.service.Stop(c.ctx, c.project.Name, api.StopOptions{})
}

func (c Container) Status() ([]api.ContainerSummary, error) {
	summary, err := c.service.Ps(c.ctx, c.project.Name, api.PsOptions{
		All: true,
	})
	if err != nil {
		return []api.ContainerSummary{}, err
	}
	return summary, nil
}
