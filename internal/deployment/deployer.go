package deployment

type Deployer interface {
	Deploy(deployment Deployment) error
	Remove(deploymentName string) error
}
