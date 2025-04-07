package main

import "github.com/glacius-labs/StormHeart/internal/model"

var staticDeployments = []model.Deployment{
	{
		Name:  "stormlink",
		Image: "eclipse-mosquitto:2",
		Labels: map[string]string{
			"system": "true",
			"role":   "broker",
		},
		Environment: map[string]string{},
	},
}
