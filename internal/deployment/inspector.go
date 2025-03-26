package deployment

type DeploymentInspector interface {
	ListActiveDeployments() ([]Deployment, error)
}
