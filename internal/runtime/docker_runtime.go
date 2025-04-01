package runtime

import (
	"context"
	"fmt"
	"strings"

	"maps"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/glacius-labs/StormHeart/internal/deployment/model"
)

const (
	LabelOwner      string = "owner"
	OwnerStormHeart string = "stormheart"
)

type DockerRuntime struct {
	cli *client.Client
}

func NewDockerRuntime() (*DockerRuntime, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {
		return nil, err
	}

	return &DockerRuntime{cli: cli}, nil
}

func (r *DockerRuntime) Deploy(deployment model.Deployment) error {
	ctx := context.Background()

	_ = r.cli.ContainerRemove(ctx, deployment.Name, container.RemoveOptions{Force: true})

	_, err := r.cli.ImagePull(ctx, deployment.Image, image.PullOptions{})

	if err != nil {
		return fmt.Errorf("failed to pull image: %w", err)
	}

	resp, err := r.cli.ContainerCreate(ctx, &container.Config{
		Image:  deployment.Image,
		Labels: generateRuntimeLabels(deployment),
		Env:    generateRuntimeEnvironment(deployment.Environment),
	}, nil, nil, nil, deployment.Name)

	if err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}

	if err := r.cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	return nil
}

func (r *DockerRuntime) Remove(deployment model.Deployment) error {
	ctx := context.Background()
	return r.cli.ContainerRemove(ctx, deployment.Name, container.RemoveOptions{Force: true})
}

func (r *DockerRuntime) List() ([]model.Deployment, error) {
	ctx := context.Background()

	containers, err := r.cli.ContainerList(ctx, container.ListOptions{
		Filters: filters.NewArgs(
			filters.Arg(LabelOwner, OwnerStormHeart),
		),
	})

	if err != nil {
		return nil, err
	}

	var result []model.Deployment

	for _, container := range containers {
		name := strings.TrimPrefix(container.Names[0], "/")

		inspect, err := r.cli.ContainerInspect(ctx, container.ID)
		if err != nil {
			continue
		}

		env := parseRuntimeEnvironment(inspect.Config.Env)

		result = append(result, model.Deployment{
			Name:        name,
			Image:       container.Image,
			Labels:      filterUserLabels(container.Labels),
			Environment: env,
		})
	}

	return result, nil
}

func generateRuntimeLabels(deployment model.Deployment) map[string]string {
	labels := make(map[string]string)

	maps.Copy(labels, deployment.Labels)

	labels[LabelOwner] = OwnerStormHeart
	return labels
}

func filterUserLabels(all map[string]string) map[string]string {
	userLabels := make(map[string]string)

	for k, v := range all {
		if k == LabelOwner {
			continue
		}
		userLabels[k] = v
	}

	return userLabels
}

func generateRuntimeEnvironment(envMap map[string]string) []string {
	env := make([]string, 0, len(envMap))
	for k, v := range envMap {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	return env
}

func parseRuntimeEnvironment(envList []string) map[string]string {
	result := make(map[string]string)
	for _, entry := range envList {
		parts := strings.SplitN(entry, "=", 2)
		if len(parts) == 2 {
			result[parts[0]] = parts[1]
		}
	}
	return result
}
