package deployment

import (
	"errors"
)

type Deployment struct {
	Name        string
	Image       string
	Labels      map[string]string
	Environment map[string]string
}

func NewDeployment(name, image string, opts DeploymentOptions) (Deployment, error) {
	if name == "" {
		return Deployment{}, errors.New("deployment name cannot be empty")
	}
	if image == "" {
		return Deployment{}, errors.New("deployment image cannot be empty")
	}

	if opts.Labels == nil {
		opts.Labels = map[string]string{}
	}

	if opts.Environment == nil {
		opts.Environment = map[string]string{}
	}

	return Deployment{
		Name:        name,
		Image:       image,
		Labels:      opts.Labels,
		Environment: opts.Environment,
	}, nil
}

func (d Deployment) Equals(other Deployment) bool {
	if d.Name != other.Name || d.Image != other.Image {
		return false
	}

	if len(d.Labels) != len(other.Labels) || len(d.Environment) != len(other.Environment) {
		return false
	}

	for k, v := range d.Labels {
		if other.Labels[k] != v {
			return false
		}
	}

	for k, v := range d.Environment {
		if other.Environment[k] != v {
			return false
		}
	}

	return true
}
