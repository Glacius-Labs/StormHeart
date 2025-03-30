package deployment

import "context"

type DeploymentProvider interface {
	Start(ctx context.Context) error
}
