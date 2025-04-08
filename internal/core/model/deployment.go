package model

import (
	"errors"
)

type Deployment struct {
	Name         string
	Image        string
	Labels       map[string]string
	Environment  map[string]string
	PortMappings []PortMapping
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

	if opts.PortMappings == nil {
		opts.PortMappings = []PortMapping{}
	}

	return Deployment{
		Name:         name,
		Image:        image,
		Labels:       opts.Labels,
		Environment:  opts.Environment,
		PortMappings: opts.PortMappings,
	}, nil
}

func (d Deployment) Equals(other Deployment) bool {
	if d.Name != other.Name || d.Image != other.Image {
		return false
	}

	if len(d.Labels) != len(other.Labels) || len(d.Environment) != len(other.Environment) || len(d.PortMappings) != len(other.PortMappings) {
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

	for i := range d.PortMappings {
		if !d.PortMappings[i].Equals(other.PortMappings[i]) {
			return false
		}
	}

	return true
}
