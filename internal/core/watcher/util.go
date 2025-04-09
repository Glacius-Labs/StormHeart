package watcher

import (
	"context"
	"time"

	"github.com/glacius-labs/StormHeart/internal/core/model"
)

func PushEmptyDeployments(handlerFunc HandlerFunc, sourceName string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	handlerFunc(ctx, sourceName, []model.Deployment{})
}
