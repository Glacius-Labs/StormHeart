package runtime

import "github.com/glacius-labs/StormHeart/internal/model"

type Runtime interface {
	Deploy(deployment model.Deployment) error
	Remove(deployment model.Deployment) error
	List() ([]model.Deployment, error)
}
