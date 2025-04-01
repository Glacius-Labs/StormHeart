package provider

import "github.com/glacius-labs/StormHeart/internal/deployment/model"

type PushFunc func(name string, deployments []model.Deployment)
