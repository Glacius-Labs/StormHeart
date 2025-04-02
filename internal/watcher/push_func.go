package watcher

import "github.com/glacius-labs/StormHeart/internal/model"

type PushFunc func(name string, deployments []model.Deployment)
