package build

import (
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
)

func (b *Build) ToMobyContainerConfig() *container.Config {
	return &container.Config{}
}

func (b *Build) ToMobyHostConfig() *container.HostConfig {
	return &container.HostConfig{}
}

func (b *Build) ToMobyNetworkConfig() *network.NetworkingConfig {
	return &network.NetworkingConfig{}
}
