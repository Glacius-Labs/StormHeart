package model

type PortMapping struct {
	HostPort      int
	ContainerPort int
}

func (pm PortMapping) Equals(other PortMapping) bool {
	return pm.HostPort == other.HostPort && pm.ContainerPort == other.ContainerPort
}
