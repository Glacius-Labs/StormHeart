package model

type DeploymentRepository interface {
	GetAll() []Deployment
	Get(source string) (Deployment, bool)
	Add(deployment Deployment)
	Remove(source string)
	RemoveAll()
}
