package provider

import "context"

type FileDeploymentProvider struct {
	path string
}

func NewFileDeploymentProvider(path string) *FileDeploymentProvider {
	return &FileDeploymentProvider{
		path: path,
	}
}

func (p *FileDeploymentProvider) Start(ctx context.Context) error {
	panic("not implemented")
}
