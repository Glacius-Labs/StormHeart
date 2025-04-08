package main

import "github.com/glacius-labs/StormHeart/internal/core/model"

var staticDeployments = []model.Deployment{
	{
		Name:  "stormlink",
		Image: "eclipse-mosquitto:2",
		Labels: map[string]string{
			"system": "true",
			"role":   "broker",
		},
		Environment: map[string]string{},
		PortMappings: []model.PortMapping{
			{
				HostPort:      1883,
				ContainerPort: 1883,
			},
		},
	},
}
